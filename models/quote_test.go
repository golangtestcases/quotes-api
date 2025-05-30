package models

import (
	"encoding/json"
	"testing"
)

func TestQuoteJSON(t *testing.T) {
	q := Quote{
		ID:     1,
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	data, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `{"id":1,"author":"Test Author","quote":"Test Quote"}`
	if string(data) != expected {
		t.Errorf("Got %s, expected %s", string(data), expected)
	}
}
