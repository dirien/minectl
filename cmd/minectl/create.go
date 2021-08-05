package minectl

import (
	"fmt"
	"log"
	"os"

	"github.com/minectl/pkg/common"
	"github.com/minectl/pkg/provisioner"
	"github.com/olekukonko/tablewriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	minectlCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("filename", "f", "", "Contains the configuration for minectl")
	createCmd.Flags().BoolP("wait", "w", true, "Wait for Minecraft Server is started")
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
	if len(filename) == 0 {
		return errors.New("Please provide a valid MinecraftServer manifest file")
	}
	if err != nil {
		return errors.Wrap(err, "Please provide a valid MinecraftServer manifest file")
	}
	p, err := provisioner.NewProvisioner(filename)
	if err != nil {
		log.Fatal(err)
	}
	wait := true
	if cmd.Flags().Changed("wait") {
		wait, _ = cmd.Flags().GetBool("wait")
	}
	res, err := p.CreateServer(wait)
	if err != nil {
		log.Fatal(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "NAME", "REGION", "TAGS", "IP"})

	table.Append([]string{res.ID, res.Name, res.Region, res.Tags, res.PublicIP})

	table.SetBorder(false)
	fmt.Println("")
	table.Render()

	common.PrintMixedGreen("\n🪓 To delete the server type:\n\n %s", fmt.Sprintf("minectl delete -f %s --id %s\n", filename, res.ID))
	common.PrintMixedGreen("\n🆙 To update the server type:\n\n %s", fmt.Sprintf("minectl update -f %s --id %s\n", filename, res.ID))
	fmt.Println("🚧 Beta features:")
	common.PrintMixedGreen("⤴️ To upload a plugin type:\n\n %s",
		fmt.Sprintf("minectl plugins -f %s --id %s --plugin <folder>/x.jar --destination /minecraft/plugins\n", filename, res.ID))
	common.PrintMixedGreen("\n🔌 Connected to RCON type:\n\n %s", fmt.Sprintf("minectl rcon -f %s --id %s\n", filename, res.ID))
	return err
}
