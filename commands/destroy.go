package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	forge "github.com/nathandines/forge/forgelib"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a CloudFormation Stack",
	Run: func(cmd *cobra.Command, args []string) {
		if assumeRoleArn != "" {
			if err := forge.AssumeRole(assumeRoleArn); err != nil {
				log.Fatal(err)
			}
		}

		// Populate Stack ID
		if err := stack.GetStackInfo(); err != nil {
			log.Fatal(err)
		}

		after, err := stack.GetLastEventTime()
		if err != nil {
			log.Fatal(err)
		}

		if err := stack.Destroy(); err != nil {
			log.Fatal(err)
		}

		for {
			// Refresh Stack State
			if err := stack.GetStackInfo(); err != nil {
				log.Fatal(err)
			}

			printStackEvents(&stack, after)

			status := *stack.StackInfo.StackStatus
			switch {
			case stackInProgressRegexp.MatchString(status):
			case status == cloudformation.StackStatusDeleteComplete:
				return
			default:
				fmt.Print("\n")
				log.Fatal(fmt.Errorf("Stack destroy failed! Stack Status: %s", status))
			}

			time.Sleep(5 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
