package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hay-kot/hookfeed/hookfeed"
	"github.com/hay-kot/hookfeed/internal/console"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type ValidateCmd struct {
	flags struct {
		input string
	}
}

func NewValidateCommand() *ValidateCmd {
	return &ValidateCmd{}
}

func (i *ValidateCmd) Register(app *cli.Command) *cli.Command {
	cmd := &cli.Command{
		Name:      "validate",
		Usage:     "Validate a Lua script by transforming sample JSON",
		UsageText: "hookfeed validate <script>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "Path to JSON input file (webhook simulation)",
				Required:    false,
				Destination: &i.flags.input,
			},
		},
		Action: i.validate,
	}

	app.Commands = append(app.Commands, cmd)
	return app
}

func (i *ValidateCmd) validate(ctx context.Context, cmd *cli.Command) error {
	const defaultWebhookJSON = `{
  "body": {
    "event": "example.event",
    "timestamp": "2025-10-16T00:00:00Z",
    "data": {
      "id": "12345",
      "status": "active"
    }
  }
}`

	scriptPath := cmd.Args().First()
	if scriptPath == "" {
		return fmt.Errorf("lua script path is required")
	}

	// Get input data - use file if provided, otherwise use default
	var inputData []byte
	var err error

	if inputPath := cmd.String("input"); inputPath != "" {
		inputData, err = os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		log.Info().Str("input", inputPath).Msg("using input file")
	} else {
		inputData = []byte(defaultWebhookJSON)
		log.Info().Msg("using default webhook JSON")
	}

	// Create transformer
	transformer := hookfeed.NewTransformer(scriptPath)

	// Transform the input
	output, err := transformer.Transform(inputData)
	if err != nil {
		return fmt.Errorf("transformation failed: %w", err)
	}

	// Print input and output with headers
	fmt.Println(console.SectionTitle("Input"))
	fmt.Print(console.PrettyJSON(inputData))
	fmt.Println(console.SectionTitle("Output"))
	fmt.Print(console.PrettyJSON(output))
	fmt.Println()

	log.Info().Msg("validation successful")
	return nil
}
