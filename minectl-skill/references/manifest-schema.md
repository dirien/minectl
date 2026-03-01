# Manifest Schema Reference

## MinecraftServer Manifest

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer       # or MinecraftProxy
metadata:
  name: <server-name>       # lowercase alphanumeric and hyphens only
spec:
  monitoring:
    enabled: <bool>         # optional, enables Prometheus + Node Exporter
  server:
    cloud: <provider>       # required: civo|do|hetzner|aws|gce|azure|scaleway|akamai|ovh|vultr|oci|ionos|vexxhost|exoscale|fuga|multipass
    region: <region>        # provider-specific region/zone
    size: <instance-type>   # provider-specific instance size
    volumeSize: <int>       # optional, additional volume in GB
    arm: <bool>             # optional, use ARM instances
    spot: <bool>            # optional, use spot/preemptible instances (AWS/Azure/GCE)
    ssh:
      port: <int>           # SSH port (22 or 1024-65535)
      publickeyfile: <path> # path to SSH public key file
      # OR
      publickey: <string>   # inline SSH public key
      fail2ban:
        bantime: <int>      # ban duration in seconds
        maxretry: <int>     # max failed attempts before ban
        ignoreip: <string>  # optional, IPs to whitelist
    port: <int>             # Minecraft server port (default 25565)
  minecraft:
    java:                   # required for Java editions
      openjdk: <int>        # OpenJDK version: 8, 16, 17, 21
      xmx: <string>        # max heap (e.g., "2G")
      xms: <string>        # initial heap (e.g., "2G")
      options: [<string>]   # optional, extra JVM flags
      rcon:
        password: <string>  # RCON password
        port: <int>         # RCON port (default 25575)
        enabled: <bool>     # enable RCON
        broadcast: <bool>   # broadcast RCON to ops
    edition: <string>       # required: java|papermc|spigot|craftbukkit|fabric|forge|purpur|bedrock|nukkit|powernukkit
    version: <string>       # required: Minecraft version (e.g., "1.21", "1.20.4-388")
    eula: <bool>            # required: must be true
    properties: |           # optional: server.properties content
      key=value
```

## MinecraftProxy Manifest

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftProxy
metadata:
  name: <proxy-name>
spec:
  # same server block as MinecraftServer
  proxy:
    edition: <string>       # bungeecord|waterfall|velocity
    version: <string>
    java:
      openjdk: <int>
      xmx: <string>
      xms: <string>
```

## Editions

| Edition | Type | Notes |
|---------|------|-------|
| `java` | Vanilla | Official Mojang server |
| `papermc` | Fork | High-performance Spigot fork |
| `spigot` | Fork | CraftBukkit fork with plugin API |
| `craftbukkit` | Fork | Modified server with Bukkit API |
| `fabric` | Modloader | Lightweight mod loader |
| `forge` | Modloader | Popular mod platform |
| `purpur` | Fork | Paper fork with extra features |
| `bedrock` | Bedrock | Cross-platform (no Java config needed) |
| `nukkit` | Bedrock | Java-based Bedrock server |
| `powernukkit` | Bedrock | Nukkit fork |

## Proxy Editions

| Edition | Notes |
|---------|-------|
| `bungeecord` | Original Minecraft proxy |
| `waterfall` | BungeeCord fork by PaperMC |
| `velocity` | Modern, high-performance proxy |

## Validation Rules

- `metadata.name`: must match `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
- `spec.server.ssh.port`: must be 22 or between 1024-65535
- `spec.minecraft.eula`: must be `true`
- `spec.server.cloud`: must be a valid provider name
