package cisco

import (
	"strings"
	"testing"
)

func TestSplitTranscript(t *testing.T) {
	input := []byte(`SW1#show version
Cisco IOS XE Software, Version 17.12.4
SW1#show vlan brief
VLAN Name                             Status    Ports
---- -------------------------------- --------- -------------------------------
1    default                          active    Gi1/0/1
10   USERS                            active    Gi1/0/2, Gi1/0/3
SW1#show ip arp
Protocol  Address          Age (min)  Hardware Addr   Type   Interface
Internet  10.10.10.1              2   0011.2233.4455  ARPA   Vlan10
`)
	commands, err := SplitTranscript(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(commands))
	}
	if commands[0].DevicePrompt != "SW1" || commands[0].Command != "show version" {
		t.Fatalf("unexpected first command: %+v", commands[0])
	}
	if commands[0].StartLine != 1 || commands[0].EndLine != 2 {
		t.Fatalf("unexpected first command line range: %d-%d", commands[0].StartLine, commands[0].EndLine)
	}
}

func TestParseCommandBundlePreservesRepeatedCommands(t *testing.T) {
	bundle := `===== COMMAND: show arp =====
first
===== END COMMAND =====
===== COMMAND: show arp =====
second
===== END COMMAND =====
`
	commands, err := ParseCommandBundle(strings.NewReader(bundle))
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 2 {
		t.Fatalf("expected repeated commands to be preserved, got %d", len(commands))
	}
	if string(commands[0].Output) != "first" || string(commands[1].Output) != "second" {
		t.Fatalf("unexpected command outputs: %q / %q", commands[0].Output, commands[1].Output)
	}
}
