package main

import (
  "fmt"
  "os"
  "strings"
  "reflect"
  "path/filepath"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "log"
//  "github.com/goinggo/tracelog"
)

var curParentDir string = ""
var curFilename string = ""
var fullFilepath string = ""
var pl = fmt.Println
var pf = fmt.Printf
var sp = fmt.Sprint


/* grab YAML file and decode into structs */

func input() (map[string]string, []map[string]string) {
  var cfgMap map[string]interface{}
  var options map[string]string
  var plugins []map[string]string

  // count and parse CLI args
  if len(os.Args) < 1 {
    pl("Missing path to config.yaml")
    os.Exit(1)
  }

  // open config
  cfgFile, err := ioutil.ReadFile(os.Args[1])
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
//      pl("key: ", kOp, "val: ", ta)
      options[kOp] = ta
    } else if ta, ok := vOp.([]interface{}); ok {
      // iterate over "Plugins:" sub-config array
      // each value is a map of plugin config values
      configArr := ta

      for _, cMap := range configArr {
        obj := reflect.ValueOf(cMap)
        tmpMap := make(map[string]string, obj.Len() - 1)
//        pl("ValueOf: ", reflect.ValueOf(cMap))
        k := cMap.(map[interface{}]interface{})

        for kPlu, vPlu := range k {
//          pl("I/V: ", reflect.ValueOf(kPlu), reflect.ValueOf(vPlu))
          tmpMap[kPlu.(string)] = vPlu.(string)
        }
        plugins = append(plugins, tmpMap)
      }
    } else {
      // catch anything else
      pl("what am i: ", reflect.TypeOf(vOp))
    }
  }

  return options, plugins
}


/*  returns 2D array of dumpfiles in the format of 
    [[path, filename] [path, filename]] */

func findDumps(options map[string]string) [][]string {
  // [[path, filename] [path, filename]]
  dumpFiles := make([][]string, 0)
  var fullFilepath string
  var curFilename string
  var curParentDir string

  // crawler func for folders (used below)
  var walkFunc = func (path string, f os.FileInfo, err error) error {
    fullFilepath = path
    curFilename = filepath.Base(path)
    fi, err := os.Stat(path)

    if err != nil {
      pl("os.Stat(): Error on path ", path)
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
  subFolders := strings.Fields(options["subfolders"])

  for _, s := range subFolders {
    memDumpPath := filepath.Join(options["memdumps"], s)
    err := filepath.Walk(memDumpPath, walkFunc)

    if err != nil {
      log.Fatalf("Error: %v", err)
    }
  }

  return dumpFiles
}


func fixField(f string) string {
  hyphens := "--" + f
  return hyphens
}


/*  build volatility command for each memory dump */

func buildCommands(dumpFiles [][]string, opt map[string]string, plu []map[string]string) []string {
  var commands []string

  for _, pathArr := range dumpFiles {
    var optString []string

    // basic command info
    optString = append(optString,
      "vol.py",
      " --profile=", opt["profile"],
      " --filename=", strings.Join(pathArr, "/"),
      " --verbose",
    )

    // iterate through supplied plugin maps in 'plu' slice
    for _, pMap := range plu {
      var pluString []string
      // create plugin string place plugin name in front of its options
      pluString = append(pluString, pMap["plugin"])
      delete(pMap, "plugin")

      // add plugin options to command string
      for field, val := range pMap {
        hField := fixField(field)
        var newStr string

        if val != "" {
          str := []string{hField, val}
          newStr = strings.Join(str, "=")
        } else {
          newStr = hField
        }
        pluString = append(pluString, newStr)
        pl("pluString: ", pluString)
      }
      // need intermediate optString for each pMap
      commands = append(commands, strings.Join(optString, ""))
//      optString = append(optString, strings.Join(pluString, " "))
//      pl("optString: ", optString)
    }
//    commands = append(commands, strings.Join(optString, ""))
  }

  return commands
}


func main() {
  options, plugins := input()
  dumpFiles := findDumps(options)
  cmds := buildCommands(dumpFiles, options, plugins)

  for _, v := range cmds {
    pl("")
    pl(v)
  }
}


/*
volatility example:
python vol/vol.py --verbose --profile=Win10x64 -f /media/folder/dumps --output-file=/nhome/me/results malfind
*/







