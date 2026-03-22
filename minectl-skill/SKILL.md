---
name: minectl
description: "Use this skill for every message containing the word \"minectl\". Claude has zero built-in knowledge of minectl — it is an obscure, third-party CLI for deploying Minecraft servers to cloud providers. Without this skill, any response about minectl will be fabricated. This covers all minectl topics: manifests, create/delete/update commands, RCON troubleshooting, authentication errors, cloud provider selection, instance sizing, budget optimization, TPS tuning, plugin uploads, and edition choice. Do NOT trigger for Minecraft servers managed via Docker, Kubernetes, Ansible, Terraform, or Pulumi — only for minectl-based workflows."
---

# minectl

minectl creates and manages Minecraft servers on 15 cloud providers via declarative YAML manifests.

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

### Edition-Specific Tuning

Different editions need different server.properties and resource allocations. Getting these wrong causes crashes, kicks, or poor performance — this is where most new server admins struggle.

**Forge / Fabric (modded servers):**
- Set `max-tick-time=90000` (or `-1` to disable). Modded servers regularly exceed the default 60s watchdog, especially during chunk generation or mod initialization. Leaving the default causes the server to kill itself under normal modded load.
- Set `allow-flight=true`. Many mods (Botania, Mekanism jetpacks, Elytra mods, etc.) trigger the vanilla flight detection. Without this, modded players get kicked constantly.
- Allocate generous heap: 8G minimum for light mods, 12G+ for heavy packs (100+ mods). Set `xms` equal to `xmx` to avoid expensive heap resizing under load.
- Always add `volumeSize: 50` (or more). Mod packs, world data with modded chunks, and backups consume significant disk. Running out of disk silently corrupts worlds.
- Upload mods to `/minecraft/mods` (not `/minecraft/plugins`).
- Enable monitoring (`monitoring.enabled: true`) — modded servers are resource-hungry and benefit from Prometheus metrics to catch memory/CPU issues early.

**PaperMC / Spigot / CraftBukkit / Purpur (plugin servers):**
- Default `max-tick-time=60000` is usually fine.
- Keep `allow-flight=false` unless a specific plugin requires it.
- 2-4G heap is sufficient for up to 20 players with typical plugins.
- Upload plugins to `/minecraft/plugins`.
- PaperMC is the recommended default — it's the most performant fork and supports all Bukkit/Spigot plugins.

**Bedrock / Nukkit / PowerNukkit:**
- No `java` block needed in the manifest (Bedrock is a native binary, not Java-based).
- Default port is `19132` (UDP), not `25565`.
- RCON is not available for Bedrock edition.
- Lower resource requirements: 1-2G RAM is typically sufficient.

### Instance Sizing Guide

| Use case | Players | Recommended RAM | Example sizes |
|----------|---------|----------------|---------------|
| Vanilla / Paper (casual) | 1-10 | 2-4 GB | Hetzner cx21, DO s-2vcpu-4gb |
| Paper with plugins | 10-20 | 4-8 GB | Hetzner cpx31, DO s-4vcpu-8gb, AWS t3.large |
| Forge/Fabric (light mods) | 1-10 | 8 GB | AWS t3.xlarge, Hetzner cpx41 |
| Forge/Fabric (heavy mods) | 5-15 | 16 GB+ | AWS m5.xlarge, Hetzner cpx51 |
| Proxy (BungeeCord/Velocity) | N/A | 1-2 GB | Smallest available |

When using spot instances (`spot: true`), prefer regions with large capacity pools (e.g., `us-east-1` for AWS) to minimize interruption risk. Always warn users that spot instances can be terminated with 2 minutes notice — world backups are essential.

### Security Recommendations

Always include these in every manifest:
- `fail2ban` block under SSH config (prevents brute-force attacks)
- Strong RCON password (never deploy with "changeme")
- Consider `white-list=true` and `enforce-whitelist=true` for private servers
- `online-mode=true` to verify player accounts (prevents unauthorized access)

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
