package minectl

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dirien/minectl-sdk/common"
	"github.com/dirien/minectl-sdk/model"
	"github.com/dirien/minectl-sdk/template"
	"github.com/dirien/minectl/internal/ui"
	"github.com/spf13/cobra"
)

var wizardCmd = &cobra.Command{
	Use:           "wizard",
	Short:         "Calls the minectl wizard to create interactively a minectl config",
	Example:       `mincetl wizard`,
	RunE:          RunFunc(runWizard),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	wizardCmd.Flags().StringP("output", "o", "", "output folder for the configuration file for minectl (default: ~/.minectl)")
	_ = wizardCmd.Flags().SetAnnotation("output", cobra.BashCompSubdirsInDir, []string{})
}

// Provider options for the wizard
var providerOptions = []huh.Option[string]{
	huh.NewOption("DigitalOcean", "DigitalOcean"),
	huh.NewOption("Civo", "Civo"),
	huh.NewOption("Scaleway", "Scaleway"),
	huh.NewOption("Hetzner", "Hetzner"),
	huh.NewOption("Akamai Connected Cloud", "Akamai Connected Cloud"),
	huh.NewOption("OVHcloud", "OVHcloud"),
	huh.NewOption("Equinix Metal", "Equinix Metal"),
	huh.NewOption("Google Compute Engine", "Google Compute Engine"),
	huh.NewOption("vultr", "vultr"),
	huh.NewOption("Azure", "Azure"),
	huh.NewOption("Oracle Cloud Infrastructure", "Oracle Cloud Infrastructure"),
	huh.NewOption("IONOS Cloud", "IONOS Cloud"),
	huh.NewOption("Amazon Web Services", "Amazon Web Services"),
	huh.NewOption("VEXXHOST", "VEXXHOST"),
	huh.NewOption("Ubuntu Multipass", "Ubuntu Multipass"),
	huh.NewOption("Exoscale", "Exoscale"),
	huh.NewOption("Fuga Cloud", "Fuga Cloud"),
}

var javaVersionOptions = []huh.Option[string]{
	huh.NewOption("Java 8", "8"),
	huh.NewOption("Java 16", "16"),
	huh.NewOption("Java 17", "17"),
}

var editionOptions = []huh.Option[string]{
	huh.NewOption("Bedrock", "bedrock"),
	huh.NewOption("Nukkit", "nukkit"),
	huh.NewOption("PowerNukkit", "powernukkit"),
	huh.NewOption("CraftBukkit", "craftbukkit"),
	huh.NewOption("Fabric", "fabric"),
	huh.NewOption("Forge", "forge"),
	huh.NewOption("Java (Vanilla)", "java"),
	huh.NewOption("PaperMC", "papermc"),
	huh.NewOption("Spigot", "spigot"),
	huh.NewOption("Purpur", "purpur"),
}

var featureOptions = []huh.Option[string]{
	huh.NewOption("Monitoring", "Monitoring"),
	huh.NewOption("RCON", "RCON"),
}

// isBedrockEdition returns true for editions that don't use Java.
func isBedrockEdition(edition string) bool {
	return edition == "bedrock" || edition == "nukkit" || edition == "powernukkit"
}

