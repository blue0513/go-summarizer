package openai

import (
	"bufio"
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

// defaults to os.LookupEnv("OPENAI_API_KEY")
var client = openai.NewClient()

func Summarize(ctx context.Context, text string, lang string) (string, error) {
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

func Ask(ctx context.Context, page string) error {
	replies := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(fmt.Sprintf("Answer the following questions about the article\n%s", page)),
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("---")
		fmt.Println("Ask question:")

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		fmt.Println("---")

		res, err := completion(ctx, input, replies)
		if err != nil {
			return fmt.Errorf("failed to get completion: %w", err)
		}
		replies = append(replies, res)

		fmt.Println("")
	}
}

func completion(
	ctx context.Context,
	input string,
	replies []openai.ChatCompletionMessageParamUnion,
) (openai.ChatCompletionMessage, error) {
	replies = append(replies, openai.UserMessage(input))
	params := openai.ChatCompletionNewParams{
		Messages: openai.F(replies),
		Model:    openai.F(model),
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		// it's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		return openai.ChatCompletionMessage{}, fmt.Errorf("stream error: %w", err)
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	return acc.Choices[0].Message, nil
}
