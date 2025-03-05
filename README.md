# Annealing

Annealing is a smart build tool that automatically detects and builds services based on Git changes. It helps streamline the development workflow by only building services that have been modified.

## Features

- Automatically detects changed files using Git
- Builds only the services affected by changes
- Supports multiple build commands per service
- Concurrent build execution
- YAML-based configuration

## Installation

```bash
go install github.com/pyljain/annealing@latest
```

## Usage

1. Create an `annealing.yaml` configuration file in your project root:

```yaml
spec:
  services:
    - path: "service1/"
      commands:
        - "go build"
        - "docker build -t service1 ."
    - path: "service2/"
      commands:
        - "npm install"
        - "npm run build"
```

2. Run Annealing:

```bash
annealing --config annealing.yaml
```

By default, Annealing will look for `annealing.yaml` in the current directory if no config file is specified.

## Configuration

The `annealing.yaml` file defines the services and their build commands:

- `path`: The directory path of the service relative to the project root
- `commands`: List of shell commands to execute when building the service

Example configuration:

```yaml
spec:
  services:
    - path: "backend/"
      commands:
        - "go mod tidy"
        - "go build"
    - path: "frontend/"
      commands:
        - "npm install"
        - "npm run build"
```

## How It Works

1. Annealing runs `git diff` to detect changed files
2. It matches the changed files against service paths defined in the configuration
3. For each affected service, it executes the specified build commands in the service's directory
4. All builds are executed concurrently using Go's error group functionality

## Error Handling

- If a build command fails, Annealing will stop all ongoing builds and exit with an error
- Configuration errors and Git command failures are reported with appropriate error messages

## Requirements

- Go 1.16 or later
- Git installed and available in PATH
- Bash shell