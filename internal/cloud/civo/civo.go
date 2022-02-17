package civo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/minectl/internal/update"

	"github.com/civo/civogo"
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
)

type Civo struct {
	client *civogo.Client
	tmpl   *minctlTemplate.Template
}

func NewCivo(apiKey, region string) (*Civo, error) {
	client, err := civogo.NewClient(apiKey, region)
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

func (c *Civo) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	sshPubKey, err := c.client.NewSSHKey(fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()), string(pubKeyFile))
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Civo SSH Key created", "id", sshPubKey.ID)
	network, err := c.client.GetDefaultNetwork()
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Civo get default network created", "network", network)

	template, err := c.client.FindDiskImage("ubuntu-focal")
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Civo get disk image", "template", template)
	config, err := c.client.NewInstanceConfig()
	if err != nil {
		return nil, err
	}
	config.TemplateID = template.ID
	config.Size = args.MinecraftResource.GetSize()
	config.Hostname = args.MinecraftResource.GetName()
	config.Region = args.MinecraftResource.GetRegion()
	config.SSHKeyID = sshPubKey.ID
	config.PublicIPRequired = "create"
	config.InitialUser = "root"
	config.Tags = []string{common.InstanceTag, args.MinecraftResource.GetEdition()}

	script, err := c.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}
	config.Script = script

	instance, err := c.client.CreateInstance(config)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Civo create instance", "instance", instance)

	if args.MinecraftResource.GetEdition() == "bedrock" || args.MinecraftResource.GetEdition() == "nukkit" || args.MinecraftResource.GetEdition() == "powernukkit" {
		createRule := true
		firewall, err := c.client.NewFirewall(fmt.Sprintf("%s-fw", args.MinecraftResource.GetName()), network.ID, &createRule)
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
		zap.S().Infow("Civo create firewall", "firewall", firewall)
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
	zap.S().Infow("Civo instance ready", "instance", instance)
	region := c.client.Region
	if len(instance.Region) > 0 {
		region = instance.Region
	}
	return &automation.ResourceResults{
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
	zap.S().Infow("Civo delete instance", "id", id)
	pubKeyFile, err := c.client.FindSSHKey(fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()))
	if err != nil {
		return err
	}
	_, err = c.client.DeleteSSHKey(pubKeyFile.ID)
	if err != nil {
		return err
	}
	zap.S().Infow("Civo delete ssh key", "pubKeyFile", pubKeyFile)
	if args.MinecraftResource.GetEdition() == "bedrock" || args.MinecraftResource.GetEdition() == "nukkit" || args.MinecraftResource.GetEdition() == "powernukkit" {
		firewall, err := c.client.FindFirewall(fmt.Sprintf("%s-fw", args.MinecraftResource.GetName()))
		if err != nil {
			return err
		}
		_, err = c.client.DeleteFirewall(firewall.ID)
		if err != nil {
			return err
		}
		zap.S().Infow("Civo delete firewall", "firewall", firewall)
	}
	return nil
}

func (c *Civo) ListServer() ([]automation.ResourceResults, error) {
	var result []automation.ResourceResults
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
					result = append(result, automation.ResourceResults{
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
	if len(result) > 0 {
		zap.S().Infow("Civo list all instances", "list", result)
	} else {
		zap.S().Infow("No minectl instances found")
	}
	return result, nil
}

func (c *Civo) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := c.client.GetInstance(id)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.PublicIP, "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	zap.S().Infow("Civo minectl server updated", "instance", instance)
	return nil
}

func (c *Civo) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := c.client.GetInstance(id)
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
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "instance", instance)
	return nil
}

func (c *Civo) GetServer(id string, _ automation.ServerArgs) (*automation.ResourceResults, error) {
	instance, err := c.client.GetInstance(id)
	if err != nil {
		return nil, err
	}
	region := c.client.Region
	if len(instance.Region) > 0 {
		region = instance.Region
	}
	return &automation.ResourceResults{
		ID:       instance.ID,
		Name:     instance.Hostname,
		Region:   region,
		PublicIP: instance.PublicIP,
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}
