package database

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// PageData as overarching struct to give to templates
type PageData struct {
	Users          []UserPageEntry //slice containing users' data
	UserActs       []UserActsEntry //slice containing user's activity data
	AllGoals       []AllGoals      //slice containing all goal categories
	AllActs        []AllActs       //slice containing all activities
	Recommendation string          //recommended goal based on survey results
}

// UserPageEntry represents one user each
type UserPageEntry struct {
	ID    int
	UName string
	FName string
	LName string
	Karma int
	Email string
}

// UserActsEntry represents one user's list of activities
type UserActsEntry struct {
	GoalName   string
	Activity   string
	SumKarma   int //total karma for this activity
	IndivKarma int //karma value for a single instance of this activity
}

// AllGoals is for holdings lists of activities
type AllGoals struct {
	GoalID   int
	GoalName string
}

// AllActs is for holdings lists of activities
type AllActs struct {
	ActID    int
	GoalID   int
	Activity string
	Karma    int
}

//LoadAllGoals loads all possible Activities, with their associated goals, and karma
// receives nothing
// returns err, and filled PageData
func (d *PageData) LoadAllGoals() error {

	//fill AllGoals slice with ALL goals and associated data
	rows, err := db.Query(`SELECT id, name FROM goals`)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		p := AllGoals{}

		if err := rows.Scan(
			&p.GoalID,
			&p.GoalName); err != nil {
			return err
		}

		d.AllGoals = append(d.AllGoals, p)
	}
	return nil
}

//LoadAllActs loads all possible Activities, with their associated goals, and karma
// receives nothing
// returns err, and filled PageData
func (d *PageData) LoadAllActs() error {

	//fill AllActs slice with ALL activities and their data
	rows, err := db.Query(`SELECT id, goal_id, activity, karma_value FROM activities`)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		p := AllActs{}

		if err := rows.Scan(
			&p.ActID,
			&p.GoalID,
			&p.Activity,
			&p.Karma); err != nil {
			return err
		}

		d.AllActs = append(d.AllActs, p)
	}
	return nil
}

// LoadUser loads a single user's data into pagedata struct (except password)
// receives user ID as ID
// returns filled d PageData struct
func (d *PageData) LoadUser(ID int) error {

	//query to get data for this user
	row := db.QueryRow(`SELECT id, username, first_name, last_name, email FROM users
	WHERE users.id=?`, ID)

	p := UserPageEntry{}

	if err := row.Scan(
		&p.ID,
		&p.UName,
		&p.FName,
		&p.LName,
		&p.Email); err != nil {
		return err
	}

	//get the user's cumulative karma
	totalKarma, err := GetTotalKarma(ID)
	if err != nil {
		return err
	}
	p.Karma = totalKarma
	d.Users = append(d.Users, p)

	return nil
}

// GetTotalKarma gets the user's cumulative karma for all activities
// returns total karma and err
func GetTotalKarma(ID int) (int, error) {
	var totalKarma int
	err := db.QueryRow(`SELECT SUM(total_karma) FROM users_have_activities WHERE user_id=?`, ID).Scan(&totalKarma)

	return totalKarma, err
}

//AddUserKarma adds the karma value of activity to the user's
//cumulative karma for that activity
func AddUserKarma(userID int, actName string) error {

	indivKarma := 0
	totalKarma := 0
	actID := 0

	//get current karma value for this activity
	err := db.QueryRow(`SELECT karma_value FROM activities WHERE activity=?`, actName).Scan(&indivKarma)

	//get the user's current cumulative value for this activity
	err = db.QueryRow(`SELECT total_karma, a.id FROM users_have_activities uha
	INNER JOIN activities a ON a.id=uha.activity_id 
	WHERE activity=? and uha.user_id=?;`, actName, userID).Scan(&totalKarma, &actID)

	//add one activity instance worth of karma
	totalKarma = totalKarma + indivKarma

	//update current total for this activity
	_, err = db.Exec("UPDATE users_have_activities SET total_karma=? WHERE user_id=? AND activity_id=?", totalKarma, userID, actID)

	return err
}

//AddUserAct adds a chosen activity to a user's list of activities
func AddUserAct(userID int, actName string) error {

	//get the actID from the actName
	var actID int
	err := db.QueryRow("SELECT id FROM activities WHERE activity=?", actName).Scan(&actID)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users_have_activities (user_id, activity_id) VALUES (?, ?)`, userID, actID)

	//1062 is duplicate entry error - that's OK, it just means they already have this activity
	if thisErr, ok := err.(*mysql.MySQLError); ok {
		if thisErr.Number == 1062 {
			return nil
		}
	}
	return err
}

// LoadUserActs loads a single user's activity data into the PageData struct
// receives user ID as ID
// returns filled d PageData struct
func (d *PageData) LoadUserActs(ID int) error {

	//query to get data for this user
	rows, err := db.Query(`SELECT g.name, a.activity, uha.total_karma, a.karma_value
	FROM activities a
	INNER JOIN goals g ON g.id=a.goal_id
	INNER JOIN users_have_activities uha ON uha.activity_id=a.id AND uha.user_id=?`, ID)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		p := UserActsEntry{}

		if err := rows.Scan(
			&p.GoalName,
			&p.Activity,
			&p.SumKarma,
			&p.IndivKarma); err != nil {
			return err
		}

		d.UserActs = append(d.UserActs, p)
	}
	return nil
}

//UserExists returns result of query on username
// receives username string as username
// returns true if user is in database, false if not
func UserExists(username string) bool {

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&username)

	return err != sql.ErrNoRows

}

//CreateUser creates a new user in the users table
// receives username and password (password assumed to be already hashed)
// returns error - error is nil if user was added successfully
func CreateUser(username string, password []byte, firstname string, lastname string, email string) error {

	_, err := db.Exec("INSERT INTO users (username, password, first_name, last_name, email) VALUES(?, ?, ?, ?, ?)", username, password, firstname, lastname, email)

	return err
}

//GetPassword returns password associated with user
// receives username
// returns password and error; error is nil if successfully retrieved
func GetPassword(username string) ([]byte, error) {

	var password []byte
	//query for the hashed password
	err := db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&password)

	return password, err
}

//GetIDFromName returns the ID associated to a username
// receives username
// returns integer of the user ID and error (nil if user successfully found)
func GetIDFromName(username string) (int, error) {

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)

	return userID, err

}

//GetIDFromUUID returns the user ID associated to a UUID
// receives UUID as string
// returns the ID as int and err
// It is caller's responsibility to validate UUID before calling
func GetIDFromUUID(UUID string) (int, error) {

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE sessionID=?", UUID).Scan(&userID)

	return userID, err

}

//GetRecc fills the recommendation field of PageData
// with a string corresponding to a goal
func (d *PageData) GetRecc(score int) {

	//this is in the database package because a better
	//design would be to do actual calculations with more survey data
	//in the db, and make more refined recommendations
	//but this is as close as we can get in the time

	if score < 100 {
		d.Recommendation = "eco-friendly eating"
	} else if score < 250 {
		d.Recommendation = "reducing your mileage"
	} else if score < 700 {
		d.Recommendation = "reducing your water usage"
	} else {
		d.Recommendation = "reducing CO2 emissions"
	}

}
