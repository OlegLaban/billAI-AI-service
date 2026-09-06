// Package ai provides a Go SDK for the [Ollama] LLM server.
//
// The package lives at pkg/ollama but is declared as package ai.
// Import it with an alias:
//
//	import ollama "github.com/OlegLaban/billAI-AI-service/app/pkg/ollama"
//
// Under the hood it uses the official client [github.com/ollama/ollama/api].
//
// # Interaction modes
//
// Three modes cover most use cases:
//
//   - [Client.Prompt] and [Client.PromptStream] — single-turn text generation
//     via the Ollama Generate API. System instructions come from [Config.System].
//   - [Client.Chat] and [Client.ChatStream] — multi-message chat via the Chat API.
//     The caller supplies message history; [Config.System] is prepended automatically.
//   - [Conversation] — stateful multi-turn dialog. History is kept inside the object
//     and sent on every [Conversation.Send].
//
// # Quick start
//
//	client, err := ollama.New(ollama.Config{
//	    Model:  "llama3.2",
//	    System: "Be concise.",
//	}.WithTemperature(0.7))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if err := client.Ping(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	answer, err := client.Prompt(ctx, "How many planets are in the Solar System?")
//
// # Configuration
//
// [Config.Model] is required. [Config.Host] overrides the OLLAMA_HOST environment
// variable (default http://127.0.0.1:11434). Model parameters such as temperature
// are passed through [Config.Options] or helper methods like [Config.WithTemperature].
//
// # Context and errors
//
// All request methods accept [context.Context] for cancellation and timeouts.
// Errors from Ollama (unreachable server, missing model, etc.) are returned as-is.
// Streaming callbacks can abort generation by returning a non-nil error.
//
// [Ollama]: https://ollama.com
package ai
