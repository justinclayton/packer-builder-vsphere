package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
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

	responseBody, _ := ioutil.ReadAll(response.Body)

	propValues := parseTaskPropertyValues(responseBody, "state", "progress", "result", "error")

	t.state = propValues["state"]
	log.Printf("task state is '%s'", t.state)
	t.progress = propValues["progress"]
	log.Printf("task progress is '%s'", t.progress)
	t.result = propValues["result"]
	log.Printf("task result is '%s'", t.result)
	t.errorDesc = propValues["error"]
	log.Printf("task errorDesc is '%s'", t.errorDesc)

	return

}

func parseTaskPropertyValues(body []byte, props ...string) map[string]string {
	values := make(map[string]string)
	log.Println(props)
	for _, prop := range props {
		root, _ := xmlpath.Parse(bytes.NewBuffer(body))
		// pathString := fmt.Sprintf("//*/RetrievePropertiesResponse/returnval/propSet/val/%s", prop)
		pathString := strings.Join([]string{"//*/RetrievePropertiesResponse/returnval/propSet/val/", prop}, "")
		path := xmlpath.MustCompile(pathString)
		if value, ok := path.String(root); ok {
			// log.Printf("ok; value of '%s' is '%s'", prop, value)
			values[prop] = value
		} else {
			// log.Printf("not ok; value of '%s' is '%s'", prop, value)
			values[prop] = ""
		}
	}
	return values
}
