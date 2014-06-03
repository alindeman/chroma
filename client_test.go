package chroma

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthorize(t *testing.T) {
	apiInvoked := false

	handler := http.NewServeMux()
	handler.HandleFunc("/api/", func(writer http.ResponseWriter, req *http.Request) {
		apiInvoked = true
		if req.Method != "POST" {
			t.Fatalf("Method was %v, but POST expected", req.Method)
		}

		var body map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatal(err)
		}

		if body["devicetype"] != "mySpecialDevice" {
			t.Fatalf("devicetype was %v, but \"mySpecialDevice\" expected", body["devicetype"])
		}
		if body["username"] != "mySpecialUsername" {
			t.Fatalf("username was %v, but \"mySpecialUsername\" expected", body["username"])
		}

		writer.Write([]byte(`[{"success": {"username": "mySpecialUsername"}}]`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &Client{
		BridgeHost: strings.TrimPrefix(server.URL, "http://"),
		Username:   "mySpecialUsername",
	}
	err := client.Authorize("mySpecialDevice")
	if err != nil {
		t.Fatal(err)
	}

	if !apiInvoked {
		t.Fatalf("API was not called")
	}
}
