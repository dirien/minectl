package model

// Spec
type Spec struct {
	Monitoring Monitoring `json:"monitoring"`
	Server     Server     `yaml:"server"`
	Minecraft  Minecraft  `yaml:"minecraft"`
	Proxy      Proxy      `yaml:"proxy"`
}

// Proxy
type Proxy struct {
	Java    Java   `yaml:"java"`
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
}

// Monitoring
type Monitoring struct {
	Enabled bool `json:"enabled"`
}

// Server
type Server struct {
	Size       string `yaml:"size"`
	VolumeSize int    `yaml:"volumeSize"`
	Ssh        string `yaml:"ssh"`
	Cloud      string `yaml:"cloud"`
	Region     string `yaml:"region"`
	Port       int    `yaml:"port"`
}

// Minecraft
type Minecraft struct {
	Java       Java   `yaml:"java"`
	Properties string `yaml:"properties"`
	Edition    string `yaml:"edition"`
	Version    string `yaml:"version"`
	Eula       bool   `yaml:"eula"`
}

// Java
type Java struct {
	Xmx     string `yaml:"xmx"`
	Xms     string `yaml:"xms"`
	OpenJDK int    `yaml:"openjdk"`
	Rcon    Rcon   `yaml:"rcon"`
}

// Rcon
type Rcon struct {
	Password  string `yaml:"password"`
	Enabled   bool   `yaml:"enabled"`
	Port      int    `yaml:"port"`
	Broadcast bool   `yaml:"broadcast"`
}

// Metadata
type Metadata struct {
	Name string `yaml:"name"`
}

// MinecraftResource
type MinecraftResource struct {
	Spec       Spec     `yaml:"spec"`
	ApiVersion string   `yaml:"apiVersion"`
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

func (m *MinecraftResource) GetSSH() string {
	return m.Spec.Server.Ssh
}

func (m *MinecraftResource) GetRegion() string {
	return m.Spec.Server.Region
}

func (m *MinecraftResource) GetSize() string {
	return m.Spec.Server.Size
}

func (m *MinecraftResource) GetEdition() string {
	if (m.Spec.Proxy != Proxy{}) {
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
	if (m.Spec.Proxy != Proxy{}) {
		return m.Spec.Proxy.Java.Rcon.Port
	}
	return m.Spec.Minecraft.Java.Rcon.Port
}
func (m *MinecraftResource) GetRCONPassword() string {
	if (m.Spec.Proxy != Proxy{}) {
		return m.Spec.Proxy.Java.Rcon.Password
	}
	return m.Spec.Minecraft.Java.Rcon.Password
}

func (m *MinecraftResource) IsProxyServer() bool {
	return m.Spec.Proxy != Proxy{}
}
