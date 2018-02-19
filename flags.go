package main

import (
	"flag"
	"fmt"
	"os"
)

// Define command parameter flags
// I had to create a new flagset namespace because beego has it's own flags that were begin displayed by flag.PrintDefaults
var flagset = new(flag.FlagSet)
var exposeFlag = flagset.String("expose", "", "\n[optional] expose this directory or file at https://*/file\n"+
	"if a directory path is given it will provide an html file/dir listing\n"+
	"that you can navigate files and sub directories. if a file path is\n"+
	"provided, the file will be availible at the file URL\n")
var expireFlag = flagset.Int("expire", 0, "\n[optional] terminate the cmd2web server after the provide number of\n"+
	"minutes. If an expiration is not provide the server will run indefinately\n"+
	"until terminated manually\n")
var refreshFlag = flagset.Int("refresh", 0, "\n[optional] page refresh interval in seconds; only works with html\n"+
	"output format with GUI web browsers (Chrome, Firefox, etc...). each\n"+
	"refresh re-runs the command.\n")
var rawFlag = flagset.Bool("raw", false, "\n[optional] the default output is html; this flag enables raw text\n"+
	"output that is more suitable for use with curl or using as input to\n"+
	"another program or logging.\n")
var helpFlag = flagset.Bool("help", false, "\nprint usage information\n")

func init() {
	// usage function that's executed if a required flag is missing or user asks for help (-h)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [--expose <path> --expire <minutes> --refresh <seconds> --raw] <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n\n")
		flagset.PrintDefaults() // note "flagset." vs "flag."
		fmt.Fprintf(os.Stderr, "\nExample 1: list they systems process table and refresh the output every 30 seconds.\n"+
			"\n\t%s --refresh 30 ps aux\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample 2: expose the \"myproject\" directory to the web for 60 minutes.\n"+
			"\n\t%s --expire 60 --expose /home/rmcdermo/myproject /usr/bin/true\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "")
	}
	flag.Parse()

	// provide help (-h)
	if *helpFlag == true {
		flag.Usage()
		os.Exit(0)
	}
}
