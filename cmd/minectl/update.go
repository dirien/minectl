package minectl

import (
	"log"

	"github.com/minectl/pgk/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	minectlCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("filename", "f", "", "Contains the configuration for minectl")
	updateCmd.Flags().String("id", "", "contains the server id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Minecraft Server.",
	Example: `mincetl update  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx`,
	RunE:          runUpdate,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runUpdate(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if filename == "" {
		return errors.New("Please provide a valid MinecraftServer manifest file")
	}
	if err != nil {
		return errors.Wrap(err, "Please provide a valid MinecraftServer manifest file")
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	p, err := provisioner.NewProvisioner(filename, id)
	if err != nil {
		log.Fatal(err)
	}
	err = p.UpdateServer()
	if err != nil {
		return err
	}
	return err
}
