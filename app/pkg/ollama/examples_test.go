package ai_test

import (
	"fmt"
	"log"

	ollama "github.com/OlegLaban/billAI-AI-service/app/pkg/ollama"
)

func ExampleConfig_withTemperature() {
	cfg := ollama.Config{Model: "llama3.2"}.WithTemperature(0.7)
	fmt.Println(cfg.Options["temperature"])
	// Output: 0.7
}

func ExampleConfig_withOption() {
	cfg := ollama.Config{Model: "llama3.2"}.
		WithOption("top_p", 0.9).
		WithOption("num_ctx", 4096)
	fmt.Println(cfg.Options["top_p"], cfg.Options["num_ctx"])
	// Output: 0.9 4096
}

func ExampleNew_validation() {
	_, err := ollama.New(ollama.Config{})
	fmt.Println(err != nil)
	// Output: true
}

func ExampleMessage_roles() {
	msgs := []ollama.Message{
		{Role: ollama.RoleUser, Content: "Hello"},
		{Role: ollama.RoleAssistant, Content: "Hi!"},
	}
	fmt.Println(msgs[0].Role, msgs[1].Role)
	// Output: user assistant
}

// This example shows the recommended client setup. It requires a running Ollama server.
func ExampleNew() {
	client, err := ollama.New(ollama.Config{
		Model:  "llama3.2",
		System: "Answer briefly.",
	}.WithTemperature(0.7))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.Model())
}
