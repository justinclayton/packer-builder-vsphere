package main

import (
	"fmt"
	"github.com/justinclayton/packer-builder-vsphere/vsphere"
	"github.com/mitchellh/packer/packer/plugin"
	"os"
)

func main() {
	// If we were called with no arguments, we
	// assume packer is invoking us as a plugin
	args := os.Args[1:]
	if len(args) == 0 {
		server, err := plugin.Server()
		if err != nil {
			panic(err)
		}
		server.RegisterBuilder(new(vsphere.Builder))
		server.Serve()
	}

	// Otherwise give the user a pointer to wiring up this custom packer plugin
	usage()
	os.Exit(1)
}

func usage() {
	`
===============================================================================
Attention! This is intended to be used as a plugin to Packer, and
cannot be invoked as a standalone executable. To use this custom
plugin, create a file called $HOME/.packerconfig with this content:

{
  "builders": {
    "vsphere": "/path/to/executable/called/packer-builder-vsphere"
  }
}

This will allow you to add a vsphere builder stanza to your packer
build template that would look something like this:

{
  "builders": [{
    "type": "vsphere",
    "vsphere_username": "myuser",
    "vsphere_password": "mypass",
    "vsphere_host": "my.vcenter.server.fqdn",
    "source_vm_path": "/MyDatacenter/vm/path/to/my/source_template",
    "vm_name": "my_new_template_that_packer_built",
    "private_key_file": "/path/on/local/machine/to/assumed/private.key"
  }]
}

Please direct any questions to @justinclayton42 on Twitter, or submit
a problem at https://github.com/justinclayton/packer-builder-vsphere/issues.
==============================================================================`
}
