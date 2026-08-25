package cisco

import "fmt"

// Status describes how completely one Cisco command was parsed.
type Status string

const (
	StatusComplete    Status = "complete"
	StatusPartial     Status = "partial"
	StatusUnsupported Status = "unsupported"
	StatusInvalid     Status = "invalid"
)

// CommandOutput is one Cisco command and the text printed by the device for it.
// StartLine and EndLine refer to the original input when it came from a CLI transcript.
type CommandOutput struct {
	DevicePrompt string
	Command      string
	Output       []byte
	StartLine    int
	EndLine      int
}

// CommandResult is implemented by every command-specific typed result.
type CommandResult interface {
	CiscoCommand() string
}

// Warning records a recoverable parser problem without discarding valid data.
type Warning struct {
	Command string
	Message string
}

// ParseError records why a command could not be parsed safely.
type ParseError struct {
	Command string
	Code    string
	Message string
}

func (e *ParseError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s: %s", e.Command, e.Code, e.Message)
}

// ParsedCommand ties the original command to its command-specific Go result.
type ParsedCommand struct {
	Source           CommandOutput
	CanonicalCommand string
	Status           Status
	Result           CommandResult
	Warnings         []Warning
	Error            *ParseError
}

// CaptureResult is the complete parser output for a set of Cisco commands.
type CaptureResult struct {
	Commands    []ParsedCommand
	Unsupported []CommandOutput
}
