package vultr

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
	"github.com/minectl/pkg/update"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

type Vultr struct {
	client *govultr.Client
	tmpl   *minctlTemplate.Template
}

func NewVultr(apiKey string) (*Vultr, error) {
	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: apiKey})
	vultrClient := govultr.NewClient(oauth2.NewClient(ctx, ts))
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	vultr := &Vultr{
		client: vultrClient,
		tmpl:   tmpl,
	}
	return vultr, nil
}

func (v *Vultr) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}
	sshKey, err := v.client.SSHKey.Create(context.Background(), &govultr.SSHKeyReq{
		SSHKey: strings.TrimSpace(string(pubKeyFile)),
		Name:   fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
	})
	if err != nil {
		return nil, err
	}

	script, err := v.tmpl.GetTemplate(args.MinecraftResource, "", minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer()))
	if err != nil {
		return nil, err
	}
	startupScript, err := v.client.StartupScript.Create(context.Background(), &govultr.StartupScriptReq{
		Script: base64.StdEncoding.EncodeToString([]byte(script)),
		Name:   fmt.Sprintf("%s-stackscript", args.MinecraftResource.GetName()),
		Type:   "boot",
	})
	if err != nil {
		return nil, err
	}

	ubuntu2004Id := 387
	opts := &govultr.InstanceCreateReq{
		SSHKeys:  []string{sshKey.ID},
		ScriptID: startupScript.ID,
		Hostname: args.MinecraftResource.GetName(),
		Label:    args.MinecraftResource.GetName(),
		Region:   args.MinecraftResource.GetRegion(),
		Plan:     args.MinecraftResource.GetSize(),
		OsID:     ubuntu2004Id,
		Tag:      fmt.Sprintf("%s,%s", common.InstanceTag, args.MinecraftResource.GetEdition()),
	}

	instance, err := v.client.Instance.Create(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	stillCreating := true
	for stillCreating {
		instance, err = v.client.Instance.Get(context.Background(), instance.ID)
		if err != nil {
			return nil, err
		}
		if instance.Status == "active" {
			stillCreating = false
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	return &automation.RessourceResults{
		ID:       instance.ID,
		Name:     instance.Label,
		Region:   instance.Region,
		PublicIP: instance.MainIP,
		Tags:     instance.Tag,
	}, err
}

func (v *Vultr) DeleteServer(id string, args automation.ServerArgs) error {
	sshKeys, _, err := v.client.SSHKey.List(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, sshKey := range sshKeys {
		if sshKey.Name == fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()) {
			err := v.client.SSHKey.Delete(context.Background(), sshKey.ID)
			if err != nil {
				return err
			}
			break
		}
	}
	err = v.client.Instance.Delete(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vultr) ListServer() ([]automation.RessourceResults, error) {
	instances, _, err := v.client.Instance.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, instance := range instances {
		if strings.Contains(instance.Tag, common.InstanceTag) {
			result = append(result, automation.RessourceResults{
				ID:       instance.ID,
				PublicIP: instance.MainIP,
				Name:     instance.Label,
				Region:   instance.Region,
				Tags:     instance.Tag,
			})
		}
	}
	return result, nil
}

func (v *Vultr) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := v.client.Instance.Get(context.Background(), id)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), instance.MainIP, "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vultr) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := v.client.Instance.Get(context.Background(), id)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), instance.MainIP, "root")
	err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)))
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service")
	if err != nil {
		return err
	}
	return nil
}

func (v *Vultr) GetServer(id string, _ automation.ServerArgs) (*automation.RessourceResults, error) {
	instance, err := v.client.Instance.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &automation.RessourceResults{
		ID:       instance.ID,
		Name:     instance.Label,
		Region:   instance.Region,
		PublicIP: instance.MainIP,
		Tags:     instance.Tag,
	}, err
}
