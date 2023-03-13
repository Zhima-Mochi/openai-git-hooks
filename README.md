# git-templates-hooks

## Git hook templates

This repository contains templates for git hooks. They are intended to be used with the [git hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) and openai features.

### Usage

To use these templates in your git repository, you need to add this repository as a submodule. Here are the steps:

1. In your local repository, navigate to the root directory.
2. Run the following command to add this repository as a submodule:

```bash
git submodule add https://github.com/Zhima-Mochi/git-templates-hooks.git .git-hooks
```
This command adds the git-templates-hooks repository as a submodule to your local repository, in a directory named .git-hooks. You can change the directory name to your liking.

3. Create the hooks directory in your repository:
```bash
mkdir .git/hooks
```

4. Create symbolic links to the hook templates:

```bash
cd .git/hooks
ln -s ../../.git-hooks/pre-commit
ln -s ../../.git-hooks/post-commit
ln -s ../../.git-hooks/prepare-commit-msg
ln -s ../../.git-hooks/commit-msg
ln -s ../../.git-hooks/post-checkout
ln -s ../../.git-hooks/post-merge
ln -s ../../.git-hooks/pre-push
ln -s ../../.git-hooks/pre-rebase
ln -s ../../.git-hooks/update
```

5. This creates symbolic links in the .git/hooks directory of your repository to the hook templates in the .git-hooks directory.

That's it! Now your repository is set up to use the hook templates from git-templates-hooks.
