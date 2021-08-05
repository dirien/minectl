package minectl

import (
	"fmt"
	"os"

	"github.com/minectl/pkg/provisioner"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {

	minectlCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("provider", "p", "", "The cloud provider - civo|scaleway|do|hetzner|linode|ovh|equinix|gce")
	listCmd.Flags().StringP("region", "r", "", "The region (gce: zone) for your cloud provider - civo|gce")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Minecraft Server.",
	Example: `mincetl list  \
    --provider civo \
    --region LON1`,
	RunE:          runList,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runList(cmd *cobra.Command, _ []string) error {

	provider, err := cmd.Flags().GetString("provider")
	if err != nil {
		return errors.Wrap(err, "failed to get 'provider' value.")
	}
	if len(provider) == 0 {
		return errors.New("Please provide a valid 'provider' value")
	}
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return errors.Wrap(err, "failed to get 'region' value.")
	}

	newProvisioner, err := provisioner.ListProvisioner(provider, region)
	if err != nil {
		return err
	}
	servers, err := newProvisioner.ListServer()
	if err != nil {
		return err
	}

	if len(servers) > 0 {
		fmt.Println("")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "NAME", "REGION", "TAGS", "IP"})

		for _, server := range servers {
			table.Append([]string{server.ID, server.Name, server.Region, server.Tags, server.PublicIP})
		}
		table.SetBorder(false)
		table.Render()
	} else {
		fmt.Println("🤷 No server found")
	}
	return nil
}
