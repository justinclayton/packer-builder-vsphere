package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"strings"
	"text/template"
	"time"
)

const waitForIpTimeoutInSeconds = 300

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

func parseVmPropertyValue(propVal string, root *xmlpath.Node) string {
	pathString := strings.Join([]string{"//*/RetrievePropertiesResponse/returnval/propSet[name='", propVal, "']/val"}, "")
	path := xmlpath.MustCompile(pathString)
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

func (v *Vm) retrieveProperties() {
	props := append(make([]string, 0), "name", "parent", "resourcePool", "guest")

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
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))

	values := make(map[string]string)
	for _, prop := range props {
		values[prop] = parseVmPropertyValue(prop, root)
	}
	v.Name = values["name"]
	v.Parent = values["parent"]
	v.ResourcePool = values["resourcePool"]

	v.Ip = parseIpProperty(root)
	return
}

func parseIpProperty(root *xmlpath.Node) string {
	path := xmlpath.MustCompile("//*/RetrievePropertiesResponse/returnval/propSet[name='guest']/val/ipAddress")
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
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

		for i := 0; i < waitForIpTimeoutInSeconds; i = i + 10 {
			newVm.retrieveProperties()
			if newVm.Ip != "" {
				fmt.Printf("New VM '%s' has IP '%s'\n", newVm.Name, newVm.Ip)
				return newVm
			}
			fmt.Printf("New VM '%s' has no IP yet, retrying in 10 seconds...\n")
			time.Sleep(10 * time.Second)
		}
	} else {
		fmt.Println("failed to get proper response back from clonevm!")
	}
	return newVm
}

func (v *Vm) MarkAsTemplate() bool {
	return true
}
