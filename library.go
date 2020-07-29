// Package gcppubsubdemo is a tool for teaching and learning Google
// Cloud Platform's pub/sub API. This package is not meant to be used
// as an API but it is documented like one for instructional purposes.
//
// Quick usage examples for the package are shown below. For more
// information, see the documentation for the individual items.
//
// Publishing
//
// To use the API for publishing, call the GetPublisher function to
// obtain a Publisher function to call. Call the Publisher function
// giving the data to publish as the only parameter. To shut down the
// Publisher, call the Publisher function with nil as the parameter.
// Failing to do so will leave pub/sub code and network operations
// initialized.
//
//  if publish, err := GetPublisher("GCP Name", "Topic Name"); err != nil {
//    fmt.Fprintln(os.Stderr, err)
//  } else {
//    publish([]byte("Payload"))
//    publish(nil)
//  }
//
// Subscribing
//
// Subscribe to a topic by callng the Subscribe function, giving it
// a callback to execute when data is received. The callback should return
// true when it no longer wants to retrieve messages and shut down the
// subscription.
//  if err := Subscribe("GCP Name", "Topic Name", "Subscription ID", func(data []byte) bool {
//    fmt.Println(string(data)
//    return true
//  }); err != nil {
//    fmt.Fprintln(os.Stderr, err)
//  }
//
// Logging
//
// Basic logging is implemented using the standard log package. Logging is
// disabled by default. When enabled with a call to LogLevel(), log output
// goes to os.Stderr but can be redirected by using LogOutput().
//
//  writer, _ := os.Create("logfile.log")
//  LogOutput(writer)
//  LogLevel(DebugLogLevel)
package gcppubsubdemo

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
)

func getPubsubClient(gcpName string) (*pubsub.Client, error) {
	logger.TRACE.Printf(`gcppubsubdemo.getPubsubClient("%s")`, gcpName)

	client, err := pubsub.NewClient(context.Background(), gcpName)
	if err != nil {
		logger.ERROR.Print("Error creating pub/sub client:", err)
		return nil, err
	}

	return client, nil
}

func getOrCreatePubsubTopic(client *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	logger.TRACE.Printf(`gcppubsubdemo.getOrCreatePubsubTopic(pubsubClient, "%s")`, topicName)

	ctx := context.Background()
	topic := client.Topic(topicName)

	// Because CreateTopic() can return an error if the top already exists,
	// there is a race condition between the time Exists() and CreateTopic() is
	// called. This loop accounts for this by following the call to CreateTopic()
	// with another call to Exist() to determine if CreateTopic() failed because
	// this race condition condition occurred. This is only a concern if there are multiple
	// processes or goroutines than can create the same subscription.

	// Simplification of this logic would be a welcome code change.
	attempts := 0
	for attempts <= 1 {
		logger.TRACE.Printf("Checking for existing %s topic", topicName)
		if ok, err := topic.Exists(ctx); err != nil {
			logger.ERROR.Print("Unable to verify topic exists:", err)
			return nil, err
		} else if ok {
			break
		}

		if attempts < 1 {
			logger.DEBUG.Printf("Topic %s does not exist. Creating it.", topicName)
			if _, err := client.CreateTopic(ctx, topicName); err != nil {
				logger.ERROR.Print("Error creating topic:", err)
			}
		} else {
			errText := "Could not create topic " + topicName
			logger.ERROR.Print(errText)
			return nil, errors.New(errText)
		}

		attempts++
	}

	return topic, nil
}
