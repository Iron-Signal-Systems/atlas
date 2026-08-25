package cisco

func rejectResponse(command string, output []byte, allowEmpty bool) *ParseError {
	switch classifyResponse(output) {
	case responseDeviceRejected:
		return &ParseError{Command: command, Code: "device_command_rejected", Message: "device rejected command"}
	case responseFeatureNotConfigured:
		return &ParseError{Command: command, Code: "feature_not_configured", Message: "feature is not configured"}
	case responseEmpty:
		if !allowEmpty {
			return &ParseError{Command: command, Code: "empty_output", Message: "command returned no output"}
		}
	}
	return nil
}

func emptyObservation(output []byte) bool {
	switch classifyResponse(output) {
	case responseEmpty, responseNoData:
		return true
	default:
		return false
	}
}
