#  OpenAI Git Hooks
## Git hooks with OpenAI

This repository contains prepare-commit-msg template for git hooks that use OpenAI's GPT-3.5 API to generate commit messages.

## Usage
1. In your local repository, navigate to the root directory.
2. Run the following command to add this repository as a submodule:

```bash
git submodule add https://github.com/Zhima-Mochi/git-templates-hooks.git .openai-git-hooks
```
This command adds the openai-git-hooks repository as a submodule to your repository. The .openai-git-hooks directory is created in your repository and contains the hook templates.

3. Create the hooks directory in your repository:
```bash
mkdir .git/hooks
```

4. Create symbolic links to the hook templates:

```bash
cd .git/hooks
ln -s ../../.openai-git-hooks/prepare-commit-msg
```

5. This creates symbolic links in the .git/hooks directory of your repository to the hook templates in the .openai-git-hooks directory.

6. execute the following command to set the openai api key as an environment variable:

```bash
export OPENAI_API_KEY=<your OpenAI API key>
```

Now, when you commit a change, the prepare-commit-msg hook is executed. If there is no commit message, the hook generates a commit message using OpenAI's GPT-3.5 API.