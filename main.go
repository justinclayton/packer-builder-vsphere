package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"text/template"
)

type auth struct {
	Username string
	Password string
	HostUrl  string
}

type customizationSpec struct {
	name    string
	network string
	ip      string
	gateway string
	dns1    string
	dns2    string
}

func login(a auth) (handle VimSession) {

	handle.hostUrl = a.HostUrl
	t := template.Must(template.New("Login").Parse(LoginTemplate))
	message := applyTemplate(t, a)
	// disable strict ssl checking
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	handle.httpClient = http.Client{Transport: tr}
	request, _ := http.NewRequest("POST",
		handle.hostUrl, bytes.NewBufferString(message))
	// send request
	response, err := handle.httpClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Errorf("Bad status code [%d] [%s]", response.StatusCode, response.Status)
	}

	// assuming cookies[] count is 1
	handle.cookie = (response.Cookies()[0].Raw)

	return handle
}

func tellPackerProvisionersToRun(ip string) bool {
	return true
}

func main() {
	fmt.Println("Next step: read input from JSON config")
	a := auth{
		"someusername",
		"somepassword",
		"somehosturl",
	}
	vc := login(a)
	templateVm := vc.getVmTemplate("somepathtovm")
	fmt.Println(templateVm)
	spec := customizationSpec{
		ip: "1.2.3.4",
	}
	vm := vc.deployVM(templateVm, spec)
	_ = tellPackerProvisionersToRun(vm.Ip)
	ok := vc.markAsTemplate(vm)
	if ok {
		// println("you did it")
	}
}
