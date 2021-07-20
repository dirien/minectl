package cloud

//CloudProvider
var cloudProvider = map[string]string{
	"do":       "DigitalOcean",
	"civo":     "Civo",
	"scaleway": "Scaleway",
	"hetzner":  "Hetzner",
	"linode":   "Linode",
	"ovh":      "OVHcloud",
}

func GetCloudProviderFullName(cloud string) string {
	return cloudProvider[cloud]
}
