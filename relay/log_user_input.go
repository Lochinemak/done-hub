package relay

import (
	"done-hub/common/config"
	claudeProvider "done-hub/providers/claude"
	"done-hub/types"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func setLogUserInput(c *gin.Context, userInput string) {
	c.Set(config.GinLogUserInputKey, strings.TrimSpace(userInput))
}

func extractLastUserMessage(messages []types.ChatCompletionMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role != types.ChatMessageRoleUser {
			continue
		}

		text := extractTextFromChatMessage(messages[i])
		if text != "" {
			return text
		}
	}

	return ""
}

func extractTextFromChatMessage(message types.ChatCompletionMessage) string {
	parts := message.ParseContent()
	if len(parts) == 0 {
		return strings.TrimSpace(message.StringContent())
	}

	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part.Type == types.ContentTypeText {
			text := strings.TrimSpace(part.Text)
			if text != "" {
				texts = append(texts, text)
			}
		}
	}

	return strings.Join(texts, "\n")
}

func extractCompletionPrompt(prompt any) string {
	switch value := prompt.(type) {
	case string:
		return strings.TrimSpace(value)
	case []string:
		return strings.TrimSpace(strings.Join(value, "\n"))
	case []any:
		texts := make([]string, 0, len(value))
		for _, item := range value {
			if text, ok := item.(string); ok {
				text = strings.TrimSpace(text)
				if text != "" {
					texts = append(texts, text)
				}
			}
		}
		return strings.TrimSpace(strings.Join(texts, "\n"))
	default:
		return ""
	}
}

func extractFirstStringValue(input any) string {
	switch value := input.(type) {
	case string:
		return strings.TrimSpace(value)
	case []string:
		for _, item := range value {
			item = strings.TrimSpace(item)
			if item != "" {
				return item
			}
		}
	case []any:
		for _, item := range value {
			if text, ok := item.(string); ok {
				text = strings.TrimSpace(text)
				if text != "" {
					return text
				}
			}
		}
	}

	return ""
}

func extractClaudeUserMessage(messages []claudeProvider.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role != types.ChatMessageRoleUser {
			continue
		}

		text := extractClaudeMessageText(messages[i].Content)
		if text != "" {
			return text
		}
	}

	return ""
}

func extractClaudeMessageText(content any) string {
	if text, ok := content.(string); ok {
		return strings.TrimSpace(text)
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return ""
	}

	var parts []claudeProvider.MessageContent
	if err := json.Unmarshal(contentBytes, &parts); err != nil {
		return ""
	}

	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part.Type != claudeProvider.ContentTypeText {
			continue
		}

		text := strings.TrimSpace(part.Text)
		if text != "" {
			texts = append(texts, text)
		}
	}

	return strings.Join(texts, "\n")
}

func extractLastGeminiUserText(requestBody []byte) string {
	if len(requestBody) == 0 {
		return ""
	}

	contents := gjson.GetBytes(requestBody, "contents")
	if !contents.Exists() || !contents.IsArray() {
		return ""
	}

	items := contents.Array()
	for i := len(items) - 1; i >= 0; i-- {
		role := strings.TrimSpace(items[i].Get("role").String())
		if role != "" && role != types.ChatMessageRoleUser {
			continue
		}

		parts := items[i].Get("parts").Array()
		texts := make([]string, 0, len(parts))
		for _, part := range parts {
			text := strings.TrimSpace(part.Get("text").String())
			if text != "" {
				texts = append(texts, text)
			}
		}

		if len(texts) > 0 {
			return strings.Join(texts, "\n")
		}
	}

	return ""
}
