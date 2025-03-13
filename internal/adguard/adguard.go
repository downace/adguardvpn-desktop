package adguard

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"os/exec"
	"strings"
)

type Status struct {
	Connected bool `json:"connected"`
}

type Cli struct {
	CliBin string
}

func (a *Cli) exec(args ...string) (string, error) {
	cmd := exec.Command(a.CliBin, args...)

	outputBytes, err := cmd.CombinedOutput()

	output := string(outputBytes)

	if err != nil && output != "" {
		err = fmt.Errorf("%s%s", output, err)
	}

	return stripansi.Strip(output), err
}

func (a *Cli) Version() (string, error) {
	return a.exec("--version")
}

func (a *Cli) Status() (*Status, error) {
	statusOutput, err := a.exec("status")
	if err != nil {
		return nil, err
	}

	status := Status{
		Connected: strings.Contains(statusOutput, "is connected"),
	}

	return &status, nil
}
