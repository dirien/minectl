package alicloud

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
	"go.uber.org/zap"
)

type AliCloud struct {
	client *ecs.Client
	tmpl   *minctlTemplate.Template
}

func NewAliCloud(accessKey, accessKeySecret, regionId string) (*AliCloud, error) {
	client, err := ecs.NewClientWithAccessKey(regionId, accessKey, accessKeySecret)
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	return &AliCloud{
		client: client,
		tmpl:   tmpl,
	}, nil
}

func getTagKeys(tags []ecs.Tag) []string {
	var keys []string
	for _, tag := range tags {
		keys = append(keys, tag.Key)
	}
	return keys
}

func (a *AliCloud) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}

	keyPairRequest := ecs.CreateImportKeyPairRequest()
	keyPairRequest.PublicKeyBody = string(pubKeyFile)
	keyPairRequest.KeyPairName = fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName())
	sshPubKey, err := a.client.ImportKeyPair(keyPairRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("AliCloud SSH Key created", "KeyPairFingerPrint", sshPubKey.KeyPairFingerPrint)

	zonesRequest := ecs.CreateDescribeZonesRequest()
	zonesRequest.RegionId = args.MinecraftResource.GetRegion()
	zones, err := a.client.DescribeZones(zonesRequest)
	if err != nil {
		return nil, err
	}
	zone := zones.Zones.Zone[0]

	createVpcRequest := ecs.CreateCreateVpcRequest()
	createVpcRequest.VpcName = fmt.Sprintf("%s-vpc", args.MinecraftResource.GetName())
	createVpcRequest.CidrBlock = "172.16.0.0/12"
	vpc, err := a.client.CreateVpc(createVpcRequest)
	if err != nil {
		return nil, err
	}

	stillCreating := true
	for stillCreating {
		fmt.Println("loooping arround")
		describeVpcsRequest := ecs.CreateDescribeVpcsRequest()
		describeVpcsRequest.VpcId = vpc.VpcId
		vpcs, err := a.client.DescribeVpcs(describeVpcsRequest)
		if err != nil {
			return nil, err
		}
		vpc := vpcs.Vpcs.Vpc[0]
		if err != nil {
			return nil, err
		}
		if vpc.Status == "Available" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	createVSwitchRequest := ecs.CreateCreateVSwitchRequest()
	createVSwitchRequest.VpcId = vpc.VpcId
	createVSwitchRequest.ZoneId = zone.ZoneId
	createVSwitchRequest.CidrBlock = "172.22.208.0/20"

	vSwitch, err := a.client.CreateVSwitch(createVSwitchRequest)
	if err != nil {
		return nil, err
	}

	securityGroupRequest := ecs.CreateCreateSecurityGroupRequest()
	securityGroupRequest.SecurityGroupName = fmt.Sprintf("%s-sg", args.MinecraftResource.GetName())
	securityGroupRequest.VpcId = vpc.VpcId
	securityGroupRequest.SecurityGroupType = "normal"
	securityGroup, err := a.client.CreateSecurityGroup(securityGroupRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("AliCloud SecurityGroup created", "SecurityGroupId", securityGroup.SecurityGroupId)
	authorizeSecurityGroupRequest := ecs.CreateAuthorizeSecurityGroupRequest()
	authorizeSecurityGroupRequest.SecurityGroupId = securityGroup.SecurityGroupId
	authorizeSecurityGroupRequest.PortRange = "22/22"
	authorizeSecurityGroupRequest.IpProtocol = "tcp"
	authorizeSecurityGroupRequest.Policy = "accept"
	authorizeSecurityGroupRequest.SourceCidrIp = "0.0.0.0/0"
	authorizeSecurityGroupRequest.Description = "ssh"
	_, err = a.client.AuthorizeSecurityGroup(authorizeSecurityGroupRequest)
	if err != nil {
		return nil, err
	}
	authorizeSecurityGroupRequest.PortRange = fmt.Sprintf("%d/%d", args.MinecraftResource.GetRCONPort(), args.MinecraftResource.GetRCONPort())
	authorizeSecurityGroupRequest.Description = "RCON Port"
	_, err = a.client.AuthorizeSecurityGroup(authorizeSecurityGroupRequest)
	if err != nil {
		return nil, err
	}
	if args.MinecraftResource.HasMonitoring() {
		authorizeSecurityGroupRequest.IpProtocol = "tcp"
		authorizeSecurityGroupRequest.PortRange = "9090/9090"
		authorizeSecurityGroupRequest.Description = "Default Prometheus Port"
		_, err = a.client.AuthorizeSecurityGroup(authorizeSecurityGroupRequest)
		if err != nil {
			return nil, err
		}
	}
	if args.MinecraftResource.GetEdition() == "bedrock" || args.MinecraftResource.GetEdition() == "nukkit" {
		authorizeSecurityGroupRequest.IpProtocol = "udp"
	} else {
		authorizeSecurityGroupRequest.IpProtocol = "tcp"
	}
	authorizeSecurityGroupRequest.PortRange = fmt.Sprintf("%d/%d", args.MinecraftResource.GetPort(), args.MinecraftResource.GetPort())
	authorizeSecurityGroupRequest.Description = "Minecraft Server Port"
	_, err = a.client.AuthorizeSecurityGroup(authorizeSecurityGroupRequest)
	if err != nil {
		return nil, err
	}

	describeImagesRequest := ecs.CreateDescribeImagesRequest()
	describeImages, err := a.client.DescribeImages(describeImagesRequest)
	if err != nil {
		return nil, err
	}
	var image ecs.Image
	for _, img := range describeImages.Images.Image {
		if strings.Contains(img.ImageId, "20_04_x64") {
			image = img
			break
		}
	}

	fmt.Println(image)
	script, err := a.tmpl.GetTemplate(args.MinecraftResource, "", minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer()))
	if err != nil {
		return nil, err
	}

	createInstanceRequest := ecs.CreateRunInstancesRequest()
	createInstanceRequest.ImageId = image.ImageId
	createInstanceRequest.InstanceChargeType = "PostPaid"
	createInstanceRequest.HostName = args.MinecraftResource.GetName()
	createInstanceRequest.InstanceType = args.MinecraftResource.GetSize()
	createInstanceRequest.InstanceName = args.MinecraftResource.GetName()
	createInstanceRequest.KeyPairName = sshPubKey.KeyPairName
	createInstanceRequest.UserData = base64.StdEncoding.EncodeToString([]byte(script))
	createInstanceRequest.VSwitchId = vSwitch.VSwitchId
	createInstanceRequest.DryRun = requests.NewBoolean(true)
	createInstanceRequest.InternetMaxBandwidthOut = requests.NewInteger(10)
	createInstanceRequest.Tag = &[]ecs.RunInstancesTag{
		{
			Key:   args.MinecraftResource.GetRegion(),
			Value: "true",
		},
		{
			Key:   common.InstanceTag,
			Value: "true",
		},
	}
	createInstanceRequest.SecurityGroupId = securityGroup.SecurityGroupId

	runInstances, err := a.client.RunInstances(createInstanceRequest)
	if err != nil {
		return nil, err
	}
	stillCreating = true
	var instance ecs.Instance
	for stillCreating {
		fmt.Println("loooping arround")

		describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
		describeInstancesRequest.InstanceIds = fmt.Sprintf("[\"%s\"]", runInstances.InstanceIdSets.InstanceIdSet[0])
		describeInstances, err := a.client.DescribeInstances(describeInstancesRequest)
		if err != nil {
			return nil, err
		}
		instance = describeInstances.Instances.Instance[0]
		if instance.Status == "Running" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	return &automation.RessourceResults{
		ID:       instance.InstanceId,
		Name:     instance.InstanceName,
		Region:   instance.RegionId,
		PublicIP: instance.PublicIpAddress.IpAddress[0],
		Tags:     strings.Join(getTagKeys(instance.Tags.Tag), ","),
	}, err

}

func (a *AliCloud) DeleteServer(id string, args automation.ServerArgs) error {
	deleteKeyPairsRequest := ecs.CreateDeleteKeyPairsRequest()
	deleteKeyPairsRequest.KeyPairNames = fmt.Sprintf("[\"%s-ssh\"]", args.MinecraftResource.GetName())
	_, err := a.client.DeleteKeyPairs(deleteKeyPairsRequest)
	if err != nil {
		return err
	}

	x := ecs.CreateDescribeSecurityGroupsRequest()
	x.SecurityGroupName = fmt.Sprintf("%s-sg", args.MinecraftResource.GetName())
	securityGroups, err := a.client.DescribeSecurityGroups(x)
	if err != nil {
		return err
	}
	for _, securityGroup := range securityGroups.SecurityGroups.SecurityGroup {
		y := ecs.CreateDeleteSecurityGroupRequest()
		y.SecurityGroupId = securityGroup.SecurityGroupId
		_, err = a.client.DeleteSecurityGroup(y)
		if err != nil {
			return err
		}

		z := ecs.CreateDescribeVSwitchesRequest()
		z.VpcId = securityGroup.VpcId
		vSwitches, err := a.client.DescribeVSwitches(z)
		if err != nil {
			return err
		}
		for _, vSwitch := range vSwitches.VSwitches.VSwitch {
			m := ecs.CreateDeleteVSwitchRequest()
			m.VSwitchId = vSwitch.VSwitchId
			_, err := a.client.DeleteVSwitch(m)
			if err != nil {
				return err
			}
		}
	}
	f := ecs.CreateDescribeVpcsRequest()
	vpcs, err := a.client.DescribeVpcs(f)
	if err != nil {
		return err
	}

	for _, vpc := range vpcs.Vpcs.Vpc {
		if vpc.VpcName == fmt.Sprintf("%s-vpc", args.MinecraftResource.GetName()) {
			time.Sleep(20 * time.Second)
			ri := ecs.CreateDeleteVpcRequest()
			ri.VpcId = vpc.VpcId
			_, err := a.client.DeleteVpc(ri)
			if err != nil {
				return err
			}
		}
	}

	panic("implement me")
}

func (a *AliCloud) ListServer() ([]automation.RessourceResults, error) {
	panic("implement me")
}

func (a *AliCloud) UpdateServer(id string, args automation.ServerArgs) error {
	panic("implement me")
}

func (a *AliCloud) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	panic("implement me")
}

func (a *AliCloud) GetServer(id string, args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
}
