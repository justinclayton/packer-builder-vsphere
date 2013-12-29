package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"net/http"
	"text/template"
)

type VimSession struct {
	hostUrl    string
	httpClient http.Client
	cookie     string
}

func (vim *VimSession) sendRequest(t *template.Template, data interface{}) (*http.Response, error) {
	message := applyTemplate(t, data)

	// println(message)

	request, _ := http.NewRequest("POST", vim.hostUrl, bytes.NewBufferString(message))
	if vim.cookie != "" {
		request.Header.Add("cookie", vim.cookie)
	}

	// send request
	return vim.httpClient.Do(request)
}

func (vim *VimSession) getVmTemplate(inventoryPath string) vm {
	// searchIndex.FindByInventoryPath(:inventoryPath => path)
	v, _ := vim.findByInventoryPath(inventoryPath)

	props := append(make([]string, 0), "name", "parent")

	propValues := v.retrieveProperties(props)
	v.Name = propValues["name"]
	v.Parent = propValues["parent"]

	return v
}

func (vim *VimSession) findByInventoryPath(inventoryPath string) (vm, error) {
	// searchIndex.FindByInventoryPath(:inventoryPath => path)
	data := struct {
		InventoryPath string
	}{
		inventoryPath,
	}
	t := template.Must(template.New("FindByInventoryPath").Parse(FindByInventoryPathRequestTemplate))

	response, err := vim.sendRequest(t, data)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Printf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//*/FindByInventoryPathResponse/returnval")
	if vmid, ok := path.String(root); ok {
		v := vm{
			Vim:           *vim,
			Id:            vmid,
			InventoryPath: inventoryPath,
		}
		return v, nil
	} else {
		return vm{}, errors.New("Found nothing")
	}
}

func (vim *VimSession) deployVM(templateVm vm, spec customizationSpec) (newVm vm) {
	newVm = vm{
		Name: spec.name,
		Ip:   spec.ip,
	}
	return newVm
}

func (vim *VimSession) markAsTemplate(thisVm vm) bool {
	return true
}
