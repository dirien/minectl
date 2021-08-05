package model

// Spec
type Spec struct {
	Monitoring Monitoring `json:"monitoring"`
	Server     Server     `yaml:"server"`
	Minecraft  Minecraft  `yaml:"minecraft"`
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

// MinecraftServer
type MinecraftServer struct {
	Spec       Spec     `yaml:"spec"`
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}

func (m *MinecraftServer) GetProperties() string {
	return m.Spec.Minecraft.Properties
}

func (m *MinecraftServer) GetName() string {
	return m.Metadata.Name
}

func (m *MinecraftServer) GetCloud() string {
	return m.Spec.Server.Cloud
}

func (m *MinecraftServer) GetSSH() string {
	return m.Spec.Server.Ssh
}

func (m *MinecraftServer) GetRegion() string {
	return m.Spec.Server.Region
}

func (m *MinecraftServer) GetSize() string {
	return m.Spec.Server.Size
}

func (m *MinecraftServer) GetEdition() string {
	return m.Spec.Minecraft.Edition
}

func (m *MinecraftServer) GetVolumeSize() int {
	return m.Spec.Server.VolumeSize
}

func (m *MinecraftServer) GetVersion() string {
	return m.Spec.Minecraft.Version
}

func (m *MinecraftServer) GetPort() int {
	return m.Spec.Server.Port
}

func (m *MinecraftServer) GetJDKVersion() int {
	return m.Spec.Minecraft.Java.OpenJDK
}

func (m *MinecraftServer) GetRCONPort() int {
	return m.Spec.Minecraft.Java.Rcon.Port
}
func (m *MinecraftServer) GetRCONPassword() string {
	return m.Spec.Minecraft.Java.Rcon.Password
}
