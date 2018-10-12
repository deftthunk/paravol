package main

import (
  "fmt"
  "os"
  "strings"
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
  Name      string    `yaml:"plugin"`
  Pid       string    `yaml:"pid,omitempty"`
  Address   string    `yaml:"address,omitempty"`
}

type Config struct {
  Threads   int       `yaml:"threads"`
  Profile   string    `yaml:"profile"`
  States    string    `yaml:"states"`
  Memdumps  string    `yaml:"memdumps"`
  OutPath   string    `yaml:"output"`
  ProcPid   string    `yaml:"proc_pid"`
  Modules   []Plugin  `yaml:"plugins"`
}


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
      dumpFiles = append(dumpFiles, []string{fi.Name(), curParentDir})
    }

    return nil
  }

  // each file/dir is passed to walkFunc
  states := strings.Fields(c.States)

  for _, s := range states {
    memDumpPath := filepath.Join(c.Memdumps, s)
    err := filepath.Walk(memDumpPath, walkFunc)

    if err != nil {
      log.Fatalf("Error: %v", err)
    }
  }

  return dumpFiles
}


func (c Config) buildCommands([][]string) []string {
  volBin := 'vol.py'
  filenameFlag := '--filename='
  outputFileFlag := '--output-file='
  verboseFlag := '--verbose'
  modAddressFlag := '--addr='
  modprofileFlag := '--profile='
  modpidFlag := '--pid='
  modoffsetFlag := '--offset='
  modprocNameFlag := '--name='
  moddumpDirFlag := '--dump-dir='
  modvadBaseAddrFlag := '--base='
  modDllDumpRebaseFlag := '--fix'
  modDllDumpMemoryFlag := '--memory'
}

func main() {
  c := input()
  c.buildCommands(c.findDumps())

  pl(dumpFiles)
  pl(c)
  pl()
  for _, plugin := range c.Modules {
    pl("Plugin Name: ", plugin.Name)
  }
}


/*
volatility example:
python vol/vol.py --profile=Win10x64 -f /media/folder/dumps --output-file=/nhome/me/results malfind
*/


