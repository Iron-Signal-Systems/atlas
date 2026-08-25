package cisco

import (
	"regexp"
	"strings"
)

var (
	versionPattern        = regexp.MustCompile(`(?i)\bVersion\s+([^,\s]+)`)
	platformPattern       = regexp.MustCompile(`(?i)^cisco\s+([^\s(]+)\s+\(`)
	interfaceCountPattern = regexp.MustCompile(`(?i)^(\d+)\s+(.+?)\s+interfaces$`)
)

type ShowVersionResult struct {
	SoftwareFamily            string
	Version                   string
	Platform                  string
	Image                     string
	Uptime                    string
	ReloadReason              string
	SystemImage               string
	ProcessorBoardID          string
	BaseMACAddress            string
	MotherboardSerial         string
	ModelNumber               string
	SystemSerial              string
	PowerSupplyPartNumber     string
	PowerSupplySerial         string
	VersionID                 string
	InterfaceCounts           map[string]string
	Licenses                  []VersionLicensePackage
	AIRLicenseLevel           string
	NextReloadAIRLicenseLevel string
	SmartLicensingStatus      string
}

func (ShowVersionResult) CiscoCommand() string { return "show version" }

type VersionLicensePackage struct {
	Package    string
	Type       string
	NextReboot string
}

func parseShowVersion(output []byte) (CommandResult, Status, []Warning, *ParseError) {
	const command = "show version"
	if problem := rejectResponse(command, output, false); problem != nil {
		return nil, StatusInvalid, nil, problem
	}
	lines, err := scanInputLines(output, 1<<20)
	if err != nil {
		status, problem := invalid(command, "input_line_limit", "output contains an overlong or unreadable line")
		return nil, status, nil, problem
	}

	result := ShowVersionResult{InterfaceCounts: make(map[string]string)}
	licenseTable := false
	recognized := false
	for _, line := range lines {
		text := strings.TrimSpace(line)
		lower := strings.ToLower(text)
		switch {
		case strings.HasPrefix(lower, "technology package license information:"):
			licenseTable = true
			recognized = true
		case strings.HasPrefix(lower, "air license level:"):
			licenseTable = false
			result.AIRLicenseLevel = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "next reload air license level:"):
			result.NextReloadAIRLicenseLevel = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "smart licensing status:"):
			result.SmartLicensingStatus = valueAfterColon(text)
			recognized = true
		case strings.Contains(lower, "ios xe software"):
			result.SoftwareFamily = "ios_xe"
			result.Version = firstSubmatch(versionPattern, text, result.Version)
			recognized = true
		case strings.Contains(lower, "ios software"):
			if result.SoftwareFamily == "" {
				result.SoftwareFamily = "ios"
			}
			result.Version = firstSubmatch(versionPattern, text, result.Version)
			recognized = true
		case strings.Contains(lower, "uptime is "):
			result.Uptime = strings.TrimSpace(text[strings.Index(lower, "uptime is ")+len("uptime is "):])
			recognized = true
		case strings.HasPrefix(lower, "system image file is "):
			result.SystemImage = strings.Trim(strings.TrimSpace(text[len("System image file is "):]), `"`)
			recognized = true
		case strings.HasPrefix(lower, "last reload reason:"):
			result.ReloadReason = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "processor board id"):
			result.ProcessorBoardID = valueAfterColon(text)
			recognized = true
		case strings.Contains(lower, "base ethernet mac address"):
			result.BaseMACAddress = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "motherboard serial number"):
			result.MotherboardSerial = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "model number"):
			result.ModelNumber = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "system serial number"):
			result.SystemSerial = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "power supply part number"):
			result.PowerSupplyPartNumber = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "power supply serial number"):
			result.PowerSupplySerial = valueAfterColon(text)
			recognized = true
		case strings.HasPrefix(lower, "version id"):
			result.VersionID = valueAfterColon(text)
			recognized = true
		}
		if match := interfaceCountPattern.FindStringSubmatch(text); len(match) == 3 {
			result.InterfaceCounts[strings.TrimSpace(match[2])] = match[1]
			recognized = true
		}
		if licenseTable && text == "" && len(result.Licenses) > 0 {
			licenseTable = false
		}
		if licenseTable {
			parseVersionLicenseLine(text, &result.Licenses)
		}
		if result.Platform == "" {
			result.Platform = firstSubmatch(platformPattern, text, "")
		}
		if result.Image == "" && strings.Contains(lower, "catalyst") {
			result.Image = text
		}
	}
	if !recognized || result.SoftwareFamily == "" && result.Version == "" {
		status, problem := invalid(command, "unrecognized_output", "output has no recognized software identity")
		return nil, status, nil, problem
	}

	missing := make([]string, 0, 3)
	if result.SoftwareFamily == "" {
		missing = append(missing, "software family")
	}
	if result.Version == "" {
		missing = append(missing, "version")
	}
	if result.Platform == "" {
		missing = append(missing, "platform")
	}
	if len(missing) > 0 {
		return result, StatusPartial, []Warning{warning(command, "output is missing "+strings.Join(missing, ", "))}, nil
	}
	return result, StatusComplete, nil, nil
}

func parseVersionLicenseLine(text string, packages *[]VersionLicensePackage) {
	if text == "" || strings.HasPrefix(text, "-") {
		return
	}
	lower := strings.ToLower(text)
	if strings.HasPrefix(lower, "current") && strings.Contains(lower, "next reboot") ||
		strings.Contains(lower, "technology package license information") || strings.Contains(lower, "license level:") ||
		strings.Contains(lower, "licensing status:") {
		return
	}
	fields := strings.Fields(text)
	if len(fields) < 3 {
		return
	}
	*packages = append(*packages, VersionLicensePackage{
		Package: fields[0], Type: strings.Join(fields[1:len(fields)-1], " "), NextReboot: fields[len(fields)-1],
	})
}

func valueAfterColon(text string) string {
	if index := strings.Index(text, ":"); index >= 0 {
		return strings.TrimSpace(text[index+1:])
	}
	return ""
}

func firstSubmatch(pattern *regexp.Regexp, text, fallback string) string {
	match := pattern.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return fallback
}
