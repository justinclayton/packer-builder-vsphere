package vsphere

import (
	"bytes"
	"crypto/tls"
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

func (vim *VimSession) sendRequest(t *template.Template, data interface{}) (response *http.Response, err error) {
	message := applyTemplate(t, data)

	// println(message)

	request, _ := http.NewRequest("POST", vim.hostUrl, bytes.NewBufferString(message))
	if vim.cookie != "" {
		request.Header.Add("cookie", vim.cookie)
	}

	// send request
	response, err = vim.httpClient.Do(request)
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
		return
	}
	return
}

func (vim *VimSession) NewVm(vmId string) (Vm, error) {
	v := Vm{
		Vim: *vim,
		Id:  vmId,
	}
	err := v.retrieveProperties()
	if err != nil {
		err := fmt.Errorf("Failed to retrieve properties for '%s' VM: %s", v.Name, err)
		return v, err
	}
	return v, err
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
	if err != nil {
		err = fmt.Errorf("Error sending request: '%s'", err.Error())
		return Vm{}, err
	}

	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//*/FindByInventoryPathResponse/returnval")
	if vmId, ok := path.String(root); ok {
		v, err := vim.NewVm(vmId)
		return v, err
	} else {
		err := fmt.Errorf("Found nothing", nil)
		return Vm{}, err
	}
}

func (vim *VimSession) DeleteVm(inventoryPath string) {
	fmt.Printf("Would delete VM %s", inventoryPath)
	return
}
