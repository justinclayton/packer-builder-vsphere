package main

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

func applyTemplate(t *template.Template, data interface{}) string {
	var b bytes.Buffer
	err := t.Execute(&b, data)
	if err != nil {
		println(err.Error())
	}
	return b.String()
}
