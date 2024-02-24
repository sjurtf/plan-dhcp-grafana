package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type Participant struct {
	ComputerName string
	SwitchName   string
	Port         string
	RefId        string
}

type Data struct {
	Version      int
	Participants []Participant
}

func main() {
	panelId := "acd5b3fe-e0f3-41fb-9698-2f782b60939c"

	leases, err := getPlanDHCPLeases()
	if err != nil {
		log.Fatalf("unable to get DHCP leases: %s", err)
	}

	version, err := getGrafanaLibraryElementVersion(panelId)
	if err != nil {
		log.Fatalf("unable to current Grafana panel version: %s", err)
	}

	participants := mapGrafanaParticipants(leases)
	payload := templatePanelJson(version, participants)
	updateGrafanaLibraryItem(payload, panelId)
}

func mapGrafanaParticipants(leases DHCPLeasesResponseModel) []Participant {
	var participants []Participant
	cnt := 0
	for s, dhcp := range leases {
		// Ignore plan clients from WiFi
		if slices.Contains([]string{"sw.pub", "sw.garage"}, dhcp.Device) {
			continue
		}

		// Handle Cisco switches with wrong port info in option82
		if slices.Contains([]string{"sw6.plan"}, dhcp.Device) {
			portParts := strings.Split(dhcp.Port, "/")
			dhcp.Port = fmt.Sprintf("Gi%s/0/%s", portParts[0], portParts[1])
		}

		if s == "H_LOQ15" {
			s = "Harald"
		}

		p := Participant{
			ComputerName: s,
			SwitchName:   strings.ReplaceAll(dhcp.Device, ".ex2200", ".plan"),
			Port:         unEscapePort(dhcp.Port),
			RefId:        fmt.Sprintf("%c", toChar(cnt)),
		}
		participants = append(participants, p)
		cnt = cnt + 1
	}
	return participants
}

func unEscapePort(port string) string {
	return strings.ReplaceAll(port, "\\/", "")
}

func toChar(i int) rune {
	return rune('A' - 1 + i)
}
