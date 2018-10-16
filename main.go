package main

import (
  "fmt"
  "os"
//  "strings"
  "reflect"
//  "path/filepath"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "log"
//  "github.com/goinggo/tracelog"
)

// globals
var curParentDir string = ""
var curFilename string = ""
var fullFilepath string = ""
var pl = fmt.Println
var pf = fmt.Printf
var sp = fmt.Sprint


/*
  grab YAML file and decode into structs
*/
func input() {
  var cfgMap map[string]interface{}

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

  convertToMap(cfgMap)
//  return
}

func convertToMap(c map[string]interface{}) {
  for k, v := range c {
    if s, ok := v.(string); ok && s != "" {
      pl("key: ", k, "val: ", s)
    } else if s, ok := v.(int); ok {
      pl("key: ", k, "val: ", s)
    } else if s, ok := v.([]interface{}); ok {
      test2(s)
    } else {
      pl("what am i: ", reflect.TypeOf(v))
    }
  }
}

func test2(c []interface{}) {
  // iterate over "Plugins:" sub-config
  for _, j := range c {
    pl("ValueOf: ", reflect.ValueOf(j))
    k := j.(map[interface{}]interface{})

    for g, h := range k {
      pl("I/V: ", reflect.ValueOf(g), reflect.ValueOf(h))
    }
  }
}

func convertToArray(c []interface{}) {
  for i, v := range c {
    if s, ok := v.(string); ok {
      pl("idx: ", i, "val: ", s)
    } else if s, ok := v.(int); ok {
      pl("idx: ", i, "val: ", s)
    } else {
      pl("what is this: ", reflect.TypeOf(v))
    }
  }
}




/*
  returns 2D array of dumpfiles in the format of 
  [[filename, path] [filename, path]]
*/
/*
func findDumps() [][]string {
  // [[file, path], [file, path]]
  dumpFiles := make([][]string, 0)

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
  subFolders := strings.Fields(c.SubFolders)

  for _, s := range subFolders {
    memDumpPath := filepath.Join(c.Memdumps, s)
    err := filepath.Walk(memDumpPath, walkFunc)

    if err != nil {
      log.Fatalf("Error: %v", err)
    }
  }

  return dumpFiles
}
*/

/*
  Convert Plugin struct fields into an array of values. Array conversion
  allows for iteration and automatic building of command string using config
  values.
func convertStruct(p Plugin) []interface{} {
  v := reflect.ValueOf(p)
  values := make([]interface{}, v.NumField())

  for i:=0; i < v.NumField(); i++ {
    values[i] = v.Field(i).Interface()
  }

  //pl(values)
  return values
}
*/


/*
  build volatility command for each 'state' memory dump listed
*/
func buildCommands(dumpFiles [][]string) []string {
  var commands []string

/*  for _, pathArr := range dumpFiles {
    // basic command info
    cmdString := []string {"vol.py",
      " --profile=", c.Profile,
      " --filename=", strings.Join(pathArr, "/"),
      " --verbose",
    }
*/

    // per plugin command strings
    //for _, plugin := range c.Modules {
    //  pluginArr := convertStruct(plugin)
    //  pl(pluginArr)
      //pl(plugin)
    //}

//  }

  return commands
}


// test struct-less yaml parsing
func test() {
  var cfgMap map[string]interface{}

  // open config
  cfgFile, err := ioutil.ReadFile(os.Args[1])
  if err != nil {
    log.Fatal(err)
  }

  // parse yaml into struct
  decode_err := yaml.Unmarshal(cfgFile, &cfgMap)
  if decode_err != nil {
    log.Fatalf("error: %v", decode_err)
  }

  for k, v := range cfgMap {
    if s, ok := v.(string); ok && s != "" {
      pl("key: ", k, "val: ", s)
    } else if s, ok := v.(int); ok {
      pl("key: ", k, "val: ", s)
    }
  }
}


func main() {
//  c := input()
//  verifyConfig(c)
//  c.buildCommands(c.findDumps())
input()
//  test()
}


/*
volatility example:
python vol/vol.py --verbose --profile=Win10x64 -f /media/folder/dumps --output-file=/nhome/me/results malfind
*/







