# pman
pman is a command line tool to keep track of all your side projects.


## Why?
I needed something to keep track of all my side projects.

## Install using the go package manager

```
go install github.com/theredditbandit/pman@latest
```

## Usage 

```
A cli project manager

Usage:
  pman [flags]
  pman [command]

Available Commands:
  add         Adds a project directory to the index
  alias       Sets the alias for a project , whose name might be too big
  completion  Generate the autocompletion script for the specified shell
  delete      Deletes a project from the index database. This does not delete the project from the filesystem
  help        Help about any command
  i           Launches pman in interactive mode
  info        The info command pretty prints the README.md file present at the root of the specified project.
  init        Takes exactly 1 argument, a directory name, and initializes it as a project directory.
  ls          List all indexed projects along with their status
  reset       Deletes the current indexed projects , run pman init to reindex the projects
  set         Set the status of a project
  status      Get the status of a project

Flags:
  -h, --help      help for pman
  -v, --version   version for pman

Use "pman [command] --help" for more information about a command.
```

## How does it work.
when you run `pman init .` in any directory , it will look for subdirectories that contain a README.md or a .git folder and consider it as a project director.


