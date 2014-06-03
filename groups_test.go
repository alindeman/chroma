package chroma

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGroups(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/mySpecialUsername/groups", func(writer http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatalf("Method was %v, but GET expected", req.Method)
		}

		writer.Write([]byte(`{"1": {"name": "Group 1"}, "2": {"name": "VRC 2"}}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &Client{
		BridgeHost: strings.TrimPrefix(server.URL, "http://"),
		Username:   "mySpecialUsername",
	}
	groups, err := client.Groups()
	if err != nil {
		t.Fatal(err)
	}

	if groups[0].Id != "1" || groups[0].Name != "Group 1" {
		t.Fatal("Group #1 was not returned correctly")
	}
	if groups[1].Id != "2" || groups[1].Name != "VRC 2" {
		t.Fatal("Group #2 was not returned correctly")
	}
}

func TestGroupAttributes(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/mySpecialUsername/groups/1", func(writer http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatalf("Method was %v, but GET expected", req.Method)
		}

		writer.Write([]byte(`
			{
					"action": {
							"on": true,
							"hue": 0,
							"effect": "none",
							"bri": 100,
							"sat": 100,
							"ct": 500,
							"xy": [0.5, 0.5]
					},
					"lights": [
							"1",
							"2"
					],
					"name": "bedroom",
					"scenes": {
					}
				}
		`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &Client{
		BridgeHost: strings.TrimPrefix(server.URL, "http://"),
		Username:   "mySpecialUsername",
	}
	attributes, err := client.GroupAttributes("1")
	if err != nil {
		t.Fatal(err)
	}

	if attributes.Lights[0] != "1" && attributes.Lights[1] != "2" {
		t.Fatalf("Lights was %v, but [\"1\", \"2\"] was expected", attributes.Lights)
	}
	if attributes.Name != "bedroom" {
		t.Fatalf("Name was %v, but \"bedroom\" was expected", attributes.Name)
	}
}

func TestSetGroupState(t *testing.T) {
	apiInvoked := false

	handler := http.NewServeMux()
	handler.HandleFunc("/api/mySpecialUsername/groups/1/action", func(writer http.ResponseWriter, req *http.Request) {
		apiInvoked = true
		if req.Method != "PUT" {
			t.Fatalf("Method was %v, but PUT expected", req.Method)
		}

		var body map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatal(err)
		}

		if body["on"] != true {
			t.Fatalf("on was %v, but true expected", body["on"])
		}
		if body["hue"] != float64(12345) {
			t.Fatalf("hue was %v, but 12345 expected", body["hue"])
		}
		if body["sat"] != nil {
			t.Fatalf("sat was not expected to be part of the body")
		}

		writer.Write([]byte(`[]`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &Client{
		BridgeHost: strings.TrimPrefix(server.URL, "http://"),
		Username:   "mySpecialUsername",
	}

	change := &GroupStateChange{
		On:  new(bool),
		Hue: new(int),
	}
	*change.On = true
	*change.Hue = 12345

	err := client.SetGroupState("1", change)
	if err != nil {
		t.Fatal(err)
	}

	if !apiInvoked {
		t.Fatalf("API was not called")
	}
}
