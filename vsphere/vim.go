package vsphere

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
	"net/http"
	"text/template"
)

type VimSession struct {
	hostUrl    string
	httpClient http.Client
	cookie     string
}

func NewVimSession(user, pass, hosturl string) (vim VimSession) {
	auth := struct {
		Username string
		Password string
		HostUrl  string
	}{
		user,
		pass,
		hosturl,
	}

	vim.hostUrl = auth.HostUrl
	t := template.Must(template.New("Login").Parse(LoginTemplate))
	message := applyTemplate(t, auth)
	// disable strict ssl checking
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	vim.httpClient = http.Client{Transport: tr}
	request, _ := http.NewRequest("POST",
		vim.hostUrl, bytes.NewBufferString(message))
	// send request
	log.Println("About to submit login request to vSphere")
	response, err := vim.httpClient.Do(request)
	log.Println("Got a response back from vSphere")
	defer response.Body.Close()

	if err != nil {
		println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Errorf("Bad status code [%d] [%s]", response.StatusCode, response.Status)
	}

	// assuming cookies[] count is 1
	vim.cookie = (response.Cookies()[0].Raw)

	return vim
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

func (vim *VimSession) NewVm(vmId string) Vm {
	v := Vm{
		Vim: *vim,
		Id:  vmId,
	}
	v.retrieveProperties()
	return v
}

func (vim *VimSession) GetVmTemplate(inventoryPath string) Vm {
	v, _ := vim.FindByInventoryPath(inventoryPath)

	return v
}

func (vim *VimSession) FindByInventoryPath(inventoryPath string) (Vm, error) {
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
	if vmId, ok := path.String(root); ok {
		v := vim.NewVm(vmId)
		return v, nil
	} else {
		return Vm{}, errors.New("Found nothing")
	}
}

func (vim *VimSession) DeleteVm(inventoryPath string) {
	fmt.Printf("Would delete VM %s", inventoryPath)
	return
}
