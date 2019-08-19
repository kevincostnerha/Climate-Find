package database

import (
	"testing"
)

func TestQueries(t *testing.T) {

	userID := 3
	actName := "Recycle aluminum cans"
	//goalID := 1
	//userID := 1

	Connect()                 //connect to databse
	defer DestroyConnection() //destroy connection when done testing

	var d PageData

	d.GetRecc(1)

	if err := d.LoadUser(userID); err != nil {
		t.Fatalf("Failed to load info for userID %d", userID)
	}

	if err := AddUserAct(userID, actName); err != nil {
		t.Fatalf("Failed to add activity %s for user %d", actName, userID)
	}

	if err := d.LoadUserActs(userID); err != nil {
		t.Fatalf("Failed to load activities for userID %d", userID)
	}

	if err := AddUserKarma(userID, actName); err != nil {
		t.Fatalf("Failed to add karma for user %d activity %s", userID, actName)
	}

}
