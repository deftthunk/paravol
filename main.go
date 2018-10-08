package main

import (
  "fmt"
  "os"
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
// aliases
var pl = fmt.Println
var pf = fmt.Printf
var sp = fmt.Sprint


type dirWalk struct {
  memdumpPath  string
  resultsPath  string
  profile     string
  fileArr     [][]string
}

type Plugin struct {
  Name      string    `yaml:"name"`
  Pid       string    `yaml:"pid,omitempty"`
  Address   string    `yaml:"address,omitempty"`
}

type Config struct {
  Profile   string    `yaml:"profile"`
	State			string		`yaml:"state"`
  Memdumps  string    `yaml:"memdumps"`
  OutPath   string    `yaml:"output"`
  ProcPid   string    `yaml:"proc_pid"`
  Modules   []Plugin  `yaml:"plugins"`
}

func newDirwalk() dirWalk {
  dw := dirWalk {
    memdumpPath:  "",
    resultsPath:  "",
    profile:      "",
    fileArr:      make([][]string, 0),
  }
  return dw
}

func input(dw *dirWalk) {
  var argFail string = "Missing path to config.yaml"
	cfg := Config{}

  // count and parse CLI args
  if len(os.Args) < 1 {
    pl(argFail)
    os.Exit(1)
  }

  cfgFile, err := ioutil.ReadFile(os.Args[1])
  if err != nil {
    log.Fatal(err)
  }

  file_err := yaml.Unmarshal(cfgFile, &cfg)
  if file_err != nil {
    log.Fatalf("error: %v", file_err)
  }

  // crawler func for folders (used below)
  var walkFunc = func (path string, f os.FileInfo, err error) error {
    fullFilepath = path
    curFilename = filepath.Base(path)
    fi, err := os.Stat(path)

    if fi.IsDir() {
      curParentDir = path
    } else {
      dw.fileArr = append(dw.fileArr, []string{fi.Name(), curParentDir})
    }
    return nil
  }

  // each file/dir is passed to walkFunc
  filepath.Walk(dw.memdumpPath, walkFunc)

/*
  pl("arch: ", dw.profile)
  pl("curFilename: ", curFilename)
  pl("fullFilepath: ", fullFilepath)
  pl("curParentDir: ", curParentDir)
  pl("resultsPath: ", dw.resultsPath)
  pl("")
*/
}

func main() {
  dw := newDirwalk()
  input(&dw)

  pl(dw.fileArr)
/*
  for {
    select {
      case arr := <-inputCh:
        pl("inputCh: ", arr)
    }
  }
*/
}




/*
take input
  -path to memdumps
  -path to results folder
  -vol profile
  -config file (yaml?)

find memdump files in folder path

kick off multiple vol processes
  -make sure results are named according to memdump context
  -handle reported errors gracefully

show progress/status

volatility example:
python vol/vol.py --profile=Win10x64 -f /media/folder/dumps --output-file=/nhome/me/results malfind
*/





