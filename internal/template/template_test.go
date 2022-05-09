package template

import (
	"testing"

	"github.com/minectl/internal/model"
	"github.com/stretchr/testify/assert"
)

var (
	bedrock = model.MinecraftResource{
		Spec: model.Spec{
			Minecraft: model.Minecraft{
				Edition:    "bedrock",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.10.04",
				Eula:       false,
			},
			Server: model.Server{
				Port: 19132,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	bedrockNoMon = model.MinecraftResource{
		Spec: model.Spec{
			Minecraft: model.Minecraft{
				Edition:    "bedrock",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.10.04",
				Eula:       false,
			},
			Server: model.Server{
				Port: 19132,
			},
		},
	}
	java = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					OpenJDK: 16,
					Xms:     "2G",
					Xmx:     "2G",
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "java",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	javaAdditionalOptions = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					OpenJDK: 16,
					Xms:     "2G",
					Xmx:     "2G",
					Options: []string{
						"-XX:+UseG1GC",
						"-XX:+ParallelRefProcEnabled",
					},
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "java",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	papermc = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 16,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "papermc",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.1-157",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	craftbukkit = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 16,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "craftbukkit",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.1-138",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	fabric = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 16,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "fabric",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.1-138",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	forge = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 16,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "forge",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.1-138",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}
	spigot = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 16,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "spigot",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.1-138",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: true,
			},
		},
	}

	nukkit = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 19132,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 8,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "nukkit",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.0-SNAPSHOT",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: false,
			},
		},
	}

	fabricNoMon = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 25565,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 16,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "fabric",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.17.1-138",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: false,
			},
		},
	}

	powerNukkit = model.MinecraftResource{
		Spec: model.Spec{
			Server: model.Server{
				Port: 19132,
			},
			Minecraft: model.Minecraft{
				Java: model.Java{
					Xms:     "2G",
					Xmx:     "2G",
					OpenJDK: 8,
					Rcon: model.Rcon{
						Port:      2,
						Password:  "test",
						Enabled:   true,
						Broadcast: true,
					},
				},
				Edition:    "powernukkit",
				Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false\n",
				Version:    "1.5.1.0-PN",
				Eula:       true,
			},
			Monitoring: model.Monitoring{
				Enabled: false,
			},
		},
	}

	bedrockBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=19132
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/bin/sh -c "LD_LIBRARY_PATH=. ./bedrock_server"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl unzip fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter



sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
URL=$(curl -s https://bedrock-version.minectl.ediri.online/binary/1.17.10.04)
curl -sLSf $URL > /tmp/bedrock-server.zip
unzip -o /tmp/bedrock-server.zip -d /minecraft
chmod +x /minecraft/bedrock_server
echo "eula=false" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	javaBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
URL=$(curl -s https://java-version.minectl.ediri.online/binary/1.17)
curl -sLSf $URL > /minecraft/server.jar
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	bedrockCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - unzip
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       
       server-port=19132
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/bin/sh -c "LD_LIBRARY_PATH=. ./bedrock_server"
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL=$(curl -s https://bedrock-version.minectl.ediri.online/binary/1.17.10.04)
  - curl -sLSf $URL > /tmp/bedrock-server.zip
  - unzip -o /tmp/bedrock-server.zip -d /minecraft
  - chmod +x /minecraft/bedrock_server
  - echo "eula=false" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	javaCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  - name: minecraft_exporter
    shell: /bin/false
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
        - job_name: 'minecraft_exporter'
          scrape_interval: 1m
          static_configs:
            - targets: ['localhost:9150']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/minecraft-exporter.service
    content: |
      [Unit]
      Description=Minecraft Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=minecraft_exporter
      Group=minecraft_exporter
      Type=simple
      ExecStart=/usr/local/bin/minecraft-exporter \
          --mc.rcon-password=test
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - export MINECRAFT_EXPORTER_VERSION=0.13.0
  - curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp minecraft-exporter /usr/local/bin
  - chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
  - systemctl start minecraft-exporter.service
  - systemctl enable minecraft-exporter.service
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL=$(curl -s https://java-version.minectl.ediri.online/binary/1.17)
  - curl -sLSf $URL > /minecraft/server.jar
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	paperMCCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  - name: minecraft_exporter
    shell: /bin/false
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
        - job_name: 'minecraft_exporter'
          scrape_interval: 1m
          static_configs:
            - targets: ['localhost:9150']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/minecraft-exporter.service
    content: |
      [Unit]
      Description=Minecraft Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=minecraft_exporter
      Group=minecraft_exporter
      Type=simple
      ExecStart=/usr/local/bin/minecraft-exporter \
          --mc.rcon-password=test
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - export MINECRAFT_EXPORTER_VERSION=0.13.0
  - curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp minecraft-exporter /usr/local/bin
  - chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
  - systemctl start minecraft-exporter.service
  - systemctl enable minecraft-exporter.service
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL="https://papermc.io/api/v2/projects/paper/versions/1.17.1/builds/157/downloads/paper-1.17.1-157.jar"
  - curl -sLSf $URL > /minecraft/server.jar
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	bedrockBashMountWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=19132
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/bin/sh -c "LD_LIBRARY_PATH=. ./bedrock_server"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl unzip fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter



sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sdc
mount /dev/sdc /minecraft
echo "/dev/sdc /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL=$(curl -s https://bedrock-version.minectl.ediri.online/binary/1.17.10.04)
curl -sLSf $URL > /tmp/bedrock-server.zip
unzip -o /tmp/bedrock-server.zip -d /minecraft
chmod +x /minecraft/bedrock_server
echo "eula=false" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	javaBashMountWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sdc
mount /dev/sdc /minecraft
echo "/dev/sdc /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL=$(curl -s https://java-version.minectl.ediri.online/binary/1.17)
curl -sLSf $URL > /minecraft/server.jar
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	javaBashAdditionalOptionsMountWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -XX:+UseG1GC -XX:+ParallelRefProcEnabled -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sdc
mount /dev/sdc /minecraft
echo "/dev/sdc /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL=$(curl -s https://java-version.minectl.ediri.online/binary/1.17)
curl -sLSf $URL > /minecraft/server.jar
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	paperMCBashMountWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sdc
mount /dev/sdc /minecraft
echo "/dev/sdc /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL="https://papermc.io/api/v2/projects/paper/versions/1.17.1/builds/157/downloads/paper-1.17.1-157.jar"
curl -sLSf $URL > /minecraft/server.jar
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	craftbukkitCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  - name: minecraft_exporter
    shell: /bin/false
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
        - job_name: 'minecraft_exporter'
          scrape_interval: 1m
          static_configs:
            - targets: ['localhost:9150']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/minecraft-exporter.service
    content: |
      [Unit]
      Description=Minecraft Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=minecraft_exporter
      Group=minecraft_exporter
      Type=simple
      ExecStart=/usr/local/bin/minecraft-exporter \
          --mc.rcon-password=test
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - export MINECRAFT_EXPORTER_VERSION=0.13.0
  - curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp minecraft-exporter /usr/local/bin
  - chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
  - systemctl start minecraft-exporter.service
  - systemctl enable minecraft-exporter.service
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - export HOME=/tmp/
  - apt-get install -y git
  - git config --global user.email "minectl@github.com"
  - git config --global user.name "minectl"
  - URL="https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
  - mkdir /tmp/build
  - cd /tmp/build
  - curl -sLSf $URL > BuildTools.jar
  - git config --global --unset core.autocrlf
  - java -jar BuildTools.jar --rev 1.17.1-138 --compile craftbukkit
  - cp craftbukkit-1.17.1-138.jar /minecraft/server.jar
  - rm -rf /tmp/build
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	craftbukkitBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
apt-get install -y git
URL="https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
mkdir /tmp/build
cd /tmp/build
curl -sLSf $URL > BuildTools.jar
git config --global --unset core.autocrlf
java -jar BuildTools.jar --rev 1.17.1-138 --compile craftbukkit
cp craftbukkit-1.17.1-138.jar /minecraft/server.jar
rm -rf /tmp/build
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	fabricCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  - name: minecraft_exporter
    shell: /bin/false
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
        - job_name: 'minecraft_exporter'
          scrape_interval: 1m
          static_configs:
            - targets: ['localhost:9150']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/minecraft-exporter.service
    content: |
      [Unit]
      Description=Minecraft Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=minecraft_exporter
      Group=minecraft_exporter
      Type=simple
      ExecStart=/usr/local/bin/minecraft-exporter \
          --mc.rcon-password=test
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - export MINECRAFT_EXPORTER_VERSION=0.13.0
  - curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp minecraft-exporter /usr/local/bin
  - chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
  - systemctl start minecraft-exporter.service
  - systemctl enable minecraft-exporter.service
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL="https://maven.fabricmc.net/net/fabricmc/fabric-installer/0.10.2/fabric-installer-0.10.2.jar"
  - mkdir /tmp/build
  - cd /tmp/build
  - curl -sLSf $URL > fabric-installer.jar
  - java -jar fabric-installer.jar server -downloadMinecraft
  - echo "serverJar=minecraft-server.jar" > /minecraft/fabric-server-launcher.properties
  - cp /tmp/build/fabric-server-launch.jar /minecraft/server.jar
  - cp /tmp/build/server.jar /minecraft/minecraft-server.jar
  - rm -rf /tmp/build
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	fabricBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL="https://maven.fabricmc.net/net/fabricmc/fabric-installer/0.10.2/fabric-installer-0.10.2.jar"
mkdir /tmp/build
cd /tmp/build
curl -sLSf $URL > fabric-installer.jar
java -jar fabric-installer.jar server -downloadMinecraft -mcversion 1.17.1-138
echo "serverJar=minecraft-server.jar" > /minecraft/fabric-server-launcher.properties
cp /tmp/build/fabric-server-launch.jar /minecraft/server.jar
cp /tmp/build/server.jar /minecraft/minecraft-server.jar
rm -rf /tmp/build
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	forgeCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  - name: minecraft_exporter
    shell: /bin/false
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
        - job_name: 'minecraft_exporter'
          scrape_interval: 1m
          static_configs:
            - targets: ['localhost:9150']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/minecraft-exporter.service
    content: |
      [Unit]
      Description=Minecraft Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=minecraft_exporter
      Group=minecraft_exporter
      Type=simple
      ExecStart=/usr/local/bin/minecraft-exporter \
          --mc.rcon-password=test
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - export MINECRAFT_EXPORTER_VERSION=0.13.0
  - curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp minecraft-exporter /usr/local/bin
  - chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
  - systemctl start minecraft-exporter.service
  - systemctl enable minecraft-exporter.service
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL="https://maven.minecraftforge.net/net/minecraftforge/forge/1.17.1-138/forge-1.17.1-138-installer.jar"
  - mkdir /tmp/build
  - cd /tmp/build
  - mkdir minecraft
  - curl -sLSf $URL > forge-installer.jar
  - java -jar forge-installer.jar --installServer /minecraft
  - cp /minecraft/forge-1.17.1-138.jar /minecraft/server.jar
  - rm -rf /tmp/build
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	forgeBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL="https://maven.minecraftforge.net/net/minecraftforge/forge/1.17.1-138/forge-1.17.1-138-installer.jar"
mkdir /tmp/build
cd /tmp/build
mkdir minecraft
curl -sLSf $URL > forge-installer.jar
java -jar forge-installer.jar --installServer /minecraft
cp /minecraft/forge-1.17.1-138.jar /minecraft/server.jar
rm -rf /tmp/build
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	spigotCloudInitWant = `#cloud-config
users:
  - default
  - name: prometheus
    shell: /bin/false
  - name: node_exporter
    shell: /bin/false
  - name: minecraft_exporter
    shell: /bin/false
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /tmp/prometheus.yml
    content: |
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9090']
        - job_name: 'node_exporter'
          scrape_interval: 5s
          static_configs:
            - targets: ['localhost:9100']
        - job_name: 'minecraft_exporter'
          scrape_interval: 1m
          static_configs:
            - targets: ['localhost:9150']
  - path: /etc/systemd/system/prometheus.service
    content: |
      [Unit]
      Description=Prometheus
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=prometheus
      Group=prometheus
      Type=simple
      ExecStart=/usr/local/bin/prometheus \
          --config.file /etc/prometheus/prometheus.yml \
          --storage.tsdb.path /var/lib/prometheus/ \
          --web.console.templates=/etc/prometheus/consoles \
          --web.console.libraries=/etc/prometheus/console_libraries
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/node_exporter.service
    content: |
      [Unit]
      Description=Node Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=node_exporter
      Group=node_exporter
      Type=simple
      ExecStart=/usr/local/bin/node_exporter
      [Install]
      WantedBy=multi-user.target
  - path: /etc/systemd/system/minecraft-exporter.service
    content: |
      [Unit]
      Description=Minecraft Exporter
      Wants=network-online.target
      After=network-online.target
      [Service]
      User=minecraft_exporter
      Group=minecraft_exporter
      Type=simple
      ExecStart=/usr/local/bin/minecraft-exporter \
          --mc.rcon-password=test
      [Install]
      WantedBy=multi-user.target
  
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - export ARCH=amd64
  - MACHINE_TYPE=$(uname -i)
  - if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi
  - export PROM_VERSION=2.35.0
  - mkdir /etc/prometheus
  - mkdir /var/lib/prometheus
  - curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
  - cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
  - chown prometheus:prometheus /usr/local/bin/prometheus
  - chown prometheus:prometheus /usr/local/bin/promtool
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
  - cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
  - chown -R prometheus:prometheus /var/lib/prometheus
  - chown -R prometheus:prometheus /etc/prometheus/consoles
  - chown -R prometheus:prometheus /etc/prometheus/console_libraries
  - mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
  - chown prometheus:prometheus /etc/prometheus/prometheus.yml
  - systemctl daemon-reload
  - systemctl start prometheus
  - systemctl enable prometheus

  - export NODE_EXPORTER_VERSION=1.3.1
  - curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
  - chown node_exporter:node_exporter /usr/local/bin/node_exporter
  - systemctl daemon-reload
  - systemctl start node_exporter
  - systemctl enable node_exporter
  - export MINECRAFT_EXPORTER_VERSION=0.13.0
  - curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
  - cp minecraft-exporter /usr/local/bin
  - chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
  - systemctl start minecraft-exporter.service
  - systemctl enable minecraft-exporter.service
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - export HOME=/tmp/
  - apt-get install -y git
  - git config --global user.email "minectl@github.com"
  - git config --global user.name "minectl"
  - URL="https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
  - mkdir /tmp/build
  - cd /tmp/build
  - curl -sLSf $URL > BuildTools.jar
  - git config --global --unset core.autocrlf
  - java -jar BuildTools.jar --rev 1.17.1-138 
  - cp spigot-1.17.1-138.jar /minecraft/server.jar
  - rm -rf /tmp/build
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	spigotBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'minecraft_exporter'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9150']
EOF
tee /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF
tee /etc/systemd/system/minecraft-exporter.service <<EOF
[Unit]
Description=Minecraft Exporter
Wants=network-online.target
After=network-online.target
[Service]
User=minecraft_exporter
Group=minecraft_exporter
Type=simple
ExecStart=/usr/local/bin/minecraft-exporter \
  --mc.rcon-password=test
[Install]
WantedBy=multi-user.target
EOF

tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban
useradd prometheus -s /bin/false
useradd node_exporter -s /bin/false
useradd minecraft_exporter -s /bin/false

export ARCH=amd64
MACHINE_TYPE=$(uname -i)
if test "$MACHINE_TYPE" = 'aarch64'; then export ARCH=arm64; fi

export PROM_VERSION=2.35.0
mkdir /etc/prometheus
mkdir /var/lib/prometheus
curl -sSL https://github.com/prometheus/prometheus/releases/download/v$PROM_VERSION/prometheus-$PROM_VERSION.linux-$ARCH.tar.gz | tar -xz
cp prometheus-$PROM_VERSION.linux-$ARCH/prometheus /usr/local/bin/
cp prometheus-$PROM_VERSION.linux-$ARCH/promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r prometheus-$PROM_VERSION.linux-$ARCH/consoles /etc/prometheus
cp -r prometheus-$PROM_VERSION.linux-$ARCH/console_libraries /etc/prometheus
chown -R prometheus:prometheus /var/lib/prometheus
chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries
mv /tmp/prometheus.yml /etc/prometheus/prometheus.yml
chown prometheus:prometheus /etc/prometheus/prometheus.yml
systemctl daemon-reload
systemctl start prometheus
systemctl enable prometheus

export NODE_EXPORTER_VERSION=1.3.1
curl -sSL https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp node_exporter-$NODE_EXPORTER_VERSION.linux-$ARCH/node_exporter /usr/local/bin
chown node_exporter:node_exporter /usr/local/bin/node_exporter
systemctl daemon-reload
systemctl start node_exporter
systemctl enable node_exporter

export MINECRAFT_EXPORTER_VERSION=0.13.0
curl -sSL https://github.com/dirien/minecraft-prometheus-exporter/releases/download/v$MINECRAFT_EXPORTER_VERSION/minecraft-exporter_$MINECRAFT_EXPORTER_VERSION.linux-$ARCH.tar.gz | tar -xz
cp minecraft-exporter /usr/local/bin
chown minecraft_exporter:minecraft_exporter /usr/local/bin/minecraft-exporter
systemctl start minecraft-exporter.service
systemctl enable minecraft-exporter.service


sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
apt-get install -y git
URL="https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
mkdir /tmp/build
cd /tmp/build
curl -sLSf $URL > BuildTools.jar
git config --global --unset core.autocrlf
java -jar BuildTools.jar --rev 1.17.1-138 
cp spigot-1.17.1-138.jar /minecraft/server.jar
rm -rf /tmp/build
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	bedrockBashNoMonWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=19132
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

EOF
tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/bin/sh -c "LD_LIBRARY_PATH=. ./bedrock_server"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl unzip fail2ban

sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
URL=$(curl -s https://bedrock-version.minectl.ediri.online/binary/1.17.10.04)
curl -sLSf $URL > /tmp/bedrock-server.zip
unzip -o /tmp/bedrock-server.zip -d /minecraft
chmod +x /minecraft/bedrock_server
echo "eula=false" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	fabricCloudInitNoMonWant = `#cloud-config
users:
  - default
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-16-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=25565
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL="https://maven.fabricmc.net/net/fabricmc/fabric-installer/0.10.2/fabric-installer-0.10.2.jar"
  - mkdir /tmp/build
  - cd /tmp/build
  - curl -sLSf $URL > fabric-installer.jar
  - java -jar fabric-installer.jar server -downloadMinecraft
  - echo "serverJar=minecraft-server.jar" > /minecraft/fabric-server-launcher.properties
  - cp /tmp/build/fabric-server-launch.jar /minecraft/server.jar
  - cp /tmp/build/server.jar /minecraft/minecraft-server.jar
  - rm -rf /tmp/build
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	fabricBashNoMonWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=25565
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-16-jre-headless fail2ban

sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL="https://maven.fabricmc.net/net/fabricmc/fabric-installer/0.10.2/fabric-installer-0.10.2.jar"
mkdir /tmp/build
cd /tmp/build
curl -sLSf $URL > fabric-installer.jar
java -jar fabric-installer.jar server -downloadMinecraft -mcversion 1.17.1-138
echo "serverJar=minecraft-server.jar" > /minecraft/fabric-server-launcher.properties
cp /tmp/build/fabric-server-launch.jar /minecraft/server.jar
cp /tmp/build/server.jar /minecraft/minecraft-server.jar
rm -rf /tmp/build
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	nukkitBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=19132
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui --language eng

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-8-jre-headless fail2ban

sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL="https://ci.opencollab.dev/job/NukkitX/job/Nukkit/job/master/lastSuccessfulBuild/artifact/target/nukkit-1.0-SNAPSHOT.jar"
curl -sLSf $URL > /minecraft/server.jar
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	nukkitCloudInitWant = `#cloud-config
users:
  - default
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-8-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=19132
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui --language eng
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL="https://ci.opencollab.dev/job/NukkitX/job/Nukkit/job/master/lastSuccessfulBuild/artifact/target/nukkit-1.0-SNAPSHOT.jar"
  - curl -sLSf $URL > /minecraft/server.jar
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`

	powerNukkitBashWant = `#!/bin/bash
iptables -I INPUT -j ACCEPT
tee /tmp/server.properties <<EOF
server-port=19132
level-seed=stackitminecraftrocks
view-distance=10
enable-jmx-monitoring=false

broadcast-rcon-to-ops=true
rcon.port=2
enable-rcon=true
rcon.password=test
EOF
tee /etc/systemd/system/minecraft.service <<EOF
[Unit]
Description=Minecraft Server
Documentation=https://www.minecraft.net/en-us/download/server
DefaultDependencies=no
After=network.target

[Service]
WorkingDirectory=/minecraft
Type=simple
ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui --language eng

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
apt update
apt-get install -y apt-transport-https ca-certificates curl openjdk-8-jre-headless fail2ban

sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
service sshd restart

tee /etc/fail2ban/jail.local <<EOF
[sshd]
port = 0
enabled = true
maxretry = 0
bantime = 0
ignoreip = 
EOF

systemctl restart fail2ban
mkdir -p /minecraft
mkfs.ext4  /dev/sda
mount /dev/sda /minecraft
echo "/dev/sda /minecraft ext4 defaults,noatime,nofail 0 2" >> /etc/fstab
URL="https://github.com/PowerNukkit/PowerNukkit/releases/download/v1.5.1.0-PN/powernukkit-1.5.1.0-PN-shaded.jar"
curl -sLSf $URL > /minecraft/server.jar
echo "eula=true" > /minecraft/eula.txt
mv /tmp/server.properties /minecraft/server.properties
chmod a+rwx /minecraft
systemctl restart minecraft.service
systemctl enable minecraft.service`

	powerNukkitCloudInitWant = `#cloud-config
users:
  - default
package_update: true

packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - openjdk-8-jre-headless
  - fail2ban
fs_setup:
  - label: minecraft
    device: /dev/sda
    filesystem: xfs
    overwrite: false

mounts:
  - [/dev/sda, /minecraft]
# Enable ipv4 forwarding, required on CIS hardened machines
write_files:
  - path: /etc/sysctl.d/enabled_ipv4_forwarding.conf
    content: |
      net.ipv4.conf.all.forwarding=1
  - path: /tmp/server.properties
    content: |
       level-seed=stackitminecraftrocks
       view-distance=10
       enable-jmx-monitoring=false
       broadcast-rcon-to-ops=true
       rcon.port=2
       enable-rcon=true
       rcon.password=test
       server-port=19132
  - path: /etc/systemd/system/minecraft.service
    content: |
      [Unit]
      Description=Minecraft Server
      Documentation=https://www.minecraft.net/en-us/download/server
      DefaultDependencies=no
      After=network.target
      [Service]
      WorkingDirectory=/minecraft
      Type=simple
      ExecStart=/usr/bin/java -Xmx2G -Xms2G -jar server.jar nogui --language eng
      
      Restart=on-failure
      RestartSec=5
      [Install]
      WantedBy=multi-user.target
  - path: /etc/fail2ban/jail.local
    content: |
      [sshd]
      port = 0
      enabled = true
      maxretry = 0
      bantime = 0
      ignoreip = 

runcmd:
  - iptables -I INPUT -j ACCEPT
  - sed -i 's/#Port 22/Port 0/g' /etc/ssh/sshd_config
  - service sshd restart
  - systemctl restart fail2ban
  - URL="https://github.com/PowerNukkit/PowerNukkit/releases/download/v1.5.1.0-PN/powernukkit-1.5.1.0-PN-shaded.jar"
  - curl -sLSf $URL > /minecraft/server.jar
  - echo "eula=true" > /minecraft/eula.txt
  - mv /tmp/server.properties /minecraft/server.properties
  - chmod a+rwx /minecraft
  - systemctl restart minecraft.service
  - systemctl enable minecraft.service`
)

func TestCivoBedrockTemplate(t *testing.T) {
	t.Run("Test Template Bedrock for Civo bash", func(t *testing.T) {
		civo, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&bedrock, &CreateUpdateTemplateArgs{Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, bedrockBashWant, got)
	})
}

func TestCivoBedrockNoMonTemplate(t *testing.T) {
	t.Run("Test Template Bedrock for Civo bash", func(t *testing.T) {
		civo, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&bedrockNoMon, &CreateUpdateTemplateArgs{Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, bedrockBashNoMonWant, got)
	})
}

func TestCivoJavaTemplate(t *testing.T) {
	t.Run("Test Template Java for Civo bash", func(t *testing.T) {
		civo, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&java, &CreateUpdateTemplateArgs{Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, javaBashWant, got)
	})
}

func TestCloudInitBedrockTemplate(t *testing.T) {
	t.Run("Test Template Bedrock for Cloud-Init", func(t *testing.T) {
		civo, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&bedrock, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, bedrockCloudInitWant, got)
	})
}

func TestCloudInitJavaTemplate(t *testing.T) {
	t.Run("Test Template Java for Cloud-Init", func(t *testing.T) {
		civo, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&java, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, javaCloudInitWant, got)
	})
}

func TestBedrockBashMountTemplate(t *testing.T) {
	t.Run("Test Template Bedrock for bash", func(t *testing.T) {
		civo, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&bedrock, &CreateUpdateTemplateArgs{Mount: "sdc", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, bedrockBashMountWant, got)
	})
}

func TestJavaBashMountTemplate(t *testing.T) {
	t.Run("Test Template Bedrock for bash", func(t *testing.T) {
		civo, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := civo.GetTemplate(&java, &CreateUpdateTemplateArgs{Mount: "sdc", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, javaBashMountWant, got)
	})
}

func TestJavaAdditionalOptionsMountTemplate(t *testing.T) {
	t.Run("Test Template Java with additional options for bash", func(t *testing.T) {
		civo, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}

		got, err := civo.GetTemplate(&javaAdditionalOptions, &CreateUpdateTemplateArgs{Mount: "sdc", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, javaBashAdditionalOptionsMountWant, got)
	})
}

func TestCloudInitPaperMCTemplate(t *testing.T) {
	t.Run("Test Template PaperMC for Cloud-Init", func(t *testing.T) {
		paper, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&papermc, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, paperMCCloudInitWant, got)
	})
}

func TestBashPaperMCTemplate(t *testing.T) {
	t.Run("Test Template PaperMC for Bash", func(t *testing.T) {
		paper, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&papermc, &CreateUpdateTemplateArgs{Mount: "sdc", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, paperMCBashMountWant, got)
	})
}

func TestCloudInitCraftBukkitTemplate(t *testing.T) {
	t.Run("Test Template CraftBukkit for Cloud-Init", func(t *testing.T) {
		paper, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&craftbukkit, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, craftbukkitCloudInitWant, got)
	})
}

