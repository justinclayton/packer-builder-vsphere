package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"strings"
	"text/template"
	"time"
)

type Task struct {
	Id    string
	Vim   VimSession
	State string
}

func (t *Task) WaitForCompletion() (result string, err error) {
	for {
		state, progress, result, err := t.GetState()
		// fmt.Printf("\nlooping: state is '%s' and progress is '%s'\n", state, progress)
		if err != nil {
			return "", err
		}
		if state == "success" {
			return result, nil
		}
		time.Sleep(2 * time.Second)
	}
}

func (t *Task) GetState() (state string, progress string, result string, err error) {
	data := struct {
		TaskId string
	}{
		t.Id,
	}
	tmpl := template.Must(template.New("TaskStatus").Parse(TaskStatusRequestTemplate2))
	response, err := t.Vim.sendRequest(tmpl, data)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Printf(
			"Bad status code [%d] [%s]\n",
			response.StatusCode,
			response.Status)
	}
	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	// fmt.Println(string(body))

	state = parseTaskPropertyValue("state", root)
	progress = parseTaskPropertyValue("progress", root)
	result = parseTaskPropertyValue("result", root)

	return

}

func parseTaskPropertyValue(propVal string, root *xmlpath.Node) string {
	pathString := strings.Join([]string{"//*/RetrievePropertiesResponse/returnval/propSet/val/", propVal}, "")
	path := xmlpath.MustCompile(pathString)
	if value, ok := path.String(root); ok {
		return value
	} else {
		return ""
	}
}
