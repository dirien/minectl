---
name: minectl
description: "Deploy and manage Minecraft servers across 16 cloud providers using the minectl CLI. Use when the user wants to create, delete, update, or list Minecraft servers, write minectl YAML manifests, connect via RCON, upload plugins, or choose cloud providers and server editions for Minecraft hosting. Triggers on: Minecraft server deployment, minectl commands, Minecraft manifest files, cloud-hosted game servers, or Minecraft edition selection."
---

# minectl

minectl creates and manages Minecraft servers on 16 cloud providers via declarative YAML manifests.

## Creating a Manifest

Every operation starts with a manifest file:

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: my-minecraft
spec:
  server:
    cloud: hetzner
    region: fsn1
    size: cx21
    ssh:
      port: 22
      publickeyfile: ~/.ssh/id_rsa.pub
      fail2ban:
        bantime: 1000
        maxretry: 3
    port: 25565
  minecraft:
    java:
      openjdk: 21
      xmx: 2G
      xms: 2G
      rcon:
        password: changeme
        port: 25575
        enabled: true
    edition: papermc
    version: "1.21"
    eula: true
    properties: |
      max-players=20
      difficulty=normal
      motd=My Minecraft Server
```

### Choosing an Edition

| Edition | Best for |
|---------|----------|
| `java` | Vanilla experience, no mods |
| `papermc` | Performance + plugin support (recommended) |
| `spigot` | Bukkit plugin ecosystem |
| `craftbukkit` | Legacy Bukkit plugins |
| `fabric` | Lightweight mods |
| `forge` | Large mod packs |
| `purpur` | Paper + extra gameplay tweaks |
| `bedrock` | Cross-platform (mobile/console/PC) |
| `nukkit` | Java-based Bedrock server |
| `powernukkit` | Extended Nukkit |

For proxy setups, use `kind: MinecraftProxy` with editions: `bungeecord`, `waterfall`, or `velocity`.

### Optional Features

Add to `spec.server`:
- `spot: true` - Spot/preemptible instances (AWS, Azure, GCE) to save cost
- `arm: true` - ARM instances (Hetzner, AWS, GCE)
- `volumeSize: 50` - Extra storage in GB

Add to `spec`:
- `monitoring: { enabled: true }` - Prometheus + Node Exporter

### Naming Rules

- `metadata.name` must be lowercase alphanumeric with hyphens only
- SSH port must be 22 or between 1024-65535
- `eula` must be `true`

## Server Lifecycle

```bash
# Generate a manifest interactively
minectl wizard

# Create a server (waits for startup by default)
minectl create -f server.yaml

# Verify the server is running after creation
minectl list -p hetzner -r fsn1
# Expected: server appears with a valid ID and status. If absent, check credentials and quota.

# List running servers
minectl list -p hetzner -r fsn1

# Connect via RCON for remote commands
minectl rcon -f server.yaml --id <server-id>

# Upload a plugin
minectl plugins -f server.yaml --id <server-id> -k ~/.ssh/id_rsa \
  -p myplugin.jar -d /minecraft/plugins

# Update Minecraft version (edit version in manifest first)
minectl update -f server.yaml --id <server-id> -k ~/.ssh/id_rsa

# Delete the server
minectl delete -f server.yaml --id <server-id> -y
```

### Error Recovery

| Symptom | Likely cause | Action |
|---------|--------------|--------|
| `create` exits without a server ID | Authentication failure | Verify provider environment variables are exported (see [references/cloud-providers.md](references/cloud-providers.md)) |
| `create` fails with quota/limit error | Cloud account quota exceeded | Switch region, use a smaller `size`, or request a quota increase from the provider |
| `create` succeeds but server absent from `list` | Wrong region or provider flag | Re-run `minectl list` with the exact `-p` and `-r` values from the manifest |
| Plugin upload fails | SSH key mismatch | Confirm `-k` points to the private key matching `spec.server.ssh.publickeyfile` |
| RCON connection refused | RCON not enabled or wrong port | Ensure `rcon.enabled: true` and the port in the manifest matches `rcon.port` |

## Cloud Provider Selection

When helping a user choose a provider, consider:
- **Budget**: Civo, Hetzner, and DigitalOcean are cheapest. Use `spot: true` on AWS/Azure/GCE for savings.
- **Region needs**: Check [references/cloud-providers.md](references/cloud-providers.md) for available regions.
- **Local testing**: Use `multipass` for local VMs with no cloud account needed.

Each provider requires specific environment variables for authentication. See [references/cloud-providers.md](references/cloud-providers.md) for the full list.

## References

- **Cloud providers, credentials, and sizes**: [references/cloud-providers.md](references/cloud-providers.md)
- **Full manifest schema**: [references/manifest-schema.md](references/manifest-schema.md)
- **CLI command flags and details**: [references/cli-commands.md](references/cli-commands.md)
