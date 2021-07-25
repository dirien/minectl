package update

import (
	"fmt"
	"strings"

	"github.com/melbahja/goph"
	"github.com/minectl/pgk/model"
	"golang.org/x/crypto/ssh"
)

type ServerOperations interface {
	UpdateServer(*model.MinecraftServer) error
}

type RemoteServer struct {
	ip            string
	privateSSHKey string
	user          string
}

func NewRemoteServer(privateKey, ip, user string) *RemoteServer {
	ssh := &RemoteServer{
		ip:            ip,
		privateSSHKey: privateKey,
		user:          user,
	}
	return ssh
}

func (r *RemoteServer) UpdateServer(args *model.MinecraftServer) error {
	var url string
	if args.GetEdition() == "java" {
		url = fmt.Sprintf("URL=$(curl -s https://java-version.minectl.ediri.online/binary/%s)", args.GetVersion())
		url = fmt.Sprintf("%s\nsudo rm -rf server.jar\nsudo curl -sLSf $URL > /tmp/server.jar\nsudo mv /tmp/server.jar /minecraft/", url)
	} else if args.GetEdition() == "bedrock" {
		url = fmt.Sprintf("URL=$(curl -s https://bedrock-version.minectl.ediri.online/binary/%s)", args.GetVersion())
		url = fmt.Sprintf("%s\nsudo rm -rf bedrock_server\nsudo curl -sLSf $URL > /tmp/bedrock-server.zip\nsudo unzip -o /tmp/bedrock-server.zip -d /minecraft\nsudo chmod +x /minecraft/bedrock_server\n", url)
	}

	cmd := `
	cd /minecraft
    sudo systemctl stop minecraft.service
	` + url + `
	ls -la
	sudo systemctl start minecraft.service
	`
	_, err := r.executeCommand(strings.TrimSpace(cmd))
	if err != nil {
		return err
	}
	return nil
}

func (r *RemoteServer) executeCommand(cmd string) (string, error) {
	//fmt.Printf("Running remote command %s\n", color.GreenString(cmd))
	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		return "", err
	}
	client, err := goph.NewConn(&goph.Config{
		User:     r.user,
		Addr:     r.ip,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return "", err
	}

	defer client.Close()
	out, err := client.Run(cmd)
	return string(out), err
}
