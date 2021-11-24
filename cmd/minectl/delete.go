package minectl

import (
	"github.com/minectl/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringP("filename", "f", "", "that contains the configuration for minectl")
	deleteCmd.Flags().String("id", "", "contains the server id")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Minecraft Server.",
	Example: `mincetl delete  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx
	`,
	RunE:          RunFunc(runDelete),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runDelete(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return errors.Wrap(err, "failed to get 'filename' value.")
	}
	if len(filename) == 0 {
		return errors.New("Please provide a valid MinecraftResource manifest file via -f|--filename flag")
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("Please provide a valid id")
	}
	newProvisioner, err := provisioner.NewProvisioner(&provisioner.MinectlProvisionerOpts{
		ManifestPath: filename,
		ID:           id,
	}, minectlLog)
	if err != nil {
		return err
	}
	err = newProvisioner.DeleteServer()
	if err != nil {
		return err
	}
	return nil
}
