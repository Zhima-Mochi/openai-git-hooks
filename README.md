# OpenAI Git Hooks

This repository contains code for an `openai-git-hook` tool, which helps execute Git hooks with the additional power of OpenAI's GPT-3.5. It allows you to leverage the OpenAI API to enhance your Git hooks using natural language processing capabilities.

## Usage

To use this tool, follow the steps below:

1. Install the package using `go install github.com/Zhima-Mochi/openai-git-hooks@latest`.
2. Make sure to add the `openai-git-hook [hookname] "$@"` command in your Git hooks file (e.g., .git/hooks/prepare-commit-msg) and make it executable.

## Hooks Supported

The following Git hooks are supported by the `openai-git-hook` tool:

- `prepare-commit-msg`

## Configuration

The `openai-git-hook` tool can be configured by creating a `.openai-git-hook` file in your home directory. The file should contain your OpenAI API key, which will be used to access the OpenAI API.

## Contributing

If you want to contribute to this project, feel free to fork the repository and submit a pull request with your changes. We welcome any improvements or new features.