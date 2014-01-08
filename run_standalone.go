package main

import (
	"fmt"
	"github.com/justinclayton/packer-builder-vsphere/vsphere"
)

func RunStandalone(user, pass, hosturl, pathToSourceVm string) {

	fmt.Printf("Connecting to vSphere server at '%s' as user '%s'...", hosturl, user)
	vc := vsphere.NewVimSession(user, pass, hosturl)
	fmt.Printf("done.\n")

	fmt.Printf("Looking for VM '%s'...", pathToSourceVm)
	sourceVm := vc.GetVmTemplate(pathToSourceVm)
	fmt.Printf("found source VM '%s'\n", sourceVm.Name)

	newVmName := "packer_vsphere_builder_test_vm"
	fmt.Printf("Creating new VM '%s'...", newVmName)
	spec := vsphere.CustomizationSpec{
		Ip: "1.2.3.4",
	}
	newVm := sourceVm.DeployVM(newVmName, spec)
	fmt.Printf("'%s' created.\n", newVm.Name)

	fmt.Printf("Marking new VM '%s' as template...", newVm.Name)
	// _ = newVm.MarkAsTemplate()
	fmt.Printf("done.\n")

}
