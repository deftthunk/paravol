package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "reflect"
  "path/filepath"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "runtime"
  "strconv"
  "log"
  "regexp"
)

var pl = fmt.Println
var pf = fmt.Printf
var spf = fmt.Sprintf


/* grab YAML file and decode into structs */

func input() (map[string]string, []map[string]string, map[string]bool) {
  var cfgMap map[string]interface{}
  var options map[string]string
  var plugins []map[string]string
  var flags   map[string]bool

  flags = make(map[string]bool, 10)

  // count and parse CLI args
  if len(os.Args) < 2 {
    pl("\nUsage: paravol [OPTIONS] configfile\n")
    pl("Options:")
    pl("\t-p, --print-commands")
    pl("\t\tdoes not execute; instead prints out vol commands")

    os.Exit(1)
  } else {
    if os.Args[1] == "-p" || os.Args[1] == "--print-commands" {
      flags["print"] = true
    }
  }

  // open config
  cfgFile, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
  if err != nil {
    log.Fatal(err)
  }

  // parse yaml into map
  decode_err := yaml.Unmarshal(cfgFile, &cfgMap)
  if decode_err != nil {
    log.Fatalf("error: %v", decode_err)
  }

  // initialize options and plugins maps
  vOptions := reflect.ValueOf(cfgMap)
  options = make(map[string]string, vOptions.Len()-1)

/*  walk through unmarshaled map; data is a 
    map[string]interface{} with one member "plugins" 
    being of type []interface{}, which is an array of
    plugins in the form of map[interface{}]interface{}.
    each structure is type asserted (ta) into concrete type */

  for kOp, vOp := range cfgMap {
    // distinguish string values from arrays
    if ta, ok := vOp.(string); ok && ta != "" {
      options[kOp] = ta
    } else if ta, ok := vOp.([]interface{}); ok {
      // iterate over "Plugins:" sub-config array
      // each value is a map of plugin config values
      configArr := ta

      // array of interfaces
      for _, cMap := range configArr {
        obj := reflect.ValueOf(cMap)
        tmpMap := make(map[string]string, obj.Len() - 1)
        k := cMap.(map[interface{}]interface{})

        // map[interface]interface
        for kPlu, vPlu := range k {
          if vPlu == nil {
            tmpMap[kPlu.(string)] = ""
          } else {
            tmpMap[kPlu.(string)] = vPlu.(string)
          }
        }
        plugins = append(plugins, tmpMap)
      }
    } else if ta, ok := vOp.(int); ok {
        options[kOp] = strconv.Itoa(ta)
    } else if vOp == nil {
      switch {
        case kOp == "_filename":
          log.Fatalf("Error: must specify %s in Yaml config", kOp)
        case kOp == "_vol-name":
          log.Fatalf("Error: must specify %s in Yaml config", kOp)
        case kOp == "profile":
          log.Fatalf("Error: must specify %s in Yaml config", kOp)
        case kOp == "_subfolders":
          log.Fatalf("Error: must specify %s in Yaml config", kOp)
        default:
          options[kOp] = ""
      }
    } else {
      // catch anything else
      pl("what am i: ", reflect.TypeOf(vOp))
    }
  }

  return options, plugins, flags
}


/*  returns 2D array of dumpfiles in the format of 
    [[path, filename] [path, filename]] */

func findDumps(options map[string]string) ([][]string, []string) {
  // [[path, filename] [path, filename]]
  dumpFiles := make([][]string, 0)
  outputPath := make([]string, 0)
  var fullFilepath string
  var curFilename string
  var curParentDir string

  // crawler func for folders (used below)
  var walkFunc = func (path string, f os.FileInfo, err error) error {
    fullFilepath = path
    curFilename = filepath.Base(path)
    fi, err := os.Stat(path)

    if err != nil {
      log.Fatalf("os.Stat(): Error with path %v", path)
      pl("Error: ", err.Error())
      if os.IsNotExist(err) {
        pl("Folder does not exist")
      }
      os.Exit(1)
    }

    if fi.IsDir() {
      curParentDir = path
    } else {
      dumpFiles = append(dumpFiles, []string{curParentDir, fi.Name()})
    }

    return nil
  }

  // each file/dir is passed to walkFunc
  subFolders := strings.Fields(options["_subfolders"])

  for _, s := range subFolders {
    memDumpPath := filepath.Join(options["_memdumps"], s)
    err := filepath.Walk(memDumpPath, walkFunc)

    if err != nil {
      log.Fatalf("Error: %v", err)
    }

    // build (partial) output paths while we're at it
    out := filepath.Join(options["_output_folder"], s)
    outputPath = append(outputPath, out)
  }

  return dumpFiles, outputPath
}


