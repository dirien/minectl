+ [Supported cloud provider ☁](#supported-cloud-provider-)
+ [TL;DR 🚀](#tldr-)
    - [Installing `minectl 🗺`](#installing-minectl-)
        * [Installation Script](#installation-script)
        * [Mac OS X](#mac-os-x)
        * [Windows](#windows)
        * [Source install](#source-install)
        * [GoFish (deprecated 🕸️)](#gofish-deprecated-)
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
        * [Vultr](#vultr)
        * [Azure](#azure)
        * [Oracle Cloud Infrastructure](#oracle-cloud-infrastructure)
        * [Ionos Cloud](#ionos-cloud-limited-support)
        * [Amazon AWS](#amazon-aws)
        * [VEXXHOST](#vexxhost)
        * [Multipass](#multipass)
        * [Exoscale](#exoscale)
        * [Fuga Cloud](#fuga-cloud)
    - [Minecraft Server Versions 📚](#minecraft-server-versions-)
    - [Minecraft Proxy Versions 📚](#minecraft-proxy-versions-)
    - [Server Configs 📋](#server-configs-)
        * [Spot Instances](#spot-instances)
        * [MinecraftProxy Config 📡](#minecraftproxy-config-)
        * [MinecraftServer Config 🕹](#mincraftserver-config-)
    - [EULA ⚖️️](#eula-)
    - [`minectl 🗺` Configuration File Wizard 🧙](#minectl--configuration-file-wizard-)
    - [Create Minecraft Server 🏗](#create-minecraft-server-)
    - [Delete Minecraft Server 🗑](#delete-minecraft-server-)
    - [List Minecraft Server 📒](#list-minecraft-server-)
    - [Update Minecraft Server 🆙](#update-minecraft-server-)
    - [Plugins Minecraft Server ⤴️](#plugins-minecraft-server-)
    - [RCON Minecraft Server 🔌](#rcon-minecraft-server-)
    - [Monitoring 📊](#monitoring-)
    - [Volumes 💽](#volumes-)
    - [Headless Mode 👻](#headless-mode-)
    - [Security 🔒](#security-)
    - [Getting Started 🎫](#getting-started-)
+ [Known Limitation 😵](#known-limitation-)
+ [Contributing 🤝](#contributing-)
    - [Contributing via GitHub](#contributing-via-github)
    - [License](#license)
+ [Roadmap 🛣️](#roadmap-)
+ [Libraries & Tools 🔥](#libraries--tools-)
+ [Legal Disclaimer 👮](#legal-disclaimer-)

# `minectl 🗺`

![minectl](https://dirien.github.io/minectl/img/minectl.png)

![Minecraft](https://img.shields.io/badge/Minecraft-62B47A?style=for-the-badge&logo=Minecraft&logoColor=white)
![Go](https://img.shields.io/badge/go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Ubuntu](https://img.shields.io/badge/ubuntu_22.04-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Scaleway](https://img.shields.io/badge/scaleway-4F0599?style=for-the-badge&logo=scaleway&logoColor=white)
![DigitalOcean](https://img.shields.io/badge/DigitalOcean-0080FF?style=for-the-badge&logo=DigitalOcean&logoColor=white)
![Civo](https://img.shields.io/badge/Civo-239DFF?style=for-the-badge&logo=Civo&logoColor=white)
![Linode](https://img.shields.io/badge/linode-00A95C?style=for-the-badge&logo=linode&logoColor=white)
![Hetzner](https://img.shields.io/badge/hetzner-d50c2d?style=for-the-badge&logo=hetzner&logoColor=white)
![OVH](https://img.shields.io/badge/ovh-123F6D?style=for-the-badge&logo=ovh&logoColor=white)
![Equinix Metal](https://img.shields.io/badge/equinix_metal-d10810?style=for-the-badge&logo=equinixmetal&logoColor=white)
![Google Cloud](https://img.shields.io/badge/google_cloud-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white)
![Vultr](https://img.shields.io/badge/vultr-007BFC?style=for-the-badge&logo=vultr&logoColor=white)
![Microsoft Azure](https://img.shields.io/badge/Microsoft_Azure-0078D4?style=for-the-badge&logo=microsoft-azure&logoColor=white)
![Oracle Cloud Infrastructure](https://img.shields.io/badge/Oracle_Cloud_Infrastructure-F80000?style=for-the-badge&logo=oracle&logoColor=white)
![Ionos Cloud](https://img.shields.io/badge/ionos--cloud-003D8F?style=for-the-badge&logo=ionos&logoColor=white)
![Amazon AWS](https://img.shields.io/badge/Amazon_AWS-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white)
![VEXXHOST](https://img.shields.io/badge/VEXXHOST-2A1659?style=for-the-badge&logo=vexxhost&logoColor=white)
![Multipass](https://img.shields.io/badge/Multipass-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Exoscale](https://img.shields.io/badge/Exoscale-DA291C?style=for-the-badge&logo=exoscale&logoColor=white)
![Fuga Cloud](https://img.shields.io/badge/fuga_cloud-242F4B?style=for-the-badge&logo=fugacloud&logoColor=white)


![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/dirien/minectl/Build%20Binary/main?logo=github&style=for-the-badge)
![GitHub](https://img.shields.io/github/license/dirien/minectl?style=for-the-badge)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/dirien/minectl?style=for-the-badge)

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/minectl&style=for-the-badge)](https://artifacthub.io/packages/search?repo=minectl)

`minectl 🗺`️️ is a cli for creating Minecraft server on different cloud provider.

It is a private side project of me, to learn more about Go, CLI and multi-cloud environments.

### Supported cloud provider ☁

+ Civo (https://www.civo.com/)
+ Scaleway (https://www.scaleway.com)
+ DigitalOcean (https://www.digitalocean.com/)
+ Hetzner (https://www.hetzner.com/)
+ Linode (https://www.linode.com/)
+ OVHCloud (https://www.ovh.com/)
+ Equinix Metal (https://metal.equinix.com/)
+ Google Compute Engine (GCE) (https://cloud.google.com/compute)
+ Azure (https://azure.microsoft.com/en-us/)
+ Oracle Cloud Infrastructure (https://www.oracle.com/cloud/)
+ Ionos Cloud (https://cloud.ionos.de/)
+ Amazon AWS (https://aws.amazon.com/)
+ VEXXHOST (https://vexxhost.com/)
+ Multipass (https://multipass.run/)
+ Exoscale (https://www.exoscale.com/)
+ Fuga Cloud (https://fuga.cloud/)

### TL;DR 🚀

#### Installing `minectl 🗺`

Download the latest binary executable for your operating system.

##### Installation Script

```bash
curl -sLS https://get.minectl.dev | sudo sh
```

or without `sudo`

```bash
curl -sLS https://get.minectl.dev | sh
```

This will install the `minectl 🗺` to `~/.minctl/` and add it to your path. When it can’t automatically add `minectl 🗺`
to your path, you will be prompted to add it manually.

##### Mac OS X

- Use [Homebrew](https://brew.sh/)
  ```bash
  brew tap dirien/homebrew-dirien
  brew install minectl
  ```

##### Windows

- Use Powershell

  ```powershell
  #Create directory
  New-Item -Path "$HOME/minectl/cli" -Type Directory
  # Download file
  Start-BitsTransfer -Source https://github.com/dirien/minectl/releases/download/v0.7.0/minectl_0.7.0_windows_amd64.zip -Destination "$HOME/minectl/cli/."
  # Uncompress zip file
  Expand-Archive $HOME/minectl/cli/*.zip -DestinationPath C:\Users\Developer\minectl\cli\.
  #Add to Windows `Environment Variables`
  [Environment]::SetEnvironmentVariable("Path",$($env:Path + ";$Home\minectl\cli"),'User')
  ```

##### Source install

You need to have [go](https://golang.org/) installed, and need to checkout
the [Git repository](https://github.com/dirien/minectl) and run the following commands:

```bash
make build
 ```

This will output the `minectl 🗺` binary in the `bin/minectl` folder.

##### GoFish (deprecated 🕸️)

GoFish works across all three major operating systems (Windows, MacOS, and Linux). It installs packages into its own
directory and symlinks their files into /usr/local (or C:\ProgramData for Windows). You can think of it as the
cross-platform Homebrew.

To install `minectl 🗺` just type

```
gofish install minectl
```

As `minectl 🗺` is already in the main [rig](https://github.com/fishworks/fish-food) of Gofish.

### Architectural overview

You can find a high level architectural overview [here](docs/architecture.md)

### Usage ⚙

#### Access Token 🔑

`minectl 🗺` is completely build on zero-trust. It does not save any API Tokens, instead it looks them up in the ENV
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

For, example, Europe visit https://eu.api.ovh.com/createToken to create your API keys for `minectl 🗺`

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
export METAL_AUTH_TOKEN=xxx
export EQUINIX_PROJECT=yyy
```

#### Google Compute Engine (GCE)

```bash
export GCE_KEY=<pathto>/key.json
```

See [Getting Started - GCE edition](docs/getting-started-gce.md) for details on how to create a GCP service account for
minectl 🗺`

#### Vultr

```bash
export VULTR_API_KEY=xxx
```

#### Azure

```bash
az login
az ad sp create-for-rbac --sdk-auth --role 'Contributor' > azure.auth

export AZURE_AUTH_LOCATION=azure.auth
```

#### Oracle Cloud Infrastructure

To keep things simple for the moment, the authentication uses OCI config file. And there the default.

Example:

```bash
cat  /Users/user/.oci/config

[DEFAULT]
user=<ocid>
fingerprint=<SSH fingerprint>
key_file=>path to PEM file>
tenancy=<ocid>
region=<region>
```

Please follow the instructions under -> https://docs.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm

#### Ionos Cloud (limited support)

> I can offer only limited support for Ionos Cloud, as I don't have access to the API anymore. Ionos Cloud is a B2B only
> cloud service.

```bash
export IONOS_USERNAME=xxx
export IONOS_PASSWORD=yyy
export IONOS_TOKEN=<optional>
```

#### Amazon AWS

`minectl 🗺` looks for credentials in the following order:

- Environment variables.
- Shared credentials file.

##### Credentials file

The credentials file is most often located in the `~/.aws/credentials` and contains following content:

```bash
cat ~/.aws/credentials
[default]
aws_access_key_id = xxxx
aws_secret_access_key = zzzz
```

##### Environment variables can be set in the following way:

```bash
export AWS_ACCESS_KEY_ID=<aws_access_key_id>
export AWS_SECRET_ACCESS_KEY=<aws_secret_access_key>
export AWS_REGION=<aws_region>
```

#### VEXXHOST

It is recommended to store OpenStack credentials as environment variables because it decouples credential information
from source code:

So download the `OpenStack RC File` from the Horizon UI by click on the "Download OpenStack RC File" button at the top
right-hand corner.

To execute the file, run source `xxxx-openrc.sh` and you will be prompted for your password.

Thats all.

#### Multipass

> ⚠️ Set the plan to cpu-memG. For example: 1-2G

Multipass is a mini-cloud on your workstation using native hypervisors of all the supported platforms (Windows, macOS
and Linux), it will give you an Ubuntu command line in just a click (“Open shell”) or a simple multipass shell command,
or even a keyboard shortcut. Find what images are available with multipass find and create new instances with multipass
launch.

To install multipass, just follow the instructions on [multipass.run](https://multipass.run/) for your platform.

#### Exoscale

Go to the IAM section in the Exoscale Console and create a new API key. You can restricte the key to just perform
operations on the `compute` service.

```bash
export EXOSCALE_API_KEY=<key>
export EXOSCALE_API_SECRET=<secret>
```

#### Fuga Cloud

To get the `OpenStack RC File` from the Fuga Cloud UI, follow this steps:

1. Log in to the Fuga Cloud Dashboard
2. Go to Account → Access → Credentials
3. You can choose a user credential or team credential.
4. If you haven’t already, you should create one of these credentials. Hold on to the password.
5. Click on download OpenRC. This file contains all necessary configurations for the client.

```bash
source  fuga-openrc.sh
```

Enter the password which matches the username of the contents of the OpenRC file.



#### Minecraft Server Versions 📚

> ⚠️ `minectl 🗺` is not(!) providing any pre-compiled binaries of Minecraft or download a pre-compiled version.
>
> Every _non-vanilla_ version will be compiled during the build phase of your server.

Following Minecraft versions is `minectl 🗺` supporting.

##### Vanilla (Mincraft: Java Edition or Bedrock Edition)

The Vanilla software is the original, untouched, unmodified Minecraft server software created and distributed directly
by Mojang.

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

#### Minecraft Proxy Versions 📚

Network proxy server is what you set up and use as the controller of a network of server - this is the server that
connects all of your playable servers together so people can log in through one server IP, and then teleport between the
separate servers ingame rather than having to log out and into each different one.

A server network typically consist of the following servers:

1. The proxy server itself running the desired software (BungeeCord being the most widely used and supported). This is
   the server that you would be advertising the IP of, as all players should be logging in through the proxy server at
   all times

2. The hub (or main) server. When users connect to the network proxy server's IP, it will re-route those users to this
   server.

3. All additional servers beyond the main server. Once you have at least one server running the proxy and one as the
   hub, any other servers will be considered extra servers that you can teleport to from the hub.

##### Bungee Cord

BungeeCord is a useful software written in-house by the team at SpigotMC. It acts as a proxy between the player's client
and the connected Minecraft servers. End-users of BungeeCord see no difference between it and a normal Minecraft server.

##### Waterfall

Waterfall is a fork of BungeeCord, a proxy used primarily to teleport players between multiple Minecraft servers.

Waterfall focuses on three main areas:

- Stability
- Features
- Scalability

#### Velocity

A Minecraft server proxy with unparalleled server support, scalability, and flexibility. Velocity is licensed under the
GPLv3 license.

- A codebase that is easy to dive into and consistently follows best practices for Java projects as much as reasonably
  possible.
- High performance: handle thousands of players on one proxy.
- A new, refreshing API built from the ground up to be flexible and powerful whilst avoiding design mistakes and
  suboptimal designs from other proxies.
- First-class support for Paper, Sponge, and Forge. (Other implementations may work, but we make every endeavor to
  support these server implementations specifically.)

#### Server Configs 📋

##### Spot Instances

When you want to run a Minecraft server on a spot instance, you can use the following configuration options:

``` yaml
...
spot: <true |false>
...
```

This will enable the server to be run on a spot instance. At the moment, this is only supported by AWS, Azure and GCP.

##### MinecraftProxy Config 📡

If you want to start a server with a Minecraft Proxy, you need to define a MinecraftProxy proxy.

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftProxy
metadata:
  name: minecraft-proxy
spec:
  server:
    cloud: <cloud provider name civo|scaleway|do|hetzner|linode|ovh|equinix|gce|vultr|azure|oci|ionos|aws|vexxhost|multipass|exoscale>
    region: <cloud provider region>
    size: <cloud provider plan>
    ssh:
      port: 22 | or your custom port
      keyfolder: "<path to ssh public and private key>/ssh"
      fail2ban:
        bantime: "<ban time in seconds>"
        maxretry: "<max retry>"
    port: <server port>
    spot: <true |false>
  proxy:
    java:
      openjdk: <jdk version>
      xmx: <xmx memory for the vm>
      xms: <xms memory for the vm>
      options:
        - "-XX:+UseG1GC"
        - "-XX:+ParallelRefProcEnabled"
        - "-XX:MaxGCPauseMillis=200"
      rcon:
        password: <RCON server password>
        port: <RCON server port >
        enabled: <RCON enabled true|false>
        broadcast: <RCON broadcase true|false
    type: "bungeecord|waterfall|velocity"
    version: <version>
```

##### MincraftServer Config 🕹

You need a MinecraftServer manifest file, to describe the underlying compute instance and the Minecraft Server:

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  monitoring:
    enabled: true|false
  server:
    cloud: "provider: civo|scaleway|do|hetzner|linode|ovh|equinix|gce|vultr|azure|oci|ionos|aws|vexxhost|multipass|exoscale"
    region: "region see cloud provider for details eg. fra1"
    size: "see cloud provider docs for details eg. g3.large"
    volumeSize: 100
    ssh:
      port: 22 | or your custom port
      keyfolder: "<path to ssh public and private key>/ssh"
      fail2ban:
        bantime: "<ban time in seconds>"
        maxretry: "<max retry>"
    port: "25565|19132 are the defaults for tcp/udp"
    spot: <true |false>
  minecraft:
    java:
      openjdk: "8|16 use jdk 8 for <1.17 java server version"
      xmx: 2G
      xms: 2G
      options:
        - "-XX:+UseG1GC"
        - "-XX:+ParallelRefProcEnabled"
        - "-XX:MaxGCPauseMillis=200"
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

#### `minectl 🗺` Configuration File Wizard 🧙

```bash
Calls the minectl wizard to create interactively a minectl 🗺 config

Usage:
  minectl wizard [flags]

Examples:
mincetl wizard

Flags:
  -h, --help            help for wizard
  -o, --output string   output folder for the configuration file for minectl 🗺 (default: /Users/dirien/.minectl)

Global Flags:
      --headless              Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)
      --log-encoding string   Set the log encoding: console|json (default: console) (default "console")
      --verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
```

Watch the demo, for more details
[![asciicast](https://asciinema.org/a/439572.svg)](https://asciinema.org/a/439572)

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

Global Flags:
      --headless              Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)
      --log-encoding string   Set the log encoding: console|json (default: console) (default "console")
      --verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
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

Global Flags:
      --headless              Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)
      --log-encoding string   Set the log encoding: console|json (default: console) (default "console")
      --verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
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
  -p, --provider string   The cloud provider - civo|scaleway|do|hetzner|linode|ovh|equinix|gce|vultr|azure|oci|ionos|aws|vexxhost|multipass|exoscale
  -r, --region string     The region for your cloud provider

Global Flags:
      --headless              Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)
      --log-encoding string   Set the log encoding: console|json (default: console) (default "console")
      --verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
```

#### Update Minecraft Server 🆙

Update the Minecraft version. The function uses  `ssh` (port 22).

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

Global Flags:
      --headless              Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)
      --log-encoding string   Set the log encoding: console|json (default: console) (default "console")
      --verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
```

#### RCON Minecraft Server 🔌

Use this function, to connect to the RCON port of your Minecraft Server. RCON is a protocol that allows server
administrators to remotely execute Minecraft commands.

```bash
minectl rcon -h
RCON client to your Minecraft server.

Usage:
  minectl rcon [flags]

Examples:
mincetl rcon  \
    --filename server-do.yaml / \
    --id xxxx

Flags:
  -f, --filename string   Contains the configuration for minectl
  -h, --help              help for rcon
      --id string         contains the server id
```

#### Plugins Minecraft Server ⤴️

> 🚧 Plugins feature is still in beta.

Raw mode, to upload a local plugin file to your server. The function uses  `ssh` (port 22).

```bash
minectl plugins  -h
Manage your plugins for a specific server

Usage:
  minectl plugins [flags]

Examples:
mincetl plugins  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx
        --plugin plugin.jar
    --destination /minecraft/mods

Flags:
  -d, --destination string   Plugin destination location
  -f, --filename string      Contains the configuration for minectl
  -h, --help                 help for plugins
      --id string            contains the server id
  -p, --plugin string        Local plugin file location

Global Flags:
      --headless              Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)
      --log-encoding string   Set the log encoding: console|json (default: console) (default "console")
      --verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
```

#### Monitoring 📊

Monitoring is optional and disabled by default. It can be enabled with simply adding following fields to the
MinecraftServer manifest:

```yaml
...
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  monitoring:
    enabled: true|false
  server:
...
```

Every instance of `minectl 🗺`, has following monitoring components included:

- Prometheus (https://github.com/prometheus/prometheus)
- Node exporter (https://github.com/prometheus/node_exporter)

The `edition:java` has on top following exporter included:

- Minecraft exporter (https://github.com/dirien/minecraft-prometheus-exporter)

You can acces the `prometheus` via

```bash
http://<ip>:9090/graph
```

#### Volumes 💽

With the `volumeSize` property, you can provision an extra volume during the creation phase of the server.

It is always recommended using the provided volume of the server, but in some cases (large mod packs, community server,
etc.) it makes sense to provision a bigger volume separately.

When a separate volume is defined, `minectl 🗺` is automatically installing the Minecraft binaries on this volume.

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  server:
    cloud: linode
    region: eu-central
    size: g6-standard-4
    volumeSize: 100
    ssh:
      port: 22 | or your custom port
      keyfolder: "<path to ssh public and private key>/ssh"
      fail2ban:
        bantime: "<ban time in seconds>"
        maxretry: "<max retry>"
    port: 25565
  minecraft:
...
```

#### Headless Mode 👻

With the global flag `headless`, it is now possible to run `minectl 🗺` in a less human-readable output version. This is
very helpful, if you want to run `minectl 🗺` in workflow.

The flag `verbose` sets the level of logging and with `log-encoding` you can decide between `json` and `console` as
encoding format.

#### Security 🔒

##### SSH Port

Now you can use the `port` property under the `ssh` object, to define the SSH port of the server. Per default, the SSH
port is 22. This helps a lot, to avoid hackers to bruteforce your server.

##### SSH Key

With the `keyfolder` property, you can define the location of your SSH public and private key on your local machine.

##### Fail2Ban

Fail2Ban is an intrusion prevention software framework that protects computer servers from brute-force attacks. With the
property `bantime` you can define the ban time in seconds. With the property `maxretry` you can define the max retry.

If `maxretry` is reached, the IP is banned for the defined time (`bantime`).

#### Getting Started 🎫

- [Civo Java Edition](docs/getting-started-civo.md)
- [Civo Bedrock Edition](docs/getting-started-civo-bedrock.md)
- [Scaleway Java Edition](docs/getting-started-scaleway.md)
- [How to monitor your multi-cloud minectl 🗺 server?](docs/multi-server-monitoring-civo.md)
- [Running a modded LuckyBlocks Minecraft Server on budget 💰 - Scaleway edition](docs/running-minecraft-luckyblocks-budget-scaleway.md)
- [Running a PaperMC SkyBlock server - Hetzner edition](docs/skyblock-papermc-hetzner.md)

### Known Limitation 😵

`minectl 🗺` is still under development. There will be the possibility for breaking changes.

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
- [x] New cloud provider - Vultr [#90](https://github.com/dirien/minectl/issues/90)
- [x] Add Suport for Proxy Server - bungeecord and waterfall [#95](https://github.com/dirien/minectl/issues/95)
- [x] New cloud provider - Azure [#56](https://github.com/dirien/minectl/issues/56)
- [x] New cloud provider - Oracle/OCI [#107](https://github.com/dirien/minectl/issues/107)
- [x] New cloud provider - Ionos Cloud [#218](https://github.com/dirien/minectl/issues/218)
- [x] New cloud provider - AWS [#210](https://github.com/dirien/minectl/pull/210)
- [ ] Much more to come...

### Libraries & Tools 🔥

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
- https://github.com/packethost/packngo
- https://github.com/hashicorp/go-retryablehttp
- https://github.com/melbahja/goph
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
- https://github.com/ionos-cloud/sdk-go
- https://github.com/AlecAivazis/survey
- https://github.com/aws/aws-sdk-go
- https://github.com/gophercloud/gophercloud
- https://github.com/exoscale/egoscale

### Legal Disclaimer 👮

This project is not affiliated with Mojang Studios, XBox Game Studios, Double Eleven or the Minecraft brand.

"Minecraft" is a trademark of Mojang Synergies AB.

Other trademarks referenced herein are property of their respective owners.

Source:

##### [1] https://www.spigotmc.org/wiki/what-is-spigot-craftbukkit-bukkit-vanilla-forg/
