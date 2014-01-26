package vsphere

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/multistep"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
	"strings"
	"text/template"
)

type StepGetTemplatePath struct{}

func (s *StepGetTemplatePath) Run(state multistep.StateBag) multistep.StepAction {
	newVm := state.Get("new_vm").(*Vm)
	data := struct {
		VmId string
	}{
		newVm.Id,
	}
	tmpl := template.Must(template.New("RetrievePropertiesPathTraversal").Parse(RetrievePropertiesPathTraversalRequestTemplate))
	response, err := newVm.Vim.sendRequest(tmpl, data)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Printf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//RetrievePropertiesResponse//val")
	iter := path.Iter(root)
	values := make([]string, 0)

	iter.Next() // skip top element "Datacenters"
	for {
		iter.Next() // skip id element
		ok := iter.Next()
		if ok == false {
			break
		} else {
			newVal := []string{iter.Node().String()}
			// add new value to the beginning of the slice
			// TODO: get rid of ids
			// current end value: Datacenters/group-d1/Tukwila/datacenter-2/vm/group-v3/1-templates/group-v53287/Lower/group-v54541/my_new_template_that_packer_built
			values = append(newVal, values...)
		}
	}
	templatePath := strings.Join(values, "/")
	log.Printf("New template path is '%s'", templatePath)
	//
	state.Put("template_path", templatePath)
	return multistep.ActionContinue
}

func (s *StepGetTemplatePath) Cleanup(state multistep.StateBag) {
}
