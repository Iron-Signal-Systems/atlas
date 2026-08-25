package cisco

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	commandPrefix = "===== COMMAND: "
	commandSuffix = " ====="
	endMarker     = "===== END COMMAND ====="
)

var promptPattern = regexp.MustCompile(`^([[:alnum:]_.-]+)[#>]\s*(.*)$`)

// NormalizeCommand returns the canonical command spelling used by the dispatcher.
func NormalizeCommand(command string) string {
	normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
	switch normalized {
	case "show version", "show ver", "sh version", "sh ver":
		return "show version"
	case "show vlan brief", "sh vlan brief":
		return "show vlan brief"
	case "show interfaces status", "show int status", "sh interfaces status", "sh int status":
		return "show interfaces status"
	case "show arp", "show ip arp", "sh arp", "sh ip arp":
		return "show arp"
	case "show ip route", "sh ip route":
		return "show ip route"
	default:
		return normalized
	}
}

// SplitTranscript separates a pasted Cisco CLI session into command/output pairs.
func SplitTranscript(input []byte) ([]CommandOutput, error) {
	lines, err := scanInputLines(input, 32<<20)
	if err != nil {
		return nil, err
	}

	commands := make([]CommandOutput, 0, 32)
	var current *CommandOutput
	var body strings.Builder

	flush := func(endLine int) {
		if current == nil {
			return
		}
		current.Output = []byte(strings.TrimRight(body.String(), "\n"))
		current.EndLine = endLine
		commands = append(commands, *current)
		current = nil
		body.Reset()
	}

	for index, line := range lines {
		lineNumber := index + 1
		trimmed := strings.TrimSpace(line)
		match := promptPattern.FindStringSubmatch(trimmed)
		if match != nil && strings.TrimSpace(match[2]) != "" {
			command := strings.TrimSpace(match[2])
			if looksLikeCommand(command) {
				flush(lineNumber - 1)
				current = &CommandOutput{
					DevicePrompt: match[1],
					Command:      command,
					StartLine:    lineNumber,
				}
				continue
			}
		}
		if current != nil {
			body.WriteString(line)
			body.WriteByte('\n')
		}
	}
	flush(len(lines))
	return commands, nil
}

// ParseCommandBundle parses the explicit command-bundle format used by offline captures.
// Unlike the old map representation, this preserves order and repeated commands.
func ParseCommandBundle(reader io.Reader) ([]CommandOutput, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 32*1024*1024)

	commands := make([]CommandOutput, 0, 32)
	var current *CommandOutput
	var body strings.Builder
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasPrefix(line, commandPrefix) && strings.HasSuffix(line, commandSuffix) {
			if current != nil {
				return nil, fmt.Errorf("line %d: nested command section", lineNumber)
			}
			current = &CommandOutput{
				Command:   strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, commandPrefix), commandSuffix)),
				StartLine: lineNumber,
			}
			body.Reset()
			continue
		}
		if line == endMarker {
			if current == nil {
				return nil, fmt.Errorf("line %d: end marker without command", lineNumber)
			}
			current.Output = []byte(strings.TrimRight(body.String(), "\n"))
			current.EndLine = lineNumber
			commands = append(commands, *current)
			current = nil
			body.Reset()
			continue
		}
		if current != nil {
			body.WriteString(line)
			body.WriteByte('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if current != nil {
		return nil, fmt.Errorf("command section %q is not closed", current.Command)
	}
	return commands, nil
}

func looksLikeCommand(command string) bool {
	normalized := strings.Join(strings.Fields(strings.ToLower(command)), " ")
	return strings.HasPrefix(normalized, "show ") || strings.HasPrefix(normalized, "sh ")
}

func scanInputLines(input []byte, maxToken int) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Buffer(make([]byte, 64*1024), maxToken)
	lines := make([]string, 0, 256)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
