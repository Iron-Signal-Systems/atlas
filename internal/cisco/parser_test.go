package cisco

import "testing"

func TestParseCommandsProducesCommandSpecificResults(t *testing.T) {
	commands := []CommandOutput{
		{Command: "show version", Output: []byte(`Cisco IOS XE Software, Version 17.12.4
SW1 uptime is 10 days, 3 hours
cisco C9300-48P (X86) processor with 123456K bytes of memory.
Processor board ID FDO1234ABC
System image file is "flash:packages.conf"`)},
		{Command: "show vlan brief", Output: []byte(`VLAN Name                             Status    Ports
---- -------------------------------- --------- -------------------------------
1    default                          active    Gi1/0/1
10   USERS                            active    Gi1/0/2, Gi1/0/3`)},
		{Command: "show interfaces status", Output: []byte(`Port         Name               Status       Vlan       Duplex  Speed Type
Gi1/0/1      Server             connected    10         a-full a-1000 10/100/1000BaseTX`)},
		{Command: "show ip arp", Output: []byte(`Protocol  Address          Age (min)  Hardware Addr   Type   Interface
Internet  10.10.10.1              2   0011.2233.4455  ARPA   Vlan10`)},
		{Command: "show ip route", Output: []byte(`Codes: C - connected, S - static
Gateway of last resort is 10.0.0.1 to network 0.0.0.0
C    10.0.0.0/24 is directly connected, Vlan10
S    192.168.50.0/24 [1/0] via 10.0.0.1`)},
		{Command: "show something future", Output: []byte("future")},
	}

	result := ParseCommands(commands)
	if len(result.Commands) != 6 {
		t.Fatalf("expected 6 parsed command envelopes, got %d", len(result.Commands))
	}
	if len(result.Unsupported) != 1 {
		t.Fatalf("expected 1 unsupported command, got %d", len(result.Unsupported))
	}

	version, ok := result.Commands[0].Result.(ShowVersionResult)
	if !ok {
		t.Fatalf("expected ShowVersionResult, got %T", result.Commands[0].Result)
	}
	if version.Version != "17.12.4" || version.Platform != "C9300-48P" {
		t.Fatalf("unexpected show version result: %+v", version)
	}

	vlans, ok := result.Commands[1].Result.(ShowVLANBriefResult)
	if !ok || len(vlans.VLANs) != 2 || len(vlans.VLANs[1].Ports) != 2 {
		t.Fatalf("unexpected VLAN result: %#v", result.Commands[1].Result)
	}

	interfaces, ok := result.Commands[2].Result.(ShowInterfacesStatusResult)
	if !ok || len(interfaces.Interfaces) != 1 || interfaces.Interfaces[0].VLAN != "10" {
		t.Fatalf("unexpected interface result: %#v", result.Commands[2].Result)
	}

	arp, ok := result.Commands[3].Result.(ShowARPResult)
	if !ok || len(arp.Entries) != 1 || arp.Entries[0].MAC != "0011.2233.4455" {
		t.Fatalf("unexpected ARP result: %#v", result.Commands[3].Result)
	}

	routes, ok := result.Commands[4].Result.(ShowIPRouteResult)
	if !ok || len(routes.Routes) != 2 || routes.Routes[1].Interface != "Vlan10" {
		t.Fatalf("unexpected route result: %#v", result.Commands[4].Result)
	}
}

func TestShowARPPartialPreservesValidRows(t *testing.T) {
	parsed := ParseCommand(CommandOutput{Command: "show arp", Output: []byte(`Protocol  Address          Age (min)  Hardware Addr   Type   Interface
Internet  10.10.10.1              2   0011.2233.4455  ARPA   Vlan10
Internet  not-an-ip               5   bad-mac         ARPA   Vlan10`)})
	if parsed.Status != StatusPartial {
		t.Fatalf("expected partial, got %s", parsed.Status)
	}
	arp, ok := parsed.Result.(ShowARPResult)
	if !ok || len(arp.Entries) != 1 {
		t.Fatalf("expected valid ARP row to survive, got %#v", parsed.Result)
	}
	if len(parsed.Warnings) != 1 {
		t.Fatalf("expected one warning, got %d", len(parsed.Warnings))
	}
}

func TestShowARPNoDataIsCompleteEmptyObservation(t *testing.T) {
	parsed := ParseCommand(CommandOutput{Command: "show arp", Output: []byte("No entries found")})
	if parsed.Status != StatusComplete {
		t.Fatalf("expected complete, got %s", parsed.Status)
	}
	arp, ok := parsed.Result.(ShowARPResult)
	if !ok || len(arp.Entries) != 0 {
		t.Fatalf("expected empty ARP result, got %#v", parsed.Result)
	}
}
