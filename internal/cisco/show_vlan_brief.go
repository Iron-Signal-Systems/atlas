package cisco

import "strings"

type ShowVLANBriefResult struct {
	VLANs []VLANBrief
}

func (ShowVLANBriefResult) CiscoCommand() string { return "show vlan brief" }

type VLANBrief struct {
	ID     string
	Name   string
	Status string
	Ports  []string
}

func parseShowVLANBrief(output []byte) (CommandResult, Status, []Warning, *ParseError) {
	const command = "show vlan brief"
	if problem := rejectResponse(command, output, true); problem != nil {
		return nil, StatusInvalid, nil, problem
	}
	if emptyObservation(output) {
		return ShowVLANBriefResult{VLANs: []VLANBrief{}}, StatusComplete, nil, nil
	}
	lines, err := scanInputLines(output, 1<<20)
	if err != nil {
		status, problem := invalid(command, "input_line_limit", "output contains an overlong or unreadable line")
		return nil, status, nil, problem
	}

	result := ShowVLANBriefResult{VLANs: make([]VLANBrief, 0, 32)}
	header := false
	for _, line := range lines {
		text := strings.TrimSpace(line)
		if strings.HasPrefix(text, "VLAN Name") {
			header = true
			continue
		}
		if !header || text == "" || strings.HasPrefix(text, "----") {
			continue
		}
		fields := strings.Fields(text)
		if len(fields) >= 3 && digits(fields[0]) {
			result.VLANs = append(result.VLANs, VLANBrief{
				ID: fields[0], Name: fields[1], Status: fields[2], Ports: splitPorts(strings.Join(fields[3:], " ")),
			})
			continue
		}
		if len(result.VLANs) > 0 {
			result.VLANs[len(result.VLANs)-1].Ports = append(result.VLANs[len(result.VLANs)-1].Ports, splitPorts(text)...)
		}
	}
	if !header {
		status, problem := invalid(command, "unrecognized_output", "output has no recognized VLAN table header")
		return nil, status, nil, problem
	}
	return result, StatusComplete, nil, nil
}

func digits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func splitPorts(value string) []string {
	ports := make([]string, 0)
	for _, port := range strings.Split(value, ",") {
		if port = strings.TrimSpace(port); port != "" {
			ports = append(ports, port)
		}
	}
	return ports
}
