# Cloud Provider Reference

## Supported Providers and Environment Variables

| Provider | `cloud` value | Required Environment Variables |
|----------|--------------|-------------------------------|
| Civo | `civo` | `CIVO_TOKEN` |
| Scaleway | `scaleway` | `ACCESS_KEY`, `SECRET_KEY`, `ORGANISATION_ID` |
| DigitalOcean | `do` | `DIGITALOCEAN_TOKEN` |
| Hetzner | `hetzner` | `HCLOUD_TOKEN` |
| Akamai (Linode) | `akamai` | `LINODE_TOKEN` |
| OVHcloud | `ovh` | `OVH_ENDPOINT`, `APPLICATION_KEY`, `APPLICATION_SECRET`, `CONSUMER_KEY`, `SERVICENAME` |
| Google Compute Engine | `gce` | `GCE_KEY` |
| Vultr | `vultr` | `VULTR_API_KEY` |
| Azure | `azure` | `AZURE_AUTH_LOCATION` |
| Oracle Cloud (OCI) | `oci` | OCI SDK config (default) |
| IONOS Cloud | `ionos` | `IONOS_USERNAME`, `IONOS_PASSWORD`, `IONOS_TOKEN` |
| AWS | `aws` | Standard AWS SDK credentials |
| VEXXHOST | `vexxhost` | OpenStack env vars |
| Exoscale | `exoscale` | `EXOSCALE_API_KEY`, `EXOSCALE_API_SECRET` |
| Fuga Cloud | `fuga` | OpenStack env vars |
| Ubuntu Multipass | `multipass` | None (local) |

## Region Examples

| Provider | Example Regions |
|----------|----------------|
| Civo | `LON1`, `NYC1`, `FRA1` |
| DigitalOcean | `fra1`, `nyc1`, `sfo1` |
| Hetzner | `fsn1`, `nbg1`, `hel1` |
| AWS | `eu-central-1`, `us-east-1`, `ap-southeast-1` |
| GCE | `europe-west6-a`, `us-central1-a` |
| Azure | `westeurope`, `eastus` |
| Scaleway | `fr-par-1`, `nl-ams-1` |

## Size Examples

| Provider | Example Sizes |
|----------|--------------|
| Civo | `g3.large`, `g3.xlarge` |
| DigitalOcean | `s-4vcpu-8gb` |
| Hetzner | `cx21`, `cpx31` |
| AWS | `t3.xlarge`, `m5.large` |
| GCE | `e2-standard-2`, `e2-standard-4` |
| Azure | `Standard_B2s` |

## Features by Provider

- **Spot instances**: AWS (`spot: true`), Azure (`spot: true`), GCE (`spot: true`)
- **ARM support**: Hetzner (`arm: true`), AWS (`arm: true`), GCE (`arm: true`)
- **Volumes**: All providers support `volumeSize` (in GB)
