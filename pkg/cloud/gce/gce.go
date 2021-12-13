package gce

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/minectl/pkg/update"

	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/oslogin/v1"
)

const doneStatus = "DONE"

type Credentials struct {
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	ClientID    string `json:"client_id"`
}

type GCE struct {
	client             *compute.Service
	user               *oslogin.Service
	projectID          string
	serviceAccountName string
	serviceAccountID   string
	zone               string
	tmpl               *minctlTemplate.Template
}

func NewGCE(keyfile, zone string) (*GCE, error) {
	file, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	var cred Credentials
	err = json.Unmarshal(file, &cred)
	if err != nil {
		return nil, err
	}
	computeService, err := compute.NewService(context.Background(), option.WithCredentialsJSON(file))
	if err != nil {
		return nil, err
	}

	userService, err := oslogin.NewService(context.Background(), option.WithCredentialsJSON(file))
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	return &GCE{
		client:             computeService,
		projectID:          cred.ProjectID,
		user:               userService,
		serviceAccountName: cred.ClientEmail,
		serviceAccountID:   cred.ClientID,
		zone:               zone,
		tmpl:               tmpl,
	}, nil
}

func (g *GCE) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	imageURL := "projects/ubuntu-os-cloud/global/images/ubuntu-2004-focal-v20210720"

	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}

	_, err = g.user.Users.ImportSshPublicKey(fmt.Sprintf("users/%s", g.serviceAccountName), &oslogin.SshPublicKey{
		Key:                string(pubKeyFile),
		ExpirationTimeUsec: 0,
	}).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	stillCreating := true
	var mount string
	if args.MinecraftResource.GetVolumeSize() > 0 {
		diskInsertOp, err := g.client.Disks.Insert(g.projectID, args.MinecraftResource.GetRegion(), &compute.Disk{
			Name:   fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()),
			SizeGb: int64(args.MinecraftResource.GetVolumeSize()),
			Type:   fmt.Sprintf("zones/%s/diskTypes/pd-standard", args.MinecraftResource.GetRegion()),
		}).Context(context.Background()).Do()
		if err != nil {
			return nil, err
		}

		for stillCreating {
			diskInsertOps, err := g.client.ZoneOperations.Get(g.projectID, args.MinecraftResource.GetRegion(), diskInsertOp.Name).Context(context.Background()).Do()
			if err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			if diskInsertOps.Status == doneStatus {
				stillCreating = false
			} else {
				time.Sleep(2 * time.Second)
			}
		}
		mount = "sdb"
	}

	userData, err := g.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Mount: mount, Name: minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	oslogin := "TRUE"
	autoRestart := true
	instance := &compute.Instance{
		Name:        args.MinecraftResource.GetName(),
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", args.MinecraftResource.GetRegion(), args.MinecraftResource.GetSize()),
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				DiskSizeGb: 10,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: imageURL,
				},
			},
		},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "enable-oslogin",
					Value: &oslogin,
				},
				{
					Key:   "startup-script",
					Value: &userData,
				},
			},
		},
		Scheduling: &compute.Scheduling{
			AutomaticRestart:  &autoRestart,
			OnHostMaintenance: "MIGRATE",
			Preemptible:       false,
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: "/global/networks/default",
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: g.serviceAccountName,
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				},
			},
		},
		Labels: map[string]string{
			common.InstanceTag: "true",
		},
		Tags: &compute.Tags{
			Items: []string{common.InstanceTag, args.MinecraftResource.GetEdition()},
		},
	}
	if args.MinecraftResource.GetVolumeSize() > 0 {
		instance.Disks = append(instance.Disks, &compute.AttachedDisk{
			Source: fmt.Sprintf("zones/%s/disks/%s-vol", args.MinecraftResource.GetRegion(),
				args.MinecraftResource.GetName()),
		})
	}

	insertInstanceOp, err := g.client.Instances.Insert(g.projectID, args.MinecraftResource.GetRegion(), instance).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	stillCreating = true
	for stillCreating {
		insertInstanceOp, err := g.client.ZoneOperations.Get(g.projectID, args.MinecraftResource.GetRegion(), insertInstanceOp.Name).Context(context.Background()).Do()
		if err != nil {
			return nil, err
		}
		if insertInstanceOp.Status == doneStatus {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	firewallRule := &compute.Firewall{
		Name:        fmt.Sprintf("%s-fw", args.MinecraftResource.GetName()),
		Description: "Firewall rule created by minectl",
		Network:     fmt.Sprintf("projects/%s/global/networks/default", g.projectID),
		Allowed: []*compute.FirewallAllowed{
			{
				IPProtocol: "tcp",
			},
		},
		SourceRanges: []string{"0.0.0.0/0"},
		Direction:    "INGRESS",
		TargetTags:   []string{common.InstanceTag},
	}
	_, err = g.client.Firewalls.Insert(g.projectID, firewallRule).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	instanceListOp, err := g.client.Instances.List(g.projectID, args.MinecraftResource.GetRegion()).
		Filter(fmt.Sprintf("(name=%s)", args.MinecraftResource.GetName())).
		Context(context.Background()).
		Do()
	if err != nil {
		return nil, err
	}

	if len(instanceListOp.Items) == 1 {
		instance := instanceListOp.Items[0]
		ip := instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
		return &automation.RessourceResults{
			ID:       strconv.Itoa(int(instance.Id)),
			Name:     instance.Name,
			Region:   instance.Zone,
			PublicIP: ip,
			Tags:     strings.Join(instance.Tags.Items, ","),
		}, err
	}
	return nil, errors.New("no instances created")
}

func (g *GCE) DeleteServer(id string, args automation.ServerArgs) error {
	profileGetOp, err := g.user.Users.GetLoginProfile(fmt.Sprintf("users/%s", g.serviceAccountName)).Context(context.Background()).Do()
	if err != nil {
		return err
	}
	for _, posixAccount := range profileGetOp.PosixAccounts {
		_, err := g.user.Users.Projects.Delete(posixAccount.Name).Context(context.Background()).Do()
		if err != nil {
			return err
		}
	}
	for _, publicKey := range profileGetOp.SshPublicKeys {
		_, err = g.user.Users.SshPublicKeys.Delete(publicKey.Name).Context(context.Background()).Do()
		if err != nil {
			return err
		}
	}
	instancesListOp, err := g.client.Instances.List(g.projectID, args.MinecraftResource.GetRegion()).
		Filter(fmt.Sprintf("(id=%s)", id)).
		Context(context.Background()).
		Do()
	if err != nil {
		return err
	}
	if len(instancesListOp.Items) == 1 {
		instanceDeleteOp, err := g.client.Instances.Delete(g.projectID, args.MinecraftResource.GetRegion(), instancesListOp.Items[0].Name).
			Context(context.Background()).
			Do()
		if err != nil {
			return err
		}
		stillDeleting := true
		for stillDeleting {
			instanceDeleteOp, err := g.client.ZoneOperations.Get(g.projectID, args.MinecraftResource.GetRegion(), instanceDeleteOp.Name).Context(context.Background()).Do()
			if err != nil {
				return err
			}
			if instanceDeleteOp.Status == doneStatus {
				stillDeleting = false
			} else {
				time.Sleep(2 * time.Second)
			}
		}

	}

	diskListOp, err := g.client.Disks.List(g.projectID, args.MinecraftResource.GetRegion()).
		Filter(fmt.Sprintf("(name=%s)", fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()))).
		Context(context.Background()).
		Do()
	if err != nil {
		return err
	}
	for _, disk := range diskListOp.Items {
		_, err := g.client.Disks.Delete(g.projectID, args.MinecraftResource.GetRegion(), disk.Name).Context(context.Background()).Do()
		if err != nil {
			return err
		}
	}

	firewallListOps, err := g.client.Firewalls.List(g.projectID).Filter(fmt.Sprintf("(name=%s)", fmt.Sprintf("%s-fw", args.MinecraftResource.GetName()))).Context(context.Background()).Do()
	if err != nil {
		return err
	}
	for _, firewall := range firewallListOps.Items {
		_, err := g.client.Firewalls.Delete(g.projectID, firewall.Name).Context(context.Background()).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GCE) ListServer() ([]automation.RessourceResults, error) {
	instanceListOp, err := g.client.Instances.List(g.projectID, g.zone).
		Filter(fmt.Sprintf("(labels.%s=true)", common.InstanceTag)).
		Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, instance := range instanceListOp.Items {
		result = append(result, automation.RessourceResults{
			ID:       strconv.Itoa(int(instance.Id)),
			Name:     instance.Name,
			Region:   instance.Zone,
			PublicIP: instance.NetworkInterfaces[0].AccessConfigs[0].NatIP,
			Tags:     strings.Join(instance.Tags.Items, ","),
		})
	}
	return result, nil
}

func (g *GCE) UpdateServer(id string, args automation.ServerArgs) error {
	instancesListOp, err := g.client.Instances.List(g.projectID, args.MinecraftResource.GetRegion()).
		Filter(fmt.Sprintf("(id=%s)", id)).
		Context(context.Background()).
		Do()
	if err != nil {
		return err
	}
	if len(instancesListOp.Items) == 1 {
		instance := instancesListOp.Items[0]
		remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), instance.NetworkInterfaces[0].AccessConfigs[0].NatIP, fmt.Sprintf("sa_%s", g.serviceAccountID))
		err = remoteCommand.UpdateServer(args.MinecraftResource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GCE) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instancesListOp, err := g.client.Instances.List(g.projectID, args.MinecraftResource.GetRegion()).
		Filter(fmt.Sprintf("(id=%s)", id)).
		Context(context.Background()).
		Do()
	if err != nil {
		return err
	}
	if len(instancesListOp.Items) == 1 {
		instance := instancesListOp.Items[0]
		remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), instance.NetworkInterfaces[0].AccessConfigs[0].NatIP, fmt.Sprintf("sa_%s", g.serviceAccountID))
		err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)))
		if err != nil {
			return err
		}
		_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service")
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GCE) GetServer(id string, args automation.ServerArgs) (*automation.RessourceResults, error) {
	instancesListOp, err := g.client.Instances.List(g.projectID, args.MinecraftResource.GetRegion()).
		Filter(fmt.Sprintf("(id=%s)", id)).
		Context(context.Background()).
		Do()
	if err != nil {
		return nil, err
	}
	if len(instancesListOp.Items) == 1 {
		instance := instancesListOp.Items[0]
		ip := instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
		return &automation.RessourceResults{
			ID:       strconv.Itoa(int(instance.Id)),
			Name:     instance.Name,
			Region:   instance.Zone,
			PublicIP: ip,
			Tags:     strings.Join(instance.Tags.Items, ","),
		}, err
	}
	return nil, nil
}
