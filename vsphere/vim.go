package vsphere

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type VimClient struct {
	hostUrl    string
	httpClient http.Client
	cookie     string
}

func NewVimClient(user, pass, host string) (vim *VimClient, err error) {
	auth := struct {
		Username string
		Password string
		Host     string
	}{
		user,
		pass,
		host,
	}

	hostUrl := fmt.Sprintf("https://%s/sdk/vimService", host)

	vim = &VimClient{
		hostUrl: hostUrl,
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	t := template.Must(template.New("Login").Parse(LoginTemplate))
	message := applyTemplate(t, auth)
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

func getAssumedTemplatePath(sourceVmPath string, templateName string) (templatePath string) {
	a := strings.Split(sourceVmPath, "/")
	a = a[1 : len(a)-1]
	a = append(a, templateName)
	templatePath = strings.Join(a, "/")
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
	body, _ := ioutil.ReadAll(response.Body)
	// log.Printf("RESPONSE BODY BELOW:\n============\n%s\n===========\nEND RESPONSE BODY\n", string(body))
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
