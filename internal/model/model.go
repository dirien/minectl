package model

import "reflect"

// Wizard represents a wizard configuration.
type Wizard struct {
	Name       string
	Provider   string
	Plan       string
	Region     string
	SSH        string `survey:"ssh"`
	SSHPort    string `survey:"ssh_port"`
	BanTime    string `survey:"fail2ban_bantime"`
	MaxRetry   string `survey:"fail2ban_maxretry"`
	Features   []string
	Java       string
	Heap       string
	RconPw     string `survey:"rconpw"`
	Edition    string
	Version    string
	Properties string
}

// Spec represents a specification configuration.
type Spec struct {
	Monitoring Monitoring `json:"monitoring"`
	Server     Server     `yaml:"server"`
	Minecraft  Minecraft  `yaml:"minecraft"`
	Proxy      Proxy      `yaml:"proxy"`
}

// Proxy represents a proxy configuration.
type Proxy struct {
	Java    Java   `yaml:"java"`
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
}

// Monitoring represents a monitoring configuration.
type Monitoring struct {
	Enabled bool `json:"enabled"`
}

// Server represents a server configuration.
type Server struct {
	Size       string `yaml:"size"`
	SSH        SSH    `yaml:"ssh"`
	Cloud      string `yaml:"cloud"`
	Region     string `yaml:"region"`
	Port       int    `yaml:"port"`
	VolumeSize int    `yaml:"volumeSize"`
	Spot       bool   `yaml:"spot"`
	Arm        bool   `yaml:"arm"`
}

// SSH represents a SSH configuration.
type SSH struct {
	Port      int      `yaml:"port"`
	KeyFolder string   `yaml:"keyFolder"`
	Fail2ban  Fail2ban `yaml:"fail2ban"`
}

// Fail2ban represents a fail2ban configuration.
type Fail2ban struct {
	Bantime  int    `yaml:"bantime"`
	Maxretry int    `yaml:"maxretry"`
	Ignoreip string `yaml:"ignoreip"`
}

// Minecraft represents a minecraft configuration.
type Minecraft struct {
	Java       Java   `yaml:"java"`
	Properties string `yaml:"properties"`
	Edition    string `yaml:"edition"`
	Version    string `yaml:"version"`
	Eula       bool   `yaml:"eula"`
}

// Java represents a java configuration.
type Java struct {
	Xmx     string   `yaml:"xmx"`
	Xms     string   `yaml:"xms"`
	Options []string `yaml:"options"`
	OpenJDK int      `yaml:"openjdk"`
	Rcon    Rcon     `yaml:"rcon"`
}

// Rcon represents a rcon configuration.
type Rcon struct {
	Password  string `yaml:"password"`
	Enabled   bool   `yaml:"enabled"`
	Port      int    `yaml:"port"`
	Broadcast bool   `yaml:"broadcast"`
}

// Metadata represents a metadata configuration.
type Metadata struct {
	Name string `yaml:"name"`
}

// MinecraftResource represents a minecraft resource.
type MinecraftResource struct {
	Spec       Spec     `yaml:"spec"`
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}

func (m *MinecraftResource) GetProperties() string {
	return m.Spec.Minecraft.Properties
}

func (m *MinecraftResource) GetName() string {
	return m.Metadata.Name
}

func (m *MinecraftResource) GetCloud() string {
	return m.Spec.Server.Cloud
}

func (m *MinecraftResource) GetSSHPort() int {
	return m.Spec.Server.SSH.Port
}

func (m *MinecraftResource) GetSSHKeyFolder() string {
	return m.Spec.Server.SSH.KeyFolder
}

func (m *MinecraftResource) GetFail2Ban() Fail2ban {
	return m.Spec.Server.SSH.Fail2ban
}

func (m *MinecraftResource) GetRegion() string {
	return m.Spec.Server.Region
}

func (m *MinecraftResource) GetSize() string {
	return m.Spec.Server.Size
}

func (m *MinecraftResource) GetEdition() string {
	if m.IsProxyServer() {
		return m.Spec.Proxy.Type
	}
	return m.Spec.Minecraft.Edition
}

func (m *MinecraftResource) GetVolumeSize() int {
	return m.Spec.Server.VolumeSize
}

func (m *MinecraftResource) GetVersion() string {
	return m.Spec.Minecraft.Version
}

func (m *MinecraftResource) GetPort() int {
	return m.Spec.Server.Port
}

func (m *MinecraftResource) GetJDKVersion() int {
	return m.Spec.Minecraft.Java.OpenJDK
}

func (m *MinecraftResource) GetRCONPort() int {
	if m.IsProxyServer() {
		return m.Spec.Proxy.Java.Rcon.Port
	}
	return m.Spec.Minecraft.Java.Rcon.Port
}

func (m *MinecraftResource) HasRCON() bool {
	if m.IsProxyServer() {
		return m.Spec.Proxy.Java.Rcon.Enabled
	}
	return m.Spec.Minecraft.Java.Rcon.Enabled
}

func (m *MinecraftResource) HasMonitoring() bool {
	return m.Spec.Monitoring.Enabled
}

func (m *MinecraftResource) GetRCONPassword() string {
	if m.IsProxyServer() {
		return m.Spec.Proxy.Java.Rcon.Password
	}
	return m.Spec.Minecraft.Java.Rcon.Password
}

func (m *MinecraftResource) IsProxyServer() bool {
	return reflect.DeepEqual(m.Spec.Minecraft, Minecraft{})
}

func (m *MinecraftResource) IsSpot() bool {
	return m.Spec.Server.Spot
}

func (m *MinecraftResource) IsArm() bool {
	return m.Spec.Server.Arm
}
