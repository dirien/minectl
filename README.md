## New - 1.21.11 support

`minectl 🗺`️️ supports the `Minecraft 1.21.11` version

<img alt="logo" src="docs/img/1_21_11_logo.png" width="40%"/>

# `minectl 🗺`

![minectl](https://dirien.github.io/minectl/img/minectl.png)

![Minecraft](https://img.shields.io/badge/Minecraft-62B47A?style=for-the-badge&logo=Minecraft&logoColor=white)
![Go](https://img.shields.io/badge/go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Ubuntu](https://img.shields.io/badge/ubuntu_22.04-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Scaleway](https://img.shields.io/badge/scaleway-4F0599?style=for-the-badge&logo=scaleway&logoColor=white)
![DigitalOcean](https://img.shields.io/badge/DigitalOcean-0080FF?style=for-the-badge&logo=DigitalOcean&logoColor=white)
![Civo](https://img.shields.io/badge/Civo-239DFF?style=for-the-badge&logo=Civo&logoColor=white)
![Akamai Connected Cloud](https://img.shields.io/badge/akamai_connected_cloud-0096D6?style=for-the-badge&logo=akamai&logoColor=white)
![Hetzner](https://img.shields.io/badge/hetzner-d50c2d?style=for-the-badge&logo=hetzner&logoColor=white)
![OVH](https://img.shields.io/badge/ovh-123F6D?style=for-the-badge&logo=ovh&logoColor=white)
![Google Cloud](https://img.shields.io/badge/google_cloud-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white)
![Vultr](https://img.shields.io/badge/vultr-007BFC?style=for-the-badge&logo=vultr&logoColor=white)
![Microsoft Azure](https://img.shields.io/badge/Microsoft_Azure-0078D4?style=for-the-badge&logo=microsoft-azure&logoColor=white)
![Oracle Cloud Infrastructure](https://img.shields.io/badge/Oracle_Cloud_Infrastructure-F80000?style=for-the-badge&logo=oracle&logoColor=white)
![Amazon AWS](https://img.shields.io/badge/Amazon_AWS-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white)
![VEXXHOST](https://img.shields.io/badge/VEXXHOST-2A1659?style=for-the-badge&logo=vexxhost&logoColor=white)
![Multipass](https://img.shields.io/badge/Multipass-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Exoscale](https://img.shields.io/badge/Exoscale-DA291C?style=for-the-badge&logo=exoscale&logoColor=white)
![Fuga Cloud](https://img.shields.io/badge/fuga_cloud-242F4B?style=for-the-badge&logo=fugacloud&logoColor=white)

[![Go Reference](https://pkg.go.dev/badge/github.com/dirien/minectl.svg)](https://pkg.go.dev/github.com/dirien/minectl)
[![Go Report Card](https://goreportcard.com/badge/github.com/dirien/minectl)](https://goreportcard.com/report/github.com/dirien/minectl)

![GitHub Workflow Status (main)](https://img.shields.io/github/actions/workflow/status/dirien/minectl/ci.yaml?branch=main&logo=github&style=for-the-badge)
![GitHub](https://img.shields.io/github/license/dirien/minectl?style=for-the-badge)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=dirien_minectl&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=dirien_minectl)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/dirien/minectl/badge?style=for-the-badge)](https://api.securityscorecards.dev/projects/github.com/dirien/minectl)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5238/badge)](https://bestpractices.coreinfrastructure.org/projects/5238)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/dirien/minectl?style=for-the-badge)

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/minectl&style=for-the-badge)](https://artifacthub.io/packages/search?repo=minectl)

`minectl 🗺`️️ is a CLI for creating Minecraft servers on different cloud providers.

It is a private side project of me, to learn more about Go, CLI and multi-cloud environments.

## Table of Contents

- [Supported Cloud Providers](#supported-cloud-providers)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Getting Started Guides](#getting-started-guides)
- [Known Limitations](#known-limitations)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [Libraries & Tools](#libraries--tools)
- [Legal Disclaimer](#legal-disclaimer)

## Supported Cloud Providers

| Provider | Website |
|----------|---------|
| Civo | https://www.civo.com/ |
| Scaleway | https://www.scaleway.com |
| DigitalOcean | https://www.digitalocean.com/ |
| Hetzner | https://www.hetzner.com/ |
| Akamai Connected Cloud | https://www.linode.com/ |
| OVHCloud | https://www.ovh.com/ |
| Google Compute Engine (GCE) | https://cloud.google.com/compute |
| Azure | https://azure.microsoft.com/en-us/ |
| Oracle Cloud Infrastructure | https://www.oracle.com/cloud/ |
| Amazon AWS | https://aws.amazon.com/ |
| VEXXHOST | https://vexxhost.com/ |
| Multipass | https://multipass.run/ |
| Exoscale | https://www.exoscale.com/ |
| Fuga Cloud | https://fuga.cloud/ |

## Quick Start

### Installation

#### Linux/macOS (Installation Script)

```bash
curl -sLS https://get.minectl.dev | sudo sh
```

or without `sudo`:

```bash
curl -sLS https://get.minectl.dev | sh
```

#### macOS (Homebrew)

```bash
brew tap dirien/homebrew-dirien
brew install minectl
```

#### Windows (PowerShell)

```powershell
# Create directory
New-Item -Path "$HOME/minectl/cli" -Type Directory
# Download file
Start-BitsTransfer -Source https://github.com/dirien/minectl/releases/download/v0.7.0/minectl_0.7.0_windows_amd64.zip -Destination "$HOME/minectl/cli/."
# Uncompress zip file
Expand-Archive $HOME/minectl/cli/*.zip -DestinationPath $HOME/minectl/cli/.
# Add to Windows Environment Variables
[Environment]::SetEnvironmentVariable("Path",$($env:Path + ";$Home\minectl\cli"),'User')
```

#### From Source

```bash
git clone https://github.com/dirien/minectl
cd minectl
make build
```

### Create Your First Server

1. **Set up authentication** for your cloud provider ([see docs](docs/authentication.md))

2. **Create a config file** using the wizard:
   ```bash
   minectl wizard
   ```

3. **Create the server**:
   ```bash
   minectl create --filename server.yaml
   ```

4. **Connect and play!**

[![asciicast](https://asciinema.org/a/439572.svg)](https://asciinema.org/a/439572)

## Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | High-level architectural overview |
| [Authentication](docs/authentication.md) | Cloud provider credentials setup |
| [Configuration](docs/configuration.md) | Server and proxy configuration reference |
| [CLI Reference](docs/cli-reference.md) | All commands and flags |
| [Editions](docs/editions.md) | Supported Minecraft server and proxy editions |
| [Features](docs/features.md) | Monitoring, volumes, security, and more |

## Getting Started Guides

- [Civo Java Edition](docs/getting-started-civo.md)
- [Civo Bedrock Edition](docs/getting-started-civo-bedrock.md)
- [Scaleway Java Edition](docs/getting-started-scaleway.md)
- [GCE Edition](docs/getting-started-gce.md)
- [Exoscale Edition](docs/getting-started-exoscale.md)
- [Multi-cloud server monitoring](docs/multi-server-monitoring-civo.md)
- [LuckyBlocks on Scaleway](docs/running-minecraft-luckyblocks-budget-scaleway.md)
- [PaperMC SkyBlock on Hetzner](docs/skyblock-papermc-hetzner.md)

## Known Limitations

`minectl 🗺` is still under development. There will be the possibility for breaking changes.

## Contributing

Feel free to join! See our [contribution guidelines](CONTRIBUTING.md).

**License:** Apache License, Version 2.0

## Roadmap

- [x] Support Bedrock edition [#10](https://github.com/dirien/minectl/issues/10)
- [x] Add monitoring capabilities [#21](https://github.com/dirien/minectl/issues/21)
- [x] List Minecraft Server [#11](https://github.com/dirien/minectl/issues/11)
- [x] Update Minecraft Server [#12](https://github.com/dirien/minectl/issues/12)
- [x] Hetzner support [#26](https://github.com/dirien/minectl/issues/26)
- [x] Linode support [#31](https://github.com/dirien/minectl/issues/31)
- [x] OVHCloud support [#43](https://github.com/dirien/minectl/issues/43)
- [x] GCE support [#55](https://github.com/dirien/minectl/issues/55)
- [x] Modded editions [#20](https://github.com/dirien/minectl/issues/20)
- [x] Vultr support [#90](https://github.com/dirien/minectl/issues/90)
- [x] Proxy servers (BungeeCord/Waterfall) [#95](https://github.com/dirien/minectl/issues/95)
- [x] Azure support [#56](https://github.com/dirien/minectl/issues/56)
- [x] Oracle/OCI support [#107](https://github.com/dirien/minectl/issues/107)
- [x] AWS support [#210](https://github.com/dirien/minectl/pull/210)
- [ ] Much more to come...

## Libraries & Tools

<details>
<summary>Click to expand</summary>

- https://github.com/fatih/color
- https://github.com/melbahja/goph
- https://github.com/spf13/cobra
- https://github.com/goreleaser
- https://github.com/civo/civogo
- https://github.com/digitalocean/godo
- https://github.com/scaleway/scaleway-sdk-go
- https://github.com/linode/linodego
- https://github.com/hetznercloud/hcloud-go
- https://github.com/olekukonko/tablewriter
- https://github.com/sethvargo/go-password
- https://github.com/ovh/go-ovh
- https://github.com/dirien/ovh-go-sdk
- https://github.com/hashicorp/go-retryablehttp
- https://github.com/googleapis/google-api-go-client
- https://github.com/Masterminds/sprig
- https://github.com/Tnze/go-mc
- https://github.com/c-bata/go-prompt
- https://github.com/vultr/govultr
- https://github.com/Azure/azure-sdk-for-go
- https://github.com/blang/semver
- https://github.com/tcnksm/go-latest
- https://github.com/uber-go/zap
- https://github.com/oracle/oci-go-sdk
- https://github.com/AlecAivazis/survey
- https://github.com/aws/aws-sdk-go
- https://github.com/gophercloud/gophercloud
- https://github.com/exoscale/egoscale

</details>

## Legal Disclaimer

This project is not affiliated with Mojang Studios, XBox Game Studios, Double Eleven or the Minecraft brand.

"Minecraft" is a trademark of Mojang Synergies AB.

Other trademarks referenced herein are property of their respective owners.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/dirien/minectl.svg)](https://starchart.cc/dirien/minectl)
