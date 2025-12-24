# Configuration

`minectl` uses YAML manifest files to describe Minecraft servers and proxies. You can create these manually or use the interactive wizard.

## Configuration Wizard

Use the wizard to create configuration files interactively:

```bash
minectl wizard
```

[![asciicast](https://asciinema.org/a/439572.svg)](https://asciinema.org/a/439572)

## MinecraftServer Config

You need a MinecraftServer manifest file to describe the underlying compute instance and the Minecraft Server:

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  monitoring:
    enabled: true|false
  server:
    cloud: "civo|scaleway|do|hetzner|akamai|ovh|gce|vultr|azure|oci|aws|vexxhost|multipass|exoscale"
    region: "region see cloud provider for details eg. fra1"
    size: "see cloud provider docs for details eg. g3.large"
    volumeSize: 100
    ssh:
      port: 22 # or your custom port
      publickeyfile: "<path to ssh public key>.pub"
      fail2ban:
        bantime: "<ban time in seconds>"
        maxretry: "<max retry>"
    port: "25565|19132 are the defaults for tcp/udp"
    spot: true|false
    arm: true|false
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
    edition: "java|bedrock|craftbukkit|fabric|forge|papermc|spigot|purpur"
    version: "<version>"
    eula: true
    properties: |
      level-seed=minectlrocks
      broadcast-rcon-to-ops=true
      ...
```

> **Attention:** Please lookup the correct service size if you are setting the `arm` attribute to `true`.

Example configs are available in the [config](https://github.com/dirien/minectl/tree/main/config) folder for all supported cloud providers and Minecraft editions.

## MinecraftProxy Config

If you want to start a server with a Minecraft Proxy, you need to define a MinecraftProxy manifest:

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftProxy
metadata:
  name: minecraft-proxy
spec:
  server:
    cloud: civo|scaleway|do|hetzner|akamai|ovh|gce|vultr|azure|oci|aws|vexxhost|multipass|exoscale
    region: <cloud provider region>
    size: <cloud provider plan>
    ssh:
      port: 22 # or your custom port
      publickeyfile: "<path to ssh public key>.pub"
      fail2ban:
        bantime: "<ban time in seconds>"
        maxretry: "<max retry>"
    port: <server port>
    spot: true|false
    arm: true|false
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
        port: <RCON server port>
        enabled: true|false
        broadcast: true|false
    type: "bungeecord|waterfall|velocity"
    version: <version>
```

## Configuration Options

### Spot Instances

When you want to run a Minecraft server on a spot instance, use:

```yaml
spec:
  server:
    spot: true
```

This is currently supported by AWS, Azure, and GCP.

### EULA

You need to explicitly set the EULA property in the MinecraftServer manifest to indicate your agreement with the [Minecraft End User License](https://minecraft.net/terms):

```yaml
spec:
  minecraft:
    eula: true
```

### SSH Configuration

```yaml
spec:
  server:
    ssh:
      port: 22                    # Custom SSH port (default: 22)
      publickeyfile: "~/.ssh/id_rsa.pub"  # Path to SSH public key
      fail2ban:
        bantime: "600"            # Ban time in seconds
        maxretry: "5"             # Max failed attempts before ban
```

You can also inline the public key content:

```yaml
spec:
  server:
    ssh:
      publickey: "ssh-rsa AAAAB3 ... xxx"
```

### Java Options

Configure JVM settings for optimal performance:

```yaml
spec:
  minecraft:
    java:
      openjdk: "17"
      xmx: 4G
      xms: 4G
      options:
        - "-XX:+UseG1GC"
        - "-XX:+ParallelRefProcEnabled"
        - "-XX:MaxGCPauseMillis=200"
        - "-XX:+UnlockExperimentalVMOptions"
        - "-XX:+DisableExplicitGC"
        - "-XX:+AlwaysPreTouch"
        - "-XX:G1NewSizePercent=30"
        - "-XX:G1MaxNewSizePercent=40"
        - "-XX:G1HeapRegionSize=8M"
        - "-XX:G1ReservePercent=20"
        - "-XX:G1HeapWastePercent=5"
        - "-XX:G1MixedGCCountTarget=4"
        - "-XX:InitiatingHeapOccupancyPercent=15"
        - "-XX:G1MixedGCLiveThresholdPercent=90"
        - "-XX:G1RSetUpdatingPauseTimePercent=5"
        - "-XX:SurvivorRatio=32"
        - "-XX:+PerfDisableSharedMem"
        - "-XX:MaxTenuringThreshold=1"
```

### Server Properties

Configure Minecraft server properties directly in the manifest:

```yaml
spec:
  minecraft:
    properties: |
      level-seed=minectlrocks
      broadcast-rcon-to-ops=true
      view-distance=10
      max-players=20
      difficulty=normal
      gamemode=survival
      pvp=true
      spawn-monsters=true
      spawn-animals=true
```
