package operations

import (
	"fmt"
	"os"
	"regexp"

	"github.com/charmbracelet/huh"
	"golang.org/x/term"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// validNamePattern validates template names: lowercase alphanumeric, starting with a-z or 0-9.
var validNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// supportedLanguages are the languages supported by the init scaffolding.
var supportedLanguages = []string{"go", "typescript", "python"}

// InitInfo holds parameters for initializing a template project.
type InitInfo struct {
	Name     string // Template project name
	Language string // Programming language
	Path     string // Output directory (defaults to ./<name>)
}

// Init initializes a new template project with scaffolded files.
// When parameters are not provided, uses interactive prompts.
func Init(info InitInfo) {
	name := info.Name
	language := info.Language
	path := info.Path

	// Interactive prompts if args are missing
	if name == "" || language == "" {
		// Require TTY for interactive mode
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			sbClient.PrintError("--name and --language are required in non-interactive mode")
			return
		}

		var fields []huh.Field

		if name == "" {
			fields = append(fields,
				huh.NewInput().
					Title("Template name").
					Description("Lowercase alphanumeric, hyphens and underscores allowed").
					Value(&name).
					Validate(func(s string) error {
						if !validNamePattern.MatchString(s) {
							return fmt.Errorf("name must match pattern: [a-z0-9][a-z0-9_-]*")
						}
						return nil
					}),
			)
		}

		if language == "" {
			langOptions := make([]huh.Option[string], 0, len(supportedLanguages))
			for _, lang := range supportedLanguages {
				langOptions = append(langOptions, huh.NewOption(lang, lang))
			}
			fields = append(fields,
				huh.NewSelect[string]().
					Title("Programming language").
					Options(langOptions...).
					Value(&language),
			)
		}

		if len(fields) > 0 {
			form := huh.NewForm(huh.NewGroup(fields...))
			if fErr := form.Run(); fErr != nil {
				sbClient.PrintError("cancelled: %v", fErr)
				return
			}
		}
	}

	// Validate name
	if !validNamePattern.MatchString(name) {
		sbClient.PrintError("invalid template name %q (must match: [a-z0-9][a-z0-9_-]*)", name)
		return
	}

	// Validate language
	validLang := false
	for _, l := range supportedLanguages {
		if language == l {
			validLang = true
			break
		}
	}
	if !validLang {
		sbClient.PrintError("unsupported language %q (supported: go, typescript, python)", language)
		return
	}

	if path == "" {
		path = "./" + name
	}

	fmt.Printf("Initializing %s template %q in %s...\n", language, name, path)
	if err := scaffold(name, language, path); err != nil {
		sbClient.PrintError("scaffold failed: %v", err)
		return
	}
	sbClient.PrintSuccess("Template %s initialized successfully!", name)
}
