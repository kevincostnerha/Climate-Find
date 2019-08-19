package database

import (
	"database/sql"
	"log"

	uuid "github.com/satori/go.uuid"
)

//This is a separate package because this will be public to the client
//The purpose of this package is to handle our database interactions
//...but for now all it does is deal with sessions

//StartSession starts a session by generting UUID for user and adding to user_data
// receives userID for the user that needs a session
// returns the UUID as a string, and error (nil if added successfully)
func StartSession(userID int) (string, error) {

	if db == nil {
		log.Fatal("db is not connected")
	}

	// create the UUID - no timestamp because not necessary
	sID := uuid.NewV4()

	// stringify the UUID because that's what the db expects
	UUID := sID.String()

	// save the UUID to the db for this user to indicate active session
	_, err := db.Exec("UPDATE users SET sessionID=? WHERE id=?", UUID, userID)

	return UUID, err

}

//ValidateSession validates a session by checking for UUID in user_data
// receives UUID to check in user_data table
// returns TRUE if session is valid, FALSE if session invalid
func ValidateSession(UUID string) bool {

	var tmp int // so we can use Scan
	// check db for this UUID, i.e., active session
	err := db.QueryRow("SELECT id FROM users WHERE sessionID=?", UUID).Scan(&tmp)

	// if no row results, return FALSE to indicate session does not exist
	return err != sql.ErrNoRows

}

//EndSession ends a session
// receives UUID to terminate
// returns error (nil if session ended successfully)
func EndSession(UUID string) {

	// set UUID to null to indicate no active session
	db.Exec("UPDATE users SET sessionID=NULL WHERE sessionID=?", UUID)

	//don't need to check for err, because if UUID is not there, update still succeeds (trivially)

}
