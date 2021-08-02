package update

import (
	"fmt"
	"strings"

	minctlTemplate "github.com/minectl/pgk/template"

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
	tmpl := minctlTemplate.GetUpdateTemplate()
	var update string
	var err error

	if args.GetEdition() == "java" {
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateJavaBinary)
	} else if args.GetEdition() == "bedrock" {
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateBedrockBinary)
	} else if args.GetEdition() == "craftbukkit" || args.GetEdition() == "spigot" {
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplatesSigotbukkitBinary)
	} else if args.GetEdition() == "fabric" {
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplatesFabricBinary)
		update = fmt.Sprintf("\nrm -rf /minecraft/minecraft-server.jar%s", update)
	} else if args.GetEdition() == "forge" {
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplatesForgeBinary)
	} else if args.GetEdition() == "papermc" {
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplatesPaperMCBinary)
	}
	if args.GetEdition() != "bedrock" {
		update = fmt.Sprintf("%s\napt-get install -y openjdk-%d-jre-headless\n", update, args.GetJDKVersion())
	}
	if err != nil {
		return err
	}

	cmd := `
cd /minecraft
sudo systemctl stop minecraft.service` + update +
		`ls -la
sudo systemctl start minecraft.service
	`

	_, err = r.ExecuteCommand(strings.TrimSpace(cmd))
	if err != nil {
		return err
	}
	return nil
}

func (r *RemoteServer) TransferFile(src, dstPath string) error {
	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		return err
	}

	client, err := goph.NewConn(&goph.Config{
		User:     r.user,
		Addr:     r.ip,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return err
	}

	defer client.Close()
	err = client.Upload(src, dstPath)
	if err != nil {
		return err
	}
	return nil
}

func (r *RemoteServer) ExecuteCommand(cmd string) (string, error) {
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
