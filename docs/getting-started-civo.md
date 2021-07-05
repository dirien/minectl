![Civo](img/civo.png)

# Getting Started - Civo edition

## API Key

Get your API Key via https://www.civo.com/account/security 

![key](img/key.png)

## Create SSH Keys

```
ssh-keygen -t rsa -f ./minecraft
```

## Create MinecraftServer config

```bash
apiVersion: ediri.io/minectl/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft-server
spec:
  server:
    cloud: civo
    region: LON1
    size: g3.large
    volumeSize: 100
    ssh: "xxx/ssh/minecraft.pub"
  minecraft:
    java:
      xmx: 2G
      xms: 2G
    properties: |
      level-seed=randomseed
      broadcast-rcon-to-ops=true
      view-distance=10
      enable-jmx-monitoring=false
      server-ip=
      resource-pack-prompt=
      rcon.port=25575
      gamemode=survival
      server-port=25565
      allow-nether=true
      enable-command-block=false
      enable-rcon=false
      sync-chunk-writes=true
      enable-query=false
      op-permission-level=4
      prevent-proxy-connections=false
      resource-pack=
      entity-broadcast-range-percentage=100
      level-name=world
      rcon.password=
      player-idle-timeout=0
      motd=Civo Minecraft
      query.port=25565
      force-gamemode=false
      rate-limit=0
      hardcore=false
      white-list=false
      broadcast-console-to-ops=true
      pvp=true
      spawn-npcs=true
      spawn-animals=true
      snooper-enabled=true
      difficulty=easy
      function-permission-level=2
      network-compression-threshold=256
      text-filtering-config=
      require-resource-pack=false
      spawn-monsters=true
      max-tick-time=60000
      enforce-whitelist=false
      use-native-transport=true
      max-players=100
      resource-pack-sha1=
      spawn-protection=16
      online-mode=true
      enable-status=true
      allow-flight=false
      max-world-size=29999984
```

## minectl 🗺

```bash
minectl create --filename config/server-civo.yaml 

Using cloud provider Civo
Creating instance (minecraft-server)... ⣟ 
Instance (minecraft-server) created
Minecraft Server IP: 74.220.20.145
Minecraft Server ID: a7ad735a-d1e9-4951-9f9b-83221efd945e
```

![instance](img/instance.png)

## Minecraft Client

### Download
Download a Minecraft Client (Java Edition) under https://www.minecraft.net/en-us/get-minecraft

Start your Minecraft Client

![img.png](img/multi.png)

Add your server

![img.png](img/add-server.png)

Join the server

![img.png](img/join.png)

Play the game

![game.png](img/game.png)

## minectl 🗺 

Feed up with your server? Deleting is as easy as creating the server

```bash
minectl delete --filename config/server-civo.yaml --id a7ad735a-d1e9-4951-9f9b-83221efd945e
Using cloud provider Civo
Delete instance (a7ad735a-d1e9-4951-9f9b-83221efd945e)...
```

### Legal Disclaimer 👮

This project is not affiliated with Mojang Studios, XBox Game Studios, Double Eleven or the Minecraft brand.

"Minecraft" is a trademark of Mojang Synergies AB.

Other trademarks referenced herein are property of their respective owners.