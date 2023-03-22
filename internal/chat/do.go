package chat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	CmdExist = "exit"

	EnvKey     = "OPENAI_API_KEY" // sk-xxx
	EnvBaseURL = "OPENAI_BASEURL" // https://api.openai.com/v1
)

func Do() {
	if len(os.Getenv(EnvKey)) == 0 || len(os.Getenv(EnvBaseURL)) == 0 {
		fmt.Println("Set env OPENAI_API_KEY && OPENAI_BASEURL first.")
		return
	}

	var ctx = context.Background()
	var session = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant.",
		},
	}

	for {
		var input string
		if input = readInput(); len(input) == 0 {
			fmt.Println("empty input!")
			continue
		}

		session = append(session, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		errChan, respChan := chat(ctx, &session)
		response := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: "",
		}
		fmt.Print("[Assistant]")

		timer := time.NewTimer(time.Second)
	loop:
		for {
			select {
			case err, ok := <-errChan:
				if ok {
					fmt.Println(err)
					return
				}
				break loop

			case resp, ok := <-respChan:
				if !ok {
					break loop
				}
				if len(resp.Choices) == 0 {
					fmt.Println("No response. End of session.")
					break loop
				}

				timer.Stop()
				msg := resp.Choices[0].Delta
				fmt.Print(msg.Content)
				response.Content += msg.Content
			case <-timer.C:
				fmt.Print(".") // waiting response.
				timer.Reset(time.Second)
			}
		}
		fmt.Println()
		session = append(session, response)
	}

}

func readInput() string {
	fmt.Print("[User]")
	reader := bufio.NewReader(os.Stdin)
	var allStr string
	for {
		str, _ := reader.ReadString('\n')
		str = strings.TrimSpace(str)

		if len(str) == 0 {
			break
		}
		allStr = allStr + str
	}

	if strings.EqualFold(allStr, CmdExist) {
		os.Exit(0)
	}
	return allStr
}

func chat(ctx context.Context, session *[]openai.ChatCompletionMessage) (
	chan error, chan openai.ChatCompletionStreamResponse) {

	client := openai.NewClientWithConfig(newConfig())

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		// MaxTokens: 20,
		Messages: *session,
		Stream:   true,
	}

	var errChan = make(chan error)
	var respChan = make(chan openai.ChatCompletionStreamResponse)

	go stream(ctx, client, req, errChan, respChan)
	return errChan, respChan
}

func newConfig() openai.ClientConfig {
	cfg := openai.DefaultConfig(os.Getenv(EnvKey))
	cfg.BaseURL = os.Getenv(EnvBaseURL)
	return cfg
}

func stream(ctx context.Context, client *openai.Client, req openai.ChatCompletionRequest,
	errChan chan error, respChan chan openai.ChatCompletionStreamResponse) {

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		errChan <- errors.New(fmt.Sprintf("ChatCompletionStream error: %v\n", err))
		close(errChan)
		return
	}
	defer stream.Close()

	for {
		response, respError := stream.Recv()
		if errors.Is(respError, io.EOF) {
			close(respChan)
			close(errChan)
			return
		}

		if respError != nil {
			errChan <- errors.New(fmt.Sprintf("\nStream error: %v\n", err))
			close(respChan)
			close(errChan)
			return
		}

		respChan <- response
	}
}
