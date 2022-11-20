package minectl

import (
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.Flags().StringP("filename", "f", "", "Location of the manifest file")
	pluginCmd.Flags().SetAnnotation("filename", cobra.BashCompFilenameExt, []string{"yaml"}) //nolint:errcheck
	pluginCmd.Flags().String("id", "", "contains the server id")
	pluginCmd.Flags().StringP("plugin", "p", "", "Location of the plugin")
	pluginCmd.Flags().SetAnnotation("plugin", cobra.BashCompFilenameExt, []string{"jar"}) //nolint:errcheck
	pluginCmd.Flags().StringP("destination", "d", "", "Plugin destination folder")
	pluginCmd.Flags().SetAnnotation("destination", cobra.BashCompSubdirsInDir, []string{}) //nolint:errcheck
	pluginCmd.Flags().StringP("ssh-key", "k", "", "specify a specific path for the SSH key")
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
	DownloadURL string
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
	RunE:          RunFunc(runPlugin),
	SilenceUsage:  true,
	SilenceErrors: true,
}

var _ = []Plugin{
	{
		Name:        "Fabric API",
		Mod:         Fabric,
		Version:     []string{"0.37.1+1.16", "0.37.1+1.17"},
		DownloadURL: "https://github.com/FabricMC/fabric/releases/download/{{ .Version }}/fabric-api-{{ .Version }}.jar",
		Destination: "/mincraft/mods",
	},
}

func runPlugin(cmd *cobra.Command, _ []string) error {
	p, err := createUpdatePluginProvisioner(cmd)
	if err != nil {
		return err
	}
	plugin, _ := cmd.Flags().GetString("plugin")
	destination, _ := cmd.Flags().GetString("destination")
	err = p.UploadPlugin(plugin, destination)
	if err != nil {
		return err
	}
	return err
}
