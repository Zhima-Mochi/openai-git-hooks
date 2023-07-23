package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Zhima-Mochi/openai-git-hooks/handler"
	"github.com/Zhima-Mochi/openai-git-hooks/hooks"
	"golang.org/x/term"
)

var (
	availableHooks = map[string]bool{
		"prepare-commit-msg": true,
		"pre-commit":         false,
		"commit-msg":         false,
		"post-commit":        false,
		"pre-rebase":         false,
		"post-rewrite":       false,
		"post-checkout":      false,
		"post-merge":         false,
		"pre-push":           false,
	}

	hooksFunc = map[string]func(r io.Reader, w io.Writer, handler *handler.Handler, args ...string) error{
		"prepare-commit-msg": hooks.PrepareCommitMsg,
	}

	ErrNoOpenAIKey = fmt.Errorf("no OpenAI key found")
)

type config struct {
	hookName string
	args     []string
}

func validateArgs(c config) error {
	if _, ok := availableHooks[c.hookName]; !ok {
		return fmt.Errorf("invalid hook name: %s", c.hookName)
	}
	return nil
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("openai-git-hook", flag.ExitOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		var usageString = `
openai-git-hook helps to execute hooks with the additional power of OpenAI's GPT-3.5.

Usage: %s [hook-name] [options]
`
		fmt.Fprintf(w, usageString, fs.Name())
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if fs.NArg() == 0 {
		return c, fmt.Errorf("invalid number of arguments")
	}
	c.hookName = fs.Arg(0)
	c.args = fs.Args()[1:]

	return c, nil
}

func getOpenAIKeyFromFile() (string, error) {
	filePath := filepath.Join(os.Getenv("HOME"), ".openai-git-hook")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", ErrNoOpenAIKey
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("os.Open failed: %w", err)
	}
	defer file.Close()

	var key string
	_, err = fmt.Fscanln(file, &key)
	if err != nil {
		return "", fmt.Errorf("fmt.Fscanln failed: %w", err)
	}
	if key == "" {
		return "", ErrNoOpenAIKey
	}
	return key, nil
}

func writeOpenAIKeyToFile(key string) error {
	filePath := filepath.Join(os.Getenv("HOME"), ".openai-git-hook")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("os.Create failed: %w", err)
	}
	defer file.Close()
	_, err = fmt.Fprintln(file, key)
	if err != nil {
		return fmt.Errorf("fmt.Fprintln failed: %w", err)
	}
	return nil
}

func getOpenAIKey(r io.Reader, w io.Writer) (string, error) {
	if key, err := getOpenAIKeyFromFile(); err == nil {
		return key, nil
	} else if err != ErrNoOpenAIKey {
		return "", fmt.Errorf("getOpenAIKeyFromFile failed: %w", err)
	}

	fmt.Fprintln(w, "Please enter your OpenAI API key:")
	var key string
	if r == os.Stdin {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("term.ReadPassword failed: %w", err)
		}
		key = string(bytePassword)
	} else {
		scanner := bufio.NewScanner(r)
		scanner.Scan()
		key = scanner.Text()
	}

	if key == "" {
		return "", ErrNoOpenAIKey
	}

	if err := writeOpenAIKeyToFile(key); err != nil {
		return "", fmt.Errorf("writeOpenAIKeyToFile failed: %w", err)
	}

	return key, nil
}

func runCmd(r io.Reader, w io.Writer, c config) error {
	key, err := getOpenAIKey(r, w)
	if err != nil {
		return fmt.Errorf("getOpenAIKey failed: %w", err)
	}
	handler := handler.NewHandler(key)
	return hooksFunc[c.hookName](r, w, handler, c.args...)
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := validateArgs(c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := runCmd(os.Stdin, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
