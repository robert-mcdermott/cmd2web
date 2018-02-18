package main

import (
	"flag"
	"fmt"
	"os"
)

// define and set default command parameter flags
var exposeFlag = flag.String("expose", "", "[optional] expose this directory via web at https://hostname/accesscode/files")
var expireFlag = flag.Int("expire", 0, "[optional] terminate the cmd2web server after this many minutes")
var helpFlag = flag.Bool("h", false, "print usage information")

func init() {
	// usage function that's executed if a required flag is missing or user asks for help (-h)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [--expose <path-to-dir> --expire <minutes>] <command>>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample: %s uptime\n\n", os.Args[0])
		//flag.PrintDefaults()
		fmt.Println()
	}
	flag.Parse()

	// provide help (-h)
	if *helpFlag == true {
		flag.Usage()
		os.Exit(0)
	}
}
