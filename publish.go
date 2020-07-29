package gcppubsubdemo

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// Publisher is a function that is returned to caller from
// GetPublisher() and is used to by the caller to publish data
type Publisher func([]byte) error

func publishPubsubData(topic *pubsub.Topic, data []byte) error {
	logger.TRACE.Printf(`gcppubsubdemo.publishPubsubData("%s", data)`, topic.ID())

	ctx := context.Background()

	msg := &pubsub.Message{
		Data: data,
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		logger.ERROR.Print("Error publishing data:", err)
		return err
	}
	logger.DEBUG.Printf("Published: %s", msg.Data)

	return nil
}

func closePubsub(client *pubsub.Client, topic *pubsub.Topic) {
	logger.TRACE.Print(`gcppubsubdemo.closePubsub(client, topic)`)

	if topic != nil {
		logger.DEBUG.Print("Deleting topic")
		topic.Delete(context.Background())
	}

	if client != nil {
		logger.DEBUG.Print("Closing client")
		client.Close()
	}
}

// GetPublisher initializes a pub/sub client and topic returning a
// Publisher function to the caller. The caller can then call the Publisher
// function with the data to be published. The caller can shut down the
// Publisher by calling it with a nil parameter. Failing to do so will
// leave pub/sub code and network operations initialized.
//
// If an error is encountered during initialization, it is returned and
// the caller should not attempt to use the Publisher function.
func GetPublisher(gcpName, topicName string) (Publisher, error) {
	logger.TRACE.Printf(`gcppubsubdemo.GetPublisher("%s", "%s")`, gcpName, topicName)

	var client *pubsub.Client
	var topic *pubsub.Topic
	var err error
	var skipDefer bool = false

	logger.DEBUG.Print("Starting pub/sub publisher")

	defer func() {
		if !skipDefer {
			closePubsub(client, topic)
		}
	}()

	if client, err = getPubsubClient(gcpName); err != nil {
		return nil, fmt.Errorf("Error getting pubsub client: %w", err)
	}

	logger.DEBUG.Print("Pub/sub client created")

	if topic, err = getOrCreatePubsubTopic(client, topicName); err != nil {
		return nil, fmt.Errorf("Error setting up topic %s: %w", topicName, err)
	}

	skipDefer = true

	return func(data []byte) error {
			logger.TRACE.Print("gcppubsubdemo.Publisher([]byte])")

			if data != nil && topic != nil {
				if err = publishPubsubData(topic, data); err != nil {
					closePubsub(client, topic)
					return fmt.Errorf("Error publishing data: %w", err)
				}
			} else {
				closePubsub(client, topic)
			}

			return nil
		},
		nil
}
