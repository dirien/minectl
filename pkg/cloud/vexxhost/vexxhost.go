package vexxhost

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/pkg/update"
	"go.uber.org/zap"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	flavors2 "github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
)

type VEXXHOST struct {
	tmpl          *minctlTemplate.Template
	computeClient *gophercloud.ServiceClient
	networkClient *gophercloud.ServiceClient
	region        string
}

func getTags(edition string) map[string]string {
	return map[string]string{
		common.InstanceTag: "true",
		edition:            "true",
	}
}

func getTagKeys(tags map[string]string) []string {
	var keys []string
	for key := range tags {
		keys = append(keys, key)
	}
	return keys
}

func NewVEXXHOST() (*VEXXHOST, error) {
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainID:         os.Getenv("OS_PROJECT_DOMAIN_ID"),
		UserID:           os.Getenv("OS_USERID"),
		Passcode:         os.Getenv("OS_PASSCODE"),
		TenantID:         os.Getenv("OS_PROJECT_ID"),
		TenantName:       os.Getenv("OS_PROJECT_NAME"),
	}
	if err != nil {
		return nil, err
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}
	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		return nil, err
	}
	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		return nil, err
	}
	return &VEXXHOST{
		tmpl:          tmpl,
		computeClient: computeClient,
		networkClient: networkClient,
		region:        os.Getenv("OS_REGION_NAME"),
	}, nil
}

