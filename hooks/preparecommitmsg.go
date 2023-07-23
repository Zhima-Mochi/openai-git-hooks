package hooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Zhima-Mochi/openai-git-hooks/handler"
	"github.com/sashabaranov/go-openai"
)

var (
	excludeFilePatterns = []string{
		"*.sum",
		"*.md",
		"*.txt",
		"*.json",
		"*.yaml",
		"*.yml",
	}
)

func PrepareCommitMsg(r io.Reader, w io.Writer, handler *handler.Handler, args ...string) error {
	if len(args) >= 2 && args[1] == "message" {
		return nil
	}

	commitMsg, err := getCommitMessage(r, w, handler)
	if err != nil {
		return fmt.Errorf("getCommitMessage failed: %w", err)
	}

	// write commit message to file
	commitEditMsgPath := args[0]
	file, err := os.OpenFile(commitEditMsgPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile failed: %w", err)
	}
	defer file.Close()
	file.WriteString(commitMsg)

	// open editor
	editor := strings.Split(os.Getenv("EDITOR"), " ")[0]
	args = strings.Split(os.Getenv("EDITOR"), " ")[1:]
	if editor == "" {
		editor = "vim"
	}
	args = append(args, commitEditMsgPath)
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCommitMessage(r io.Reader, w io.Writer, handler *handler.Handler) (string, error) {
	args := []string{"diff", "--cached", "--ignore-all-space", "--ignore-space-change", "--ignore-submodules", "--word-diff", "--word-diff-regex=.", "HEAD", "--"}
	for _, pattern := range excludeFilePatterns {
		args = append(args, fmt.Sprintf(":(exclude)%s", pattern))
	}
	diff, err := getDiff(args...)
	if err != nil {
		return "", fmt.Errorf("getDiff failed: %w", err)
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo0613,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Compose commit message: 50-char summary & 200-word description based on diff user offers.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: diff,
			},
		},
	}

	resp, err := handler.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("handler.CreateChatCompletion failed: %w", err)
	}
	fmt.Fprintf(w, "Prompt Tokens: %d, Completion Tokens: %d\n", resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
	fmt.Fprintf(w, "Cost: %f $NTD\n", (float64(resp.Usage.PromptTokens)/1000*0.0015+float64(resp.Usage.CompletionTokens)/1000*0.002)*30)
	return resp.Choices[0].Message.Content, nil
}

func getDiff(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd.Run failed: %w", err)
	}
	content := out.String()
	return content, nil
}
