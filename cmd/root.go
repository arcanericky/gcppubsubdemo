package main

import (
	"fmt"
	"os"

	"github.com/arcanericky/gcppubsubdemo"
	"github.com/spf13/cobra"
)

func outputOptionError(option string, err error) {
	fmt.Fprintf(os.Stderr, "Error parsing %s option: %s\n", option, err)
}

func setLogLevel(cmd *cobra.Command) error {
	var verbose bool
	var err error

	if verbose, err = cmd.Flags().GetBool(optionVerbose); err != nil {
		outputOptionError(optionVerbose, err)
		return err
	}

	if verbose == true {
		gcppubsubdemo.LogLevel(gcppubsubdemo.TraceLogLevel)
	} else {
		gcppubsubdemo.LogLevel(gcppubsubdemo.ErrorLogLevel)
	}

	return nil
}

func getPersistentFlags(cmd *cobra.Command) (string, string, bool, error) {
	var gcpName,
		topicName string
	var once bool
	var err error

	if gcpName, err = cmd.Flags().GetString(optionGCP); err != nil {
		outputOptionError(optionGCP, err)
		return gcpName, topicName, once, err
	}

	if topicName, err = cmd.Flags().GetString(optionTopic); err != nil {
		outputOptionError(optionTopic, err)
		return gcpName, topicName, once, err
	}

	if once, err = cmd.Flags().GetBool(optionOnce); err != nil {
		outputOptionError(optionOnce, err)
		return gcpName, topicName, once, err
	}

	return gcpName, topicName, once, nil
}

var rootCmd = &cobra.Command{
	Use:   "gcppubsubdemo",
	Short: "Google Cloud Platform Pubsub Demo",
}

func init() {
	rootCmd.PersistentFlags().StringP(optionGCP, "", "google-cloud-project-id", helpGCP)
	rootCmd.PersistentFlags().StringP(optionTopic, "", "mytopicname", helpTopic)
	rootCmd.PersistentFlags().BoolP(optionOnce, "", false, helpOnce)
	rootCmd.PersistentFlags().BoolP(optionVerbose, "", false, helpVerbose)
}
