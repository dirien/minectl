package update

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	minctlTemplate "github.com/minectl/pkg/template"

	"github.com/melbahja/goph"
	"github.com/minectl/pkg/model"
	"golang.org/x/crypto/ssh"
)

type ServerOperations interface {
	UpdateServer(*model.MinecraftResource) error
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

func (r *RemoteServer) UpdateServer(args *model.MinecraftResource) error {
	tmpl := minctlTemplate.GetUpdateTemplate()
	var update string
	var err error

	switch args.GetEdition() {
	case "java":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateJavaBinary)
	case "bedrock":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateBedrockBinary)
	case "craftbukkit":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateSpigotBukkitBinary)
	case "spigot":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateSpigotBukkitBinary)
	case "fabric":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateFabricBinary)
		update = fmt.Sprintf("\nrm -rf /minecraft/minecraft-server.jar%s", update)
	case "forge":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateForgeBinary)
	case "papermc":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplatePaperMCBinary)
	case "bungeecord":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateBungeeCordBinary)
	case "waterfall":
		update, err = tmpl.DoUpdate(args, minctlTemplate.TemplateWaterfallBinary)
	}

	if args.GetEdition() != "bedrock" {
		update = fmt.Sprintf("%s\napt-get install -y openjdk-%d-jre-headless\n", update, args.GetJDKVersion())
	}
	if err != nil {
		return err
	}

	cmd := `
cd /minecraft
sudo systemctl stop minecraft.service
sudo bash -c '` + update + `'
ls -la
sudo systemctl start minecraft.service
	`
	zap.S().Infof("server updated cmd %s", cmd)
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
