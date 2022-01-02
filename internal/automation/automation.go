package automation

import "github.com/minectl/internal/model"

type Automation interface {
	CreateServer(args ServerArgs) (*ResourceResults, error)
	DeleteServer(id string, args ServerArgs) error
	ListServer() ([]ResourceResults, error)
	UpdateServer(id string, args ServerArgs) error
	UploadPlugin(id string, args ServerArgs, plugin, destination string) error
	GetServer(id string, args ServerArgs) (*ResourceResults, error)
}

type Rcon struct {
	Password  string
	Enabled   bool
	Port      int
	Broadcast bool
}

type ServerArgs struct {
	ID                string
	MinecraftResource *model.MinecraftResource
}

type ResourceResults struct {
	ID       string
	Name     string
	Region   string
	PublicIP string
	Tags     string
}
