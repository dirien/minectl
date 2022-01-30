package minectl

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/minectl/internal/common"

	"github.com/AlecAivazis/survey/v2"
	"github.com/minectl/internal/model"
	"github.com/minectl/internal/template"
	"github.com/spf13/cobra"
)

var wizardQuestions = []*survey.Question{
	{
		Name:   "name",
		Prompt: &survey.Input{Message: "Enter the name of your Minecraft server:"},
		Validate: func(val interface{}) error {
			str, err := val.(string)
			if !err {
				return errors.New("please enter a the name of your Minecraft server")
			}
			match, _ := regexp.MatchString(common.NameRegex, str)
			if !match {
				return errors.New("the name of your Minecraft server must consist of lower case alphanumeric characters or '-'")
			}
			return nil
		},
	},
	{
		Name:     "provider",
		Validate: survey.Required,
		Prompt: &survey.Select{
			Message: "Choose a cloud provider:",
			Options: []string{
				"DigitalOcean",
				"Civo",
				"Scaleway",
				"Hetzner",
				"Linode",
				"OVHcloud",
				"Equinix Metal",
				"Google Compute Engine",
				"vultr",
				"Azure",
				"Oracle Cloud Infrastructure",
				"IONOS Cloud",
				"Amazon WebServices",
				"VEXXHOST",
				"Ubuntu Multipass",
				"Exoscale",
			},
			PageSize: 15,
		},
	},
	{
		Name:     "plan",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter the plan/size for the server:",
			Help:    "Plans could have dedicated names, like s-4vcpu-8gb for DO or e2-standard-2 for GCE. Check your provider for details.",
		},
	},
	{
		Name:     "region",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter the region/datacenter:",
			Help:    "Regions could have dedicated names, like fra1 for DO or europe-west6-a for GCE. Check your provider for details",
		},
	},
	{
		Name:     "ssh",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter a full path to ssh private key (like /f/a/key):",
			Help:    "Please enter the full path to ssh private key like this -> /Users/dirien/Tools/repos/stackit-minecraft/minecraft/ssh/minecraft-be",
		},
	},
	{
		Name:     "ssh_port",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter the ssh port (default 22)",
			Default: "22",
		},
	},
	{
		Name:     "fail2ban_bantime",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter the fail2ban bantime (default 600):",
			Default: "600",
		},
	},
	{
		Name:     "fail2ban_maxretry",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter the fail2ban maxretry (default 6):",
			Default: "6",
		},
	},
	{
		Name: "features",
		Prompt: &survey.MultiSelect{
			Message: "Which additional features do you want for your server?:",
			Options: []string{"Monitoring", "RCON"},
		},
	},
	{
		Name: "java",
		Prompt: &survey.Select{
			Message: "Choose the Java version:",
			Help:    "If you are not running bedrock, you have to select the java version.",
			Options: []string{"8", "16", "17"},
			Default: "16",
		},
	},
	{
		Name: "heap",
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || !strings.Contains(str, "G") {
				return errors.New("enter the Java heap size in following notation <size>G, eg 2G")
			}
			return nil
		},
		Prompt: &survey.Input{
			Message: "Enter the Java heap size in following notation <size>G:",
			Help:    "Rule of thumb is to set the heap size to around half of your available system RAM",
			Default: "2G",
		},
	},
	{
		Name: "rconpw",
		Prompt: &survey.Input{
			Message: "Enter a RCON password (if you selected that feature):",
			Help:    "The password for the RCON function",
		},
	},
	{
		Name: "edition",
		Prompt: &survey.Select{
			Message:  "Select a Minecraft edition:",
			Options:  []string{"bedrock", "nukkit", "powernukkit", "craftbukkit", "fabric", "forge", "java", "papermc", "spigot"},
			Default:  "java",
			PageSize: 10,
		},
	},
	{
		Name:     "version",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "Enter the Minecraft version number:",
			Help:    "The version of the Minecraft edition, you want to install.",
		},
	},
	{
		Name: "properties",
		Prompt: &survey.Multiline{
			Message: "Add additional Minecraft server properties:",
			Help:    "Properties like level-seed,max-players etc. Enter as key=value",
		},
	},
}

var wizardCmd = &cobra.Command{
	Use:           "wizard",
	Short:         "Calls the minectl wizard to create interactively a minectl 🗺 config",
	Example:       `mincetl wizard`,
	RunE:          RunFunc(runWizard),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	usage := fmt.Sprintf("output folder for the configuration file for minectl 🗺 (default: %s)", GetHomeFolder())
	wizardCmd.Flags().StringP("output", "o", "", usage)
}

func runWizard(cmd *cobra.Command, _ []string) error {
	minectlLog.RawMessage("🧙 minectl configuration file wizard\n")
	wizard := model.Wizard{}
	err := survey.Ask(wizardQuestions, &wizard)
	if err != nil {
		return err
	}

	config, err := template.NewTemplateConfig(wizard)
	if err != nil {
		return err
	}

	outputFolder := GetHomeFolder()
	output, err := cmd.Flags().GetString("output")
	if err == nil {
		if len(output) > 0 {
			outputFolder = output
		}
	}

	filename := fmt.Sprintf("%s/config-%s.yaml", outputFolder, wizard.Name)
	minectlLog.PrintMixedGreen("\n📄 Writing configuration file to %s", filename)
	err = os.WriteFile(filename, []byte(config), 0o600)
	if err != nil {
		return err
	}
	minectlLog.PrintMixedGreen("\n🏗 To create the server type:\n\n %s", fmt.Sprintf("minectl create -f %s \n", filename))
	return nil
}
