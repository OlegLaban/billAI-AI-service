package ai

import "github.com/ollama/ollama/api"

// Role identifies the author of a chat message.
type Role string

const (
	RoleSystem    Role = "system"    // System instruction.
	RoleUser      Role = "user"      // End-user input.
	RoleAssistant Role = "assistant" // Model response.
)

// Message is one turn in a chat history passed to [Client.Chat] or stored in [Conversation].
type Message struct {
	Role    Role
	Content string
}

func toAPIMessages(msgs []Message) []api.Message {
	out := make([]api.Message, len(msgs))
	for i, m := range msgs {
		out[i] = api.Message{Role: string(m.Role), Content: m.Content}
	}
	return out
}
