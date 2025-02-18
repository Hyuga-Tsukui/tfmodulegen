# tfmodulegen

A CLI tool **tfmodulegen** that generates boilerplate files for Terraform modules.

## Overview

`tfmodulegen` is a tool that automatically generates basic Terraform module files (`versions.tf`, `main.tf`, `output.tf`, `variable.tf`, `README.md`) through an interactive prompt.  
Additionally, by using a JSON configuration file (`tfmodulegen.config.json`), you can predefine default settings for Terraform versions and providers, which can be overridden during the interactive prompt.

## Features

- **Interactive Prompt**  
  Configure module name, description, Terraform version, and provider information based on user input.

- **Automatic Directory Generation**  
  Automatically creates a directory with the same name as the input module name and outputs necessary files.

- **Default Values via Configuration File**  
  Use `tfmodulegen.config.json` to predefine default settings for Terraform versions and providers, which can be overridden during the interactive prompt.

## Requirements

- [Go](https://golang.org/) 1.16 or higher (recommended)
- Terraform (※Knowledge of Terraform module creation is required)

## Installation

1. Clone this repository or obtain the source code.

2. Run the following command in your terminal to build the binary:

   ```bash
   go build -o tfmodulegen main.go
   ```

This command will generate an executable tfmodulegen binary.

## Using the Configuration File

Optionally, you can use a configuration file named `tfmodulegen.config.json`.
This configuration file automatically loads default values for Terraform versions and provider information, which are displayed as initial values during the interactive prompt.

### Configuration File Example

```json
{
  "terraform_version": "~> 1.9.6",
  "providers": [
    {
      "name": "google",
      "source": "hashicorp/google",
      "version": "6.4.0"
    }
  ]
}
```

Place the file in the project root, and it will be automatically loaded when the tool starts.

## Usage

1. Run the `tfmodulegen` command in your terminal:

```bash
./tfmodulegen
```

2. Follow the interactive prompt to input the module name, description, Terraform version, and provider information.
3. Once input is complete, a directory with the same name as the module will be generated, containing all necessary files.

## Cautions

- Overwrite Warning

If a directory with the same name already exists, existing files may be overwritten.
Please backup as needed before execution.
