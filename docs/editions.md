# Minecraft Editions

> `minectl` is not(!) providing any pre-compiled binaries of Minecraft or downloading a pre-compiled version.
>
> Every _non-vanilla_ version will be compiled during the build phase of your server.

## Server Editions

### Vanilla (Minecraft: Java Edition or Bedrock Edition)

The Vanilla software is the original, untouched, unmodified Minecraft server software created and distributed directly by Mojang.

### CraftBukkit

CraftBukkit is a lightly modified version of the Vanilla software allowing it to be able to run Bukkit plugins.

### Spigot

Spigot is the most popular Minecraft server software in the world. Spigot is a modified version of CraftBukkit with hundreds of improvements and optimizations that can only make CraftBukkit shrink in shame.

### PaperMC

Paper (formerly known as PaperSpigot, distributed via the Paperclip patch utility) is a high performance fork of Spigot.

### Purpur

Purpur is a drop-in replacement for Paper servers designed for configurability and new, fun, exciting gameplay features.

### Forge

Forge is well known for being able to use Forge Mods which are direct modifications to the Minecraft program code. In doing so, Forge Mods can change the gaming-feel drastically as a result of this.

### Fabric

Fabric is also a mod loader like Forge with some improvements. It's lightweight and faster and it may be the best mod loader in the future because it's doing very well.

## Proxy Editions

Network proxy server is what you set up and use as the controller of a network of servers - this is the server that connects all of your playable servers together so people can log in through one server IP, and then teleport between the separate servers in-game rather than having to log out and into each different one.

A server network typically consists of the following servers:

1. **The proxy server** itself running the desired software (BungeeCord being the most widely used and supported). This is the server that you would be advertising the IP of, as all players should be logging in through the proxy server at all times.

2. **The hub (or main) server**. When users connect to the network proxy server's IP, it will re-route those users to this server.

3. **Additional servers** beyond the main server. Once you have at least one server running the proxy and one as the hub, any other servers will be considered extra servers that you can teleport to from the hub.

### BungeeCord

BungeeCord is a useful software written in-house by the team at SpigotMC. It acts as a proxy between the player's client and the connected Minecraft servers. End-users of BungeeCord see no difference between it and a normal Minecraft server.

### Waterfall

Waterfall is a fork of BungeeCord, a proxy used primarily to teleport players between multiple Minecraft servers.

Waterfall focuses on three main areas:

- Stability
- Features
- Scalability

### Velocity

A Minecraft server proxy with unparalleled server support, scalability, and flexibility. Velocity is licensed under the GPLv3 license.

- A codebase that is easy to dive into and consistently follows best practices for Java projects as much as reasonably possible
- High performance: handle thousands of players on one proxy
- A new, refreshing API built from the ground up to be flexible and powerful whilst avoiding design mistakes and suboptimal designs from other proxies
- First-class support for Paper, Sponge, and Forge (other implementations may work, but we make every endeavor to support these server implementations specifically)

---

Source: [SpigotMC Wiki](https://www.spigotmc.org/wiki/what-is-spigot-craftbukkit-bukkit-vanilla-forg/)
