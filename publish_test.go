package gcppubsubdemo

func ExampleGetPublisher() {
	publish, _ := GetPublisher("GCP Name", "Topic Name")
	publish([]byte("Payload"))
}
