package minectl

import (
	"fmt"
	"os"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

func init() {
	minectlCmd.AddCommand(versionCmd)
}

var minectlCmd = &cobra.Command{
	Use:   "minectl",
	Short: "Create Minecraft Server on different cloud provider.",
	Run:   runMineCtl,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the clients version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	os.Exit(0)
}

func Execute(version, gitCommit string) error {

	Version = version
	GitCommit = gitCommit

	if err := minectlCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func runMineCtl(cmd *cobra.Command, args []string) {
	printLogo()
	err := cmd.Help()
	if err != nil {
		os.Exit(0)
	}
}

func printLogo() {
	minectlLogo := aec.WhiteF.Apply(minectlFigletStr)
	fmt.Println(minectlLogo)
}

const minectlFigletStr = `
 _______ _____ __   _ _______ _______ _______       
 |  |  |   |   | \  | |______ |          |    |     
 |  |  | __|__ |  \_| |______ |_____     |    |_____
`
