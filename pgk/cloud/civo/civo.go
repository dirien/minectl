package civo

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/pgk/update"

	"github.com/civo/civogo"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/common"
	minctlTemplate "github.com/minectl/pgk/template"
)

type Civo struct {
	client *civogo.Client
	tmpl   *minctlTemplate.Template
}

func NewCivo(APIKey, region string) (*Civo, error) {
	client, err := civogo.NewClient(APIKey, region)
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	do := &Civo{
		client: client,
		tmpl:   tmpl,
	}
	return do, nil
}

func (c *Civo) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftServer.GetSSH()))
	if err != nil {
		return nil, err
	}
	sshPubKey, err := c.client.NewSSHKey(fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()), string(pubKeyFile))
	if err != nil {
		return nil, err
	}

	network, err := c.client.GetDefaultNetwork()
	if err != nil {
		return nil, err
	}

	template, err := c.client.FindDiskImage("ubuntu-focal")
	if err != nil {
		return nil, err
	}
	config, err := c.client.NewInstanceConfig()
	if err != nil {
		return nil, err
	}
	config.TemplateID = template.ID
	config.Size = args.MinecraftServer.GetSize()
	config.Hostname = args.MinecraftServer.GetName()
	config.Region = args.MinecraftServer.GetRegion()
	config.SSHKeyID = sshPubKey.ID
	config.PublicIPRequired = "create"
	config.InitialUser = "root"
	config.Tags = []string{common.InstanceTag, args.MinecraftServer.GetEdition()}

	script, err := c.tmpl.GetTemplate(args.MinecraftServer, "", minctlTemplate.TemplateBash)
	if err != nil {
		return nil, err
	}
	config.Script = script

	instance, err := c.client.CreateInstance(config)
	if err != nil {
		return nil, err
	}

	if args.MinecraftServer.GetEdition() == "bedrock" {
		firewall, err := c.client.NewFirewall(fmt.Sprintf("%s-fw", args.MinecraftServer.GetName()), network.ID)
		if err != nil {
			return nil, err
		}
		_, err = c.client.NewFirewallRule(&civogo.FirewallRuleConfig{
			FirewallID: firewall.ID,
			Protocol:   "udp",
			StartPort:  "19132",
			EndPort:    "19133",
			Cidr: []string{
				"0.0.0.0/0",
			},
			Label: "Minecraft Bedrock UDP",
		})
		if err != nil {
			return nil, err
		}
		_, err = c.client.SetInstanceFirewall(instance.ID, firewall.ID)
		if err != nil {
			return nil, err
		}
	}

	stillCreating := true
	for stillCreating {
		instance, err = c.client.FindInstance(instance.ID)
		if err != nil {
			return nil, err
		}
		if instance.Status == "ACTIVE" {
			stillCreating = false
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	instance, err = c.client.FindInstance(instance.ID)
	if err != nil {
		return nil, err
	}
	region := c.client.Region
	if len(instance.Region) > 0 {
		region = instance.Region
	}
	return &automation.RessourceResults{
		ID:       instance.ID,
		Name:     instance.Hostname,
		Region:   region,
		PublicIP: instance.PublicIP,
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}

func (c *Civo) DeleteServer(id string, args automation.ServerArgs) error {
	_, err := c.client.DeleteInstance(id)
	if err != nil {
		return err
	}
	pubKeyFile, err := c.client.FindSSHKey(fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()))
	if err != nil {
		return err
	}
	_, err = c.client.DeleteSSHKey(pubKeyFile.ID)
	if err != nil {
		return err
	}
	if args.MinecraftServer.GetEdition() == "bedrock" {
		firewall, err := c.client.FindFirewall(fmt.Sprintf("%s-fw", args.MinecraftServer.GetName()))
		if err != nil {
			return err
		}
		_, err = c.client.DeleteFirewall(firewall.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Civo) ListServer() ([]automation.RessourceResults, error) {
	var result []automation.RessourceResults
	instances, err := c.client.ListAllInstances()
	if err != nil {
		return nil, err
	}
	for _, instance := range instances {
		if instance.Tags[0] == common.InstanceTag {
			for _, tag := range instance.Tags {
				if tag == common.InstanceTag {
					region := c.client.Region
					if len(instance.Region) > 0 {
						region = instance.Region
					}
					result = append(result, automation.RessourceResults{
						ID:       instance.ID,
						PublicIP: instance.PublicIP,
						Name:     instance.Hostname,
						Region:   region,
						Tags:     strings.Join(instance.Tags, ","),
					})
				}
			}
		}
	}
	return result, nil
}

func (c *Civo) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := c.client.GetInstance(id)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftServer.GetSSH(), instance.PublicIP, "root")
	err = remoteCommand.UpdateServer(args.MinecraftServer)
	if err != nil {
		return err
	}
	return nil
}

func (c *Civo) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := c.client.GetInstance(id)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftServer.GetSSH(), instance.PublicIP, "root")
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
