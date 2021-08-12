package minectl

import (
	"github.com/minectl/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	minectlCmd.AddCommand(rconCmd)
	rconCmd.Flags().StringP("filename", "f", "", "Contains the configuration for minectl")
	rconCmd.Flags().String("id", "", "contains the server id")
}

var rconCmd = &cobra.Command{
	Use:   "rcon",
	Short: "RCON client to your Minecraft server.",
	Example: `mincetl rcon  \
    --filename server-do.yaml \
    --id xxxx`,
	RunE:          RunFunc(runRCON),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runRCON(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if len(filename) == 0 {
		return errors.New("Please provide a valid MinecraftResource manifest file")
	}
	if err != nil {
		return errors.Wrap(err, "Please provide a valid MinecraftResource manifest file")
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("Please provide a valid id")
	}
	p, err := provisioner.NewProvisioner(filename, id)
	if err != nil {
		return err
	}
	err = p.DoRCON()
	if err != nil {
		return err
	}
	return nil
}
