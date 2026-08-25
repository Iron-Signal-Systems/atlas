package cisco

// ParseCommands parses already-separated Cisco command outputs.
func ParseCommands(commands []CommandOutput) CaptureResult {
	result := CaptureResult{
		Commands:    make([]ParsedCommand, 0, len(commands)),
		Unsupported: make([]CommandOutput, 0),
	}
	for _, command := range commands {
		parsed := ParseCommand(command)
		result.Commands = append(result.Commands, parsed)
		if parsed.Status == StatusUnsupported {
			result.Unsupported = append(result.Unsupported, command)
		}
	}
	return result
}

// ParseCommand dispatches one Cisco command to its command-specific parser.
func ParseCommand(source CommandOutput) ParsedCommand {
	canonical := NormalizeCommand(source.Command)
	parsed := ParsedCommand{Source: source, CanonicalCommand: canonical}

	switch canonical {
	case "show version":
		parsed.Result, parsed.Status, parsed.Warnings, parsed.Error = parseShowVersion(source.Output)
	case "show vlan brief":
		parsed.Result, parsed.Status, parsed.Warnings, parsed.Error = parseShowVLANBrief(source.Output)
	case "show interfaces status":
		parsed.Result, parsed.Status, parsed.Warnings, parsed.Error = parseShowInterfacesStatus(source.Output)
	case "show arp":
		parsed.Result, parsed.Status, parsed.Warnings, parsed.Error = parseShowARP(source.Output)
	case "show ip route":
		parsed.Result, parsed.Status, parsed.Warnings, parsed.Error = parseShowIPRoute(source.Output)
	default:
		parsed.Status = StatusUnsupported
	}
	return parsed
}

func invalid(command, code, message string) (Status, *ParseError) {
	return StatusInvalid, &ParseError{Command: command, Code: code, Message: message}
}

func warning(command, message string) Warning {
	return Warning{Command: command, Message: message}
}
