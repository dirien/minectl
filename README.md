+ [Supported cloud provider ☁](#supported-cloud-provider-)
+ [TL;DR 🚀](#tldr-)
+ [Usage ⚙](#usage-)
    - [Access Token 🔑](#access-token-)
        * [Civo](#civo)
        * [Digital Ocean](#digital-ocean)
        * [Scaleway](#scaleway)
        * [Hetzner](#hetzner)
        * [Linode](#linode)
        * [OVHCloud](#ovhcloud)
        * [Equinix Metal](#equinix-metal)
        * [Google Compute Engine (GCE)](#google-compute-engine-gce)
    - [Minecraft Server Versions 📚](#minecraft-server-versions-)
    - [Server Config 📋](#server-config-)
    - [EULA ⚖️️](#eula-)
    - [Create Minecraft Server 🏗](#create-minecraft-server-)
    - [Delete Minecraft Server 🗑](#delete-minecraft-server-)
    - [List Minecraft Server 📒](#list-minecraft-server-)
    - [Update Minecraft Server 🆙](#update-minecraft-server-)
    - [Monitoring 📊](#monitoring-)
    - [Volumes 💽](#volumes-)
    - [Getting Started 🎫](#getting-started-)
+ [Known Limitation 😵](#known-limitation-)
+ [Contributing 🤝](#contributing-)
    - [Contributing via GitHub](#contributing-via-github)
    - [License](#license)
+ [Roadmap 🛣️](#roadmap-)
+ [Libraries & Tools 🔥](#libraries--tools-)
+ [Legal Disclaimer 👮](#legal-disclaimer-)

# minectl 🗺

![Minecraft](https://img.shields.io/badge/Minecraft-62B47A?style=for-the-badge&logo=Minecraft&logoColor=white)
![Go](https://img.shields.io/badge/go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Scaleway](https://img.shields.io/badge/scaleway-4F0599?style=for-the-badge&logo=scaleway&logoColor=white)
![DigitalOcean](https://img.shields.io/badge/DigitalOcean-0080FF?style=for-the-badge&logo=DigitalOcean&logoColor=white)
![Civo](https://img.shields.io/badge/Civo-239DFF?style=for-the-badge&logo=Civo&logoColor=white)
![Linode](https://img.shields.io/badge/linode-00A95C?style=for-the-badge&logo=linode&logoColor=white)
![Hetzner](https://img.shields.io/badge/hetzner-d50c2d?style=for-the-badge&logo=hetzner&logoColor=white)
![OVH](https://img.shields.io/badge/ovh-123F6D?style=for-the-badge&logo=ovh&logoColor=white)
![Equinix Metal](https://img.shields.io/badge/equinix--metal-d10810?style=for-the-badge&logo=equinix-metal&logoColor=white)
![Google Cloud](https://img.shields.io/badge/google--cloud-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white)

[![Build Binary](https://github.com/dirien/minectl/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/dirien/minectl/actions/workflows/ci.yaml)

`minectl`️️ is a cli for creating Minecraft (java or bedrock) server on different cloud provider.

It is a private side project of me, to learn more about Go, CLI and multi cloud.

### Supported cloud provider ☁

+ Civo (https://www.civo.com/)
+ Scaleway (https://www.scaleway.com)
+ DigitalOcean (https://www.digitalocean.com/)
+ Hetzner (https://www.hetzner.com/)
+ Linode (https://www.linode.com/)
+ OVHCloud (https://www.ovh.com/)
+ Equinix Metal (https://metal.equinix.com/)
+ Google Compute Engine (GCE) (https://cloud.google.com/compute)

### TL;DR 🚀

Install via homebrew:

```bash
brew tap dirien/homebrew-dirien
brew install minectl
```

Linux or Windows user, can directly download (or use `curl`/`wget`) the binary via
the [release page](https://github.com/dirien/minectl/releases).

### Usage ⚙

#### Access Token 🔑

`minectl` is completely build on zero-trust. It does not save any API Tokens, instead it looks them up in the ENV
variables.

##### Civo

```bash
export CIVO_TOKEN=xx
```

##### Digital Ocean

```bash
export DIGITALOCEAN_TOKEN=xxx
```

##### Scaleway

```bash
export ACCESS_KEY=xxx
export SECRET_KEY=yyy
export ORGANISATION_ID=zzz
```

##### Hetzner

```bash
export HCLOUD_TOKEN=yyyy
```

##### Linode

```bash
export LINODE_TOKEN=xxxx
```

#### OVHCloud

You need to create API keys per endpoint. For an overview of available endpoint
check [supported-apis](https://github.com/ovh/go-ovh#supported-apis) documentation

For, example, Europe visit https://eu.api.ovh.com/createToken to create your API keys for minectl 🗺

![img.png](docs/img/ovh_create_token.png)

For the proper `rights` choose all HTTP Verbs (GET,PUT,DELETE, POST), and we need only the `/cloud/` API.

```bash
export OVH_ENDPOINT=ovh-eu
export APPLICATION_KEY=xxx
export APPLICATION_SECRET=yyy
export CONSUMER_KEY=zzz
export SERVICENAME=<projectid>
```

#### Equinix Metal

```bash
export PACKET_AUTH_TOKEN=xxx
export EQUINIX_PROJECT=yyy
```

#### Google Compute Engine (GCE)

```bash
export GCE_KEY=<pathto>/key.json
```

See [Getting Started - GCE edition](docs/getting-started-gce.md) for details on how to create a GCP service account for
minectl 🗺

#### Minecraft Server Versions 📚

> ⚠️ `minectl 🗺` is not(!) providing any pre-compiled binaries of minecraft or use a page to download a pre-compiled version.
>
> Every non-vanilla version will be compiled during the creation phase of the server

Following Minecraft versions is `minectl 🗺` supporting.

##### Vanilla (Mincraft: Java Edition or Bedrock Edition)

The Vanilla software is the original, untouched, unmodified Minecraft server software created and distributed directly
by Mojang.[2]

##### CraftBukkit

CraftBukkit is lightly modified version of the Vanilla software allowing it to be able to run Bukkit plugins.

##### Spigot

Spigot is the most popular used Minecraft server software in the world. Spigot is a modified version of CraftBukkit with
hundreds of improvements and optimizations that can only make CraftBukkit shrink in shame.

##### PaperMC

Paper (formerly known as PaperSpigot, distributed via the Paperclip patch utility) is a high performance fork* of
Spigot.

##### Forge

Forge is well known for being able to use Forge Mods which are direct modifications to the Minecraft program code. In
doing so, Forge Mods can change the gaming-feel drastically as a result of this.

##### Fabric

Fabric is also an mod loader like Forge is with some improvements. Its lightweight and faster and it may is being the
best mod loader in the future because its doing very good.

Source: [[1]](#1-httpswwwspigotmcorgwikiwhat-is-spigot-craftbukkit-bukkit-vanilla-forg)

#### Server Config 📋

You need a MinecraftServer manifest file, to describe your VM and the Minecraft Server:

```yaml
apiVersion: ediri.io/minectl/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  monitoring:
    enabled: true|false
  server:
    cloud: "provider: civo|scaleway|do|hetzner|linode|ovh|equinix|gce"
    region: "region see cloud provider for details eg. fra1"
    size: "see cloud provider docs for details eg. g3.large"
    volumeSize: 100
    ssh: "<path to ssh public and private key>/ssh"
    port: "25565|19132 are the defaults for tcp/udp"
  minecraft:
    java:
      openjdk: "8|16 use jdk 8 for <1.17 java server version"
      xmx: 2G
      xms: 2G
      rcon:
        password: test
        port: 25575
        enabled: true
        broadcast: true
    edition: "java|bedrock|craftbukkit|fabric|forge|papermc|spigot"
    version: "<version>"
    eula: true
    properties: |
      level-seed=stackitminecraftrocks
      broadcast-rcon-to-ops=true
      ...
```

I created some example configs in the [config](config) folder for currently supported cloud provider and Minecraft
editions.

#### EULA ⚖

You need to set explicitly the EULA as new property in the MinecraftServer manifest to indicate your agreement with the
Minecraft End User License. See -> https://minecraft.net/terms for the details.

#### Create Minecraft Server 🏗

```bash
minectl create -h

Create an Minecraft Server.

Usage:
  minectl create [flags]

Examples:
mincetl create  \
    --filename server-do.yaml

Flags:
  -f, --filename string   Contains the configuration for minectl
  -h, --help              help for create
  -w, --wait              Wait for Minecraft Server is started (default true)  
```

#### Delete Minecraft Server 🗑

```bash
minectl delete -h

Delete an Minecraft Server.

Usage:
  minectl delete [flags]

Examples:
mincetl delete  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx
        

Flags:
  -f, --filename string   that contains the configuration for minectl
  -h, --help              help for delete
      --id string         contains the server id
```

#### List Minecraft Server 📒

```bash
minectl list -h

List all Minecraft Server.

Usage:
  minectl list [flags]

Examples:
mincetl list  \
    --provider civo \
    --region LON1

Flags:
  -h, --help              help for list
  -p, --provider string   The cloud provider - civo|scaleway|do|hetzner|linode|ovh|equinix|gce
  -r, --region string     The region for your cloud provider
```

#### Update Minecraft Server 🆙

```bash
minectl update -h
Update an Minecraft Server.

Usage:
  minectl update [flags]

Examples:
mincetl update  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx

Flags:
  -f, --filename string   Contains the configuration for minectl
  -h, --help              help for update
      --id string         contains the server id
```

#### Monitoring 📊

Monitoring is optional and disabled by default. It can be enabled with simply adding following fields to the
MinecraftServer manifest:

```yaml
...
apiVersion: ediri.io/minectl/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  monitoring:
    enabled: true|false
  server:
...
```

Every instance of minectl 🗺, has following monitoring components included:

- Prometheus (https://github.com/prometheus/prometheus)
- Node exporter (https://github.com/prometheus/node_exporter)

The `edition:java` has on top following exporter included:

- Minecraft exporter (https://github.com/dirien/minecraft-prometheus-exporter)

You can acces the `prometheus` via

```bash
http://<ip>:9090/graph
```

#### Volumes 💽

With the `volumeSize` tag, you are able to provision an extra volume during the creation phase of the server.

It is always recommended using the provided volume of the server, but in some cases (large mod packs, community server,
etc.) it make sense to provision a bigger volume separately.

When a seperate volume is defined, `minectl` is automatically installing Minecraft binaries on this volume.

```yaml
apiVersion: ediri.io/minectl/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  server:
    cloud: linode
    region: eu-central
    size: g6-standard-4
    volumeSize: 100
    ssh: "/Users/dirien/Tools/repos/stackit-minecraft/minecraft/ssh/minecraft"
    port: 25565
  minecraft:
...
```

#### Getting Started 🎫

- [Civo Java Edition](docs/getting-started-civo.md)
- [Civo Bedrock Edition](docs/getting-started-civo-bedrock.md)
- [Scaleway Java Edition](docs/getting-started-scaleway.md)
- [How to monitor your multi-cloud minectl 🗺 server?](docs/multi-server-monitoring-civo.md)

### Known Limitation 😵

`minectl` is still under development. There will be the possibility for breaking changes.

### Contributing 🤝

#### Contributing via GitHub

Feel free to join.

#### License

Apache License, Version 2.0

### Roadmap 🛣

- [x] Support Bedrock edition [#10](https://github.com/dirien/minectl/issues/10)
- [x] Add monitoring capabilities to minectl server [#21](https://github.com/dirien/minectl/issues/21)
- [x] List Minecraft Server [#11](https://github.com/dirien/minectl/issues/11)
- [x] New Command - Update Minecraft Server [#12](https://github.com/dirien/minectl/issues/12)
- [x] New cloud provider - Hetzner [#26](https://github.com/dirien/minectl/issues/26)
- [x] New cloud provider - Linode [#31](https://github.com/dirien/minectl/issues/31)
- [x] New cloud provider - OVHCloud [#43](https://github.com/dirien/minectl/issues/43)
- [x] New Cloud Provider Equinix Metal [#49](https://github.com/dirien/minectl/issues/49)
- [x] New cloud provider - GCE [#55](https://github.com/dirien/minectl/issues/55)
- [x] Add modded versions as new edition [#20](https://github.com/dirien/minectl/issues/20)
- [ ] ...

### Libraries & Tools 🔥

- https://github.com/fatih/color
- https://github.com/melbahja/goph
- https://github.com/spf13/cobra
- https://github.com/goreleaser
- https://github.com/briandowns/spinner
- https://github.com/civo/civogo
- https://github.com/digitalocean/godo
- https://github.com/scaleway/scaleway-sdk-go
- https://github.com/linode/linodego
- https://github.com/hetznercloud/hcloud-go
- https://github.com/olekukonko/tablewriter
- https://github.com/sethvargo/go-password
- https://github.com/ovh/go-ovh
- https://github.com/dirien/ovh-go-sdk
- https://github.com/packethost/packngo
- https://github.com/hashicorp/go-retryablehttp
- https://github.com/melbahja/goph
- https://github.com/googleapis/google-api-go-client

### Legal Disclaimer 👮

This project is not affiliated with Mojang Studios, XBox Game Studios, Double Eleven or the Minecraft brand.

"Minecraft" is a trademark of Mojang Synergies AB.

Other trademarks referenced herein are property of their respective owners.

Source:

##### [1] https://www.spigotmc.org/wiki/what-is-spigot-craftbukkit-bukkit-vanilla-forg/