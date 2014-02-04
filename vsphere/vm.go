package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func (vim *VimClient) getVmBasicInfo(vmId string) (values map[string]string, err error) {
	// static set of what constitutes "basic info"
	props := []string{"name", "parent", "resourcePool"}
	response, err := vim.retrieveProperties(vmId, "name", "parent", "resourcePool")
	if err != nil {
		err = fmt.Errorf("Failed to get VM basic info: '%s'", err.Error())
		return
	}
	values = make(map[string]string)

	body, _ := ioutil.ReadAll(response.Body)
	for _, prop := range props {
		values[prop] = parseVmPropertyValue(prop, bytes.NewBuffer(body))
		if values[prop] == "" {
			err = fmt.Errorf("Failed to get value for VM property '%s'", prop)
		}
	}
	return
}

func (vim *VimClient) getVmIp(vmId string) (ip string, err error) {
	response, err := vim.retrieveProperties(vmId, "guest")
	if err != nil {
		err = fmt.Errorf("Failed to get VM IP: '%s'", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	ip = parseIpProperty(bytes.NewBuffer(body))
	return
}

func (vim *VimClient) getVmPowerState(vmId string) (state string, err error) {
	response, err := vim.retrieveProperties(vmId, "runtime")
	if err != nil {
		err = fmt.Errorf("Failed to get VM power state: '%s'", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	state = parseVmPowerStateProperty(bytes.NewBuffer(body))
	return
}

func (vim *VimClient) getVmPath(vmId string) (path string, err error) {

	response, err := vim.retrievePropertiesPathTraversal(vmId)
	if err != nil {
		err = fmt.Errorf("Failed to get VM Path: '%s'", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	path = parsePathProperty(bytes.NewBuffer(body))
	return
}

func (vim *VimClient) retrieveProperties(vmId string, props ...string) (response *http.Response, err error) {

	data := struct {
		VmId       string
		Properties []string
	}{
		vmId,
		props,
	}
	t := template.Must(template.New("RetrieveProperties").Parse(RetrievePropertiesRequestTemplate))

	request, _ := vim.prepareRequest(t, data)
	response, err = vim.Do(request)
	if err != nil {
		err = fmt.Errorf("Error retrieving VM properties: '%s'", err.Error())
	}
	return
}

func (vim *VimClient) retrievePropertiesPathTraversal(vmId string) (response *http.Response, err error) {
	data := struct {
		VmId string
	}{
		vmId,
	}
	t := template.Must(template.New("RetrievePropertiesPathTraversal").Parse(RetrievePropertiesPathTraversalRequestTemplate))

	request, _ := vim.prepareRequest(t, data)
	response, err = vim.Do(request)
	if err != nil {
		err = fmt.Errorf("Error retrieving VM path properties: '%s'", err.Error())
	}
	return
}

func parseVmPropertyValue(prop string, body *bytes.Buffer) (value string) {
	root, _ := xmlpath.Parse(body)
	pathString := strings.Join([]string{"//*/RetrievePropertiesResponse/returnval/propSet[name='", prop, "']/val"}, "")
	path := xmlpath.MustCompile(pathString)
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

func parseIpProperty(body *bytes.Buffer) (value string) {
	root, _ := xmlpath.Parse(bytes.NewBuffer(body.Bytes()))
	path := xmlpath.MustCompile("//*/RetrievePropertiesResponse/returnval/propSet[name='guest']/val/ipAddress")
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

func parseVmPowerStateProperty(body *bytes.Buffer) (value string) {
	root, _ := xmlpath.Parse(bytes.NewBuffer(body.Bytes()))
	path := xmlpath.MustCompile("//*/RetrievePropertiesResponse/returnval/propSet[name='runtime']/val/powerState")
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

func parsePathProperty(body *bytes.Buffer) (value string) {

	root, _ := xmlpath.Parse(bytes.NewBuffer(body.Bytes()))
	path := xmlpath.MustCompile("//RetrievePropertiesResponse//val")
	iter := path.Iter(root)
	values := make([]string, 0)

	// iter.Next() // skip top element "Datacenters"
	for {
		// iter.Next() // skip id element
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
	value = strings.Join(values, "/")
	return value
}

func parseTaskIdFromResponse(response *http.Response) (value string) {
	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))

	path := xmlpath.MustCompile("//*/CloneVM_TaskResponse/returnval")
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}

// cloneVmTask() calls the CloneVM_Task vSphere API method and returns a
// task id that can be used to track progress of the clone operation.
func (vim *VimClient) cloneVmTask(sourceVmId, folder, name string) (taskId string, err error) {

	data := struct {
		SourceVmId string
		Folder     string
		Name       string
	}{
		sourceVmId,
		folder,
		name,
	}

	t := template.Must(template.New("CloneVMTask").Parse(CloneVMTaskRequestTemplate))

	request, _ := vim.prepareRequest(t, data)
	response, err := vim.Do(request)
	if err != nil {
		err = fmt.Errorf("Error calling CloneVM_Task: '%s'", err.Error())
		return
	}
	taskId = parseTaskIdFromResponse(response)
	if taskId == "" {
		err = fmt.Errorf("Error calling CloneVM_Task: COULDNT FIND TASK ID IN RESPONSE. BUG?", nil)
		return
	}
	return
}

func (vim *VimClient) waitForIp(resultCh chan<- string, errCh chan<- error, vmId string) {
	for {
		ip, err := vim.getVmIp(vmId)
		if err != nil {
			errCh <- err
			return
		}
		if ip != "" {
			log.Printf("VM has IP '%s'\n", ip)
			resultCh <- ip
			return
		}
		log.Println("No IP yet, retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}

func (vim *VimClient) waitUntilVmShutdownComplete(errCh chan<- error, vmId string) {
	for {
		powerState, err := vim.getVmPowerState(vmId)
		if err != nil {
			errCh <- err
			return
		}
		if powerState == "poweredOff" {
			log.Printf("VM's power state is '%s'", powerState)
			errCh <- nil
			return
		}
		log.Printf("VM's power state is '%s', retrying in 5 seconds...", powerState)
		time.Sleep(5 * time.Second)
	}
}

// shutdownGuest() calls the ShutdownGuest vSphere API method.
func (vim *VimClient) shutdownGuest(vmId string) (err error) {

	data := struct {
		VmId string
	}{
		vmId,
	}

	t := template.Must(template.New("ShutdownGuest").Parse(ShutdownGuestRequestTemplate))

	request, _ := vim.prepareRequest(t, data)
	_, err = vim.Do(request)
	if err != nil {
		err = fmt.Errorf("Error calling ShutdownGuest: '%s'", err.Error())
		return
	}
	return
}

func (vim *VimClient) markAsTemplate(vmId string) (err error) {
	data := struct {
		VmId string
	}{
		vmId,
	}

	t := template.Must(template.New("MarkAsTemplate").Parse(MarkAsTemplateRequestTemplate))

	request, _ := vim.prepareRequest(t, data)
	_, err = vim.Do(request)
	if err != nil {
		err = fmt.Errorf("Error calling MarkAsTemplate: '%s'", err.Error())
		return
	}
	return
}
