package main

import (
	"climatefind/database"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//gets the UUID out of the cookie and returns it
func getUUID(r *http.Request) (string, error) {
	c, err := r.Cookie("UUID")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func validateLogin(r *http.Request) bool {
	//get the cookie
	UUID, err := getUUID(r)
	if err != nil {
		return false
	}

	//validate the session
	if !database.ValidateSession(UUID) {
		return false
	}
	return true
}

// function to signup for an account
func signupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "signup.html", "")
		return
	}

	// initialize the values of username and password
	username := r.FormValue("username")
	password := r.FormValue("password")
	fname := r.FormValue("firstname")
	lname := r.FormValue("lastname")
	email := r.FormValue("emailaddr")

	// if blank entries, go to catchall error
	if len(username) == 0 || len(password) == 0 {
		http.Error(w, "Invalid input.", 500)
		return
	}
	// if user exists, don't make a new account
	if database.UserExists(username) {
		http.Error(w, "That username has been taken.", 500)
		return
	}
	// make an account
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error, unable to create your account.", 500)
		return
	}
	// insert username and hashed password into database
	err = database.CreateUser(username, hashedPassword, fname, lname, email)
	if err != nil {
		http.Error(w, "Database error, unable to create your account.", 500)
		return
	}

	// start a session for the user
	// get this user's ID
	userID, err := database.GetIDFromName(username)
	if err != nil {
		http.Error(w, "Could not find user.", 500)
		return
	}
	// start the session for this user
	UUID, err := database.StartSession(userID)
	if err != nil {
		http.Error(w, "Could not start session.", 500)
		return
	}

	// set the cookie
	c := http.Cookie{
		Name:  "UUID",
		Value: UUID,
	}
	http.SetCookie(w, &c)

	// direct the new user to the video tutorial
	http.Redirect(w, r, "/video", http.StatusFound)

}

// function to login to account
func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "login.html", "")
		return
	}

	// initialize the values of username and password
	username := r.FormValue("username")
	password := r.FormValue("password")

	// get the hashed password from the database
	databasePassword, err := database.GetPassword(username)

	// if both values are NULL redirect user to error message
	if err != nil {
		http.Error(w, "NULL username or password.", 500)
		//http.Redirect(w, r, "/login", 301)
		return
	}

	// incorrect input will print error message
	err = bcrypt.CompareHashAndPassword(databasePassword, []byte(password))
	if err != nil {
		http.Error(w, "Wrong username or password.", 500)
		return
	}

	// user has entered correct username and password, so start a session for them
	// get this user's ID
	userID, err := database.GetIDFromName(username)
	if err != nil {
		http.Error(w, "Could not find user.", 500)
		return
	}

	// start the session for this user
	UUID, err := database.StartSession(userID)
	if err != nil {
		http.Error(w, "Could not start session.", 500)
		return
	}

	// set the cookie
	c := http.Cookie{
		Name:  "UUID",
		Value: UUID,
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/user", 301)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {

	//get the cookie
	UUID, err := getUUID(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}

	//end the session
	database.EndSession(UUID)

	//redirect to login
	http.Redirect(w, r, "/login", 302)

}

func userPage(w http.ResponseWriter, r *http.Request) {
	if !validateLogin(r) {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "user.html", "")
		return
	}
}

// just executes the video template - assumed only new users are directed here
func videoPage(w http.ResponseWriter, r *http.Request) {

	//pageData is empty, if data needed from db
	//fill by calling a function from queries.go
	//in the GET or POST below before executing template
	var pageData database.PageData

	if !validateLogin(r) {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	tmpl.ExecuteTemplate(w, "video.html", pageData) //user is NEW, send to tutorial video
}

func surveyPage(w http.ResponseWriter, r *http.Request) {
	if !validateLogin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "survey.html", "")
		return
	}
}
