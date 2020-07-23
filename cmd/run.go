package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the sweepstakes' rules",
	Long: `Start the process of execute the rules that defines a sweepstakes (it needs a logged account).
Example: sorteiagram run RULES_ID`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if sessionHash == "" {
			log.Panicln("Unable to execute the rules because there is no valid session")
		}

		instagram, err := LoadSession(sessionHash)

		if err != nil {
			log.WithError(err).Panicln("Unable to load the Instagram's session")
		}

		user, err := instagram.Profiles.ByName("carronovodasthe")

		log.Debug(user.Feed("1595532929"))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
