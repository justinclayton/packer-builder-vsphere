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

type Task struct {
	Id        string
	Vim       *VimClient
	state     string
	progress  string
	result    string
	errorDesc string
}

func (vim *VimClient) NewTask(taskId string) (t *Task) {
	t = &Task{Id: taskId, Vim: vim}
	return
}

func (t *Task) WaitForCompletion(resultCh chan<- string, errCh chan<- error) {
	for {
		err := t.RefreshState()
		if err != nil {
			errCh <- fmt.Errorf("Unable to refresh task status: '%s'", err.Error())
			return
		}
		if t.state == "error" {
			errCh <- fmt.Errorf("Task ended with an error: '%s'", t.errorDesc)
			return
		}
		if t.state == "success" {
			log.Printf("Task ended with state '%s' and result '%s'", t.state, t.result)
			resultCh <- t.result
			return
		}
		log.Printf("Task %s: %s percent\n", t.state, t.progress)
		time.Sleep(2 * time.Second)
	}
}

func (t *Task) RefreshState() (err error) {
	data := struct {
		TaskId string
	}{
		t.Id,
	}
	tmpl := template.Must(template.New("TaskStatus").Parse(TaskStatusRequestTemplate2))

	request, _ := t.Vim.prepareRequest(tmpl, data)
	response, err := t.Vim.Do(request)

	t.state = parseTaskPropertyValue("state", response)
	t.progress = parseTaskPropertyValue("progress", response)
	t.result = parseTaskPropertyValue("result", response)
	t.errorDesc = parseTaskPropertyValue("error", response)

	return

}

func parseTaskPropertyValue(propVal string, response *http.Response) string {
	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))

	pathString := strings.Join([]string{"//*/RetrievePropertiesResponse/returnval/propSet/val/", propVal}, "")
	path := xmlpath.MustCompile(pathString)
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}
