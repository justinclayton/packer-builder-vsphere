package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
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

func (v *Vm) retrieveProperties() error {
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
	if err != nil {
		err = fmt.Errorf("Error sending request: '%s'", err.Error())
		return err
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

	// uses guest property collected in the request
	v.Ip = parseIpProperty(root)
	return nil
}

func parseIpProperty(root *xmlpath.Node) string {
	path := xmlpath.MustCompile("//*/RetrievePropertiesResponse/returnval/propSet[name='guest']/val/ipAddress")
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

func (v *Vm) DeployVM(newVmName string) (newVm Vm, err error) {

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
	if err != nil {
		err = fmt.Errorf("Error sending request: '%s'", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//*/CloneVM_TaskResponse/returnval")
	if taskId, ok := path.String(root); ok {
		tsk := Task{Id: taskId, Vim: v.Vim}
		newVmId, err := tsk.WaitForCompletion()
		if err != nil {
			err = fmt.Errorf("Error waiting for task to complete: '%s'", err.Error())
			return newVm, err
		}
		newVm, err = v.Vim.NewVm(newVmId)
		if err != nil {
			err = fmt.Errorf("Error creating new VM: '%s'", err.Error())
			return newVm, err
		}

		c := make(chan error, 1)
		go newVm.waitForIp(c)

		select {
		case err := <-c:
			return newVm, err
		case <-time.After(waitForIpTimeoutInSeconds * time.Second):
			err := fmt.Errorf("timed out waiting for '%s' to get an IP", newVm.Name)
			return newVm, err
		}
	} else {

		log.Println("failed to get proper response back from clonevm!")
	}
	return newVm, err
}

func (v *Vm) waitForIp(c chan<- error) {
	for {
		err := v.retrieveProperties()
		if err != nil {
			c <- err
			return
		}
		if v.Ip != "" {
			log.Printf("New VM '%s' has IP '%s'\n", v.Name, v.Ip)
			c <- nil
			return
		}

		log.Printf("New VM '%s' has no IP yet, retrying in 10 seconds...\n", v.Name)
		time.Sleep(10 * time.Second)
	}
}

func (v *Vm) MarkAsTemplate() bool {
	return true
}
