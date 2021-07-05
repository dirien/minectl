package automation

type Automation interface {
	CreateServer(args ServerArgs) (*RessourceResults, error)
	DeleteServer(id string, args ServerArgs) error
	ListServer(args ServerArgs) (*[]RessourceResults, error)
	UpdateServer(args ServerArgs) (*RessourceResults, error)
}

type ServerArgs struct {
	ID         string
	StackName  string
	Size       string
	Region     string
	SSH        string
	VolumeSize int
	Properties string
}

type RessourceResults struct {
	ID       string
	PublicIP string
}
