package minectl

import (
	"github.com/dirien/minectl/internal/provisioner"
	"github.com/dirien/minectl/internal/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringP("filename", "f", "", "Location of the manifest file")
	_ = deleteCmd.Flags().SetAnnotation("filename", cobra.BashCompFilenameExt, []string{"yaml"})
	deleteCmd.Flags().String("id", "", "Contains the server id")
	deleteCmd.Flags().BoolP("yes", "y", false, "Automatically delete the server")
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
		return errors.Wrap(err, "failed to get 'filename' value")
	}
	if filename == "" {
		return errors.New("Please provide a valid manifest file via -f|--filename flag")
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	if id == "" {
		return errors.New("Please provide a valid id")
	}
	newProvisioner, err := provisioner.NewProvisioner(&provisioner.MinectlProvisionerOpts{
		ManifestPath: filename,
		ID:           id,
	}, minectlUI)
	if err != nil {
		return err
	}

	yes := cmd.Flag("yes").Changed
	if yes {
		err = newProvisioner.DeleteServer()
		if err != nil {
			return err
		}
	} else {
		confirmed, err := ui.Confirm("Do you want to delete the Minecraft server?")
		if err != nil {
			return err
		}
		if confirmed {
			err = newProvisioner.DeleteServer()
			if err != nil {
				return err
			}
		} else {
			minectlUI.Warn("Delete canceled.")
		}
	}

	return nil
}