/* add hyphens to option flags */

func fixField(f string) string {
  hyphens := "--" + f
  return hyphens
}


/*  build volatility command for each memory dump */

func buildCommands(dumpFiles [][]string, outputPath []string,
  opt map[string]string, plu []map[string]string,
  flags map[string]bool) [][]string {

  var commands [][]string

  for dIndex,pathArr := range dumpFiles {
    var optString []string

    // make plugins first argument if it exists
    if val,ok := opt["_comm_plugins"]; ok {
      optString = append(optString, spf("%s%s", "--plugins=", val))
    }
    optString = append(optString,
      spf("%s%s", "--filename=", strings.Join(pathArr, "/")),
    )

    // iterate throguh supplied vol options map in 'opt' slice
    for field,val := range opt {
      // catch and skip over script config params that start with an '_'
      match,err := regexp.MatchString("^_.*", field)
      if err != nil {
        log.Fatalf("error: %v", err)
      }
      if match { continue }
      hyField := fixField(field)
      var newStr string

      // handle empty string values
      if val != "" {
        str := []string{hyField, val}
        newStr = strings.Join(str, "=")
      } else {
        newStr = hyField
      }
      optString = append(optString, newStr)
    }

    // iterate through supplied plugin maps in 'plu' slice
    for _, pMap := range plu {
      var pluString []string
      var interString []string

      // finish output path using plugin name for destination folder
      fullOutPath := filepath.Join(outputPath[dIndex], pMap["plugin"])
      outStr := spf("%s%s", "--output-file=", fullOutPath)
      pluString = append(pluString, outStr)
      // ensure base path exists or create it
      if ! flags["print"] {
        os.MkdirAll(filepath.Dir(fullOutPath), 0755)
      }
      // create plugin string; place plugin name in front of its options
      pluString = append(pluString, pMap["plugin"])

      // add plugin options to command string
      for field, val := range pMap {
        // skip re-adding the plugin name
        if field == "plugin" { continue }
        hyField := fixField(field)
        var newStr string

        if val != "" {
          str := []string{hyField, val}
          newStr = strings.Join(str, "=")
        } else {
          newStr = hyField
        }

        // append each plugin option to pluString array
        pluString = append(pluString, newStr)
      }
      interString = append(interString, optString...)
      interString = append(interString, pluString...)

      commands = append(commands, interString)
    }
  }

  return commands
}


/* act as a Go thread to execute shell commands */

func manager(ch chan string, volPath string, cStr []string) {
  result, err := exec.Command(volPath, cStr...).CombinedOutput()

  if err != nil {
    log.Fatalf("error: %v", err)
  } else {
    ch<-string(result)
  }
}


func main() {
  options, plugins, flags := input()
  dumpFiles, outputPath := findDumps(options)
  cmds := buildCommands(dumpFiles, outputPath, options, plugins, flags)

  var threads int = 0
  volPath := options["_vol-name"]
  ch := make(chan string, len(cmds))
  _ = runtime.GOMAXPROCS(threads)

  // if flag set on CLI, print out commands rather than executing them
  if flags["print"] {
    for i:=0; i<(len(cmds)); i++ {
      pl(volPath, strings.Join(cmds[i], " "))
      pl("")
    }
    os.Exit(0)
  }

  // set go thread support
  if options["_threads"] != "" {
    threads, _ = strconv.Atoi(options["_threads"])
  } else {
    threads = runtime.NumCPU()
  }

  cmdIndex, kickoff := 0, 0
  cmdCount := len(cmds)

  // ensure enough work for starting batch of threads
  if threads >= cmdCount {
    kickoff = cmdCount
  } else {
    kickoff = threads
  }

  // start workers
  for i:=0; i<kickoff; i++ {
    go manager(ch, volPath, cmds[cmdIndex])
    cmdIndex++
  }

/*  for each iteration, wait for thread return. if there
    is another command waiting for execution, kick it off,
    and adjust the counter */

  for i:=cmdCount; i>0; i-- {
    select {
      case ret := <-ch:
        pl(ret)

        if cmdIndex < cmdCount {
          go manager(ch, volPath, cmds[cmdIndex])
          cmdIndex++
        }
    }
  }
}


