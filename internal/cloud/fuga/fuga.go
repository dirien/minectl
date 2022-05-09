package fuga

import (
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/cloud/openstack"
)

type Fuga struct {
	openshift *openstack.OpenStack
}

const imageName = "Ubuntu 20.04 LTS"

func NewFuga() (*Fuga, error) {
	client, err := openstack.NewOpenStack(imageName)
	if err != nil {
		return nil, err
	}
	return &Fuga{
		openshift: client,
	}, nil
}

func (f *Fuga) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	return f.openshift.CreateServer(args)
}

func (f *Fuga) DeleteServer(id string, args automation.ServerArgs) error {
	return f.openshift.DeleteServer(id, args)
}

func (f *Fuga) ListServer() ([]automation.ResourceResults, error) {
	return f.openshift.ListServer()
}

func (f *Fuga) UpdateServer(id string, args automation.ServerArgs) error {
	return f.openshift.UpdateServer(id, args)
}

func (f *Fuga) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	return f.openshift.UploadPlugin(id, args, plugin, destination)
}

func (f *Fuga) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	return f.openshift.GetServer(id, args)
}
