package main

import (
	"fmt"
	"os"
	"time"

	"github.com/arcanericky/gcppubsubdemo"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish to a topic",
	Run: func(cmd *cobra.Command, args []string) {
		var gcpName,
			topicName string
		var once bool
		var interval time.Duration
		var err error
		var publish gcppubsubdemo.Publisher

		if setLogLevel(cmd) != nil {
			return
		}

		fmt.Println("Executing publish command")

		if gcpName,
			topicName,
			once,
			err = getPersistentFlags(cmd); err != nil {
			exitCode = 1
			return
		}

		if interval, err = cmd.Flags().GetDuration(optionInterval); err != nil {
			outputOptionError(optionInterval, err)
			exitCode = 1
			return
		}

		if publish, err = gcppubsubdemo.GetPublisher(gcpName, topicName); err != nil {
			fmt.Fprint(os.Stderr, "Error getting publisher:", err)
			exitCode = 1
			return
		}

		if once {
			publish([]byte(getNow()))
			publish(nil)
			return
		}

		for publish([]byte(getNow())) == nil {
			time.Sleep(interval)
		}
	},
}

func init() {
	publishCmd.Flags().DurationP(optionInterval, "", time.Duration(5)*time.Second, helpInterval)
	rootCmd.AddCommand(publishCmd)
}

func getNow() string {
	now := time.Now().String()
	fmt.Println("Publishing:", now)
	return now
}