func runWizard(cmd *cobra.Command, _ []string) error {
	minectlUI.Info("minectl configuration file wizard")

	wizard := model.Wizard{
		SSHPort:  "22",
		BanTime:  "600",
		MaxRetry: "6",
		Java:     "16",
		Heap:     "2G",
		Edition:  "java",
	}
	var features []string

	// Validation functions
	nameValidator := func(s string) error {
		if s == "" {
			return errors.New("please enter the name of your Minecraft server")
		}
		match, _ := regexp.MatchString(common.NameRegex, s)
		if !match {
			return errors.New("the name must consist of lower case alphanumeric characters or '-'")
		}
		return nil
	}

	requiredValidator := func(s string) error {
		if s == "" {
			return errors.New("this field is required")
		}
		return nil
	}

	portValidator := func(s string) error {
		if s == "" {
			return errors.New("this field is required")
		}
		port, err := strconv.Atoi(s)
		if err != nil || port < 1 || port > 65535 {
			return errors.New("enter a valid port number (1-65535)")
		}
		return nil
	}

	heapValidator := func(s string) error {
		if s == "" {
			return errors.New("this field is required")
		}
		if !strings.Contains(s, "G") {
			return errors.New("enter the Java heap size in following notation <size>G, eg 2G")
		}
		return nil
	}

	// Group 1: Server basics
	basicGroup := huh.NewGroup(
		huh.NewInput().
			Title("Server Name").
			Description("Enter the name of your Minecraft server").
			Value(&wizard.Name).
			Validate(nameValidator),
		huh.NewSelect[string]().
			Title("Cloud Provider").
			Description("Choose a cloud provider").
			Options(providerOptions...).
			Value(&wizard.Provider),
		huh.NewInput().
			Title("Plan/Size").
			Description("Enter the plan/size for the server (e.g., s-4vcpu-8gb, e2-standard-2)").
			Value(&wizard.Plan).
			Validate(requiredValidator),
		huh.NewInput().
			Title("Region/Datacenter").
			Description("Enter the region/datacenter (e.g., fra1, europe-west6-a)").
			Value(&wizard.Region).
			Validate(requiredValidator),
	)

	// Group 2: SSH configuration
	sshGroup := huh.NewGroup(
		huh.NewInput().
			Title("SSH Public Key Path").
			Description("Enter the full path to your SSH public key (e.g., /home/user/.ssh/id_rsa.pub)").
			Value(&wizard.SSH).
			Validate(requiredValidator),
		huh.NewInput().
			Title("SSH Port").
			Description("Enter the SSH port").
			Value(&wizard.SSHPort).
			Validate(portValidator),
	)

	// Group 3: Security (fail2ban)
	securityGroup := huh.NewGroup(
		huh.NewInput().
			Title("Fail2ban Ban Time").
			Description("Enter the fail2ban ban time in seconds").
			Value(&wizard.BanTime).
			Validate(requiredValidator),
		huh.NewInput().
			Title("Fail2ban Max Retry").
			Description("Enter the fail2ban max retry count").
			Value(&wizard.MaxRetry).
			Validate(requiredValidator),
	)

	// Group 4: Features and edition selection
	featuresGroup := huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("Additional Features").
			Description("Select additional features for your server").
			Options(featureOptions...).
			Value(&features),
		huh.NewSelect[string]().
			Title("Minecraft Edition").
			Description("Select a Minecraft edition").
			Options(editionOptions...).
			Value(&wizard.Edition),
	)

	// Group 5: Java configuration (hidden for Bedrock-based editions)
	javaGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Java Version").
			Description("Choose the Java version").
			Options(javaVersionOptions...).
			Value(&wizard.Java),
		huh.NewInput().
			Title("Java Heap Size").
			Description("Enter the Java heap size (rule of thumb: half of available RAM)").
			Value(&wizard.Heap).
			Validate(heapValidator),
		huh.NewInput().
			Title("RCON Password").
			Description("Enter a RCON password (if you selected that feature)").
			Value(&wizard.RconPw),
	).WithHideFunc(func() bool {
		return isBedrockEdition(wizard.Edition)
	})

	// Group 6: Minecraft version and properties
	minecraftGroup := huh.NewGroup(
		huh.NewInput().
			Title("Minecraft Version").
			Description("Enter the Minecraft version number").
			Value(&wizard.Version).
			Validate(requiredValidator),
		huh.NewText().
			Title("Additional Properties").
			Description("Add additional Minecraft server properties (key=value, one per line)").
			Value(&wizard.Properties).
			CharLimit(1000),
	)

	// Create the form
	form := huh.NewForm(
		basicGroup,
		sshGroup,
		securityGroup,
		featuresGroup,
		javaGroup,
		minecraftGroup,
	)

	// Run the form with headless support
	err := ui.RunForm(form, headless)
	if err != nil {
		return err
	}

	// Set features
	wizard.Features = features

	config, err := template.NewTemplateConfig(wizard)
	if err != nil {
		return err
	}

	outputFolder := GetHomeFolder()
	output, err := cmd.Flags().GetString("output")
	if err == nil {
		if output != "" {
			outputFolder = output
		}
	}

	filename := fmt.Sprintf("%s/config-%s.yaml", outputFolder, wizard.Name)
	minectlUI.Info("Writing configuration file to " + filename)
	err = os.WriteFile(filename, []byte(config), 0o600)
	if err != nil {
		return err
	}
	minectlUI.Info("To create the server:\n\n  minectl create -f " + filename)
	return nil
}
