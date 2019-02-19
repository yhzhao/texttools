// gofind file based on regex
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type FileTool struct {
	FileContentPattern string // regexp pattern string
}

func (t FileTool) PrintFilename(path string, info os.FileInfo, err error) error {
	fmt.Println(path)
	return nil
}

func (t FileTool) PrintFilenameIfMatchContent(path string, info os.FileInfo, err error) error {
	b, err := ioutil.ReadFile(path)
	s := string(b)
	matched, err := regexp.MatchString(t.FileContentPattern, s)
	if err != nil {
		fmt.Printf("encountered error %v while accessing file %s\n", err, path)
		return err
	}
	if matched {
		fmt.Println(path)
	}
	return nil
}

func main() {
	//err := filepath.Walk(".", PrintFilename)
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v pattern\n", os.Args[0])
		os.Exit(-1)
	}
	fmt.Printf("Matching pattern %v to filecontent\n", os.Args[1])
	var t FileTool = FileTool{os.Args[1]}
	err := filepath.Walk(".", t.PrintFilenameIfMatchContent)
	if err != nil {
		fmt.Printf("encountering error %v\n", err)
	}
}
