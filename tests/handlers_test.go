package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MeherKandukuri/studioClasses_API/handlers"
)

func TestCreateClass(t *testing.T) {
	reqBody := `{
	"class_name":"Pilates",
	"start_date":"2024-12-01",
	"end_date" : "2024-12-20",
	"capacity":10
	}`
	req, err := http.NewRequest("POST", "/classes", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.PostCreateClass)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expectedResponse := map[string]string{"message": "Class created successfully"}

	var actualResponse map[string]string

	if err = json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("handler returned unexpected body: got %v want %v", actualResponse, expectedResponse)
	}

}
