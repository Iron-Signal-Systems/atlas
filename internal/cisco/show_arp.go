package cisco

import (
	"net/netip"
	"regexp"
	"strings"
)

var arpMACPattern = regexp.MustCompile(`(?i)^(?:[0-9a-f]{4}\.){2}[0-9a-f]{4}$`)

type ShowARPResult struct {
	Entries []ARPEntry
}

func (ShowARPResult) CiscoCommand() string { return "show arp" }

type ARPEntry struct {
	Protocol   string
	Address    string
	AgeMinutes string
	MAC        string
	Type       string
	Interface  string
}

func parseShowARP(output []byte) (CommandResult, Status, []Warning, *ParseError) {
	const command = "show arp"
	if problem := rejectResponse(command, output, true); problem != nil {
		return nil, StatusInvalid, nil, problem
	}
	if emptyObservation(output) {
		return ShowARPResult{Entries: []ARPEntry{}}, StatusComplete, nil, nil
	}
	lines, err := scanInputLines(output, 1<<20)
	if err != nil {
		status, problem := invalid(command, "input_line_limit", "output contains an overlong or unreadable line")
		return nil, status, nil, problem
	}

	result := ShowARPResult{Entries: make([]ARPEntry, 0, len(lines))}
	headerSeen := false
	malformed := 0
	for _, line := range lines {
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		lower := strings.ToLower(text)
		if strings.HasPrefix(lower, "protocol") && strings.Contains(lower, "hardware addr") && strings.Contains(lower, "interface") {
			headerSeen = true
			continue
		}
		if !headerSeen {
			continue
		}
		fields := strings.Fields(text)
		if len(fields) < 5 {
			malformed++
			continue
		}
		address, parseErr := netip.ParseAddr(fields[1])
		if parseErr != nil || !address.Is4() || (!arpMACPattern.MatchString(fields[3]) && !strings.EqualFold(fields[3], "incomplete")) {
			malformed++
			continue
		}
		interfaceName := ""
		if len(fields) > 5 {
			interfaceName = fields[5]
		}
		result.Entries = append(result.Entries, ARPEntry{
			Protocol: fields[0], Address: fields[1], AgeMinutes: fields[2], MAC: fields[3], Type: fields[4], Interface: interfaceName,
		})
	}

	if !headerSeen {
		status, problem := invalid(command, "unrecognized_output", "output has no recognized ARP table header")
		return nil, status, nil, problem
	}
	if malformed > 0 {
		if len(result.Entries) == 0 {
			status, problem := invalid(command, "malformed_rows", "output contains no valid ARP rows")
			return nil, status, nil, problem
		}
		return result, StatusPartial, []Warning{warning(command, "one or more ARP rows were malformed")}, nil
	}
	return result, StatusComplete, nil, nil
}
