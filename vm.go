package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"strings"
	"text/template"
)

type vm struct {
	Vim           VimSession
	Id            string
	Name          string
	Parent        string
	Config        string
	InventoryPath string
	Ip            string
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

func (v *vm) retrieveProperties(props []string) map[string]string {
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

func (v *vm) deployVM(spec customizationSpec) (newVm vm) {
	newVm = vm{
		Name: spec.name,
		Ip:   spec.ip,
	}
	return newVm
}

func (v *vm) markAsTemplate() bool {
	return true
}
