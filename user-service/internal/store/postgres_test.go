package store

import (
	"testing"
)

func TestConnection(t *testing.T) {
	_, err := NewPostgre()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

}
