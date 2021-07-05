package minectl

import (
	"github.com/spf13/cobra"
)

func init() {

	minectlCmd.AddCommand(listCmd)
	listCmd.Flags().String("region", "", "that contains the configuration for minectl")

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Minecraft Server.",
	Example: `mincetl list  \
    --region LON1`,
	RunE:          runList,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runList(cmd *cobra.Command, _ []string) error {

	//_, err := cmd.Flags().GetString("filename")
	/*if err != nil {
		return errors.Wrap(err, "failed to get 'filename' value.")
	}
	do, _ := provisioner.NewProvisioner(filename)
	res, err := do.UpdateServer()
	common.PrintMixedGreen("Minecraft Server IP: %s\n", res.PublicIP)*/
	return nil
}
