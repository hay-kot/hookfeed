package console

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/hookfeed/backend/pkgs/styles"
)

var (
	jsonKeyColor      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7dd3fc")) // Light blue for keys
	jsonStringColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#86efac")) // Light green for strings
	jsonNumberColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#fbbf24")) // Yellow for numbers
	jsonBoolNullColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#c084fc")) // Purple for bool/null
	jsonPunctColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#a3a3a3")) // Gray for punctuation
	jsonBoxPrefix     = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.ColorSubtle))
)

// FatalError formats a fatal error for the CLI
// FatalError printer an error message for an unknown or unexpected error.
// This is used when an error in the system was unexpected, and the error output
// should be displayed to the user.
//
// If the error implements the ConsoleOutput interface, the ConsoleOutput method
// will be called to get the error output.
func FatalError(err error) string {
	bldr := &strings.Builder{}

	// Create error box with red border
	errorBox := styles.ErrorBox("An error occurred", err.Error())
	bldr.WriteString(errorBox)
	bldr.WriteString("\n")

	return bldr.String()
}

// SectionTitle formats a section title for the CLI
func SectionTitle(title string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a3a3a3")).
		Bold(true)

	subtleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSubtle))

	return "\n" + subtleStyle.Render("╭ ") + titleStyle.Render(title)
}

// PrettyJSON return a pretty formatted and colorized display for JSON that
// can be used in the terminal.
func PrettyJSON(v any) string {
	bldr := &strings.Builder{}

	// Parse the input based on its type
	var jsonBytes []byte
	var err error

	// If v is already []byte, use it directly, otherwise marshal it
	switch data := v.(type) {
	case []byte:
		jsonBytes = data
	case string:
		jsonBytes = []byte(data)
	default:
		// Marshal with indentation
		jsonBytes, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err.Error()
		}
	}

	// For []byte or string input, ensure it's valid JSON and format it
	if _, ok := v.([]byte); ok {
		var result any
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return string(jsonBytes) // Return as-is if not valid JSON
		}
		jsonBytes, _ = json.MarshalIndent(result, "", "  ")
	} else if _, ok := v.(string); ok {
		var result any
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return string(jsonBytes) // Return as-is if not valid JSON
		}
		jsonBytes, _ = json.MarshalIndent(result, "", "  ")
	}

	jsonStr := string(jsonBytes)

	// Apply syntax highlighting with box prefix
	prefix := jsonBoxPrefix.Render("│") + "  "

	lines := strings.Split(jsonStr, "\n")
	for _, line := range lines {
		coloredLine := colorizeJSONLine(line)
		bldr.WriteString(prefix)
		bldr.WriteString(coloredLine)
		bldr.WriteString("\n")
	}

	return bldr.String()
}

// colorizeJSONLine applies syntax highlighting to a single line of JSON
func colorizeJSONLine(line string) string {
	result := &strings.Builder{}
	i := 0

	for i < len(line) {
		ch := line[i]

		// Handle whitespace
		if ch == ' ' || ch == '\t' {
			result.WriteByte(ch)
			i++
			continue
		}

		// Handle structural characters
		if ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == ',' || ch == ':' {
			result.WriteString(jsonPunctColor.Render(string(ch)))
			i++
			continue
		}

		// Handle strings (keys and values)
		if ch == '"' {
			str, newI := colorizeJSONString(line, i)
			if newI > i {
				result.WriteString(str)
				i = newI
				continue
			}
		}

		// Handle numbers
		if (ch >= '0' && ch <= '9') || ch == '-' {
			num, newI := colorizeJSONNumber(line, i)
			result.WriteString(num)
			i = newI
			continue
		}

		// Handle booleans and null
		if keyword, newI := colorizeJSONKeyword(line, i); newI > i {
			result.WriteString(keyword)
			i = newI
			continue
		}

		// Default: output as-is
		result.WriteByte(ch)
		i++
	}

	return result.String()
}

// colorizeJSONString handles colorizing a JSON string starting at position i
func colorizeJSONString(line string, i int) (string, int) {
	end := i + 1
	for end < len(line) && line[end] != '"' {
		if line[end] == '\\' {
			end++ // Skip escaped character
		}
		end++
	}
	if end >= len(line) {
		return "", i
	}

	end++ // Include closing quote
	str := line[i:end]

	// Check if this is a key (followed by colon after whitespace)
	isKey := isJSONKey(line, end)

	if isKey {
		return jsonKeyColor.Render(str), end
	}
	return jsonStringColor.Render(str), end
}

// isJSONKey checks if a string is followed by a colon (making it a key)
func isJSONKey(line string, pos int) bool {
	for pos < len(line) && (line[pos] == ' ' || line[pos] == '\t') {
		pos++
	}
	return pos < len(line) && line[pos] == ':'
}

// colorizeJSONNumber handles colorizing a JSON number starting at position i
func colorizeJSONNumber(line string, i int) (string, int) {
	end := i + 1
	for end < len(line) && isNumberChar(line[end]) {
		end++
	}
	return jsonNumberColor.Render(line[i:end]), end
}

// isNumberChar checks if a character is valid in a JSON number
func isNumberChar(ch byte) bool {
	return (ch >= '0' && ch <= '9') || ch == '.' || ch == 'e' || ch == 'E' || ch == '-' || ch == '+'
}

// colorizeJSONKeyword handles colorizing JSON keywords (true, false, null)
func colorizeJSONKeyword(line string, i int) (string, int) {
	if strings.HasPrefix(line[i:], "true") {
		return jsonBoolNullColor.Render("true"), i + 4
	}
	if strings.HasPrefix(line[i:], "false") {
		return jsonBoolNullColor.Render("false"), i + 5
	}
	if strings.HasPrefix(line[i:], "null") {
		return jsonBoolNullColor.Render("null"), i + 4
	}
	return "", i
}
