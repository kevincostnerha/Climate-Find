package database

import "testing"

func TestConnection(t *testing.T) {
	Connect()
	defer DestroyConnection()
}
