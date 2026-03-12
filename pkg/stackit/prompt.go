package stackit

import (
	"fmt"
	"time"

	libnuke "github.com/ekristen/libnuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/utils"
)

// Prompt is a struct that contains the parameters and tenant details use to craft a unique prompt
// for the user to confirm the nuke operation.
type Prompt struct {
	Parameters *libnuke.Parameters
	ProjectID  string
}

func (p *Prompt) Prompt() error {
	forceSleep := time.Duration(p.Parameters.ForceSleep) * time.Second

	if p.Parameters.Force {
		fmt.Printf("no-prompt flag set, continuing without prompting user")
		fmt.Printf("waiting %v before continuing", forceSleep)
		time.Sleep(forceSleep)
	} else {
		fmt.Printf("Do you really want to nuke the StackIT project with the ID %d.\n", p.ProjectID)
		fmt.Printf("Do you want to continue? Enter project id %q to continue.\n", p.ProjectID)
		if err := utils.Prompt(p.ProjectID); err != nil {
			return err
		}
	}

	return nil
}
