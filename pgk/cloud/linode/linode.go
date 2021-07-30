package linode

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/minectl/pgk/update"

	"github.com/linode/linodego"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/common"
	minctlTemplate "github.com/minectl/pgk/template"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/oauth2"
)

type Linode struct {
	client linodego.Client
	tmpl   *minctlTemplate.Template
}

func NewLinode(APItoken string) (*Linode, error) {

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: APItoken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	linodeClient := linodego.NewClient(oauth2Client)
	tmpl, err := minctlTemplate.NewTemplateBash("sdc")
	if err != nil {
		return nil, err
	}
	linode := &Linode{
		client: linodeClient,
		tmpl:   tmpl,
	}
	return linode, nil
}

func (l *Linode) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	os := "linode/ubuntu20.04"
	pubKeyFile, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftServer.GetSSH()))
	if err != nil {
		return nil, err
	}
	key, err := l.client.CreateSSHKey(context.Background(), linodego.SSHKeyCreateOptions{
		SSHKey: strings.TrimSpace(string(pubKeyFile)),
		Label:  fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()),
	})
	if err != nil {
		return nil, err
	}
	volume, err := l.client.CreateVolume(context.Background(), linodego.VolumeCreateOptions{
		Label:  fmt.Sprintf("%s-vol", args.MinecraftServer.GetName()),
		Size:   args.MinecraftServer.GetVolumeSize(),
		Region: args.MinecraftServer.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	userData, err := l.tmpl.GetTemplate(args.MinecraftServer, minctlTemplate.TemplateBash)
	if err != nil {
		return nil, err
	}

	stackscript, err := l.client.CreateStackscript(context.Background(), linodego.StackscriptCreateOptions{
		IsPublic: false,
		Label:    fmt.Sprintf("%s-stackscript", args.MinecraftServer.GetName()),
		Images:   []string{os},
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
		Label:          args.MinecraftServer.GetName(),
		Region:         args.MinecraftServer.GetRegion(),
		Image:          os,
		Type:           args.MinecraftServer.GetSize(),
		AuthorizedKeys: []string{key.SSHKey},
		StackScriptID:  stackscript.ID,
		RootPass:       rootPassword,
		Tags:           []string{common.InstanceTag, args.MinecraftServer.GetEdition()},
	})
	if err != nil {
		return nil, err
	}
	_, err = l.client.AttachVolume(context.Background(), volume.ID, &linodego.VolumeAttachOptions{
		LinodeID: instance.ID,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, 600)
	if err != nil {
		return nil, err
	}
	return &automation.RessourceResults{
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
		if key.Label == fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()) {
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
		if volume.Label == fmt.Sprintf("%s-vol", args.MinecraftServer.GetName()) {
			err := l.client.DetachVolume(context.Background(), volume.ID)
			if err != nil {
				return err
			}
			//wait 40 secs for detach volume, to be sure.
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
		if stackscript.Label == fmt.Sprintf("%s-stackscript", args.MinecraftServer.GetName()) {
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

func (l *Linode) ListServer() ([]automation.RessourceResults, error) {
	servers, err := l.client.ListInstances(context.Background(), linodego.NewListOptions(0, "{\"tags\":\"minectl\"}"))
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, server := range servers {
		result = append(result, automation.RessourceResults{
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

	remoteCommand := update.NewRemoteServer(args.MinecraftServer.GetSSH(), instance.IPv4[0].String(), "root")
	err = remoteCommand.UpdateServer(args.MinecraftServer, l.tmpl)
	if err != nil {
		return err
	}
	return nil
}
