package minectl

import (
	"fmt"
	"github.com/minectl/pgk/common"
	"github.com/minectl/pgk/provisioner"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	minectlCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("filename", "f", "", "Contains the configuration for minectl")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Minecraft Server.",
	Example: `mincetl create  \
    --filename server-do.yaml`,
	RunE:          runCreate,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runCreate(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if filename == "" {
		return errors.New("Please provide a valid MinecraftServer manifest file")
	}
	if err != nil {
		return errors.Wrap(err, "Please provide a valid MinecraftServer manifest file")
	}
	p, err := provisioner.NewCreateServerProvisioner(filename)
	if err != nil {
		log.Fatal(err)
	}
	res, err := p.CreateServer()
	if err != nil {
		log.Fatal(err)
	}
	common.PrintMixedGreen("Minecraft Server IP: %s\n", res.PublicIP)
	common.PrintMixedGreen("Minecraft Server ID: %s\n", res.ID)

	common.PrintMixedGreen("\nTo delete the server type:\n\n %s", fmt.Sprintf("minectl delete -f %s --id %s\n", filename, res.ID))
	return err
}
