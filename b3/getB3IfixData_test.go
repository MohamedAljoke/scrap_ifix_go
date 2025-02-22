package b3

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_B3_GetDataSuccess(t *testing.T) {
	expected := []Asset{
		{
			Code:  "BRCR11",
			Asset: "BCFF11",
			Type:  "FII",
			Part:  "DR",
		},
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Results: expected,
		}
		respByte, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respByte)
	}))

	defer mockServer.Close()

	// t.Logf(mockServer.URL)
	// time.Sleep(10 * time.Second)

	oldURL := url
	url = mockServer.URL
	defer func() { url = oldURL }()

	result, err := GetB3IfixData()

	if err != nil {
		t.Errorf("Expected error to be nil got %s", err.Error())
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result, expected: %v, got: %v", expected, result)
	}
}

func Test_B3_GetDataErrorBodyUnmarshal(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bytes := []byte("Some data")
		w.Write(bytes)
	}))

	defer mockServer.Close()

	oldURL := url
	url = mockServer.URL
	defer func() { url = oldURL }()

	result, err := GetB3IfixData()

	if result != nil {
		t.Errorf("Expected nil for result, got: %v", result)
	}
	if err == nil {
		t.Error("Expected non-nil error, got nil")
	}
	errorMessage := "body must have only a single JSON value"
	if err.Error() != errorMessage {
		t.Errorf("Expected error message to be %s got: %s", errorMessage, err.Error())
	}
}

func Test_B3_GetDataErrorReadingResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Length", "50")
		w.Write([]byte("a"))
	}))

	defer mockServer.Close()

	oldURL := url
	url = mockServer.URL
	defer func() { url = oldURL }()

	result, err := GetB3IfixData()

	if result != nil {
		t.Errorf("Expected nil for result, got: %v", result)
	}
	if err == nil {
		t.Error("Expected non-nil error, got nil")
	}
	errorMessage := "error getting body"
	if err.Error() != errorMessage {
		t.Errorf("Expected error message to be %s got: %s", errorMessage, err.Error())
	}
}

func Test_B3_GetB3DataErrorCreatingRequest(t *testing.T) {
	originalDefaultClient := http.DefaultClient

	httpClient := &http.Client{
		Transport: errorTransport{},
	}
	http.DefaultClient = httpClient
	defer func() {
		http.DefaultClient = originalDefaultClient
	}()
	result, err := GetB3IfixData()
	if result != nil {
		t.Errorf("Expected nil for result, got: %v", result)
	}
	if err == nil {
		t.Error("Expected non-nil error, got nil")
	}
	errorMessage := "error starting client"
	if err.Error() != errorMessage {
		t.Errorf("Expected error message to be %s got: %s", errorMessage, err.Error())
	}
}

type errorTransport struct{}

func (t errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("Mocked error")
}
