package main

import (
	"fmt"
	"os"
	"path/filepath"
//	"log"
//	"gopkg.in/yaml.v2"
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
	memdumpPath	string
	resultsPath	string
	profile     string
	fileArr     [][]string
}

func newDirwalk() dirWalk {
	dw := dirWalk {
		memdumpPath:	"",
		resultsPath:	"",
		profile:			"",
		fileArr:			make([][]string, 0),
	}
	return dw
}

func input(dw *dirWalk) {
	var argFail string = ""

	// instructions for program args
	argFail += sp("Missing args:\n\t-Memory Dump Path")
	argFail += sp("\n\t-Results Destination Folder Path")
	argFail += sp("\n\t-Volatility Profile\n")

	// count and parse CLI args
	if len(os.Args) > 1 {
		dw.memdumpPath = os.Args[1]
	} else {
		pl(argFail)
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		dw.resultsPath = os.Args[2]
	} else {
		pl(argFail)
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		dw.profile = os.Args[3]
	} else {
		pl(argFail)
		os.Exit(1)
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