// CreateServer TODO: https://github.com/dirien/minectl/issues/299
func (v *VEXXHOST) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) { // nolint: gocyclo
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}

	keyPair, err := keypairs.Create(v.computeClient, keypairs.CreateOpts{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
		PublicKey: string(pubKeyFile),
	}).Extract()
	if err != nil {
		return nil, err
	}

	listOpts := images.ListOpts{
		Status: "active",
	}

	var image images.Image
	pager := images.ListDetail(v.computeClient, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}
		for _, i := range imageList {
			if strings.Contains(i.Name, "Ubuntu 20.04.3") && i.Metadata["architecture"] == "x86_64" {
				image = i
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	var flavor flavors2.Flavor

	flavorPager := flavors2.ListDetail(v.computeClient, flavors2.ListOpts{})
	err = flavorPager.EachPage(func(page pagination.Page) (bool, error) {
		flavorsList, err := flavors2.ExtractFlavors(page)
		if err != nil {
			return false, err
		}
		for _, i := range flavorsList {
			if i.Name == args.MinecraftResource.GetSize() {
				flavor = i
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	createOpts := secgroups.CreateOpts{
		Name:        fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()),
		Description: "minectl",
	}

	group, err := secgroups.Create(v.computeClient, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	err = v.createSecurityGroup(group, 22, "TCP")
	if err != nil {
		return nil, err
	}
	if args.MinecraftResource.GetEdition() == "bedrock" || args.MinecraftResource.GetEdition() == "nukkit" || args.MinecraftResource.GetEdition() == "powernukkit" {
		err = v.createSecurityGroup(group, args.MinecraftResource.GetPort(), "UDP")
		if err != nil {
			return nil, err
		}
	} else {
		err = v.createSecurityGroup(group, args.MinecraftResource.GetPort(), "TCP")
		if err != nil {
			return nil, err
		}
		if args.MinecraftResource.HasRCON() {
			err = v.createSecurityGroup(group, args.MinecraftResource.GetRCONPort(), "TCP")
			if err != nil {
				return nil, err
			}
		}
	}

	adminStateUp := true
	networkOpts := networks.CreateOpts{
		Name:         fmt.Sprintf("%s-net", args.MinecraftResource.GetName()),
		AdminStateUp: &adminStateUp,
	}

	network, err := networks.Create(v.networkClient, networkOpts).Extract()
	if err != nil {
		return nil, err
	}

	subnetOpts := subnets.CreateOpts{
		Name:      fmt.Sprintf("%s-subnet", args.MinecraftResource.GetName()),
		NetworkID: network.ID,
		CIDR:      "10.1.10.0/24",
		IPVersion: gophercloud.IPVersion(4),
		DNSNameservers: []string{
			"8.8.8.8",
			"8.8.4.4",
		},
	}

	subnet, err := subnets.Create(v.networkClient, subnetOpts).Extract()
	if err != nil {
		return nil, err
	}

	networkPager := networks.List(v.networkClient, networks.ListOpts{
		Name: "public",
	})
	var publicNetwork networks.Network
	err = networkPager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}
		for _, i := range networkList {
			publicNetwork = i
			break
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	gatewayInfo := &routers.GatewayInfo{
		NetworkID: publicNetwork.ID,
	}

	router, err := routers.Create(v.networkClient, routers.CreateOpts{
		Name:         fmt.Sprintf("%s-router", args.MinecraftResource.GetName()),
		AdminStateUp: &adminStateUp,
		GatewayInfo:  gatewayInfo,
	}).Extract()
	if err != nil {
		return nil, err
	}
	routers.AddInterface(v.networkClient, router.ID, &routers.AddInterfaceOpts{
		SubnetID: subnet.ID,
	})
	userData, err := v.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	server, err := servers.Create(v.computeClient, keypairs.CreateOptsExt{
		CreateOptsBuilder: servers.CreateOpts{
			Name: args.MinecraftResource.GetName(),
			SecurityGroups: []string{
				group.ID,
			},
			FlavorRef: flavor.ID,
			ImageRef:  image.ID,
			Networks: []servers.Network{
				{
					UUID: network.ID,
				},
			},
			Metadata: getTags(args.MinecraftResource.GetEdition()),
			UserData: []byte(base64.StdEncoding.EncodeToString([]byte(userData))),
		},
		KeyName: keyPair.Name,
	}).Extract()
	if err != nil {
		return nil, err
	}

	stillCreating := true
	for stillCreating {
		server, err = servers.Get(v.computeClient, server.ID).Extract()
		if err != nil {
			return nil, err
		}
		if server.Status == "ACTIVE" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	floatingIP, err := floatingips.Create(v.computeClient, floatingips.CreateOpts{
		Pool: "public",
	}).Extract()
	if err != nil {
		return nil, err
	}

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP.IP,
	}

	err = floatingips.AssociateInstance(v.computeClient, server.ID, associateOpts).ExtractErr()
	if err != nil {
		return nil, err
	}

	return &automation.RessourceResults{
		ID:       server.ID,
		Name:     server.Name,
		Region:   v.region,
		PublicIP: floatingIP.IP,
		Tags:     strings.Join(getTagKeys(server.Metadata), ","),
	}, err
}

func (v *VEXXHOST) createSecurityGroup(group *secgroups.SecurityGroup, port int, protocol string) error {
	ssh := secgroups.CreateRuleOpts{
		ParentGroupID: group.ID,
		FromPort:      port,
		ToPort:        port,
		IPProtocol:    protocol,
		CIDR:          "0.0.0.0/0",
	}

	_, err := secgroups.CreateRule(v.computeClient, ssh).Extract()
	if err != nil {
		return err
	}
	return nil
}

func (v *VEXXHOST) DeleteServer(id string, args automation.ServerArgs) error {
	server, err := servers.Get(v.computeClient, id).Extract()
	if err != nil {
		return err
	}
	err = keypairs.Delete(v.computeClient, fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()), &keypairs.DeleteOpts{}).Err
	if err != nil {
		return err
	}
	floatingIP, err := v.getFloatingIPByInstanceID(server.ID)
	if err != nil {
		return err
	}
	floatingips.DisassociateInstance(v.computeClient, id, floatingips.DisassociateOpts{
		FloatingIP: floatingIP.IP,
	})
	err = floatingips.Delete(v.computeClient, floatingIP.ID).Err
	if err != nil {
		return err
	}

	err = servers.Delete(v.computeClient, server.ID).Err
	if err != nil {
		return err
	}
	stillCreating := true
	for stillCreating {
		server, err = servers.Get(v.computeClient, server.ID).Extract()
		if err != nil {
			stillCreating = false
		}
		if server.Status == "DELETED" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	network, err := v.getNetworkByName(args)
	if err != nil {
		return err
	}
	subnet, err := v.getSubNetByName(args)
	if err != nil {
		return err
	}

	securityGroup, err := v.getSecurityGroupByName(args)
	if err != nil {
		return err
	}
	err = secgroups.Delete(v.computeClient, securityGroup.ID).Err
	if err != nil {
		return err
	}
	router, err := v.getRouterByName(args)
	if err != nil {
		return err
	}
	_, err = routers.RemoveInterface(v.networkClient, router.ID, routers.RemoveInterfaceOpts{
		SubnetID: subnet.ID,
	}).Extract()
	if err != nil {
		return err
	}
	err = routers.Delete(v.networkClient, router.ID).Err
	if err != nil {
		return err
	}
	err = subnets.Delete(v.networkClient, subnet.ID).Err
	if err != nil {
		return err
	}
	err = networks.Delete(v.networkClient, network.ID).Err
	if err != nil {
		return err
	}
	return nil
}

func (v *VEXXHOST) getRouterByName(args automation.ServerArgs) (*routers.Router, error) {
	var router *routers.Router
	pager := routers.List(v.networkClient, routers.ListOpts{
		Name: fmt.Sprintf("%s-router", args.MinecraftResource.GetName()),
	})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		routerList, err := routers.ExtractRouters(page)
		if err != nil {
			return false, err
		}
		for i, routerItem := range routerList {
			if routerItem.Name == fmt.Sprintf("%s-router", args.MinecraftResource.GetName()) {
				router = &routerList[i]
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return router, nil
}

func (v *VEXXHOST) getFloatingIPByInstanceID(id string) (*floatingips.FloatingIP, error) {
	var floatingIP *floatingips.FloatingIP
	pager := floatingips.List(v.computeClient)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		list, err := floatingips.ExtractFloatingIPs(page)
		if err != nil {
			return false, err
		}
		for i, item := range list {
			if item.InstanceID == id {
				floatingIP = &list[i]
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return floatingIP, nil
}

func (v *VEXXHOST) getSecurityGroupByName(args automation.ServerArgs) (*secgroups.SecurityGroup, error) {
	var securityGroup *secgroups.SecurityGroup
	pager := secgroups.List(v.computeClient)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		list, err := secgroups.ExtractSecurityGroups(page)
		if err != nil {
			return false, err
		}
		for i, item := range list {
			if item.Name == fmt.Sprintf("%s-sg", args.MinecraftResource.GetName()) {
				securityGroup = &list[i]
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return securityGroup, nil
}

func (v *VEXXHOST) getSubNetByName(args automation.ServerArgs) (*subnets.Subnet, error) {
	var subnet *subnets.Subnet
	pager := subnets.List(v.networkClient, subnets.ListOpts{
		Name: fmt.Sprintf("%s-subnet", args.MinecraftResource.GetName()),
	})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		list, err := subnets.ExtractSubnets(page)
		if err != nil {
			return false, err
		}
		for i, item := range list {
			if item.Name == fmt.Sprintf("%s-subnet", args.MinecraftResource.GetName()) {
				subnet = &list[i]
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return subnet, nil
}

func (v *VEXXHOST) getNetworkByName(args automation.ServerArgs) (*networks.Network, error) {
	var network *networks.Network
	pager := networks.List(v.networkClient, networks.ListOpts{
		Name: fmt.Sprintf("%s-net", args.MinecraftResource.GetName()),
	})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		list, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}
		for i, networkItem := range list {
			if networkItem.Name == fmt.Sprintf("%s-net", args.MinecraftResource.GetName()) {
				network = &list[i]
				break
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (v *VEXXHOST) ListServer() ([]automation.RessourceResults, error) {
	var result []automation.RessourceResults
	pager := servers.List(v.computeClient, servers.ListOpts{})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		list, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		for _, i := range list {
			for key := range i.Metadata {
				if key == common.InstanceTag {
					floatingIP, err := v.getFloatingIPByInstanceID(i.ID)
					if err != nil {
						return false, err
					}
					result = append(result, automation.RessourceResults{
						ID:       i.ID,
						Name:     i.Name,
						Region:   v.region,
						PublicIP: floatingIP.IP,
						Tags:     strings.Join(getTagKeys(i.Metadata), ","),
					})
				}
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VEXXHOST) UpdateServer(id string, args automation.ServerArgs) error {
	server, err := v.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), server.PublicIP, "ubuntu")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	zap.S().Infow("minectl server updated", "name", server.Name)
	return nil
}

func (v *VEXXHOST) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	server, err := v.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), server.PublicIP, "ubuntu")

	// as we are not allowed to login via root user, we need to add sudo to the command
	source := filepath.Join("/tmp", filepath.Base(plugin))
	err = remoteCommand.TransferFile(plugin, source)
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand(fmt.Sprintf("sudo mv %s %s\nsudo systemctl restart minecraft.service", source, destination))
	if err != nil {
		return err
	}
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "server", server.Name)
	return nil
}

func (v *VEXXHOST) GetServer(id string, _ automation.ServerArgs) (*automation.RessourceResults, error) {
	server, err := servers.Get(v.computeClient, id).Extract()
	if err != nil {
		return nil, err
	}
	floatingIP, err := v.getFloatingIPByInstanceID(server.ID)
	if err != nil {
		return nil, err
	}
	return &automation.RessourceResults{
		ID:       server.ID,
		Name:     server.Name,
		Region:   v.region,
		PublicIP: floatingIP.IP,
		Tags:     strings.Join(getTagKeys(server.Metadata), ","),
	}, nil
}
