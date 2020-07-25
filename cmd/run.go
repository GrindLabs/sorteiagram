package cmd

import (
	"reflect"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/grindlabs/sorteiagram/actions"

	"github.com/grindlabs/sorteiagram/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the sweepstakes' rules",
	Long: `Start the process of execute the rules that defines a sweepstakes (it needs a logged account).
Example: sorteiagram run RULES_ID`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if sessionHash == "" {
			log.Panicln("Unable to execute the rules because there is no valid session")
		}

		instagram, err := utils.LoadSession(sessionHash)

		if err != nil {
			log.WithError(err).Panicln("Unable to load the Instagram's session")
		}

		call(utils.LoadRules(args[0]), instagram)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// Call the respective rules
func call(rules map[string][]interface{}, instagram *goinsta.Instagram) {
	for action, params := range rules {
		if _, ok := actions.ActionsMap[action]; ok {
			fn := reflect.ValueOf(actions.ActionsMap[action])
			in := make([]reflect.Value, len(params)+1)
			in[0] = reflect.ValueOf(instagram)

			for k, param := range params {
				in[k+1] = reflect.ValueOf(param)
			}

			fn.Call(in)
		}
	}
}
