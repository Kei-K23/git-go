# Git-Go

**Git-Go** is a simplified and lightweight `Git` implementation written in **Go**. It replicates key features of `Git`, such as **creating repositories**, **adding files**, **committing changes**, and working with **branches**.

## Features

- Initialize a new Git repository (`init`)
- Add files to the staging area (`add`)
- Commit changes to the repository (`commit`)
- Show the list of commits (`log`)
- Show the list of staged files (`ls-files-stage`)
- List, create and delete branch (`branch`)

## Setup and Installation

1. Clone the repo and install modules

```bash
git clone https://github.com/Kei-K23/git-go.git
cd git-go
go mod tidy
```

2. Build the project

```bash
go build -o git-go
```

3. Test the binary

```bash
./git-go --help
```

Example output will be

```bash
git-go is a lightweight version of Git implemented in Go.

Usage:
  git-go [flags]
  git-go [command]

Available Commands:
  add            Add file contents to the index
  branch         List, create, or delete branches
  commit         Record changes to the repository
  completion     Generate the autocompletion script for the specified shell
  help           Help about any command
  init           Initialize a new Git repository
  log            Show commits log
  ls-files-stage Show information about files in staging area

Flags:
  -h, --help   help for git-go

Use "git-go [command] --help" for more information about a command.
```

## Example

1. **Init** (Initialize a new Git repository)

```bash
./git-go init
```

Example output will be

```bash
Initialized empty .git-go repository.
```

2. **Add** (Add files to the staging area)

```bash
./git-go add test.txt
```

Example output will be (hash value will be different depending on your file content)

```bash
Stored object as : .git-go/objects/98/9e5c2c4b1fc947fc3e35ea817eff7224113931
```

3. **Commit** (Commit changes to the repository)

```bash
./git-go commit -m "Add test commit"
```

Example output will be (hash value will be different depending on your file content)

```bash
[master c4ab8b32e630f904a8512b8858557c62ea6f1ed2] Add test commit
```

4. **Log** (Show the list of commits)

```bash
./git-go log
```

Example output will be (hash value will be different depending on your file content)

```bash
commit c4ab8b32e630f904a8512b8858557c62ea6f1ed2
Author author <author@gmail.com> 2024-09-26T18:09:28+06:30
Date <author@gmail.com> 2024-09-26T18:09:28+06:30
        Add test commit
```

5. **List-files-stage** (Show the list of staged files)

```bash
./git-go ls-files-stage
```

Example output will be (hash value will be different depending on your file content)

```bash
100644 989e5c2c4b1fc947fc3e35ea817eff7224113931 test.txt
```

6. **Branch** (List, create and delete branch)

```bash
./git-go branch dev

```

Example output will be

```bash
New branch call 'dev' is created
```

```bash
./git-go branch
```

Example output will be

```bash
dev
master
```

```bash
./git-go branch -d dev
```

Example output will be

```bash
Branch name 'dev' is deleted
```

## Contribution

All contributions are welcome. You can contribute to refactor the codebase, add more features, and fix the issues. Please open issues or make PR.

## License

This project is under [MIT LICENSE](/LICENSE).
