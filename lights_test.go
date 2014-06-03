package chroma

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLights(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/mySpecialUsername/lights", func(writer http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatalf("Method was %v, but GET expected", req.Method)
		}

		writer.Write([]byte(`{"1": {"name": "Bedroom"}, "2": {"name": "Kitchen"}}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &Client{
		BridgeHost: strings.TrimPrefix(server.URL, "http://"),
		Username:   "mySpecialUsername",
	}
	lights, err := client.Lights()
	if err != nil {
		t.Fatal(err)
	}

	if lights[0].Id != "1" || lights[0].Name != "Bedroom" {
		t.Fatal("Light #1 was not returned correctly")
	}
	if lights[1].Id != "2" || lights[1].Name != "Kitchen" {
		t.Fatal("Light #2 was not returned correctly")
	}
}

func TestLightAttributes(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/mySpecialUsername/lights/1", func(writer http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Fatalf("Method was %v, but GET expected", req.Method)
		}

		writer.Write([]byte(`
			{
					"state": {
							"hue": 50000,
							"on": true,
							"effect": "none",
							"alert": "none",
							"bri": 200,
							"sat": 200,
							"ct": 500,
							"xy": [0.5, 0.5],
							"reachable": true,
							"colormode": "hs"
					},
					"type": "Living Colors",
					"name": "LC 1",
					"modelid": "LC0015",
					"swversion": "1.0.3",
					"pointsymbol": {
							"1": "none",
							"2": "none",
							"3": "none",
							"4": "none",
							"5": "none",
							"6": "none",
							"7": "none",
							"8": "none"
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
	attributes, err := client.LightAttributes("1")
	if err != nil {
		t.Fatal(err)
	}

	// Spot check a few attributes
	if attributes.State.Hue != 50000 {
		t.Fatalf("Hue was %v, but 50000 was expected", attributes.State.Hue)
	}
	if attributes.Type != "Living Colors" {
		t.Fatalf("Type was %v, but \"Living Colors\" was expected", attributes.Type)
	}
}

func TestSetLightState(t *testing.T) {
	apiInvoked := false

	handler := http.NewServeMux()
	handler.HandleFunc("/api/mySpecialUsername/lights/1/state", func(writer http.ResponseWriter, req *http.Request) {
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

	change := &LightStateChange{
		On:  new(bool),
		Hue: new(int),
	}
	*change.On = true
	*change.Hue = 12345

	err := client.SetLightState("1", change)
	if err != nil {
		t.Fatal(err)
	}

	if !apiInvoked {
		t.Fatalf("API was not called")
	}
}
