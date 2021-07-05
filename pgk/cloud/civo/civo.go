package civo

import (
	_ "embed"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/common"
	"io/ioutil"
	"strings"
	"time"
)

//go:embed civo.sh
var bash string

type Civo struct {
	client *civogo.Client
}

func NewCivo(APIKey, region string) (*Civo, error) {
	client, err := civogo.NewClient(APIKey, region)
	if err != nil {
		return nil, err
	}
	do := &Civo{
		client: client,
	}
	return do, nil
}

func (c *Civo) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(args.SSH)
	if err != nil {
		return nil, err
	}
	sshPubKey, err := c.client.NewSSHKey(fmt.Sprintf("%s-ssh", args.StackName), string(pubKeyFile))
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
	config.Size = args.Size
	config.Hostname = args.StackName
	config.Region = args.Region
	config.SSHKeyID = sshPubKey.ID
	config.PublicIPRequired = "create"
	config.InitialUser = "root"
	config.Script = strings.Replace(bash, "<properties>", args.Properties, -1)

	instance, err := c.client.CreateInstance(config)
	if err != nil {
		return nil, err
	}

	stillCreating := true
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("Creating instance (%s)... ", common.Green(instance.Hostname))
	s.FinalMSG = fmt.Sprintf("\nInstance (%s) created\n", common.Green(instance.Hostname))
	s.Start()

	for stillCreating {
		instance, err = c.client.FindInstance(instance.ID)
		if err != nil {
			return nil, err
		}
		if instance.Status == "ACTIVE" {
			stillCreating = false
			s.Stop()
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	return &automation.RessourceResults{
		ID:       instance.ID,
		PublicIP: instance.PublicIP,
	}, err
}

func (c *Civo) DeleteServer(id string, args automation.ServerArgs) error {
	common.PrintMixedGreen("Delete instance (%s)... ", id)

	_, err := c.client.DeleteInstance(id)
	if err != nil {
		return err
	}
	pubKeyFile, err := c.client.FindSSHKey(fmt.Sprintf("%s-ssh", args.StackName))
	if err != nil {
		return err
	}
	_, err = c.client.DeleteSSHKey(pubKeyFile.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Civo) ListServer(args automation.ServerArgs) (*[]automation.RessourceResults, error) {
	panic("implement me")
}

func (c *Civo) UpdateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
}
