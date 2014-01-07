package vsphere

import (
	"bytes"
	"text/template"
)

const LoginTemplate = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <Login xmlns="urn:vim25">
      <_this type="SessionManager">SessionManager</_this>
      <userName>{{.Username}}</userName>
      <password>{{.Password}}</password>
    </Login>
  </env:Body>
</env:Envelope>`

const RootFolderEnvelope = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <RetrieveProperties xmlns="urn:vim25">
      <_this type="PropertyCollector">propertyCollector</_this>
      <specSet xsi:type="PropertyFilterSpec">
        <propSet xsi:type="PropertySpec">
          <type>Folder</type>
          <pathSet>childEntity</pathSet>
        </propSet>
        <objectSet xsi:type="ObjectSpec">
          <obj type="Folder">group-d1</obj>
        </objectSet>
      </specSet>
    </RetrieveProperties>
  </env:Body>
</env:Envelope>
`

const FindByInventoryPathRequestTemplate = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <FindByInventoryPath xmlns="urn:vim25">
      <_this type="SearchIndex">SearchIndex</_this>
      <inventoryPath>{{.InventoryPath}}</inventoryPath>
    </FindByInventoryPath>
  </env:Body>
</env:Envelope>`

const RetrievePropertiesRequestTemplate = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <RetrieveProperties xmlns="urn:vim25">
      <_this type="PropertyCollector">propertyCollector</_this>
      <specSet xsi:type="PropertyFilterSpec">
        <propSet xsi:type="PropertySpec">
          <type>VirtualMachine</type>
          {{range .Properties}}<pathSet>{{.}}</pathSet>{{end}}
          </propSet>
          <objectSet xsi:type="ObjectSpec">
          <obj type="VirtualMachine">{{.VmId}}</obj>
        </objectSet>
      </specSet>
    </RetrieveProperties>
  </env:Body>
</env:Envelope>`

const CloneVMTaskRequestTemplate = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <CloneVM_Task xmlns="urn:vim25">
      <_this type="VirtualMachine">vm-58499</_this>
      <folder type="Folder">{{.Folder}}</folder>
      <name>{{.Name}}</name>
      <spec xsi:type="VirtualMachineCloneSpec">
        <location xsi:type="VirtualMachineRelocateSpec"></location>
        <template>1</template>
        <powerOn>0</powerOn>
      </spec>
      </CloneVM_Task>
    </env:Body>
</env:Envelope>
`

const TaskStatusRequestTemplate = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <CreateFilter xmlns="urn:vim25">
      <_this type="PropertyCollector">propertyCollector</_this>
      <spec xsi:type="PropertyFilterSpec">
        <propSet xsi:type="PropertySpec">
          <type>Task</type>
          <all>0</all>
          <pathSet>info.state</pathSet>
        </propSet>
        <objectSet xsi:type="ObjectSpec">
          <obj type="Task">{{.TaskId}}</obj>
        </objectSet>
      </spec>
      <partialUpdates>0</partialUpdates>
    </CreateFilter>
  </env:Body>
</env:Envelope>
`

const TaskStatusRequestTemplate2 = `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <env:Body>
    <RetrieveProperties xmlns="urn:vim25">
      <_this type="PropertyCollector">propertyCollector</_this>
      <specSet xsi:type="PropertyFilterSpec">
        <propSet xsi:type="PropertySpec">
          <type>Task</type>
          <pathSet>info</pathSet>
        </propSet>
        <objectSet xsi:type="ObjectSpec">
          <obj type="Task">{{.TaskId}}</obj>
        </objectSet>
      </specSet>
    </RetrieveProperties>
  </env:Body>
</env:Envelope>
`

func applyTemplate(t *template.Template, data interface{}) string {
	var b bytes.Buffer
	err := t.Execute(&b, data)
	if err != nil {
		println(err.Error())
	}
	return b.String()
}
