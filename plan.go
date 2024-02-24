package main

import (
	"encoding/json"
	"net/http"
)

type DHCPLeasesResponseModel map[string]DHCP

type DHCP struct {
	IP      string `json:"ip"`
	MAC     string `json:"mac"`
	Device  string `json:"device"`
	Port    string `json:"port"`
	EndTime string `json:"end-time"`
}

func getPlanDHCPLeases() (DHCPLeasesResponseModel, error) {
	get, err := http.Get("http://plan.mlarsen.no/dhcp_api.php")
	if err != nil {
		return nil, err
	}

	resp := DHCPLeasesResponseModel{}
	err = json.NewDecoder(get.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
