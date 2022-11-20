package minectl

import (
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().StringP("filename", "f", "", "Location of the manifest file")
	updateCmd.Flags().StringP("ssh-key", "k", "", "specify a specific path for the SSH key")
	updateCmd.Flags().SetAnnotation("filename", cobra.BashCompFilenameExt, []string{"yaml"}) //nolint:errcheck
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
	p, err := createUpdatePluginProvisioner(cmd)
	if err != nil {
		return err
	}
	err = p.UpdateServer()
	if err != nil {
		return err
	}
	return err
}
