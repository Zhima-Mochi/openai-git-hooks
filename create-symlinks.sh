chmod +x $(pwd)/.openai-git-hooks/prepare-commit-msg
ln -s $(pwd)/.openai-git-hooks/prepare-commit-msg $(pwd)/.git/hooks/prepare-commit-msg
ln -s $(pwd)/.openai-git-hooks/.diffignore $(pwd)/.git/hooks/.diffignore