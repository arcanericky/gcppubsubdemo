package gcppubsubdemo

import "fmt"

func ExampleSubscribe() {
	Subscribe("GCP Name", "Topic Name", "Subscription ID", func(data []byte) bool {
		fmt.Println("Received Data:", string(data))
		return true
	})
}

func ExamplePubsubReceiveCallback() {
	var callback PubsubReceiveCallback = func(data []byte) bool {
		fmt.Println("Received Data:", string(data))
		return true
	}

	callback([]byte("data"))
}
