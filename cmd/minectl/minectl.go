package minectl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/blang/semver/v4"
	"github.com/morikuni/aec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
)

var (
	Version   string
	GitCommit string
	Date      string
)

func init() {
	minectlCmd.AddCommand(versionCmd)
}

var updateCheckResult chan *string

var minectlCmd = &cobra.Command{
	Use:   "minectl",
	Short: "Create Minecraft Server on different cloud provider.",
	Run:   runMineCtl,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: false,
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("PersistentPostRun")
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var waitForUpdateCheck bool
		defer func() {
			if !waitForUpdateCheck {
				close(updateCheckResult)
			}
		}()

		updateCheckResult = make(chan *string)
		waitForUpdateCheck = true
		go func() {
			updateCheckResult <- checkForUpdate()
			close(updateCheckResult)
		}()
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		checkVersionMsg, ok := <-updateCheckResult
		if ok && checkVersionMsg != nil {
			fmt.Println()
			fmt.Println(to.String(checkVersionMsg))
		}
		return nil
	},
}

func isDevVersion(s semver.Version) bool {
	if len(s.Pre) == 0 {
		return false
	}

	devStrings := regexp.MustCompile(`alpha|beta|dev|rc`)
	return !s.Pre[0].IsNum && devStrings.MatchString(s.Pre[0].VersionStr)
}

func isBrewInstall(exe string) (bool, error) {
	if runtime.GOOS != "darwin" {
		return false, nil
	}

	exePath, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return false, err
	}

	brewBin, err := exec.LookPath("brew")
	if err != nil {
		return false, err
	}

	brewPrefixCmd := exec.Command(brewBin, "--prefix", "minectl")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	brewPrefixCmd.Stdout = &stdout
	brewPrefixCmd.Stderr = &stderr
	if err = brewPrefixCmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			ee.Stderr = stderr.Bytes()
		}
		return false, errors.Wrapf(err, "'brew --prefix minectl' failed")
	}

	brewPrefixCmdOutput := strings.TrimSpace(stdout.String())
	if brewPrefixCmdOutput == "" {
		return false, errors.New("trimmed output from 'brew --prefix minectl' is empty")
	}

	brewPrefixPath, err := filepath.EvalSymlinks(brewPrefixCmdOutput)
	if err != nil {
		return false, err
	}

	brewPrefixExePath := filepath.Join(brewPrefixPath, "minectl")
	return exePath == brewPrefixExePath, nil
}

func runPostCommandHooks(c *cobra.Command, args []string) error {
	if c.PostRunE != nil {
		if err := c.PostRunE(c, args); err != nil {
			return err
		}
	} else if c.PostRun != nil {
		c.PostRun(c, args)
	}
	for p := c; p != nil; p = p.Parent() {
		if p.PersistentPostRunE != nil {
			if err := p.PersistentPostRunE(c, args); err != nil {
				return err
			}
			break
		} else if p.PersistentPostRun != nil {
			p.PersistentPostRun(c, args)
			break
		}
	}
	return nil
}

func RunFunc(run func(cmd *cobra.Command, args []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if res := run(cmd, args); res != nil {
			fmt.Println(res)
			if postRunErr := runPostCommandHooks(cmd, args); postRunErr != nil {
				fmt.Println(res)
			}
			os.Exit(1)
		}
		os.Exit(0)
		return nil
	}
}

func getUpgradeCommand() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}

	isBrew, err := isBrewInstall(exe)
	if err != nil {
		fmt.Printf("error determining if the running executable was installed with brew: %s", err)
	}
	if isBrew {
		return "$ brew upgrade minectl"
	}

	if runtime.GOOS != "windows" {
		return "$ curl -sSL https://get.minectl.dev | sh"
	}
	return ""
}

func getUpgradeMessage(latest *semver.Version, current *semver.Version) *string {
	cmd := getUpgradeCommand()
	msg := fmt.Sprintf("A new version of minectl is available. To upgrade from version '%s' to '%s', ", current, latest)
	if cmd != "" {
		msg += "run \n   " + cmd + "\n\nor "
	}

	msg += "visit https://github.com/dirien/minectl#installing-minectl- for manual instructions."
	return &msg
}

func getCLIVersionInfo(current *semver.Version) (*semver.Version, error) {
	githubTag := &latest.GithubTag{
		Owner:      "dirien",
		Repository: "minectl",
	}

	res, err := latest.Check(githubTag, current.String())
	if err != nil {
		return nil, err
	}
	version, err := semver.New(res.Current)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func checkForUpdate() *string {
	curVer, err := semver.ParseTolerant(getVersion())
	if err != nil {
		fmt.Printf("error parsing current version: %s", err)
	}
	if isDevVersion(curVer) {
		return nil
	}
	latestVer, err := getCLIVersionInfo(&curVer)
	if err != nil {
		return nil
	}
	if latestVer.GT(curVer) {
		return getUpgradeMessage(latestVer, &curVer)
	}

	return nil
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
	return "0.1.0-dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()
	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	fmt.Println("Build date:", Date)
}

func Execute(version, gitCommit, date string) error {

	Version = version
	GitCommit = gitCommit
	Date = date

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
