package vsphere

import (
	"fmt"
	"log"

	"code.google.com/p/go.crypto/ssh"
	"github.com/mitchellh/multistep"
)

// sshAddress returns the ssh address.
func sshAddress(state multistep.StateBag) (string, error) {
	log.Print("Getting SSH address...")
	config := state.Get("config").(*Config)
	ipAddress := state.Get("new_vm_ip").(string)

	ipAndPort := fmt.Sprintf("%s:%d", ipAddress, config.SSHPort)

	log.Printf("'%s'\n", ipAndPort)
	return ipAndPort, nil
}

// sshConfig returns the ssh configuration.
func sshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	log.Print("Getting SSH config...")
	config := state.Get("config").(*Config)

	log.Printf("sshConfig.User is '%s'\n", config.SSHUsername)

	privateKey := string(config.privateKeyBytes)
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("Error setting up SSH config: %s", err)
	}

	return &ssh.ClientConfig{
		User: config.SSHUsername,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}, nil
}
