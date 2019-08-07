package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/mitchellh/cli"
	"trellis-cli/trellis"
)

type DotEnvCommand struct {
	UI      cli.Ui
	Trellis *trellis.Trellis
}

func deletePlaybook(path string) {
	err := os.Remove(path)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *DotEnvCommand) Run(args []string) int {
	if err := c.Trellis.LoadProject(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var environment string

	switch len(args) {
	case 0:
		environment = "development"
	case 1:
		environment = args[0]
	default:
		c.UI.Error(fmt.Sprintf("Error: too many arguments (expected 0 or 1, got %d)\n", len(args)))
		c.UI.Output(c.Help())
		return 1
	}

	_, ok := c.Trellis.Environments[environment]
	if !ok {
		c.UI.Error(fmt.Sprintf("Error: %s is not a valid environment", environment))
		return 1
	}

	// Copy playbook file from package to Trellis
	_, filename, _, ok := runtime.Caller(0)
	playbookTemplatePath := path.Join(path.Dir(filename), "../playbooks/dotenv.yml")

	b, readFileErr := ioutil.ReadFile(playbookTemplatePath)
	if readFileErr != nil {
		log.Fatal(readFileErr)
	}

	playbookPath := "dotenv.yml"
	writeFileErr := ioutil.WriteFile(playbookPath, b, 0644)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
	defer deletePlaybook(playbookPath)

	dotEnv := execCommand("ansible-playbook", "dotenv.yml", "-e", "env=" + environment)

	env := os.Environ()
	// To allow mockExecCommand injects its environment variables
	if dotEnv.Env != nil {
		env = dotEnv.Env
	}
	dotEnv.Env = append(env, "ANSIBLE_RETRY_FILES_ENABLED=false")

	logCmd(dotEnv, c.UI, true)
	runErr := dotEnv.Run()

	if runErr != nil {
		c.UI.Error(fmt.Sprintf("Error running ansible-playbook: %s", runErr))
	}

	return 0
}

func (c *DotEnvCommand) Synopsis() string {
	return "Template .env files to local system"
}

func (c *DotEnvCommand) Help() string {
	helpText := `
Usage: trellis dotenv [options] [ENVIRONMENT=development]

Template .env files to local system

Options:
  -h, --help  show this help
`

	return strings.TrimSpace(helpText)
}
