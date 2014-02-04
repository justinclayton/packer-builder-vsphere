package vsphere

import (
	"errors"
	"fmt"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"time"
)

// Config is the configuration structure for the vsphere builder. It stores
// both the publicly settable state as well as the privately generated
// state of the config object.
type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	VsphereUsername string `mapstructure:"vsphere_username"`
	VspherePassword string `mapstructure:"vsphere_password"`
	VsphereHost     string `mapstructure:"vsphere_host"`
	SourceVmPath    string `mapstructure:"source_vm_path"`
	VmName          string `mapstructure:"vm_name"`
	Passphrase      string `mapstructure:"passphrase"`
	PrivateKeyFile  string `mapstructure:"private_key_file"`
	SSHUsername     string `mapstructure:"ssh_username"`
	SSHPort         uint   `mapstructure:"ssh_port"`
	RawSSHTimeout   string `mapstructure:"ssh_timeout"`
	RawStateTimeout string `mapstructure:"state_timeout"`

	privateKeyBytes []byte
	sshTimeout      time.Duration
	stateTimeout    time.Duration
	tpl             *packer.ConfigTemplate
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	md, err := common.DecodeConfig(c, raws...)
	if err != nil {
		return nil, nil, err
	}

	c.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, nil, err
	}

	errs := common.CheckUnusedConfig(md)

	// // Set defaults.
	if c.RawSSHTimeout == "" {
		c.RawSSHTimeout = "5m"
	}
	if c.RawStateTimeout == "" {
		c.RawStateTimeout = "5m"
	}
	if c.SSHUsername == "" {
		c.SSHUsername = "root"
	}
	if c.SSHPort == 0 {
		c.SSHPort = 22
	}

	// Process timeout settings.
	sshTimeout, err := time.ParseDuration(c.RawSSHTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing ssh_timeout: %s", err))
	}
	c.sshTimeout = sshTimeout

	stateTimeout, err := time.ParseDuration(c.RawStateTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing state_timeout: %s", err))
	}
	c.stateTimeout = stateTimeout

	// Process required parameters.
	if c.VsphereUsername == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a vsphere_username must be specified"))
	}
	if c.VspherePassword == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a vsphere_password must be specified"))
	}
	if c.VsphereHost == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a vsphere_host must be specified"))
	}
	if c.SourceVmPath == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a source_vm_path must be specified"))
	}
	if c.VmName == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a vm_name must be specified"))
	}
	if c.PrivateKeyFile == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a private_key_file"))
	} else {
		// Load the private key.
		c.privateKeyBytes, err = processPrivateKeyFile(c.PrivateKeyFile, c.Passphrase)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Failed loading private key file: %s", err))
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	return c, nil, nil

}
