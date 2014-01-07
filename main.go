package main

import (
	// "github.com/mitchellh/packer"
	// "github.com/mitchellh/packer/plugin"
	"os"
)

func main() {

	// If we were called with no arguments, we
	// assume packer is invoking us as a plugin
	args := os.Args[1:]
	if len(args) == 0 {
		// plugin.ServeBuilder(new(Builder))
		panic("BEHAVE AS A PLUGIN HERE")
		return
	}

	//// Otherwise we march onward!
	user := args[0]
	pass := args[1]
	hosturl := args[2]
	pathToSourceVm := args[3]

	RunStandalone(user, pass, hosturl, pathToSourceVm)
}

// // SOLO USAGE: ./packer-communicator-winrm cmd -user vagrant -pass vagrant command-text
// // set WINRM_DEBUG=1 for more output

// func main() {
//

//         Run(&cmd{
//                 Handle: func(user, pass, command string) {
//                         communicator := &Communicator{endpoint, user, pass}
//                         rc := &packer.RemoteCmd{
//                                 Command: command,
//                                 Stdout:  os.Stdout,
//                                 Stderr:  os.Stderr,
//                         }

//                         err := communicator.Start(rc)
//                         if err != nil {
//                                 log.Printf("unable to run command: %s", err)
//                                 return
//                         }

//                         rc.Wait()
//                 },
//         })
// }
