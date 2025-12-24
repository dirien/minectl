# Features

## Monitoring

Monitoring is optional and disabled by default. Enable it by adding the following to your MinecraftServer manifest:

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  monitoring:
    enabled: true
  server:
    ...
```

Every instance of `minectl` with monitoring enabled includes:

- [Prometheus](https://github.com/prometheus/prometheus)
- [Node exporter](https://github.com/prometheus/node_exporter)

The `edition:java` has an additional exporter included:

- [Minecraft exporter](https://github.com/dirien/minecraft-prometheus-exporter)

Access Prometheus via:

```
http://<ip>:9090/graph
```

For more details on monitoring, see [How to monitor your multi-cloud minectl server](multi-server-monitoring-civo.md).

## Volumes

With the `volumeSize` property, you can provision an extra volume during the creation phase of the server.

It is always recommended to use the provided volume of the server, but in some cases (large mod packs, community server, etc.) it makes sense to provision a bigger volume separately.

When a separate volume is defined, `minectl` automatically installs the Minecraft binaries on this volume.

```yaml
apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  server:
    cloud: akamai
    region: eu-central
    size: g6-standard-4
    volumeSize: 100
    ssh:
      port: 22
      publickeyfile: "<path to ssh public key>.pub"
      fail2ban:
        bantime: "600"
        maxretry: "5"
    port: 25565
  minecraft:
    ...
```

## Security

### SSH Port

You can use the `port` property under the `ssh` object to define the SSH port of the server. The default SSH port is 22. Changing this helps avoid brute-force attacks on your server.

```yaml
spec:
  server:
    ssh:
      port: 2222
```

### SSH Key

With the `publickeyfile` property, you can define the location of your SSH public key on your local machine:

```yaml
spec:
  server:
    ssh:
      publickeyfile: "~/.ssh/id_rsa.pub"
```

Alternatively, use the `publickey` property to define the content of your SSH public key directly:

```yaml
spec:
  server:
    ssh:
      publickey: "ssh-rsa AAAAB3 ... xxx"
```

If you need to update or upload a plugin to your server, provide the SSH private key in the command with the `--ssh-key` flag:

```bash
minectl update --filename server.yaml --id xxx --ssh-key ~/.ssh/id_rsa
```

### Fail2Ban

Fail2Ban is an intrusion prevention software framework that protects computer servers from brute-force attacks.

- `bantime` - The ban time in seconds
- `maxretry` - The maximum number of failed attempts before banning

If `maxretry` is reached, the IP is banned for the defined time (`bantime`).

```yaml
spec:
  server:
    ssh:
      fail2ban:
        bantime: "600"
        maxretry: "5"
```

## Spot Instances

Run your Minecraft server on spot instances for cost savings:

```yaml
spec:
  server:
    spot: true
```

Currently supported by:
- AWS
- Azure
- GCP

> **Note:** Spot instances can be terminated by the cloud provider with short notice. Use for non-critical or easily recoverable servers.

## ARM Support

Run on ARM-based instances for cost efficiency:

```yaml
spec:
  server:
    arm: true
```

> **Attention:** Please lookup the correct service size when setting the `arm` attribute to `true`. Not all instance types are available on ARM.
