package minectl

import (
	"github.com/minectl/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().StringP("filename", "f", "", "Contains the configuration for minectl")
	updateCmd.Flags().String("id", "", "contains the server id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Minecraft Server.",
	Example: `mincetl update  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx`,
	RunE:          RunFunc(runUpdate),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runUpdate(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if len(filename) == 0 {
		return errors.New("Please provide a valid MinecraftResource manifest file")
	}
	if err != nil {
		return errors.Wrap(err, "Please provide a valid MinecraftResource manifest file via -f|--filename flag")
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("Please provide a valid id")
	}
	p, err := provisioner.NewProvisioner(&provisioner.MinectlProvisionerOpts{
		ManifestPath: filename,
		ID:           id,
	}, minectlLog)
	if err != nil {
		return err
	}
	err = p.UpdateServer()
	if err != nil {
		return err
	}
	return err
}