func TestBashCraftBukkitTemplate(t *testing.T) {
	t.Run("Test Template CraftBukkit for Bash", func(t *testing.T) {
		paper, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&craftbukkit, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, craftbukkitBashWant, got)
	})
}

func TestCloudInitFabricTemplate(t *testing.T) {
	t.Run("Test Template fabric for Cloud-Init", func(t *testing.T) {
		paper, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&fabric, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fabricCloudInitWant, got)
	})
}

func TestCloudInitFabricNoMonTemplate(t *testing.T) {
	t.Run("Test Template fabric for Cloud-Init", func(t *testing.T) {
		paper, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&fabricNoMon, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fabricCloudInitNoMonWant, got)
	})
}

func TestBashFabricTemplate(t *testing.T) {
	t.Run("Test Template fabric for Bash", func(t *testing.T) {
		paper, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&fabric, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fabricBashWant, got)
	})
}

func TestBashFabricNoMonTemplate(t *testing.T) {
	t.Run("Test Template fabric for Bash", func(t *testing.T) {
		paper, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&fabricNoMon, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fabricBashNoMonWant, got)
	})
}

func TestCloudInitForgeTemplate(t *testing.T) {
	t.Run("Test Template forge for Cloud-Init", func(t *testing.T) {
		paper, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&forge, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, forgeCloudInitWant, got)
	})
}

