package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

type ReReplace struct {
	Regex *regexp.Regexp
	Text  string
}

func NewReReplace(regex string, replace string) ReReplace {
	return ReReplace{
		Regex: regexp.MustCompile(regex),
		Text:  replace,
	}
}

func NewReDate(dateStr string) ReReplace {
	return ReReplace{
		Regex: regexp.MustCompile(fmt.Sprintf(`%s: string`, dateStr)),
		Text:  fmt.Sprintf(`%s: Date | string`, dateStr),
	}
}

func GenerateTypescript(c *Config) error {
	for _, file := range c.Typescript.DataContracts {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Error().Msgf("File %s does not exist", file)
			return errors.New("file does not exist")
		}

		text := "/* post-processed by ./scripts/process-types.go */"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read file %s", file)
			return err
		}

		if strings.HasPrefix(string(data), text) {
			text = string(data)
		} else {
			text += "\n" + string(data)
		}

		replaces := [...]ReReplace{
			NewReReplace(` Dtos`, " "),
			NewReReplace(`\?:`, ":"),
			NewReReplace(`(\w+):\s(.*null.*)`, "$1?: $2"), // make null union types optional
			NewReDate("createdAt"),
			NewReDate("updatedAt"),
		}

		for _, replace := range replaces {
			log.Info().Msgf("Replacing '%v' -> '%s'", replace.Regex, replace.Text)
			text = replace.Regex.ReplaceAllString(text, replace.Text)
		}

		err = os.WriteFile(file, []byte(text), 0o644)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to write file %s\n", file)
			return err
		}
	}

	// ========================
	// Process swagger file

	f, err := os.ReadFile(c.Typescript.SwaggerFile)
	if err != nil {
		return err
	}

	swagger := map[string]any{}

	err = json.Unmarshal(f, &swagger)
	if err != nil {
		return err
	}

	pathsMp := swagger["paths"].(map[string]any)

	paths := make([]string, 0, len(pathsMp))
	for path := range pathsMp {
		paths = append(paths, path)
	}

	// Sort paths
	sort.Strings(paths)

	bldr := strings.Builder{}

	bldr.WriteString("export type RCPRoute =\n")

	for i, path := range paths {
		// Replace {.*} with {string}
		path = regexp.MustCompile(`{.*}`).ReplaceAllString(path, `{string}`)
		path = strings.Replace(path, "{", "${", 10)
		path += "/"
		path = strings.TrimPrefix(path, "/v1")
		bldr.WriteString(fmt.Sprintf("  | `%s`", path))

		if i+1 == len(paths) {
			bldr.WriteString(";")
		}

		bldr.WriteString("\n")
	}

	return os.WriteFile(c.Typescript.OutputRoutes, []byte(bldr.String()), 0o644)
}
