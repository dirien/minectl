package minectl

import (
	"fmt"

	"github.com/dirien/minectl/internal/provisioner"
	"github.com/dirien/minectl/internal/ui"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.Flags().StringP("filename", "f", "", "Location of the manifest file")
	_ = createCmd.Flags().SetAnnotation("filename", cobra.BashCompFilenameExt, []string{"yaml"})
	createCmd.Flags().BoolP("wait", "w", true, "Wait for Minecraft Server is started")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Minecraft Server.",
	Example: `minectl create  \
    --filename server-do.yaml`,
	RunE:          RunFunc(runCreate),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runCreate(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if filename == "" {
		return errors.New("Please provide a valid manifest file via -f|--filename flag")
	}
	if err != nil {
		return errors.Wrap(err, "Please provide a valid manifest file")
	}
	p, err := provisioner.NewProvisioner(&provisioner.MinectlProvisionerOpts{
		ManifestPath: filename,
	}, minectlUI)
	if err != nil {
		return err
	}
	wait := true
	if cmd.Flags().Changed("wait") {
		wait, _ = cmd.Flags().GetBool("wait")
	}
	res, err := p.CreateServer(wait)
	if err != nil {
		return err
	}
	if !headless {
		table := ui.NewTable(minectlUI, "ID", "NAME", "REGION", "TAGS", "IP")
		table.Append([]string{res.ID, res.Name, res.Region, res.Tags, res.PublicIP})

		fmt.Println("")
		table.Render()

		minectlUI.Info(fmt.Sprintf("To delete the server:\n\n  minectl delete -f %s --id %s", filename, res.ID))
		minectlUI.Info(fmt.Sprintf("To update the server:\n\n  minectl update -f %s --id %s", filename, res.ID))
		minectlUI.Info(fmt.Sprintf("To connect via RCON:\n\n  minectl rcon -f %s --id %s", filename, res.ID))
		minectlUI.Warn("Beta features:")
		minectlUI.Info(fmt.Sprintf("To upload a plugin:\n\n  minectl plugins -f %s --id %s --plugin <folder>/x.jar --destination /minecraft/plugins", filename, res.ID))
	}
	return nil
}
