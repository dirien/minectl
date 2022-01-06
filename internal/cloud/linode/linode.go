package linode

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/minectl/internal/update"

	"github.com/linode/linodego"
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/oauth2"
)

type Linode struct {
	client linodego.Client
	tmpl   *minctlTemplate.Template
}

func NewLinode(apiToken string) (*Linode, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	linodeClient := linodego.NewClient(oauth2Client)
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	linode := &Linode{
		client: linodeClient,
		tmpl:   tmpl,
	}
	return linode, nil
}

func (l *Linode) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	ubuntuImage := "linode/ubuntu20.04"
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	key, err := l.client.CreateSSHKey(context.Background(), linodego.SSHKeyCreateOptions{
		SSHKey: strings.TrimSpace(string(pubKeyFile)),
		Label:  fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
	})
	if err != nil {
		return nil, err
	}

	var volume *linodego.Volume
	var mount string
	if args.MinecraftResource.GetVolumeSize() > 0 {
		volume, err = l.client.CreateVolume(context.Background(), linodego.VolumeCreateOptions{
			Label:  fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()),
			Size:   args.MinecraftResource.GetVolumeSize(),
			Region: args.MinecraftResource.GetRegion(),
		})
		mount = "sdc"
	}
	if err != nil {
		return nil, err
	}

	userData, err := l.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Mount: mount, Name: minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	stackscript, err := l.client.CreateStackscript(context.Background(), linodego.StackscriptCreateOptions{
		IsPublic: false,
		Label:    fmt.Sprintf("%s-stackscript", args.MinecraftResource.GetName()),
		Images:   []string{ubuntuImage},
		Script:   userData,
	})
	if err != nil {
		return nil, err
	}
	rootPassword, err := password.Generate(16, 4, 0, false, true)
	if err != nil {
		return nil, err
	}
	instance, err := l.client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Label:          args.MinecraftResource.GetName(),
		Region:         args.MinecraftResource.GetRegion(),
		Image:          ubuntuImage,
		Type:           args.MinecraftResource.GetSize(),
		AuthorizedKeys: []string{key.SSHKey},
		StackScriptID:  stackscript.ID,
		RootPass:       rootPassword,
		Tags:           []string{common.InstanceTag, args.MinecraftResource.GetEdition()},
	})
	if err != nil {
		return nil, err
	}

	if args.MinecraftResource.GetVolumeSize() > 0 {
		_, err = l.client.AttachVolume(context.Background(), volume.ID, &linodego.VolumeAttachOptions{
			LinodeID: instance.ID,
		})
		if err != nil {
			return nil, err
		}
	}

	_, err = l.client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, 600)
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       strconv.Itoa(instance.ID),
		Name:     instance.Label,
		Region:   instance.Region,
		PublicIP: instance.IPv4[0].String(),
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}

func (l *Linode) DeleteServer(id string, args automation.ServerArgs) error {
	keys, err := l.client.ListSSHKeys(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, key := range keys {
		if key.Label == fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()) {
			err := l.client.DeleteSSHKey(context.Background(), key.ID)
			if err != nil {
				return err
			}
		}
	}
	volumes, err := l.client.ListVolumes(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, volume := range volumes {
		if volume.Label == fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()) {
			err := l.client.DetachVolume(context.Background(), volume.ID)
			if err != nil {
				return err
			}
			// wait 40 secs for detach volume, to be sure.
			time.Sleep(40 * time.Second)
			err = l.client.DeleteVolume(context.Background(), volume.ID)
			if err != nil {
				return err
			}
		}
	}

	stackscripts, err := l.client.ListStackscripts(context.Background(), &linodego.ListOptions{Filter: "{\"mine\":true}"})
	if err != nil {
		return err
	}
	for _, stackscript := range stackscripts {
		if stackscript.Label == fmt.Sprintf("%s-stackscript", args.MinecraftResource.GetName()) {
			err := l.client.DeleteStackscript(context.Background(), stackscript.ID)
			if err != nil {
				return err
			}
		}
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	err = l.client.DeleteInstance(context.Background(), intID)
	if err != nil {
		return err
	}
	return nil
}

func (l *Linode) ListServer() ([]automation.ResourceResults, error) {
	servers, err := l.client.ListInstances(context.Background(), linodego.NewListOptions(0, "{\"tags\":\"minectl\"}"))
	if err != nil {
		return nil, err
	}
	var result []automation.ResourceResults
	for _, server := range servers {
		result = append(result, automation.ResourceResults{
			ID:       strconv.Itoa(server.ID),
			Name:     server.Label,
			Region:   server.Region,
			PublicIP: server.IPv4[0].String(),
			Tags:     strings.Join(server.Tags, ","),
		})
	}
	return result, nil
}

func (l *Linode) UpdateServer(id string, args automation.ServerArgs) error {
	intID, _ := strconv.Atoi(id)
	instance, err := l.client.GetInstance(context.Background(), intID)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.IPv4[0].String(), "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (l *Linode) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	intID, _ := strconv.Atoi(id)
	instance, err := l.client.GetInstance(context.Background(), intID)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.IPv4[0].String(), "root")
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

func (l *Linode) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	intID, _ := strconv.Atoi(id)
	instance, err := l.client.GetInstance(context.Background(), intID)
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       strconv.Itoa(instance.ID),
		Name:     instance.Label,
		Region:   instance.Region,
		PublicIP: instance.IPv4[0].String(),
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}