func TestBashForgeTemplate(t *testing.T) {
	t.Run("Test Template forge for Bash", func(t *testing.T) {
		paper, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&forge, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, forgeBashWant, got)
	})
}

func TestCloudInitSpigotTemplate(t *testing.T) {
	t.Run("Test Template spigot for Cloud-Init", func(t *testing.T) {
		paper, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&spigot, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, spigotCloudInitWant, got)
	})
}

func TestBashSpigotTemplate(t *testing.T) {
	t.Run("Test Template spigot for Bash", func(t *testing.T) {
		paper, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := paper.GetTemplate(&spigot, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, spigotBashWant, got)
	})
}

func TestBashNukkitTemplate(t *testing.T) {
	t.Run("Test Template nukkit for Bash", func(t *testing.T) {
		bash, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := bash.GetTemplate(&nukkit, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, nukkitBashWant, got)
	})
}

func TestCloudInitNukkitTemplate(t *testing.T) {
	t.Run("Test Template nukkit for Cloud Config", func(t *testing.T) {
		cloudConfig, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := cloudConfig.GetTemplate(&nukkit, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, nukkitCloudInitWant, got)
	})
}

func TestBashPowerNukkitTemplate(t *testing.T) {
	t.Run("Test Template powernukkit for Bash", func(t *testing.T) {
		bash, err := NewTemplateBash()
		if err != nil {
			t.Fatal(err)
		}
		got, err := bash.GetTemplate(&powerNukkit, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateBash})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, powerNukkitBashWant, got)
	})
}

