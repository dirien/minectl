package minectl

import (
	"log"

	"github.com/minectl/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	minectlCmd.AddCommand(pluginCmd)
	pluginCmd.Flags().StringP("filename", "f", "", "Contains the configuration for minectl")
	pluginCmd.Flags().String("id", "", "contains the server id")
	pluginCmd.Flags().StringP("plugin", "p", "", "Local plugin file location")
	pluginCmd.Flags().StringP("destination", "d", "", "Plugin destination location")
}

type ModType string

const (
	Forge  ModType = "forge"
	Fabric ModType = "fabric"
)

type Plugin struct {
	Name        string
	Mod         ModType
	Version     []string
	DownloadUrl string
	Destination string
}

var pluginCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage your plugins for a specific server",
	Example: `mincetl plugins  \
    --filename server-do.yaml
    --id xxx-xxx-xxx-xxx
	--plugin plugin.jar
    --destination /minecraft/mods`,
	RunE:          runPlugin,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var _ = []Plugin{
	{
		Name:        "Fabric API",
		Mod:         Fabric,
		Version:     []string{"0.37.1+1.16", "0.37.1+1.17"},
		DownloadUrl: "https://github.com/FabricMC/fabric/releases/download/{{ .Version }}/fabric-api-{{ .Version }}.jar",
		Destination: "/mincraft/mods",
	},
}

func runPlugin(cmd *cobra.Command, _ []string) error {
	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return errors.Wrap(err, "Please provide a valid MinecraftResource manifest file")
	}
	if len(filename) == 0 {
		return errors.New("Please provide a valid MinecraftResource manifest file")
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
		log.Fatal(err)
	}
	plugin, _ := cmd.Flags().GetString("plugin")
	destination, _ := cmd.Flags().GetString("destination")
	err = p.UploadPlugin(plugin, destination)
	if err != nil {
		return err
	}
	return err
}
