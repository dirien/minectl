# CLI Command Reference

## Global Flags

```
--verbose <level>     Logging level: debug|info|warn|error|dpanic|panic|fatal
--log-encoding <fmt>  Log format: console|json (default: console)
--headless            CI mode: enables logging, disables interactive output
```

## Commands

### create

Create a Minecraft server from a manifest file.

```bash
minectl create -f <manifest.yaml> [--wait=true]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--filename` | `-f` | Path to manifest YAML file (required) |
| `--wait` | `-w` | Wait for server startup (default: true) |

Output: table with ID, NAME, REGION, TAGS, IP.

### delete

Delete an existing Minecraft server.

```bash
minectl delete -f <manifest.yaml> --id <server-id> [-y]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--filename` | `-f` | Path to manifest YAML file (required) |
| `--id` | | Server ID (required) |
| `--yes` | `-y` | Skip confirmation prompt |

### list

List all Minecraft servers for a provider/region.

```bash
minectl list -p <provider> [-r <region>]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--provider` | `-p` | Cloud provider name (required) |
| `--region` | `-r` | Region filter (required for civo, gce) |

### update

Update Minecraft version on an existing server via SSH.

```bash
minectl update -f <manifest.yaml> --id <server-id> -k <ssh-key-path>
```

| Flag | Short | Description |
|------|-------|-------------|
| `--filename` | `-f` | Path to manifest YAML file (required) |
| `--id` | | Server ID (required) |
| `--ssh-key` | `-k` | Path to SSH private key (required) |

### rcon

Connect to a server's RCON port for remote commands.

```bash
minectl rcon -f <manifest.yaml> --id <server-id>
```

| Flag | Short | Description |
|------|-------|-------------|
| `--filename` | `-f` | Path to manifest YAML file (required) |
| `--id` | | Server ID (required) |

Opens an interactive RCON shell. RCON must be enabled in the manifest.

### plugins

Upload plugins/mods to a running server (beta).

```bash
minectl plugins -f <manifest.yaml> --id <server-id> -k <ssh-key> -p <plugin.jar> -d <dest>
```

| Flag | Short | Description |
|------|-------|-------------|
| `--filename` | `-f` | Path to manifest YAML file (required) |
| `--id` | | Server ID (required) |
| `--ssh-key` | `-k` | Path to SSH private key (required) |
| `--plugin` | `-p` | Path to plugin JAR file |
| `--destination` | `-d` | Server-side destination folder (e.g., `/minecraft/plugins`) |

### wizard

Interactive questionnaire to generate a manifest file.

```bash
minectl wizard [-o <output-dir>]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output directory (default: `~/.minectl`) |

### version

Display version information.

```bash
minectl version
```

## Typical Workflow

```bash
# 1. Generate manifest interactively
minectl wizard -o ./configs

# 2. Create server
minectl create -f ./configs/config-my-server.yaml

# 3. Connect via RCON
minectl rcon -f ./configs/config-my-server.yaml --id <server-id>

# 4. Upload a plugin
minectl plugins -f ./configs/config-my-server.yaml --id <server-id> \
  -k ~/.ssh/id_rsa -p ./myplugin.jar -d /minecraft/plugins

# 5. Update Minecraft version (edit manifest version first)
minectl update -f ./configs/config-my-server.yaml --id <server-id> -k ~/.ssh/id_rsa

# 6. Delete when done
minectl delete -f ./configs/config-my-server.yaml --id <server-id> -y
```
