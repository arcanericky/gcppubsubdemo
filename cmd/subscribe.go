package main

import (
	"fmt"
	"os"

	"github.com/arcanericky/gcppubsubdemo"
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a topic",
	Run: func(cmd *cobra.Command, args []string) {
		var gcpName,
			subscriptionID,
			topicName string
		var once bool
		var err error

		if setLogLevel(cmd) != nil {
			return
		}

		fmt.Println("Executing subscribe command")

		if gcpName,
			topicName,
			once,
			err = getPersistentFlags(cmd); err != nil {
			exitCode = 1
			return
		}

		if subscriptionID, err = cmd.Flags().GetString(optionSubscriptionID); err != nil {
			outputOptionError(optionSubscriptionID, err)
			exitCode = 1
			return
		}

		if err = gcppubsubdemo.Subscribe(gcpName, topicName, subscriptionID, func(data []byte) bool {
			fmt.Printf("Received: %s\n", data)

			if once {
				return true
			}

			return false
		}); err != nil {
			fmt.Fprintln(os.Stderr, "Error with the subscribe command:", err)
			exitCode = 1
		}
	},
}

func init() {
	subscribeCmd.Flags().StringP(optionSubscriptionID, "", "mysubscriptionid", helpSubscriptionID)

	rootCmd.AddCommand(subscribeCmd)
}
