# minectl 🗺

`minectl`️️ is a cli for creating Minecraft (java or bedrock) server on different cloud provider.

It is a private side project of me, to learn more about Go, CLI and multi cloud.

### TL;DR 🚀

Install via homebrew:

```bash
brew tap dirien/homebrew-dirien
brew install minectl
```

Linux or Windows user, can directly download (or use `curl`/`wget`) the binary via
the [release page](https://github.com/dirien/minectl/releases).

### Usage ⚙

#### Access Token

`minectl` is completely build on zero-trust. It does not save any API Tokens, instead it looks them up in the ENV
variables.

##### Civo
```
export CIVO_TOKEN=xx
```

##### Digital Ocean
```
export DIGITALOCEAN_TOKEN=xxx
```

##### Scaleway
```
export ACCESS_KEY=xxx
export SECRET_KEY=yyy
export ORGANISATION_ID=zzz
```

#### Server Config 📋

You need a MinecraftServer manifest file, to define some informations regarding the VM and the Minecraft Server:

```yaml
apiVersion: ediri.io/minectl/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  server:
    cloud: "provider: civo|scaleway|do"
    region: "region see cloud provider for details eg. fra1"
    size: "see cloud provider docs for details eg. g3.large"
    volumeSize: 100
    ssh: "<path to>/ssh.pub"
  minecraft:
    java:
      xmx: 2G
      xms: 2G
    edition: "java|bedrock"
    properties: |
      level-seed=stackitminecraftrocks
      broadcast-rcon-to-ops=true
      ...
```

I created some example configs in the [config](config) folder for currently supported cloud provider and Minecraft editions.

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

#### Getting Started

- [Civo Java Edition](docs/getting-started-civo.md)
- [Civo Bedrock Edition](docs/getting-started-civo-bedrock.md)
- [Scaleway Java Edition](docs/getting-started-scaleway.md)

### Supported cloud provider ☁

+ Civo (https://www.civo.com/)
+ Scaleway (https://www.scaleway.com)
+ DigitalOcean (https://www.digitalocean.com/)

### Known Limitation 😵

`minectl` is still under development and supports only creation and deletion of server. There is no mod or plugin
functionality for the Minecraft servers.

### Contributing 🤝

#### Contributing via GitHub

Feel free to join.

#### License

Apache License, Version 2.0

### Roadmap 🛣️

- [x] Support Bedrock edition [#10](https://github.com/dirien/minectl/issues/10)
- [ ] List Minecraft Server
- [ ] Update Minecraft Server
- [ ] Support Mods and Plugins
- [ ] Add additional cloud provider
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

### Legal Disclaimer 👮

This project is not affiliated with Mojang Studios, XBox Game Studios, Double Eleven or the Minecraft brand.

"Minecraft" is a trademark of Mojang Synergies AB.

Other trademarks referenced herein are property of their respective owners.