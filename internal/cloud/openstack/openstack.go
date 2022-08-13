package openstack

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
	"github.com/minectl/internal/update"
	"go.uber.org/zap"
)

type OpenStack struct {
	tmpl          *minctlTemplate.Template
	computeClient *gophercloud.ServiceClient
	networkClient *gophercloud.ServiceClient
	region        string
	imageName     string
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

func NewOpenStack(imageName string) (*OpenStack, error) {
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}

	userID := os.Getenv("OS_USER_ID")
	domainID := os.Getenv("OS_PROJECT_DOMAIN_ID")
	if len(userID) != 0 {
		domainID = ""
		userID = os.Getenv("OS_USER_ID")
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainID:         domainID,
		UserID:           userID,
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
	return &OpenStack{
		tmpl:          tmpl,
		computeClient: computeClient,
		networkClient: networkClient,
		region:        os.Getenv("OS_REGION_NAME"),
		imageName:     imageName,
	}, nil
}

// CreateServer TODO: https://github.com/dirien/minectl/issues/299
func (o *OpenStack) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) { //nolint: gocyclo
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}

	keyPair, err := keypairs.Create(o.computeClient, keypairs.CreateOpts{
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
	pager := images.ListDetail(o.computeClient, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}
		for _, i := range imageList {
			if strings.Contains(i.Name, o.imageName) && !strings.HasSuffix(i.Name, "vGPU") {
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

	flavorPager := flavors2.ListDetail(o.computeClient, flavors2.ListOpts{})
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

	group, err := secgroups.Create(o.computeClient, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	err = o.createSecurityGroup(group, args.MinecraftResource.GetSSHPort(), "TCP")
	if err != nil {
		return nil, err
	}
	if args.MinecraftResource.GetEdition() == "bedrock" || args.MinecraftResource.GetEdition() == "nukkit" || args.MinecraftResource.GetEdition() == "powernukkit" {
		err = o.createSecurityGroup(group, args.MinecraftResource.GetPort(), "UDP")
		if err != nil {
			return nil, err
		}
	} else {
		err = o.createSecurityGroup(group, args.MinecraftResource.GetPort(), "TCP")
		if err != nil {
			return nil, err
		}
		if args.MinecraftResource.HasRCON() {
			err = o.createSecurityGroup(group, args.MinecraftResource.GetRCONPort(), "TCP")
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

	network, err := networks.Create(o.networkClient, networkOpts).Extract()
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

	subnet, err := subnets.Create(o.networkClient, subnetOpts).Extract()
	if err != nil {
		return nil, err
	}

	networkPager := networks.List(o.networkClient, networks.ListOpts{
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

	router, err := routers.Create(o.networkClient, routers.CreateOpts{
		Name:         fmt.Sprintf("%s-router", args.MinecraftResource.GetName()),
		AdminStateUp: &adminStateUp,
		GatewayInfo:  gatewayInfo,
	}).Extract()
	if err != nil {
		return nil, err
	}
	routers.AddInterface(o.networkClient, router.ID, &routers.AddInterfaceOpts{
		SubnetID: subnet.ID,
	})
	userData, err := o.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	server, err := servers.Create(o.computeClient, keypairs.CreateOptsExt{
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
		server, err = servers.Get(o.computeClient, server.ID).Extract()
		if err != nil {
			return nil, err
		}
		if server.Status == "ACTIVE" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	floatingIP, err := floatingips.Create(o.computeClient, floatingips.CreateOpts{
		Pool: "public",
	}).Extract()
	if err != nil {
		return nil, err
	}

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP.IP,
	}

	err = floatingips.AssociateInstance(o.computeClient, server.ID, associateOpts).ExtractErr()
	if err != nil {
		return nil, err
	}

	return &automation.ResourceResults{
		ID:       server.ID,
		Name:     server.Name,
		Region:   o.region,
		PublicIP: floatingIP.IP,
		Tags:     strings.Join(getTagKeys(server.Metadata), ","),
	}, err
}

func (o *OpenStack) createSecurityGroup(group *secgroups.SecurityGroup, port int, protocol string) error {
	ssh := secgroups.CreateRuleOpts{
		ParentGroupID: group.ID,
		FromPort:      port,
		ToPort:        port,
		IPProtocol:    protocol,
		CIDR:          "0.0.0.0/0",
	}

	_, err := secgroups.CreateRule(o.computeClient, ssh).Extract()
	if err != nil {
		return err
	}
	return nil
}

func (o *OpenStack) DeleteServer(id string, args automation.ServerArgs) error {
	server, err := servers.Get(o.computeClient, id).Extract()
	if err != nil {
		return err
	}
	err = keypairs.Delete(o.computeClient, fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()), &keypairs.DeleteOpts{}).Err
	if err != nil {
		return err
	}
	floatingIP, err := o.getFloatingIPByInstanceID(server.ID)
	if err != nil {
		return err
	}
	floatingips.DisassociateInstance(o.computeClient, id, floatingips.DisassociateOpts{
		FloatingIP: floatingIP.IP,
	})
	err = floatingips.Delete(o.computeClient, floatingIP.ID).Err
	if err != nil {
		return err
	}

	err = servers.Delete(o.computeClient, server.ID).Err
	if err != nil {
		return err
	}
	stillCreating := true
	for stillCreating {
		server, err = servers.Get(o.computeClient, server.ID).Extract()
		if err != nil {
			stillCreating = false
		}
		if server.Status == "DELETED" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	network, err := o.getNetworkByName(args)
	if err != nil {
		return err
	}
	subnet, err := o.getSubNetByName(args)
	if err != nil {
		return err
	}

	securityGroup, err := o.getSecurityGroupByName(args)
	if err != nil {
		return err
	}
	err = secgroups.Delete(o.computeClient, securityGroup.ID).Err
	if err != nil {
		return err
	}
	router, err := o.getRouterByName(args)
	if err != nil {
		return err
	}
	_, err = routers.RemoveInterface(o.networkClient, router.ID, routers.RemoveInterfaceOpts{
		SubnetID: subnet.ID,
	}).Extract()
	if err != nil {
		return err
	}
	err = routers.Delete(o.networkClient, router.ID).Err
	if err != nil {
		return err
	}
	err = subnets.Delete(o.networkClient, subnet.ID).Err
	if err != nil {
		return err
	}
	err = networks.Delete(o.networkClient, network.ID).Err
	if err != nil {
		return err
	}
	return nil
}

func (o *OpenStack) getRouterByName(args automation.ServerArgs) (*routers.Router, error) {
	var router *routers.Router
	pager := routers.List(o.networkClient, routers.ListOpts{
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

func (o *OpenStack) getFloatingIPByInstanceID(id string) (*floatingips.FloatingIP, error) {
	var floatingIP *floatingips.FloatingIP
	pager := floatingips.List(o.computeClient)
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

func (o *OpenStack) getSecurityGroupByName(args automation.ServerArgs) (*secgroups.SecurityGroup, error) {
	var securityGroup *secgroups.SecurityGroup
	pager := secgroups.List(o.computeClient)
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

func (o *OpenStack) getSubNetByName(args automation.ServerArgs) (*subnets.Subnet, error) {
	var subnet *subnets.Subnet
	pager := subnets.List(o.networkClient, subnets.ListOpts{
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

func (o *OpenStack) getNetworkByName(args automation.ServerArgs) (*networks.Network, error) {
	var network *networks.Network
	pager := networks.List(o.networkClient, networks.ListOpts{
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

func (o *OpenStack) ListServer() ([]automation.ResourceResults, error) {
	var result []automation.ResourceResults
	pager := servers.List(o.computeClient, servers.ListOpts{})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		list, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		for _, i := range list {
			for key := range i.Metadata {
				if key == common.InstanceTag {
					floatingIP, err := o.getFloatingIPByInstanceID(i.ID)
					if err != nil {
						return false, err
					}
					result = append(result, automation.ResourceResults{
						ID:       i.ID,
						Name:     i.Name,
						Region:   o.region,
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

func (o *OpenStack) UpdateServer(id string, args automation.ServerArgs) error {
	server, err := o.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), server.PublicIP, "ubuntu")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	zap.S().Infow("minectl server updated", "name", server.Name)
	return nil
}

func (o *OpenStack) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	server, err := o.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), server.PublicIP, "ubuntu")

	// as we are not allowed to login via root user, we need to add sudo to the command
	source := filepath.Join("/tmp", filepath.Base(plugin))
	err = remoteCommand.TransferFile(plugin, source, args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand(fmt.Sprintf("sudo mv %s %s\nsudo systemctl restart minecraft.service", source, destination), args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "server", server.Name)
	return nil
}

func (o *OpenStack) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	server, err := servers.Get(o.computeClient, id).Extract()
	if err != nil {
		return nil, err
	}
	floatingIP, err := o.getFloatingIPByInstanceID(server.ID)
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       server.ID,
		Name:     server.Name,
		Region:   o.region,
		PublicIP: floatingIP.IP,
		Tags:     strings.Join(getTagKeys(server.Metadata), ","),
	}, nil
}
