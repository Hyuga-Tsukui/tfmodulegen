package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Provider is a Terraform provider configuration.
type Provider struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Version string `json:"version"`
}

// ModuleData is the data used to generate the module files.
type ModuleData struct {
	ModuleName       string
	Description      string
	TerraformVersion string
	Providers        []Provider
}

// Config is the configuration for the module generator.
// This is loaded from a JSON file (tfmodulegen.config.json).
type Config struct {
	TerraformVersion string     `json:"terraform_version"`
	Providers        []Provider `json:"providers"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var config Config
	configPath := "tfmodulegen.config.json"
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println("Error opening config file:", err)
		} else {
			defer file.Close()
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&config); err != nil {
				fmt.Println("Error decoding config file:", err)
			} else {
				fmt.Println("Loaded configuration from", configPath)
			}
		}
	}

	// Start the interactive module generation process.

	// Input module name.
	fmt.Print("Enter module name: ")
	moduleName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading module name:", err)
		return
	}
	moduleName = strings.TrimSpace(moduleName)

	// Input module description.
	fmt.Print("Enter module description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading description:", err)
		return
	}
	description = strings.TrimSpace(description)

	// Input Terraform version.
	// if Terraform version is not set in the config, use a default value.
	defaultTFVersion := ">= 0.12"
	if config.TerraformVersion != "" {
		defaultTFVersion = config.TerraformVersion
		fmt.Printf("Enter required Terraform version (default from config: %s): ", defaultTFVersion)
	} else {
		fmt.Printf("Enter required Terraform version (default: %s): ", defaultTFVersion)
	}
	tfVersion, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading Terraform version:", err)
		return
	}
	tfVersion = strings.TrimSpace(tfVersion)
	if tfVersion == "" {
		tfVersion = defaultTFVersion
	}

	// Input providers.
	var providers []Provider
	if len(config.Providers) > 0 {
		// not empty, use the providers from the config file.
		// not override the providers from the config file.
		fmt.Println("Using provider configuration from config file:")
		for _, p := range config.Providers {
			fmt.Printf("  - %s: source=%s, version=%s\n", p.Name, p.Source, p.Version)
		}
		providers = config.Providers
	} else {
		// empty, ask the user to input providers.
		for {
			fmt.Print("Do you want to add a provider? (y/n): ")
			ans, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
			ans = strings.ToLower(strings.TrimSpace(ans))
			if ans != "y" && ans != "yes" {
				break
			}

			fmt.Print("Enter provider name (e.g. google): ")
			providerName, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading provider name:", err)
				return
			}
			providerName = strings.TrimSpace(providerName)

			fmt.Print("Enter provider source (e.g. hashicorp/google): ")
			providerSource, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading provider source:", err)
				return
			}
			providerSource = strings.TrimSpace(providerSource)

			fmt.Print("Enter provider version (e.g. 6.4.0): ")
			providerVersion, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading provider version:", err)
				return
			}
			providerVersion = strings.TrimSpace(providerVersion)

			providers = append(providers, Provider{
				Name:    providerName,
				Source:  providerSource,
				Version: providerVersion,
			})
		}
	}

	data := ModuleData{
		ModuleName:       moduleName,
		Description:      description,
		TerraformVersion: tfVersion,
		Providers:        providers,
	}

	// Create the module directory.
	if err := os.Mkdir(moduleName, 0755); err != nil {
		if !os.IsExist(err) {
			fmt.Println("Error creating directory:", err)
			return
		} else {
			fmt.Printf("Directory '%s' already exists. Files will be overwritten if they exist.\n", moduleName)
		}
	}

	files := []struct {
		filename string
		tmpl     string
	}{
		{"versions.tf", versionsTemplate},
		{"main.tf", mainTemplate},
		{"output.tf", outputTemplate},
		{"variable.tf", variableTemplate},
		{"README.md", readmeTemplate},
	}

	for _, f := range files {
		if err := generateFile(moduleName, f.filename, f.tmpl, data); err != nil {
			fmt.Printf("Failed to generate %s: %v\n", f.filename, err)
			return
		}
	}

	fmt.Printf("Terraform module boilerplate files generated successfully in the '%s' directory!\n", moduleName)
}

// generateFile generates a file with the given template and data.
func generateFile(dirName, filename, tmplStr string, data ModuleData) error {
	filePath := filepath.Join(dirName, filename)
	tmpl, err := template.New(filename).Funcs(template.FuncMap{
		"codeFence": func() string { return "```" },
	}).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

const versionsTemplate = `terraform {
  required_version = "{{.TerraformVersion}}"
  {{- if .Providers }}
  required_providers {
  {{- range .Providers }}
    {{ .Name }} = {
      source  = "{{ .Source }}"
      version = "{{ .Version }}"
    }
  {{- end }}
  }
  {{- end }}
}
`

const mainTemplate = `// Main configuration for module {{.ModuleName}}
resource "example_resource" "default" {
  provisioner "local-exec" {
    command = "echo Hello from module {{.ModuleName}}!"
  }
}
`

const outputTemplate = `// Outputs for module {{.ModuleName}}
output "example" {
  description = "An example output"
  value       = "example_value"
}
`

const variableTemplate = `// Variables for module {{.ModuleName}}
variable "example_variable" {
  description = "An example variable"
  type        = string
  default     = "default_value"
}
`

const readmeTemplate = `# {{.ModuleName}}
{{.Description}}

This Terraform module is automatically generated.

## Requirements

- Terraform version {{.TerraformVersion}}

{{- if .Providers }}
## Providers
{{- range .Providers }}
- **{{ .Name }}**: source={{ .Source }}, version={{ .Version }}
{{- end }}
{{- end }}

## Usage

{{codeFence}}hcl
module "{{.ModuleName}}" {
  source = "./{{.ModuleName}}"
  # ... module inputs
}
{{codeFence}}
`
