package automation

import "github.com/minectl/pgk/model"

type Automation interface {
	CreateServer(args ServerArgs) (*RessourceResults, error)
	DeleteServer(id string, args ServerArgs) error
	ListServer() ([]RessourceResults, error)
	UpdateServer(id string, args ServerArgs) error
	UploadPlugin(id string, args ServerArgs, plugin, destination string) error
	GetServer(id string, args ServerArgs) (*RessourceResults, error)
}

type Rcon struct {
	Password  string
	Enabled   bool
	Port      int
	Broadcast bool
}

type ServerArgs struct {
	ID              string
	MinecraftServer *model.MinecraftServer
}

type RessourceResults struct {
	ID       string
	Name     string
	Region   string
	PublicIP string
	Tags     string
}
