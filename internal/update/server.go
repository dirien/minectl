package update

import (
	"fmt"
	"strings"

	"github.com/melbahja/goph"
	minctlTemplate "github.com/minectl/internal/template"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/minectl/internal/model"
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
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateJavaBinary})
	case "bedrock":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateBedrockBinary})
	case "craftbukkit":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateSpigotBukkitBinary})
	case "spigot":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateSpigotBukkitBinary})
	case "fabric":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateFabricBinary})
		update = fmt.Sprintf("\nrm -rf /minecraft/minecraft-server.jar%s", update)
	case "forge":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateForgeBinary})
	case "papermc":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplatePaperMCBinary})
	case "bungeecord":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateBungeeCordBinary})
	case "waterfall":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateWaterfallBinary})
	case "nukkit":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateNukkitBinary})
	case "powernukkit":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplatePowerNukkitBinary})
	case "velocity":
		update, err = tmpl.DoUpdate(args, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.TemplateVelocityBinary})
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
		Callback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
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
	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		return "", err
	}
	client, err := goph.NewConn(&goph.Config{
		User:     r.user,
		Addr:     r.ip,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
	})
	if err != nil {
		return "", err
	}

	defer client.Close()
	out, err := client.Run(cmd)
	return string(out), err
}
