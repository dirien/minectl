package cloud

//CloudProvider
var cloudProvider = map[string]string{
	"do":       "DigitalOcean",
	"civo":     "Civo",
	"scaleway": "Scaleway",
}

func GetCloudProviderFullName(cloud string) string {
	return cloudProvider[cloud]
}
