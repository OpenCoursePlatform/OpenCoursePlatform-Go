package server

import (
	"testing"
)

func TestMake(t *testing.T) {
	server := Make()
	if server == nil {
		t.Errorf("Server not OK")
	}
}
