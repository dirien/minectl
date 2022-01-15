package exoscale

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/exoscale/egoscale"
	v2 "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/minectl/internal/automation"
	minctlTemplate "github.com/minectl/internal/template"
	"github.com/minectl/internal/update"
)

var userAgent string

func String(v string) *string {
	return &v
}

type Exoscale struct {
	client   *egoscale.Client
	clientv2 *v2.Client
	tmpl     *minctlTemplate.Template
}

type defaultTransport struct {
	next http.RoundTripper
}

func (t *defaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", userAgent)

	resp, err := t.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewExoscale(apiKey, apiSecret string) (*Exoscale, error) {
	httpClient := cleanhttp.DefaultPooledClient()
	httpClient.Transport = &defaultTransport{next: httpClient.Transport}

	client := egoscale.NewClient("https://api.exoscale.com/v1", apiKey, apiSecret,
		egoscale.WithHTTPClient(httpClient),
		egoscale.WithoutV2Client())

	clientv2, err := v2.NewClient(apiKey, apiSecret,
		v2.ClientOptWithAPIEndpoint("https://api.exoscale.com/v1"),
		v2.ClientOptWithTimeout(5*time.Minute))
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	es := &Exoscale{
		client:   client,
		clientv2: clientv2,
		tmpl:     tmpl,
	}
	return es, nil
}

func (e *Exoscale) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	ctx := context.Background()

	_, cidr, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return nil, err
	}

	_, err = e.client.Request(egoscale.CreateSecurityGroup{
		Name: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
	})
	if err != nil {
		return nil, err
	}

	_, err = e.client.Request(egoscale.AuthorizeSecurityGroupIngress{
		SecurityGroupName: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
		Protocol:          "TCP",
		StartPort:         uint16(args.MinecraftResource.GetSSHPort()),
		EndPort:           uint16(args.MinecraftResource.GetSSHPort()),
		CIDRList: []egoscale.CIDR{
			{
				IPNet: *cidr,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if args.MinecraftResource.GetEdition() == "bedrock" || args.MinecraftResource.GetEdition() == "nukkit" || args.MinecraftResource.GetEdition() == "powernukkit" {
		_, err = e.client.Request(egoscale.AuthorizeSecurityGroupIngress{
			SecurityGroupName: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
			Protocol:          "UDP",
			StartPort:         uint16(args.MinecraftResource.GetPort()),
			EndPort:           uint16(args.MinecraftResource.GetPort()),
			CIDRList: []egoscale.CIDR{
				{
					IPNet: *cidr,
				},
			},
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = e.client.Request(egoscale.AuthorizeSecurityGroupIngress{
			SecurityGroupName: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
			Protocol:          "TCP",
			StartPort:         uint16(args.MinecraftResource.GetPort()),
			EndPort:           uint16(args.MinecraftResource.GetPort()),
			CIDRList: []egoscale.CIDR{
				{
					IPNet: *cidr,
				},
			},
		})
		if err != nil {
			return nil, err
		}
		if args.MinecraftResource.HasRCON() {
			_, err = e.client.Request(egoscale.AuthorizeSecurityGroupIngress{
				SecurityGroupName: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
				Protocol:          "TCP",
				StartPort:         uint16(args.MinecraftResource.GetRCONPort()),
				EndPort:           uint16(args.MinecraftResource.GetRCONPort()),
				CIDRList: []egoscale.CIDR{
					{
						IPNet: *cidr,
					},
				},
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if args.MinecraftResource.HasMonitoring() {
		_, err = e.client.Request(egoscale.AuthorizeSecurityGroupIngress{
			SecurityGroupName: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
			Protocol:          "TCP",
			StartPort:         uint16(9090),
			EndPort:           uint16(9090),
			CIDRList: []egoscale.CIDR{
				{
					IPNet: *cidr,
				},
			},
		})
		if err != nil {
			return nil, err
		}
	}

	script, err := e.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	listTemplates, err := e.clientv2.ListTemplates(ctx, args.MinecraftResource.GetRegion(), v2.ListTemplatesWithFamily("ubuntu"))
	if err != nil {
		return nil, err
	}
	var templateID string
	for _, template := range listTemplates {
		if strings.Contains(*template.Name, "Ubuntu 20.04 LTS") {
			templateID = *template.ID
			break
		}
	}
	instanceTypes, err := e.clientv2.ListInstanceTypes(ctx, args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}
	var instanceTypeID string
	for _, instanceType := range instanceTypes {
		if strings.Contains(*instanceType.Size, args.MinecraftResource.GetSize()) {
			instanceTypeID = *instanceType.ID
			break
		}
	}

	zones, err := e.client.ListWithContext(ctx, &egoscale.Zone{})
	if err != nil {
		return nil, err
	}
	var zoneID *egoscale.UUID
	for _, value := range zones {
		zone := value.(*egoscale.Zone)
		if zone.Name == args.MinecraftResource.GetRegion() {
			zoneID = zone.ID
			break
		}
	}

	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	sshPubKey, err := e.clientv2.RegisterSSHKey(ctx, args.MinecraftResource.GetRegion(), fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()), string(pubKeyFile))
	if err != nil {
		return nil, err
	}

	groups, err := e.clientv2.ListSecurityGroups(ctx, args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}
	var groupID string
	for _, group := range groups {
		if *group.Name == fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()) {
			groupID = *group.ID
			break
		}
	}

	virtualMachine := egoscale.DeployVirtualMachine{
		Name:              args.MinecraftResource.GetName(),
		TemplateID:        egoscale.MustParseUUID(templateID),
		ServiceOfferingID: egoscale.MustParseUUID(instanceTypeID),
		RootDiskSize:      50,
		SecurityGroupIDs: []egoscale.UUID{
			*egoscale.MustParseUUID(groupID),
		},
		UserData: base64.StdEncoding.EncodeToString([]byte(script)),
		KeyPair:  *sshPubKey.Name,
		ZoneID:   zoneID,
	}
	resp, err := e.client.Request(virtualMachine)
	if err != nil {
		return nil, err
	}

	vm := resp.(*egoscale.VirtualMachine)
	return &automation.ResourceResults{
		ID:       vm.ID.String(),
		Name:     vm.Name,
		Region:   vm.ZoneName,
		PublicIP: vm.PublicIP,
	}, err
}

func (e *Exoscale) DeleteServer(id string, args automation.ServerArgs) error {
	ctx := context.Background()

	virtualMachine := egoscale.DestroyVirtualMachine{
		ID: egoscale.MustParseUUID(id),
	}
	_, err := e.client.Request(virtualMachine)
	if err != nil {
		return err
	}
	securityGroups, err := e.client.ListWithContext(ctx, &egoscale.SecurityGroup{
		Name: fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
	})
	if err != nil || len(securityGroups) == 0 {
		return err
	}
	securityGroup := securityGroups[0].(*egoscale.SecurityGroup)

	err = e.clientv2.DeleteSecurityGroup(ctx, args.MinecraftResource.GetRegion(), &v2.SecurityGroup{
		ID: String(securityGroup.ID.String()),
	})
	if err != nil {
		return err
	}

	sshKey, err := e.clientv2.GetSSHKey(ctx, args.MinecraftResource.GetRegion(), fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()))
	if err != nil {
		return err
	}
	err = e.clientv2.DeleteSSHKey(ctx, args.MinecraftResource.GetRegion(), sshKey)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exoscale) ListServer() ([]automation.ResourceResults, error) {
	panic("List Server is not possible with Exoscale, as it does not support labels in v1")
}

func (e *Exoscale) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := e.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.PublicIP, "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exoscale) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := e.GetServer(id, args)
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

func (e *Exoscale) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	ctx := context.Background()

	instance, err := e.clientv2.GetInstance(ctx, args.MinecraftResource.GetRegion(), id)
	if err != nil {
		return nil, err
	}

	return &automation.ResourceResults{
		ID:       *instance.ID,
		Name:     *instance.Name,
		Region:   *instance.Zone,
		PublicIP: instance.PublicIPAddress.String(),
	}, err
}
