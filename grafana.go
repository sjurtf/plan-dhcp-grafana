package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"text/template"
)

type LibraryItemResponseModel struct {
	Result LibraryItemModel `json:"result"`
}

type LibraryItemModel struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Version int    `json:"version"`
}

func templatePanelJson(version int, participants []Participant) *bytes.Buffer {
	file, err := os.ReadFile("template/library-model.json.tmpl")
	if err != nil {
		log.Fatalln(err)
	}

	var fns = template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
	}
	tmpl, err := template.New("grafana-participants").Funcs(fns).Parse(string(file))
	if err != nil {
		log.Fatalln(err)
	}

	d := make([]byte, 0)
	buffer := bytes.NewBuffer(d)

	err = tmpl.Execute(buffer, Data{
		Version:      version,
		Participants: participants,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return buffer
}

func updateGrafanaLibraryItem(buffer *bytes.Buffer, panelId string) {
	url := "http://192.168.11.8:3000/api/library-elements/" + panelId
	query := "?ds_type=influxdb"

	token := os.Getenv("GRAFANA_TOKEN")

	r, _ := http.NewRequest(http.MethodPatch, url+query, buffer)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println("unable to do request")
		log.Fatalln(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println(response.Status)
		respBody, _ := io.ReadAll(response.Body)
		fmt.Printf("%s\n", respBody)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Successfully updated participant panel with DHCP data")
	}
}

func getGrafanaLibraryElementVersion(panelId string) (int, error) {
	response, err := http.Get("http://192.168.11.8:3000/api/library-elements/" + panelId)
	if err != nil {
		return 0, err
	}

	libraryElementModel := &LibraryItemResponseModel{}
	err = json.NewDecoder(response.Body).Decode(libraryElementModel)
	if err != nil {
		return 0, err
	}
	return libraryElementModel.Result.Version, nil
}