func TestCloudInitPowerNukkitTemplate(t *testing.T) {
	t.Run("Test Template powernukkit for Cloud Config", func(t *testing.T) {
		cloudConfig, err := NewTemplateCloudConfig()
		if err != nil {
			t.Fatal(err)
		}
		got, err := cloudConfig.GetTemplate(&powerNukkit, &CreateUpdateTemplateArgs{Mount: "sda", Name: TemplateCloudConfig})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, powerNukkitCloudInitWant, got)
	})
}

var (
	fullFeatureJava = `apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft
spec:
  monitoring:
    enabled: true
  server:
    cloud: do
    region: xxx
    size: xxx
    ssh:
      port: 
      keyfolder: xxx
      fail2ban:
        bantime: 
        maxretry: 
    port: 25565
  minecraft:
    java:
      openjdk: 8
      xmx: 2G
      xms: 2G
      rcon:
        password: xxx
        port: 25575
        enabled: true
        broadcast: true
    edition: java
    version: xx
    eula: true
    properties: |
      level-seed=stackitminecraftrocks
      view-distance=10
      enable-jmx-monitoring=false`

	JavaWithoutRcon = `apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft
spec:
  monitoring:
    enabled: true
  server:
    cloud: do
    region: xxx
    size: xxx
    ssh:
      port: 
      keyfolder: xxx
      fail2ban:
        bantime: 
        maxretry: 
    port: 25565
  minecraft:
    java:
      openjdk: 8
      xmx: 2G
      xms: 2G
    edition: java
    version: xx
    eula: true
    properties: |
      level-seed=stackitminecraftrocks
      view-distance=10
      enable-jmx-monitoring=false`

	plainJavaNoThrill = `apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft
spec:
  server:
    cloud: do
    region: xxx
    size: xxx
    ssh:
      port: 
      keyfolder: xxx
      fail2ban:
        bantime: 
        maxretry: 
    port: 25565
  minecraft:
    java:
      openjdk: 8
      xmx: 2G
      xms: 2G
    edition: java
    version: xx
    eula: true
    properties: |
      level-seed=stackitminecraftrocks
      view-distance=10
      enable-jmx-monitoring=false`

	bedrockConfig = `apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft
spec:
  server:
    cloud: do
    region: xxx
    size: xxx
    ssh:
      port: 
      keyfolder: xxx
      fail2ban:
        bantime: 
        maxretry: 
    port: 19132
  minecraft:
    edition: bedrock
    version: xx
    eula: true
    properties: |
      level-seed=stackitminecraftrocks
      view-distance=10
      enable-jmx-monitoring=false`

	nukkitConfig = `apiVersion: minectl.ediri.io/v1alpha1
kind: MinecraftServer
metadata:
  name: minecraft
spec:
  monitoring:
    enabled: true
  server:
    cloud: do
    region: xxx
    size: xxx
    ssh:
      port: 
      keyfolder: xxx
      fail2ban:
        bantime: 
        maxretry: 
    port: 19132
  minecraft:
    java:
      openjdk: 8
      xmx: 2G
      xms: 2G
      rcon:
        password: xxx
        port: 25575
        enabled: true
        broadcast: true
    edition: nukkit
    version: xx
    eula: true
    properties: |
      level-seed=stackitminecraftrocks
      view-distance=10
      enable-jmx-monitoring=false`
)

