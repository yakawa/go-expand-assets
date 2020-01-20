package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/yakawa1128/go-expand-assets"
)

type Options struct {
	Verbose  []bool   `short:"v" long:"verbose" description:"Verbose output"`
	Version  bool     `long:"version" description:"Print Version"`
	Source   []string `short:"s" long:"source" description:"Source Folder" required:"true"`
	Absolute bool     `short:"a" long:"absolute" describe:"Freeze absolute path"`
}

var opts Options

func printVersion() {
	fmt.Println("Version: ")
}

func getExecPath() {
	e, _ := os.Executable()
	fmt.Println("Exec: ", e)
}

func getCurrentPath() {
	c, _ := os.Getwd()
	fmt.Println("Cwd: ", c)
}

func generate(s []string) {
	g := assets.Generator{}

	for _, d := range s {
		g.Add(d)
	}
	g.Write(os.Stdout)
}

func main() {
	//getExecPath()
	//getCurrentPath()
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return
	}
	if opts.Version == true {
		printVersion()
		return
	}
	generate(opts.Source)
}
