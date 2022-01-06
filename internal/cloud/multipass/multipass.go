package multipass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/minectl/internal/automation"
	minctlTemplate "github.com/minectl/internal/template"
	"github.com/minectl/internal/update"
)

const (
	multipassBinary = "multipass"
)

type Multipass struct {
	tmpl *minctlTemplate.Template
}

func NewMultipass() (*Multipass, error) {
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	return &Multipass{
		tmpl: tmpl,
	}, nil
}

func (m *Multipass) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	script, err := m.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{SSHPublicKey: string(pubKeyFile), Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(fmt.Sprintf("%s/cloud-config.yaml", os.TempDir()), []byte(script), 0o600)
	if err != nil {
		return nil, err
	}

	app := multipassBinary
	arg0 := "launch"
	arg1 := "-n"
	arg2 := args.MinecraftResource.GetName()
	arg3 := "--cloud-init"
	arg4 := fmt.Sprintf("%s/cloud-config.yaml", os.TempDir())
	arg5 := "-c"
	arg6 := strings.Split(args.MinecraftResource.GetSize(), "-")[0]
	arg7 := "-m"
	arg8 := strings.Split(args.MinecraftResource.GetSize(), "-")[1]

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return m.GetServer(args.MinecraftResource.GetName(), args)
}

func (m *Multipass) DeleteServer(id string, args automation.ServerArgs) error {
	cmd := exec.Command(multipassBinary, "delete", id)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command(multipassBinary, "purge")
	cmdOutput = &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (m *Multipass) ListServer() ([]automation.ResourceResults, error) {
	panic("List Server is not possible with Multipass, as it does not support labels")
}

func (m Multipass) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := m.GetServer(id, args)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.PublicIP, "ubuntu")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (m Multipass) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := m.GetServer(id, args)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.PublicIP, "root")
	err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)), args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service", args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	return nil
}

func (m Multipass) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	cmd := exec.Command(multipassBinary, "info", "--format", "json", args.MinecraftResource.GetName()) //nolint: gosec
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(cmdOutput.Bytes(), &result)
	server := result["info"].(map[string]interface{})[args.MinecraftResource.GetName()].(map[string]interface{})
	ip := server["ipv4"].([]interface{})[0].(string)
	return &automation.ResourceResults{
		ID:       args.MinecraftResource.GetName(),
		Name:     args.MinecraftResource.GetName(),
		Region:   multipassBinary,
		PublicIP: ip,
		Tags:     "",
	}, err
}
