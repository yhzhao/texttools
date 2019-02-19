// xor byte by byte
package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

type FileTool struct {
	XorPattern string // regexp pattern string
}

func (t FileTool) Xor(input, output string) error {
	b, err := ioutil.ReadFile(input)
	o := make([]byte, len(b))
	mask, err := hex.DecodeString(t.XorPattern)
	if err != nil {
		fmt.Printf("Encountered error while decode xor pattern %v\n%v\n", t.XorPattern, err)
		return err
	}
	for i, t := range b {
		o[i] = t ^ mask[0]
	}
	err = ioutil.WriteFile(output, o, 0777)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %v inputfilename outputfilenamek xor_pattern\n", os.Args[0])
		os.Exit(-1)
	}
	fmt.Printf("input: %v, output: %v, xor pattern: %v\n", os.Args[1], os.Args[2], os.Args[3])
	var t FileTool = FileTool{os.Args[3]}
	err := t.Xor(os.Args[1], os.Args[2])
	fmt.Println(err)
}
