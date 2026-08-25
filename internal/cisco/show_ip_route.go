package cisco

import (
	"net/netip"
	"regexp"
	"strings"
)

var (
	routePattern  = regexp.MustCompile(`^([A-Za-z*+%&]+)\s+(\S+)(?:\s+\[(\d+)/(\d+)\])?(?:\s+via\s+(\S+))?(?:,\s*)?(.*)$`)
	legendPattern = regexp.MustCompile(`^([A-Za-z0-9*+%&]+)\s+-\s+(.+)$`)
)

type ShowIPRouteResult struct {
	Legends  []RouteLegend
	Gateways []RouteGateway
	Subnets  []RouteSubnet
	Routes   []Route
}

func (ShowIPRouteResult) CiscoCommand() string { return "show ip route" }

type RouteLegend struct {
	Code        string
	Description string
}

type RouteGateway struct {
	Text string
}

type RouteSubnet struct {
	Prefix string
	Text   string
}

type Route struct {
	Codes       []string
	Prefix      string
	NextHop     string
	Interface   string
	Distance    string
	Metric      string
	Description string
}

func parseShowIPRoute(output []byte) (CommandResult, Status, []Warning, *ParseError) {
	const command = "show ip route"
	if problem := rejectResponse(command, output, false); problem != nil {
		return nil, StatusInvalid, nil, problem
	}
	lines, err := scanInputLines(output, 1<<20)
	if err != nil {
		status, problem := invalid(command, "input_line_limit", "output contains an overlong or unreadable line")
		return nil, status, nil, problem
	}

	result := ShowIPRouteResult{}
	recognized := false
	for _, line := range lines {
		text := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(text, "Codes:"):
			result.Legends = append(result.Legends, parseRouteLegends(strings.TrimSpace(strings.TrimPrefix(text, "Codes:")))...)
			recognized = true
		case strings.HasPrefix(text, "Gateway of last resort"), strings.HasPrefix(text, "Default gateway is"):
			result.Gateways = append(result.Gateways, RouteGateway{Text: text})
			recognized = true
		case strings.Contains(text, " is subnetted") || strings.Contains(text, " is variably subnetted"):
			result.Subnets = append(result.Subnets, RouteSubnet{Prefix: firstField(text), Text: text})
			recognized = true
		default:
			if match := legendPattern.FindStringSubmatch(text); match != nil && !strings.Contains(text, " is ") {
				result.Legends = append(result.Legends, RouteLegend{Code: match[1], Description: match[2]})
				recognized = true
				continue
			}
			if match := routePattern.FindStringSubmatch(text); match != nil && validPrefix(match[2]) {
				description := strings.TrimSpace(match[6])
				result.Routes = append(result.Routes, Route{
					Codes: strings.Split(match[1], ""), Prefix: match[2], Distance: match[3], Metric: match[4],
					NextHop: match[5], Interface: routeInterface(description), Description: description,
				})
				recognized = true
			}
		}
	}
	bindNextHopInterfaces(&result)
	if !recognized {
		status, problem := invalid(command, "unrecognized_output", "output is not recognized as an IPv4 routing table")
		return nil, status, nil, problem
	}
	return result, StatusComplete, nil, nil
}

func parseRouteLegends(value string) []RouteLegend {
	result := make([]RouteLegend, 0)
	for _, part := range strings.Split(value, ",") {
		fields := strings.SplitN(strings.TrimSpace(part), " - ", 2)
		if len(fields) == 2 {
			result = append(result, RouteLegend{Code: strings.TrimSpace(fields[0]), Description: strings.TrimSpace(fields[1])})
		}
	}
	return result
}

func validPrefix(value string) bool {
	_, err := netip.ParsePrefix(value)
	return err == nil
}

func routeInterface(description string) string {
	fields := strings.Fields(description)
	if len(fields) >= 4 && strings.EqualFold(fields[0], "is") && strings.EqualFold(fields[1], "directly") && strings.EqualFold(fields[2], "connected,") {
		return fields[3]
	}
	return ""
}

func bindNextHopInterfaces(result *ShowIPRouteResult) {
	for index := range result.Routes {
		route := &result.Routes[index]
		if route.Interface != "" || route.NextHop == "" {
			continue
		}
		address, err := netip.ParseAddr(route.NextHop)
		if err != nil {
			continue
		}
		bestPrefixLength := -1
		for _, candidate := range result.Routes {
			prefix, err := netip.ParsePrefix(candidate.Prefix)
			if err == nil && candidate.Interface != "" && prefix.Contains(address) && prefix.Bits() > bestPrefixLength {
				route.Interface = candidate.Interface
				bestPrefixLength = prefix.Bits()
			}
		}
	}
}

func firstField(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
