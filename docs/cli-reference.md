# CLI Reference

## Global Flags

All commands support these global flags:

```
--headless              Run in CI mode with logging enabled and human-readable output disabled (default: false)
--log-encoding string   Set the log encoding: console|json (default: "console")
--verbose string        Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal
```

## Commands

### wizard

Create a configuration file interactively.

```bash
minectl wizard [flags]
```

**Flags:**
- `-h, --help` - Help for wizard
- `-o, --output string` - Output folder for the configuration file (default: ~/.minectl)

**Example:**
```bash
minectl wizard
```

---

### create

Create a Minecraft Server.

```bash
minectl create [flags]
```

**Flags:**
- `-f, --filename string` - Location of the manifest file
- `-h, --help` - Help for create
- `-w, --wait` - Wait for Minecraft Server to start (default: true)

**Example:**
```bash
minectl create --filename server-do.yaml
```

---

### delete

Delete a Minecraft Server.

```bash
minectl delete [flags]
```

**Flags:**
- `-f, --filename string` - Location of the manifest file
- `-h, --help` - Help for delete
- `--id string` - Contains the server ID
- `-y, --yes` - Automatically delete the server without confirmation

**Example:**
```bash
minectl delete --filename server-do.yaml --id xxx-xxx-xxx-xxx
```

---

### list

List all Minecraft Servers.

```bash
minectl list [flags]
```

**Flags:**
- `-h, --help` - Help for list
- `-p, --provider string` - The cloud provider (civo|scaleway|do|hetzner|akamai|ovh|gce|vultr|azure|oci|aws|vexxhost|multipass|exoscale)
- `-r, --region string` - The region for your cloud provider

**Example:**
```bash
minectl list --provider civo --region LON1
```

---

### update

Update a Minecraft Server version. Uses SSH (port 22) to connect.

```bash
minectl update [flags]
```

**Flags:**
- `-f, --filename string` - Location of the manifest file
- `-h, --help` - Help for update
- `--id string` - Contains the server ID
- `-k, --ssh-key string` - Specify a specific path for the SSH key

**Example:**
```bash
minectl update --filename server-do.yaml --id xxx-xxx-xxx-xxx
```

---

### rcon

Connect to the RCON port of your Minecraft Server. RCON is a protocol that allows server administrators to remotely execute Minecraft commands.

```bash
minectl rcon [flags]
```

**Flags:**
- `-f, --filename string` - Location of the manifest file
- `-h, --help` - Help for rcon
- `--id string` - Contains the server ID

**Example:**
```bash
minectl rcon --filename server-do.yaml --id xxxx
```

---

### plugins

> This feature is still in beta.

Upload a local plugin file to your server. Uses SSH (port 22) to connect.

```bash
minectl plugins [flags]
```

**Flags:**
- `-d, --destination string` - Plugin destination folder
- `-f, --filename string` - Location of the manifest file
- `-h, --help` - Help for plugins
- `--id string` - Contains the server ID
- `-p, --plugin string` - Location of the plugin
- `-k, --ssh-key string` - Specify a specific path for the SSH key

**Example:**
```bash
minectl plugins \
    --filename server-do.yaml \
    --id xxx-xxx-xxx-xxx \
    --plugin plugin.jar \
    --destination /minecraft/mods
```

## Headless Mode

With the global flag `--headless`, you can run `minectl` in a less human-readable output version. This is helpful when running `minectl` in CI/CD workflows.

The `--verbose` flag sets the level of logging and `--log-encoding` lets you choose between `json` and `console` as the encoding format.

**Example:**
```bash
minectl create --filename server.yaml --headless --verbose info --log-encoding json
```
