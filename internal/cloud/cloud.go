package cloud

// CloudProvider mapping cloud provider from short name to full name
var cloudProvider = map[string]string{
	"do":        "DigitalOcean",
	"civo":      "Civo",
	"scaleway":  "Scaleway",
	"hetzner":   "Hetzner",
	"linode":    "Linode",
	"ovh":       "OVHcloud",
	"equinix":   "Equinix Metal",
	"gce":       "Google Compute Engine",
	"vultr":     "vultr",
	"azure":     "Azure",
	"oci":       "Oracle Cloud Infrastructure",
	"ionos":     "IONOS Cloud",
	"aws":       "Amazon WebServices",
	"vexxhost":  "VEXXHOST",
	"exoscale":  "Exoscale",
	"multipass": "Ubuntu Multipass",
	"fuga":      "Fuga Cloud",
}

func GetCloudProviderFullName(cloud string) string {
	return cloudProvider[cloud]
}

func GetCloudProviderCode(fullName string) string {
	for code, name := range cloudProvider {
		if name == fullName {
			return code
		}
	}
	return ""
}
