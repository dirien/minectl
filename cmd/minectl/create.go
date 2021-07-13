package minectl

import (
	"fmt"
	"github.com/minectl/pgk/common"
	"github.com/minectl/pgk/provisioner"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"

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
	p, err := provisioner.NewProvisioner(filename)
	if err != nil {
		log.Fatal(err)
	}
	res, err := p.CreateServer()
	if err != nil {
		log.Fatal(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "NAME", "REGION", "TAGS", "IP"})

	table.Append([]string{res.ID, res.Name, res.Region, res.Tags, res.PublicIP})

	table.SetBorder(false)
	fmt.Println("")
	table.Render()

	common.PrintMixedGreen("\nTo delete the server type:\n\n %s", fmt.Sprintf("minectl delete -f %s --id %s\n", filename, res.ID))
	return err
}
