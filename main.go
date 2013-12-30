package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
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

func login(a auth) (vim VimSession) {

	vim.hostUrl = a.HostUrl
	t := template.Must(template.New("Login").Parse(LoginTemplate))
	message := applyTemplate(t, a)
	// disable strict ssl checking
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	vim.httpClient = http.Client{Transport: tr}
	request, _ := http.NewRequest("POST",
		vim.hostUrl, bytes.NewBufferString(message))
	// send request
	response, err := vim.httpClient.Do(request)
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

func tellPackerProvisionersToRun(ip string) bool {
	return true
}

func main() {
	username := os.Args[1]
	password := os.Args[2]
	hostUrl := os.Args[3]
	pathToVm := os.Args[4]

	a := auth{
		username,
		password,
		hostUrl,
	}

	vc := login(a)
	templateVm := vc.getVmTemplate(pathToVm)
	fmt.Println(templateVm)
	spec := customizationSpec{
		ip: "1.2.3.4",
	}
	newVm := templateVm.deployVM(spec)
	_ = tellPackerProvisionersToRun(newVm.Ip)
	ok := newVm.markAsTemplate()
	if ok {
		println("you did it")
	}
}
