package minectl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"go.uber.org/zap"

	"github.com/blang/semver/v4"
	"github.com/dirien/minectl/internal/logging"
	"github.com/dirien/minectl/internal/provisioner"
	"github.com/dirien/minectl/internal/ui"
	"github.com/mitchellh/go-homedir"
	"github.com/morikuni/aec"
	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
)

var (
	Version   string
	GitCommit string
	Date      string
)

func createUpdatePluginProvisioner(cmd *cobra.Command) (provisioner.Provisioner, error) {
	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Please provide a valid manifest file")
	}
	if filename == "" {
		return nil, pkgerrors.New("Please provide a valid manifest file via -f|--filename flag")
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, pkgerrors.New("Please provide a valid id")
	}
	sshKey, err := cmd.Flags().GetString("ssh-key")
	if err != nil {
		return nil, err
	}
	if sshKey == "" {
		return nil, pkgerrors.New("Please provide a valid ssh key path")
	}
	p, err := provisioner.NewProvisioner(&provisioner.MinectlProvisionerOpts{
		ManifestPath:      filename,
		ID:                id,
		SSHPrivateKeyPath: sshKey,
	}, minectlUI)
	if err != nil {
		return nil, err
	}
	return p, nil
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

	brewPrefixCmd := exec.CommandContext(context.Background(), brewBin, "--prefix", "minectl")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	brewPrefixCmd.Stdout = &stdout
	brewPrefixCmd.Stderr = &stderr
	if err = brewPrefixCmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			ee.Stderr = stderr.Bytes()
		}
		return false, pkgerrors.Wrapf(err, "'brew --prefix minectl' failed")
	}

	brewPrefixCmdOutput := strings.TrimSpace(stdout.String())
	if brewPrefixCmdOutput == "" {
		return false, pkgerrors.New("trimmed output from 'brew --prefix minectl' is empty")
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
			minectlUI.ErrorMsg(res)
			if postRunErr := runPostCommandHooks(cmd, args); postRunErr != nil {
				minectlUI.ErrorMsg(postRunErr)
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

func getUpgradeMessage(latestVer, current *semver.Version) *string {
	cmd := getUpgradeCommand()
	msg := fmt.Sprintf("A new version of minectl is available. To upgrade from version '%s' to '%s', ", current, latestVer)
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
	if Version != "" {
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

var (
	headless          bool
	minectlLog        *logging.MinectlLogging
	minectlUI         *ui.UI
	updateCheckResult chan *string
)

var minectlCmd = &cobra.Command{
	Use:   "minectl",
	Short: "Create Minecraft Server on different cloud provider.",
	Run:   runMineCtl,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: false,
	},
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		headless, _ = cmd.Flags().GetBool("headless")
		verbose, _ := cmd.Flags().GetString("verbose")
		logEncoding, _ := cmd.Flags().GetString("log-encoding")
		if headless && verbose == "" {
			verbose = "info"
		}
		var err error
		minectlLog, err = logging.NewLogging(verbose, logEncoding, headless)
		if err != nil {
			os.Exit(0)
		}
		minectlUI = ui.NewUI(headless, minectlLog)
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
	PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
		checkVersionMsg, ok := <-updateCheckResult
		if ok && checkVersionMsg != nil {
			zap.S().Infof("Warning(%v)", checkVersionMsg)
			fmt.Println()
			fmt.Println(*checkVersionMsg)
		}
		return nil
	},
}

func makeAppDirectoryIfNotExists() {
	dir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("error determining the home directory: %s", err)
	}
	path := dir + "/.minectl"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModeDir|0o755)
		if err != nil {
			fmt.Printf("error while creating the minectl directory: %s", err)
		}
	}
}

func GetHomeFolder() string {
	dir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("error determining the home directory: %s", err)
	}
	return dir + "/.minectl"
}

func init() {
	makeAppDirectoryIfNotExists()
	minectlCmd.PersistentFlags().String("verbose", "",
		"Enable verbose logging: debug|info|warn|error|dpanic|panic|fatal")
	minectlCmd.PersistentFlags().String("log-encoding", "console",
		"Set the log encoding: console|json (default: console)")
	minectlCmd.PersistentFlags().Bool("headless", false,
		"Set this value to if mincetl is called by a CI system. Enables logging and disables human-readable output rendering (default: false)")
	minectlCmd.AddCommand(versionCmd)
	minectlCmd.AddCommand(createCmd)
	minectlCmd.AddCommand(deleteCmd)
	minectlCmd.AddCommand(listCmd)
	minectlCmd.AddCommand(wizardCmd)
	minectlCmd.AddCommand(pluginCmd)
	minectlCmd.AddCommand(rconCmd)
	minectlCmd.AddCommand(updateCmd)
}

func Execute(version, gitCommit, date string) error {
	Version = version
	GitCommit = gitCommit
	Date = date

	return minectlCmd.Execute()
}

func runMineCtl(cmd *cobra.Command, _ []string) {
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
