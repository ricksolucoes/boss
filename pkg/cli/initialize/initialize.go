package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashload/boss/internal/pkg/configuration"
	"github.com/hashload/boss/internal/pkg/models"
	"github.com/hashload/boss/internal/pkg/ui"
	"github.com/hashload/boss/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdInitialize add the command line init
func NewCmdInitialize(config *configuration.Configuration) *cobra.Command {
	var quiet bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Long:  "Initialize a new project and creates a boss.json file",
		Example: `  Initialize a new project:
  boss init

  Initialize a new project without having it ask any questions:
  boss init --quiet`,
		Run: func(cmd *cobra.Command, args []string) {
			err := initalizePackage(config, quiet)
			util.CheckErr(err)
		},
	}
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "without asking questions")
	return cmd
}

func initalizePackage(config *configuration.Configuration, quiet bool) error {
	folderName := ""
	wd, err := config.CurrentDir()
	if err != nil {
		return err
	}

	var finalFile = filepath.Join(wd, "boss.json")

	if stat, err := os.Stat(finalFile); err == nil && !stat.IsDir() {
		if ok, _ := ui.GetConfirmation("boss.json already exists, do you want to overwrite it", false); !ok {
			return nil
		}
	}

	rxp, err := regexp.Compile(`^.+\` + string(filepath.Separator) + `([^\\]+)$`)
	if err == nil {
		allString := rxp.FindAllStringSubmatch(wd, -1)
		folderName = allString[0][1]
	}

	var pkg = models.MakeBossPackage()

	if quiet {
		pkg.Name = folderName
		pkg.Version = "1.0.0"
		pkg.MainSrc = "./"
	} else {
		printHead()
		pkg.Name = ui.GetTextOrDef("Package name", folderName)
		pkg.Homepage = ui.GetTextOrDef("homepage", "")
		pkg.Version = ui.GetTextOrDef("version: (1.0.0)", "1.0.0")
		pkg.Description = ui.GetTextOrDef("description", "")
		pkg.MainSrc = ui.GetTextOrDef("source folder: (./)", "./")
	}

	json, err := pkg.SaveToFile(finalFile)
	if err == nil {
		fmt.Println("\n" + string([]byte(json)))
	}
	return err
}

func printHead() {
	fmt.Println(`
This utility will walk you through creating a boss.json file.
It only covers the most common items, and tries to guess sensible defaults.

Use 'boss install <pkg>' afterwards to install a package and
save it as a dependency in the boss.json file.
Press ^C at any time to quit.`)
}
