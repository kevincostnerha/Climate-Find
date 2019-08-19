package database

import "testing"

func TestSessions(t *testing.T) {

	userID := 14
	var UUID string
	var err error

	Connect()                 //connect to databse
	defer DestroyConnection() //destroy connection when done testing

	//create a session with userID - returns UUID
	if UUID, err = StartSession(userID); err != nil {
		t.Fatal("Failed to start session ", err)
	}

	//check the session created above is still active
	if !ValidateSession(UUID) {
		t.Fatalf("Expected valid session for %s", UUID)
	}

	//end the session by passing UUID to terminate
	EndSession(UUID)

	//try to end an already ended session - should be fine
	EndSession(UUID)

	//check that the session was annihilated
	if ValidateSession(UUID) {
		t.Fatal("Expected INVALID session")
	}

}
