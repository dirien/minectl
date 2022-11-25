package minectl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dirien/minectl/internal/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringP("filename", "f", "", "Location of the manifest file")
	deleteCmd.Flags().SetAnnotation("filename", cobra.BashCompFilenameExt, []string{"yaml"}) //nolint:errcheck
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
	if len(filename) == 0 {
		return errors.New("Please provide a valid manifest file via -f|--filename flag")
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

	yes := cmd.Flag("yes").Changed
	if yes {
		err = newProvisioner.DeleteServer()
		if err != nil {
			return err
		}
	} else {
		fmt.Print("Do you want to delete the Minecraft server? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		if strings.Compare("y", text) == 0 {
			err = newProvisioner.DeleteServer()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
