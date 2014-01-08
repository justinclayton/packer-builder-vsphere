package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"strings"
	"text/template"
)

type Vm struct {
	Vim           VimSession
	Id            string
	Name          string
	Parent        string
	ResourcePool  string
	Config        string
	InventoryPath string
	Ip            string
}

type CustomizationSpec struct {
	Name    string
	Network string
	Ip      string
	Gateway string
	Dns1    string
	Dns2    string
}

func parsePropertyValue(propVal string, root *xmlpath.Node) string {
	pathString := strings.Join([]string{"//*/RetrievePropertiesResponse/returnval/propSet[name='", propVal, "']/val"}, "")
	path := xmlpath.MustCompile(pathString)
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

func (v *Vm) retrieveProperties(props []string) map[string]string {
	data := struct {
		VmId       string
		Properties []string
	}{
		v.Id,
		props,
	}
	t := template.Must(template.New("RetrieveProperties").Parse(RetrievePropertiesRequestTemplate))

	response, err := v.Vim.sendRequest(t, data)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Printf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("BEGIN RESPONSE BODY")
	// fmt.Println(string(body))
	// fmt.Println("END RESPONSE BODY")
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	// path := xmlpath.MustCompile("//*/RetrievePropertiesResponse/returnval/propSet")

	values := make(map[string]string)
	for _, prop := range props {
		values[prop] = parsePropertyValue(prop, root)
	}
	return values
	// v.Name = parsePropertyValue("name", root)
	// v.Parent = parsePropertyValue("parent", root)
	// return nil, errors.New("Found nothing")
}

func (v *Vm) DeployVM(newVmName string, spec CustomizationSpec) (newVm Vm) {

	// Use empty relocate spec in go template, no need to create type
	// Be sure to set template to true to avoid having to find and set resource pool for new vm

	data := struct {
		SourceVmId string
		Folder     string
		Name       string
	}{
		v.Id,
		v.Parent,
		newVmName,
	}
	tmpl := template.Must(template.New("CloneVMTask").Parse(CloneVMTaskRequestTemplate))
	response, err := v.Vim.sendRequest(tmpl, data)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Printf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//*/CloneVM_TaskResponse/returnval")
	if taskId, ok := path.String(root); ok {
		tsk := Task{Id: taskId, Vim: v.Vim}
		newVmId, err := tsk.WaitForCompletion()
		if err != nil {
			fmt.Println(err.Error())
		}
		newVm := v.Vim.NewVm(newVmId)
		return newVm
	} else {
		fmt.Println("failed to get proper response back from clonevm!")
	}
	return newVm
}

func (v *Vm) MarkAsTemplate() bool {
	return true
}
