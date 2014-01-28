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

type VimClient struct {
	hostUrl    string
	httpClient http.Client
	cookie     string
}

func NewVimClient(user, pass, hosturl string) (vim *VimClient, err error) {
	auth := struct {
		Username string
		Password string
		HostUrl  string
	}{
		user,
		pass,
		hosturl,
	}

	vim = &VimClient{
		hostUrl: auth.HostUrl,
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	// vim.hostUrl = auth.HostUrl
	t := template.Must(template.New("Login").Parse(LoginTemplate))
	message := applyTemplate(t, auth)
	// disable strict ssl checking
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// vim.httpClient = http.Client{Transport: tr}
	request, _ := http.NewRequest("POST",
		vim.hostUrl, bytes.NewBufferString(message))
	// send request
	log.Println("About to submit login request to vSphere")
	response, err := vim.httpClient.Do(request)
	log.Println("Got a response back from vSphere")
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("Bad status code [%d] [%s]", response.StatusCode, response.Status)
	}

	if err != nil {
		err = fmt.Errorf("Error connecting to vSphere: '%s'", err.Error())
		return
	}

	// assuming cookies[] count is 1
	vim.cookie = (response.Cookies()[0].Raw)

	return
}

func (vim *VimClient) prepareRequest(t *template.Template, data interface{}) (request *http.Request, err error) {
	message := applyTemplate(t, data)

	request, err = http.NewRequest("POST", vim.hostUrl, bytes.NewBufferString(message))
	if err != nil {
		return
	}

	if vim.cookie != "" {
		request.Header.Add("cookie", vim.cookie)
	}
	return
}

func (vim *VimClient) Do(request *http.Request) (response *http.Response, err error) {
	// send request
	response, err = vim.httpClient.Do(request)

	if response.StatusCode != 200 {
		err = fmt.Errorf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
		return
	}
	return
}

func (vim *VimClient) FindByInventoryPath(inventoryPath string) (vmId string, err error) {
	data := struct {
		InventoryPath string
	}{
		inventoryPath,
	}
	t := template.Must(template.New("FindByInventoryPath").Parse(FindByInventoryPathRequestTemplate))

	log.Printf("Looking for '%s'", inventoryPath)

	request, _ := vim.prepareRequest(t, data)
	response, err := vim.Do(request)
	// defer response.Body.Close()
	if err != nil {
		err = fmt.Errorf("Error calling FindByInventoryPath: '%s'", err.Error())
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR IN ioutil.ReadAll: '%s'", err.Error())
	}
	log.Printf("RESPONSE BODY BELOW:\n============\n%s\n===========\nEND RESPONSE BODY\n", string(body))
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//*/FindByInventoryPathResponse/returnval")
	if vmId, ok := path.String(root); ok {
		return vmId, err
	} else {
		err := fmt.Errorf("Found nothing.")
		return vmId, err
	}
}

func (vim *VimClient) DeleteVm(inventoryPath string) (err error) {
	err = fmt.Errorf("DeleteVm() NOT IMPLEMENTED YET")
	return
}
