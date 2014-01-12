package vsphere

import (
	"fmt"
	"log"

	gossh "code.google.com/p/go.crypto/ssh"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/communicator/ssh"
)

// sshAddress returns the ssh address.
func sshAddress(state multistep.StateBag) (string, error) {
	log.Print("Getting SSH address...")
	config := state.Get("config").(*Config)
	ipAddress := state.Get("new_vm").(*Vm).Ip

	ipAndPort := fmt.Sprintf("%s:%d", ipAddress, config.SSHPort)

	log.Printf("'%s'\n", ipAndPort)
	return ipAndPort, nil
}

// sshConfig returns the ssh configuration.
func sshConfig(state multistep.StateBag) (*gossh.ClientConfig, error) {
	log.Print("Getting SSH config...")
	config := state.Get("config").(*Config)
	privateKey := string(config.privateKeyBytes)

	log.Printf("value of privateKey is '%s'\n", privateKey)
	keyring := new(ssh.SimpleKeychain)
	if err := keyring.AddPEMKey(privateKey); err != nil {
		return nil, fmt.Errorf("Error setting up SSH config: %s", err)
	}

	sshConfig := &gossh.ClientConfig{
		User: config.SSHUsername,
		Auth: []gossh.ClientAuth{gossh.ClientAuthKeyring(keyring)},
	}

	log.Printf("sshConfig.User is '%s'\n", sshConfig.User)

	return sshConfig, nil
}
