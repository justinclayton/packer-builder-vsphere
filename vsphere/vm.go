package vsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"strings"
	"text/template"
)

type Vm struct {
	Vim           VimSession
	Id            string
	Name          string
	Parent        string
	Config        string
	InventoryPath string
	Ip            string
}

type CustomizationSpec struct {
	Name    string
	Network string
	Ip      string
	Gateway string
	Dns1    string
	Dns2    string
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

func (v *Vm) retrieveProperties(props []string) map[string]string {
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

func (v *Vm) DeployVM(newVmName string, spec CustomizationSpec) (newVm Vm) {

	// Use empty relocate spec in go template, no need to create type
	// Be sure to set template to true to avoid having to find and set resource pool for new vm

	data := struct {
		Folder string
		Name   string
	}{
		v.Parent,
		newVmName,
	}
	t := template.Must(template.New("CloneVMTask").Parse(CloneVMTaskRequestTemplate))
	response, err := v.Vim.sendRequest(t, data)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Printf("Bad status code [%d] [%s]\n", response.StatusCode, response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)
	root, _ := xmlpath.Parse(bytes.NewBuffer(body))
	path := xmlpath.MustCompile("//*/CloneVM_TaskResponse/returnval")
	if taskId, ok := path.String(root); ok {
		data := struct {
			TaskId string
		}{
			taskId,
		}
		t := template.Must(template.New("TaskStatus").Parse(TaskStatusRequestTemplate2))
		response, err = v.Vim.sendRequest(t, data)
		////******
	}
	// task := v.cloneVmTask(v.Parent, newVmName)
	// task.waitForCompletion()
	// relocateSpec = VIM.VirtualMachineRelocateSpec
	// spec = RbVmomi::VIM.VirtualMachineCloneSpec(:location => relocateSpec,
	//                                    :powerOn => false,
	//                                    :template => false)
	// vm.CloneVM_Task(:folder => vm.parent, :name => vm_target, :spec => spec).wait_for_completion

	// irb(main):115:0> vm.CloneVM_Task(:folder => vm.parent, :name => "deleteme54321", :spec => spec)

	//XXXXXXXX*******XXXXX

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><CloneVM_Task xmlns="urn:vim25"><_this type="VirtualMachine">vm-58499</_this><folder type="Folder">group-v53287</folder><name>deleteme54321</name><spec xsi:type="VirtualMachineCloneSpec"><location xsi:type="VirtualMachineRelocateSpec"></location><template>1</template><powerOn>0</powerOn></spec></CloneVM_Task></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa119a8f4c8>
	// set-cookie value from Response (if any):

	// Response (in 0.012 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <CloneVM_TaskResponse xmlns="urn:vim25"><returnval type="Task">task-518718</returnval></CloneVM_TaskResponse>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// => Task("task-518718")

	//TASK stuff

	//   irb(main):118:0> t.wait_for_completion
	// Headers:
	// {"content-type"=>"text/xml; charset=utf-8", "SOAPAction"=>"urn:vim25/5.0", "cookie"=>"vmware_soap_session=\"520b60d9-0187-0c70-4a65-290c24f32a05\"; Path=/; HttpOnly;"}

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><CreateFilter xmlns="urn:vim25"><_this type="PropertyCollector">propertyCollector</_this><spec xsi:type="PropertyFilterSpec"><propSet xsi:type="PropertySpec"><type>Task</type><all>0</all><pathSet>info.state</pathSet></propSet><objectSet xsi:type="ObjectSpec"><obj type="Task">task-518718</obj></objectSet></spec><partialUpdates>0</partialUpdates></CreateFilter></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa11912ec80>
	// set-cookie value from Response (if any):

	// Response (in 0.007 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <CreateFilterResponse xmlns="urn:vim25"><returnval type="PropertyFilter">session[5250cb13-57b7-ce32-1631-4732a3d3849d]52ee2323-09bd-119a-bf7a-c1e30525bdb0</returnval></CreateFilterResponse>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// Headers:
	// {"content-type"=>"text/xml; charset=utf-8", "SOAPAction"=>"urn:vim25/5.0", "cookie"=>"vmware_soap_session=\"520b60d9-0187-0c70-4a65-290c24f32a05\"; Path=/; HttpOnly;"}

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><WaitForUpdates xmlns="urn:vim25"><_this type="PropertyCollector">propertyCollector</_this><version></version></WaitForUpdates></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa1192b4b90>
	// set-cookie value from Response (if any):

	// Response (in 0.004 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <WaitForUpdatesResponse xmlns="urn:vim25"><returnval><version>6</version><filterSet><filter type="PropertyFilter">session[5250cb13-57b7-ce32-1631-4732a3d3849d]52ee2323-09bd-119a-bf7a-c1e30525bdb0</filter><objectSet><kind>enter</kind><obj type="Task">task-518718</obj><changeSet><name>info.state</name><op>assign</op><val xsi:type="TaskInfoState">success</val></changeSet></objectSet></filterSet></returnval></WaitForUpdatesResponse>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// Headers:
	// {"content-type"=>"text/xml; charset=utf-8", "SOAPAction"=>"urn:vim25/5.0", "cookie"=>"vmware_soap_session=\"520b60d9-0187-0c70-4a65-290c24f32a05\"; Path=/; HttpOnly;"}

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><RetrieveProperties xmlns="urn:vim25"><_this type="PropertyCollector">propertyCollector</_this><specSet xsi:type="PropertyFilterSpec"><propSet xsi:type="PropertySpec"><type>Task</type><pathSet>info</pathSet></propSet><objectSet xsi:type="ObjectSpec"><obj type="Task">task-518718</obj></objectSet></specSet></RetrieveProperties></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa119530dd8>
	// set-cookie value from Response (if any):

	// Response (in 0.006 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <RetrievePropertiesResponse xmlns="urn:vim25"><returnval><obj type="Task">task-518718</obj><propSet><name>info</name><val xsi:type="TaskInfo"><key>task-518718</key><task type="Task">task-518718</task><name>CloneVM_Task</name><descriptionId>VirtualMachine.clone</descriptionId><entity type="VirtualMachine">vm-58499</entity><entityName>CentOS 6 16 GB Template</entityName><state>success</state><cancelled>false</cancelled><cancelable>false</cancelable><result type="VirtualMachine" xsi:type="ManagedObjectReference">vm-58946</result><reason xsi:type="TaskReasonUser"><userName>AMER\svcvcent</userName></reason><queueTime>2013-12-31T01:01:51.15455Z</queueTime><startTime>2013-12-31T01:01:51.15455Z</startTime><completeTime>2013-12-31T01:02:15.77955Z</completeTime><eventChainId>57514083</eventChainId></val></propSet></returnval></RetrievePropertiesResponse>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// Headers:
	// {"content-type"=>"text/xml; charset=utf-8", "SOAPAction"=>"urn:vim25/5.0", "cookie"=>"vmware_soap_session=\"520b60d9-0187-0c70-4a65-290c24f32a05\"; Path=/; HttpOnly;"}

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><DestroyPropertyFilter xmlns="urn:vim25"><_this type="PropertyFilter">session[5250cb13-57b7-ce32-1631-4732a3d3849d]52ee2323-09bd-119a-bf7a-c1e30525bdb0</_this></DestroyPropertyFilter></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa11948b388>
	// set-cookie value from Response (if any):

	// Response (in 0.006 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <DestroyPropertyFilterResponse xmlns="urn:vim25"/>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// Headers:
	// {"content-type"=>"text/xml; charset=utf-8", "SOAPAction"=>"urn:vim25/5.0", "cookie"=>"vmware_soap_session=\"520b60d9-0187-0c70-4a65-290c24f32a05\"; Path=/; HttpOnly;"}

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><RetrieveProperties xmlns="urn:vim25"><_this type="PropertyCollector">propertyCollector</_this><specSet xsi:type="PropertyFilterSpec"><propSet xsi:type="PropertySpec"><type>Task</type><pathSet>info</pathSet></propSet><objectSet xsi:type="ObjectSpec"><obj type="Task">task-518718</obj></objectSet></specSet></RetrieveProperties></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa11950d928>
	// set-cookie value from Response (if any):

	// Response (in 0.006 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <RetrievePropertiesResponse xmlns="urn:vim25"><returnval><obj type="Task">task-518718</obj><propSet><name>info</name><val xsi:type="TaskInfo"><key>task-518718</key><task type="Task">task-518718</task><name>CloneVM_Task</name><descriptionId>VirtualMachine.clone</descriptionId><entity type="VirtualMachine">vm-58499</entity><entityName>CentOS 6 16 GB Template</entityName><state>success</state><cancelled>false</cancelled><cancelable>false</cancelable><result type="VirtualMachine" xsi:type="ManagedObjectReference">vm-58946</result><reason xsi:type="TaskReasonUser"><userName>AMER\svcvcent</userName></reason><queueTime>2013-12-31T01:01:51.15455Z</queueTime><startTime>2013-12-31T01:01:51.15455Z</startTime><completeTime>2013-12-31T01:02:15.77955Z</completeTime><eventChainId>57514083</eventChainId></val></propSet></returnval></RetrievePropertiesResponse>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// Headers:
	// {"content-type"=>"text/xml; charset=utf-8", "SOAPAction"=>"urn:vim25/5.0", "cookie"=>"vmware_soap_session=\"520b60d9-0187-0c70-4a65-290c24f32a05\"; Path=/; HttpOnly;"}

	// Request:
	// <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><env:Body><RetrieveProperties xmlns="urn:vim25"><_this type="PropertyCollector">propertyCollector</_this><specSet xsi:type="PropertyFilterSpec"><propSet xsi:type="PropertySpec"><type>Task</type><pathSet>info</pathSet></propSet><objectSet xsi:type="ObjectSpec"><obj type="Task">task-518718</obj></objectSet></specSet></RetrieveProperties></env:Body></env:Envelope>

	// DEBUG:
	// Net::HTTPOK

	// #<Net::HTTPOK:0x007fa1193a6b20>
	// set-cookie value from Response (if any):

	// Response (in 0.007 s)
	// <?xml version="1.0" encoding="UTF-8"?>
	// <soapenv:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// <soapenv:Body>
	// <RetrievePropertiesResponse xmlns="urn:vim25"><returnval><obj type="Task">task-518718</obj><propSet><name>info</name><val xsi:type="TaskInfo"><key>task-518718</key><task type="Task">task-518718</task><name>CloneVM_Task</name><descriptionId>VirtualMachine.clone</descriptionId><entity type="VirtualMachine">vm-58499</entity><entityName>CentOS 6 16 GB Template</entityName><state>success</state><cancelled>false</cancelled><cancelable>false</cancelable><result type="VirtualMachine" xsi:type="ManagedObjectReference">vm-58946</result><reason xsi:type="TaskReasonUser"><userName>AMER\svcvcent</userName></reason><queueTime>2013-12-31T01:01:51.15455Z</queueTime><startTime>2013-12-31T01:01:51.15455Z</startTime><completeTime>2013-12-31T01:02:15.77955Z</completeTime><eventChainId>57514083</eventChainId></val></propSet></returnval></RetrievePropertiesResponse>
	// </soapenv:Body>
	// </soapenv:Envelope>

	// => VirtualMachine("vm-58946")
	// irb(main):119:0>

	newVm = Vm{
		Name: spec.Name,
		Ip:   spec.Ip,
	}
	return newVm
}

func (v *Vm) MarkAsTemplate() bool {
	return true
}
