package main

import (
	"github.com/justinclayton/packer-builder-vsphere/vsphere"
	"github.com/mitchellh/packer/packer/plugin"
	"os"
)

func main() {
	// If we were called with no arguments, we
	// assume packer is invoking us as a plugin
	args := os.Args[1:]
	if len(args) == 0 {
		server, err := plugin.Server() // THE ERROR IS HAPPENING HERE
		if err != nil {
			panic(err)
		}
		// log.Println("'IM NOT DEAD YET', he said to stderr")
		server.RegisterBuilder(new(vsphere.Builder))
		server.Serve()
	}

	//// Otherwise we march onward!
	user := args[0]
	pass := args[1]
	hosturl := args[2]
	pathToSourceVm := args[3]

	RunStandalone(user, pass, hosturl, pathToSourceVm)
}
