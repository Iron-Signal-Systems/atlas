package cisco

import "strings"

type ShowInterfacesStatusResult struct {
	Interfaces []InterfaceStatus
}

func (ShowInterfacesStatusResult) CiscoCommand() string { return "show interfaces status" }

type InterfaceStatus struct {
	Port   string
	Name   string
	Status string
	VLAN   string
	Duplex string
	Speed  string
	Type   string
}

func parseShowInterfacesStatus(output []byte) (CommandResult, Status, []Warning, *ParseError) {
	const command = "show interfaces status"
	if problem := rejectResponse(command, output, false); problem != nil {
		return nil, StatusInvalid, nil, problem
	}
	lines, err := scanInputLines(output, 1<<20)
	if err != nil {
		status, problem := invalid(command, "input_line_limit", "output contains an overlong or unreadable line")
		return nil, status, nil, problem
	}

	header := ""
	result := ShowInterfacesStatusResult{Interfaces: make([]InterfaceStatus, 0, 48)}
	malformedRows := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "port") && strings.Contains(lower, "status") {
			header = line
			continue
		}
		if header == "" || trimmed == "" || onlyDashes(trimmed) {
			continue
		}
		fields := statusFields(header, line)
		if len(fields) != 7 || !isInterfaceName(fields[0]) {
			malformedRows++
			continue
		}
		result.Interfaces = append(result.Interfaces, InterfaceStatus{
			Port: fields[0], Name: fields[1], Status: fields[2], VLAN: fields[3], Duplex: fields[4], Speed: fields[5], Type: fields[6],
		})
	}
	if header == "" {
		status, problem := invalid(command, "unrecognized_output", "output has no recognized interface status table header")
		return nil, status, nil, problem
	}
	if malformedRows > 0 {
		return result, StatusPartial, []Warning{warning(command, "one or more rows did not match the table columns")}, nil
	}
	return result, StatusComplete, nil, nil
}

func statusFields(header, line string) []string {
	columns := []string{"Port", "Name", "Status", "Vlan", "Duplex", "Speed", "Type"}
	positions := make([]int, 0, len(columns))
	last := -1
	for _, column := range columns {
		position := strings.Index(strings.ToLower(header[last+1:]), strings.ToLower(column))
		if position < 0 {
			return nil
		}
		position += last + 1
		positions = append(positions, position)
		last = position
	}
	fields := make([]string, len(positions))
	for index, start := range positions {
		end := len(line)
		if index+1 < len(positions) {
			end = positions[index+1]
		}
		if start >= len(line) {
			continue
		}
		if end > len(line) {
			end = len(line)
		}
		fields[index] = strings.TrimSpace(line[start:end])
	}
	return fields
}

func isInterfaceName(value string) bool {
	for _, prefix := range []string{
		"Gi", "GigabitEthernet", "Te", "TenGigabitEthernet", "Fa", "FastEthernet", "Fo", "FortyGigabitEthernet",
		"Eth", "Ethernet", "Po", "Port-channel", "Vl", "Vlan", "Ap", "AppGigabitEthernet", "Hu", "HundredGigE", "Twe", "TwentyFiveGigE",
	} {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func onlyDashes(value string) bool {
	return strings.Trim(value, "- ") == ""
}
