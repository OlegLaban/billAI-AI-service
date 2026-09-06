package ai

import (
	"context"
	"strings"

	"github.com/ollama/ollama/api"
)

// Chat sends messages via the Ollama Chat API and returns the assistant reply.
// When Config.System is set, it is prepended as the first system message.
func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	msgs := messages
	if c.cfg.System != "" {
		msgs = append([]Message{{Role: RoleSystem, Content: c.cfg.System}}, messages...)
	}

	stream := false
	req := &api.ChatRequest{
		Model:    c.cfg.Model,
		Messages: toAPIMessages(msgs),
		Stream:   &stream,
		Options:  c.cfg.Options,
	}

	var sb strings.Builder
	err := c.ollama.Chat(ctx, req, func(resp api.ChatResponse) error {
		sb.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

// ChatStream streams the Chat API response, calling onChunk for each non-empty fragment.
// Config.System is prepended the same way as in [Client.Chat].
func (c *Client) ChatStream(ctx context.Context, messages []Message, onChunk func(chunk string) error) error {
	msgs := messages
	if c.cfg.System != "" {
		msgs = append([]Message{{Role: RoleSystem, Content: c.cfg.System}}, messages...)
	}

	req := &api.ChatRequest{
		Model:    c.cfg.Model,
		Messages: toAPIMessages(msgs),
		Options:  c.cfg.Options,
	}

	return c.ollama.Chat(ctx, req, func(resp api.ChatResponse) error {
		if resp.Message.Content == "" {
			return nil
		}
		return onChunk(resp.Message.Content)
	})
}

// Conversation is a stateful multi-turn chat session.
//
// History is stored internally and sent to Ollama on every [Conversation.Send].
// On request failure the last user message is rolled back from history.
type Conversation struct {
	client   *Client
	messages []Message
}

// NewConversation starts a new session.
// Config.System, if set, becomes the first message in history.
func (c *Client) NewConversation() *Conversation {
	conv := &Conversation{client: c}
	if c.cfg.System != "" {
		conv.messages = []Message{{Role: RoleSystem, Content: c.cfg.System}}
	}
	return conv
}

// Send adds a user message, requests a reply, and appends the assistant response to history.
// On error the user message is removed from history before the error is returned.
func (conv *Conversation) Send(ctx context.Context, content string) (string, error) {
	conv.messages = append(conv.messages, Message{Role: RoleUser, Content: content})

	stream := false
	req := &api.ChatRequest{
		Model:    conv.client.cfg.Model,
		Messages: toAPIMessages(conv.messages),
		Stream:   &stream,
		Options:  conv.client.cfg.Options,
	}

	var sb strings.Builder
	err := conv.client.ollama.Chat(ctx, req, func(resp api.ChatResponse) error {
		sb.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		conv.messages = conv.messages[:len(conv.messages)-1]
		return "", err
	}

	reply := sb.String()
	conv.messages = append(conv.messages, Message{Role: RoleAssistant, Content: reply})
	return reply, nil
}

// History returns a snapshot of all messages in the session (including system, if any).
func (conv *Conversation) History() []Message {
	out := make([]Message, len(conv.messages))
	copy(out, conv.messages)
	return out
}

// Reset clears the conversation. Config.System is preserved as the only message when set.
func (conv *Conversation) Reset() {
	if conv.client.cfg.System != "" {
		conv.messages = []Message{{Role: RoleSystem, Content: conv.client.cfg.System}}
	} else {
		conv.messages = nil
	}
}
