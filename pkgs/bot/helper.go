package bot

import (
	"fmt"
	"log"
)

func (c *Context) Reply(message interface{}) error {
	api := c.GetAPIClient()
	if api == nil {
		return fmt.Errorf("API client not found in context")
	}
	
	if c.Event.IsGroupMessage() {
		_, err := api.SendGroupMessage(c.Event.GroupID, message)
		return err
	}
	
	if c.Event.IsPrivateMessage() {
		_, err := api.SendPrivateMessage(c.Event.UserID, message)
		return err
	}
	
	return fmt.Errorf("unsupported message type for reply")
}

func (c *Context) ReplyText(text string) error {
	return c.Reply(text)
}

func (c *Context) ReplyAt(text string) error {
	if !c.Event.IsGroupMessage() {
		return c.ReplyText(text)
	}
	
	atSegment := MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": fmt.Sprintf("%d", c.Event.UserID),
		},
	}
	
	textSegment := MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": " " + text,
		},
	}
	
	return c.Reply([]MessageSegment{atSegment, textSegment})
}

func (c *Context) GetAPIClient() APIClient {
	if client, exists := c.Get("api_client"); exists {
		if api, ok := client.(APIClient); ok {
			return api
		}
	}
	return nil
}

func InjectAPIClient(api APIClient) HandlerFunc {
	return func(ctx *Context) {
		ctx.Set("api_client", api)
		ctx.Next()
	}
}

func OnCommand(prefix string, command string, handler HandlerFunc) HandlerFunc {
	fullCommand := prefix + command
	
	return func(ctx *Context) {
		text := ctx.Event.GetText()
		if len(text) < len(fullCommand) {
			return
		}
		
		if text[:len(fullCommand)] == fullCommand {
			handler(ctx)
		}
	}
}

func OnKeyword(keyword string, handler HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		text := ctx.Event.GetText()
		if contains(text, keyword) {
			handler(ctx)
		}
	}
}

func OnPrefix(prefix string, handler HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		text := ctx.Event.GetText()
		if len(text) >= len(prefix) && text[:len(prefix)] == prefix {
			handler(ctx)
		}
	}
}

func OnRegex(pattern string, handler HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		handler(ctx)
	}
}

func LogError(err error, context string) {
	if err != nil {
		log.Printf("Error [%s]: %v", context, err)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
