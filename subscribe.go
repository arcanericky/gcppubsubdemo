package gcppubsubdemo

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// PubsubReceiveCallback is called with the data returned from the
// subscriber as a parameter. If the implemented callback returns true
// the subscription will be deleted and Subscribe() will return.
type PubsubReceiveCallback func([]byte) bool

func getPubsubSubscription(client *pubsub.Client, topic *pubsub.Topic, subscriptionID string) (*pubsub.Subscription, error) {
	logger.TRACE.Printf(`gcppubsubdemo.getPubsubSubscription(pubsubClient, "%s", "%s")`, topic.ID(), subscriptionID)

	var sub *pubsub.Subscription
	var err error

	ctx := context.Background()

	// Because CreateSubscription() can return an error if the subscription already exists,
	// there is a race condition between the time Exists() and CreateSubscription() is
	// called. This loop accounts for this by following the call to CreateSubscription()
	// with another call to Exist() to determine if CreateSubscription() failed because
	// this race condition condition occurred. This is only a concern if there are multiple
	// processes or goroutines than can create the same subscription.

	// Simplification of this logic would be a welcome code change.
	attempts := 0
	for attempts <= 1 {
		sub = client.Subscription(subscriptionID)

		logger.DEBUG.Printf("Checking for existing %s subscription", sub.ID())
		if ok, err := sub.Exists(ctx); err != nil {
			logger.ERROR.Printf("Error determining existence of subscription %s: %s", sub.ID(), err)
			return nil, err
		} else if ok {
			break
		}

		logger.DEBUG.Printf("Subscription %s does not exist. Creating it.", subscriptionID)

		if attempts < 1 {
			if sub, err = client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{Topic: topic}); err != nil {
				logger.ERROR.Println("Error creating subscription:", err)
			}
		} else {
			errText := "Could not create subscription " + subscriptionID
			logger.ERROR.Print(errText)
			return nil, errors.New(errText)
		}

		attempts++
	}

	return sub, nil
}

func receivePubsubData(sub *pubsub.Subscription, callback PubsubReceiveCallback) error {
	logger.TRACE.Printf(`gcppubsubdemo.receivePubsubData("%s", callback)`, sub.ID())

	cctx, cancel := context.WithCancel(context.Background())

	err := sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		logger.DEBUG.Printf("Received message: %s", m.Data)
		if callback(m.Data) {
			m.Ack()
			cancel()
			return
		}

		m.Ack()
	})

	if err != nil {
		logger.ERROR.Print("Error receiving from subscription:", err)
	}

	return err
}

// Subscribe subscribes and retrieves pub/sub messages from a given Google Cloud
// Project and topic then calls the callback so the caller can process each
// message. The callback should return true when it no longer wants to retrieve
// messages and shut down the subscription.
//
// If an error is encountered during initialization for the subscription,
// it is returned.
func Subscribe(gcpName, topicName, subscriptionID string, callback PubsubReceiveCallback) error {
	logger.TRACE.Printf(`gcppubsubdemo.Subscribe("%s", "%s", "%s", callback)`, gcpName, topicName, subscriptionID)

	var err error

	logger.DEBUG.Print("Starting pub/sub subscriber")

	var client *pubsub.Client
	if client, err = getPubsubClient(gcpName); err != nil {
		return fmt.Errorf("Error getting pub/sub client: %w", err)
	}
	defer client.Close()

	logger.DEBUG.Print("Pub/Sub client created")

	var topic *pubsub.Topic
	if topic, err = getOrCreatePubsubTopic(client, topicName); err != nil {
		return fmt.Errorf("Error setting up topic %s: %w", topicName, err)
	}

	logger.DEBUG.Println("Topic retrieved")

	var sub *pubsub.Subscription
	if sub, err = getPubsubSubscription(client, topic, subscriptionID); err != nil {
		return fmt.Errorf("Error getting pub/sub subscription: %w", err)
	}
	defer sub.Delete(context.Background())

	logger.DEBUG.Print("Subscription created, starting receiver")

	if err = receivePubsubData(sub, callback); err != nil {
		return fmt.Errorf("Error receiving pub/sub data: %w", err)
	}

	return nil
}