func crateWizardMock() model.Wizard {
	return model.Wizard{
		Name:       "minecraft",
		Provider:   "DigitalOcean",
		Plan:       "xxx",
		Region:     "xxx",
		SSH:        "xxx",
		Features:   []string{"Monitoring", "RCON"},
		Java:       "8",
		Heap:       "2G",
		RconPw:     "xxx",
		Edition:    "java",
		Version:    "xx",
		Properties: "level-seed=stackitminecraftrocks\nview-distance=10\nenable-jmx-monitoring=false",
	}
}

func TestConfigTemplate(t *testing.T) {
	t.Run("Test Config", func(t *testing.T) {
		got, err := NewTemplateConfig(crateWizardMock())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fullFeatureJava, got)
	})
}

func TestConfigMonTemplate(t *testing.T) {
	t.Run("Test Config", func(t *testing.T) {
		wizard := crateWizardMock()
		wizard.Features = []string{"Monitoring"}
		got, err := NewTemplateConfig(wizard)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, JavaWithoutRcon, got)
	})
}

func TestConfigNoMonRconTemplate(t *testing.T) {
	t.Run("Test Config", func(t *testing.T) {
		wizard := crateWizardMock()
		wizard.Features = []string{}
		got, err := NewTemplateConfig(wizard)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, plainJavaNoThrill, got)
	})
}

func TestConfigBedrockTemplate(t *testing.T) {
	t.Run("Test Config", func(t *testing.T) {
		wizard := crateWizardMock()
		wizard.Edition = "bedrock"
		got, err := NewTemplateConfig(wizard)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, bedrockConfig, got)
	})
}

func TestConfigNukkitTemplate(t *testing.T) {
	t.Run("Test Config", func(t *testing.T) {
		wizard := crateWizardMock()
		wizard.Edition = "nukkit"
		got, err := NewTemplateConfig(wizard)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, nukkitConfig, got)
	})
}
