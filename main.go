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

// globals
var curParentDir string = ""
var curFilename string = ""
var fullFilepath string = ""
var pl = fmt.Println
var pf = fmt.Printf
var sp = fmt.Sprint


type Plugin struct {
  PluginName      string    `yaml:"plugin"`
  Pid             string    `yaml:"pid,omitempty"`
  Address         string    `yaml:"address,omitempty"`
  HiveOffset      string    `yaml:"hive-offset,omitempty"`
  Offset          string    `yaml:"offset,omitempty"`
  Name            string    `yaml:"name,omitempty"`
  Fix             string    `yaml:"fix,omitempty"`
  Quick           string    `yaml:"quick,omitempty"`
  NoWhitelist     string    `yaml:"no-whitelist,omitempty"`
  SkipKernel      string    `yaml:"skip-kernel,omitempty"`
  SkipProcess     string    `yaml:"skip-process,omitempty"`
  Virtual         string    `yaml:"virtual,omitempty"`
  ShowUnallocated string    `yaml:"show-unallocated,omitempty"`
  StartAddress    string    `yaml:"start-address,omitempty"`
  Length          string    `yaml:"length,omitempty"`
  SysOffset       string    `yaml:"sys-offset,omitempty"`
  SecOffset       string    `yaml:"sec-offset,omitempty"`
  MaxHistory      string    `yaml:"max-history,omitempty"`
  PhysicalOffset  string    `yaml:"physical-offset,omitempty"`
  HistoryBuffers  string    `yaml:"history-buffers,omitempty"`
  DumpDir         string    `yaml:"dump-dir,omitempty"`
  Regex           string    `yaml:"regex,omitempty"`
  IgnoreCase      string    `yaml:"ignore-case,omitempty"`
  Base            string    `yaml:"base,omitempty"`
  Addr            string    `yaml:"addr,omitempty"`
  Ssl             string    `yaml:"ssl,omitempty"`
  Physical        string    `yaml:"physical,omitempty"`
  PhysOffset      string    `yaml:"physoffset,omitempty"`
  SummaryFile     string    `yaml:"summary-file,omitempty"`
  Unsafe          string    `yaml:"unsafe,omitempty"`
  Filter          []string  `yaml:"filter,omitempty"`
  Silent          string    `yaml:"silent,omitempty"`
  ObjectType      []string  `yaml:"object-type,omitempty"`
  SamOffset       string    `yaml:"sam-offset,omitempty"`
  BlockSize       string    `yaml:"blocksize,omitempty"`
  OutputImage     string    `yaml:"output-image,omitempty"`
  Count           string    `yaml:"count,omitempty"`
  Size            string    `yaml:"size,omitempty"`
  MaxSize         string    `yaml:"max-size,omitempty"`
  Refined         string    `yaml:"refined,omitempty"`
  Hex             string    `yaml:"hex,omitempty"`
  Hash            string    `yaml:"hash,omitempty"`
  FullHash        string    `yaml:"fullhash,omitempty"`
  Disoffset       string    `yaml:"disoffset,omitempty"`
  NoCheck         string    `yaml:"nocheck,omitempty"`
  Disk            string    `yaml:"disk,omitempty"`
  MaxDistance     string    `yaml:"maxdistance,omitempty"`
  ZeroStart       string    `yaml:"zerostart,omitempty"`
  Machine         string    `yaml:"machine,omitempty"`
  DebugOut        string    `yaml:"debugout,omitempty"`
  Memory          string    `yaml:"memory,omitempty"`
  Tag             string    `yaml:"tag,omitempty"`
  MinSize         string    `yaml:"min-size,omitempty"`
  Paged           string    `yaml:"paged,omitempty"`
  Key             string    `yaml:"key,omitempty"`
  ApplyRules      string    `yaml:"apply-rules,omitempty"`
  StringFile      string    `yaml:"string-file,omitempty"`
  Scan            string    `yaml:"scan,omitempty"`
  LookupPid       string    `yaml:"lookup-pid,omitempty"`
  ListTags        string    `yaml:"listtags,omitempty"`
  Hive            string    `yaml:"hive,omitempty"`
  User            string    `yaml:"user,omitempty"`
  Type            []string  `yaml:"type,omitempty"`
  ListHead        string    `yaml:"listhead,omitempty"`
  Free            string    `yaml:"free,omitempty"`
  All             string    `yaml:"all,omitempty"`
  Case            string    `yaml:"case,omitempty"`
  Kernel          string    `yaml:"kernel,omitempty"`
  Wide            string    `yaml:"wide,omitempty"`
  YaraRules       string    `yaml:"yara-rules,omitempty"`
  YaraFile        string    `yaml:"yara-file,omitempty"`
  Reverse         string    `yaml:"reverse,omitempty"`
}

type Config struct {
  Threads     int       `yaml:"threads"`
  Profile     string    `yaml:"profile"`
  SubFolders  string    `yaml:"subfolders"`
  Memdumps    string    `yaml:"memdumps"`
  OutPath     string    `yaml:"output"`
  ProcPid     string    `yaml:"proc_pid"`
  Modules     []Plugin    `yaml:"plugins"`
}


/*
  grab YAML file and decode into structs
*/
func input() Config {
  var argFail string = "Missing path to config.yaml"
  cfg := Config{}

  // count and parse CLI args
  if len(os.Args) < 1 {
    pl(argFail)
    os.Exit(1)
  }

  // open config
  cfgFile, err := ioutil.ReadFile(os.Args[1])
  if err != nil {
    log.Fatal(err)
  }

  // parse yaml into struct
  decode_err := yaml.Unmarshal(cfgFile, &cfg)
  if decode_err != nil {
    log.Fatalf("error: %v", decode_err)
  }

  return cfg
}


/*
  returns 2D array of dumpfiles in the format of 
  [[filename, path] [filename, path]]
*/
func (c Config) findDumps() [][]string {
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


/*
  Convert Plugin struct fields into an array of values. Array conversion
  allows for iteration and automatic building of command string using config
  values.
*/
func convertStruct(p Plugin) []interface{} {
  v := reflect.ValueOf(p)
  values := make([]interface{}, v.NumField())

  for i:=0; i < v.NumField(); i++ {
    values[i] = v.Field(i).Interface()
  }

  //pl(values)
  return values
}


/*
  build volatility command for each 'state' memory dump listed
*/
func (c Config) buildCommands(dumpFiles [][]string) []string {
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
    for _, plugin := range c.Modules {
      pluginArr := convertStruct(plugin)
      pl(pluginArr)
      //pl(plugin)
    }

//  }

  return commands
}


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

  for i,j := range cfgMap {
    if i == "plugins" {
      pl("j: ", reflect.TypeOf(j))
      arr, _ := j.([]map[string]string)
      pl("arr", reflect.TypeOf(arr))
      for w,x := range arr {
        pl("w :", reflect.TypeOf(w))
        pl("x :", reflect.TypeOf(x))
        for y,z := range x {
          pl("y - z: ", reflect.TypeOf(y), " - ", reflect.TypeOf(z))
        }
      }
    }
  }
}


func main() {
  c := input()
//  verifyConfig(c)
  c.buildCommands(c.findDumps())

  test()
}


/*
volatility example:
python vol/vol.py --verbose --profile=Win10x64 -f /media/folder/dumps --output-file=/nhome/me/results malfind
*/







