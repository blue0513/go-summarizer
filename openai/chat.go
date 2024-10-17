package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
)

const defaultModel = openai.ChatModelGPT4oMini

var model = func() string {
	model := os.Getenv("OPENAI_MODEL")
	if model != "" {
		return model
	}

	return defaultModel
}()

func Summarize(ctx context.Context, text string, lang string) (string, error) {
	// defaults to os.LookupEnv("OPENAI_API_KEY")
	client := openai.NewClient()
	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.F(model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fmt.Sprintf("I need a summary of this text in %s:\n%s", lang, text)),
		}),
	})
	if err != nil {
		return "", fmt.Errorf("openai chat completion: %w", err)
	}

	return res.Choices[0].Message.Content, nil
}
