package cisco

import "strings"

type responseKind string

const (
	responseData                 responseKind = "data"
	responseEmpty                responseKind = "empty"
	responseDeviceRejected       responseKind = "device_rejected"
	responseFeatureNotConfigured responseKind = "feature_not_configured"
	responseNoData               responseKind = "no_data"
)

func classifyResponse(output []byte) responseKind {
	lines, err := scanInputLines(output, 1<<20)
	if err != nil {
		return responseDeviceRejected
	}
	meaningful := make([]string, 0, len(lines))
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			meaningful = append(meaningful, trimmed)
		}
	}
	if len(meaningful) == 0 {
		return responseEmpty
	}
	for _, line := range meaningful {
		if isDeviceRejection(line) {
			return responseDeviceRejected
		}
	}
	if standaloneFeatureAbsence(meaningful) {
		return responseFeatureNotConfigured
	}
	if standaloneNoData(meaningful) {
		return responseNoData
	}
	return responseData
}

func isDeviceRejection(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	return strings.HasPrefix(lower, "% invalid input") ||
		strings.HasPrefix(lower, "% unrecognized command") ||
		strings.HasPrefix(lower, "% incomplete command") ||
		strings.HasPrefix(lower, "% ambiguous command") ||
		strings.HasPrefix(lower, "% unknown command") ||
		strings.Contains(lower, "command rejected")
}

func standaloneFeatureAbsence(lines []string) bool {
	if len(lines) > 4 {
		return false
	}
	matched := false
	for _, line := range lines {
		lower := normalizedMessage(line)
		if lower == "" || strings.Trim(lower, " ^") == "" {
			continue
		}
		if containsAny(lower,
			"not configured",
			"not enabled",
			"not running",
			"not active",
			"no ospf process",
			"no eigrp process",
			"no bgp process",
			"no isis process",
			"process does not exist",
			"instance does not exist",
		) {
			matched = true
			continue
		}
		return false
	}
	return matched
}

func standaloneNoData(lines []string) bool {
	if len(lines) > 4 {
		return false
	}
	matched := false
	for _, line := range lines {
		lower := normalizedMessage(line)
		if lower == "" || strings.Trim(lower, " ^") == "" {
			continue
		}
		if lower == "none" || strings.HasPrefix(lower, "no entries") ||
			strings.HasPrefix(lower, "no sessions") || strings.HasPrefix(lower, "no neighbors") ||
			strings.HasPrefix(lower, "no bindings") || strings.HasPrefix(lower, "no records") ||
			strings.HasPrefix(lower, "no matching") ||
			(strings.HasPrefix(lower, "no ") && (strings.HasSuffix(lower, " configured") ||
				strings.HasSuffix(lower, " found") || strings.HasSuffix(lower, " present"))) {
			matched = true
			continue
		}
		return false
	}
	return matched
}

func normalizedMessage(line string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.ToLower(strings.TrimSpace(line)), "%"))
}

func containsAny(value string, fragments ...string) bool {
	for _, fragment := range fragments {
		if strings.Contains(value, fragment) {
			return true
		}
	}
	return false
}
