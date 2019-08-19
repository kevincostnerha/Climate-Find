package main

import (
	"climatefind/database"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func trackersPage(w http.ResponseWriter, r *http.Request) {

	//get the UUID out of the cookie
	UUID, err := getUUID(r)
	if err != nil || UUID == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// validate the user's session
	if !database.ValidateSession(UUID) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	//fill this with data the trackers page needs
	var pageData database.PageData

	//have to figure out which user we need activities for
	userID, err := database.GetIDFromUUID(UUID)

	if err != nil {
		http.Error(w, "Unable to find that user", 500)
		return
	}

	//load it up!
	pageData.LoadUser(userID)
	pageData.LoadUserActs(userID)
	pageData.LoadAllGoals()
	pageData.LoadAllActs()
	//execute the template with pageData
	tmpl.ExecuteTemplate(w, "tracker.html", pageData)
}

func recommendationPage(w http.ResponseWriter, r *http.Request) {

	//get the UUID out of the cookie
	UUID, err := getUUID(r)
	if err != nil || UUID == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// validate the user's session
	if !database.ValidateSession(UUID) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	var score, diet, transport, water, energy int
	var pageData database.PageData

	//get values out of forms - bail out if bad values were entered
	if diet, err = strconv.Atoi(r.FormValue("r1")); err != nil {
		diet = 0
	}
	if transport, err = strconv.Atoi(r.FormValue("r2")); err != nil {
		transport = 0
	}
	if water, err = strconv.Atoi(r.FormValue("r3")); err != nil {
		water = 0
	}
	if energy, err = strconv.Atoi(r.FormValue("r4")); err != nil {
		energy = 0
	}

	score = diet + transport + water + energy

	//if they didn't fill out enough questions, or tried to go here
	//by entering it into the address bar, redirect to survey
	if score < 4 {
		http.Redirect(w, r, "/survey", http.StatusFound)
	}

	//have to figure out which user we need activities for
	userID, err := database.GetIDFromUUID(UUID)

	if err != nil {
		http.Error(w, "Unable to find that user", 500)
		return
	}

	//fill pageData recc with string of recommmended goal
	pageData.GetRecc(score)

	//fill pageData with all the activities and their goals
	pageData.LoadAllGoals()
	pageData.LoadAllActs()
	//fill pageData with all of this user's activities and karma
	pageData.LoadUserActs(userID)
	//fill pageData with the user's data
	pageData.LoadUser(userID)

	tmpl.ExecuteTemplate(w, "recommendation.html", pageData)

}

//adds karma for the activity the user enacted
//then executes the trackerTemplate
func doActivity(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/trackers", http.StatusFound)
		return
	}
	//get the UUID out of the cookie
	UUID, err := getUUID(r)
	if err != nil || UUID == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// validate the user's session
	if !database.ValidateSession(UUID) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	//have to figure out which user we need activities for
	userID, err := database.GetIDFromUUID(UUID)

	//get the activity they want to do out of the form
	actName := r.FormValue("doAct")

	//add karma using a db package func
	if err := database.AddUserKarma(userID, actName); err != nil {
		http.Error(w, "Could not add activity for user", 500)
		return
	}

	http.Redirect(w, r, "/trackers", http.StatusFound)

}

//adds the activity the user chose
//then executes the trackerTemplate
func addActivity(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/trackers", http.StatusFound)
		return
	}

	//get the UUID out of the cookie
	UUID, err := getUUID(r)
	if err != nil || UUID == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// validate the user's session
	if !database.ValidateSession(UUID) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	//have to figure out which user we need activities for
	userID, err := database.GetIDFromUUID(UUID)

	if err != nil {
		http.Error(w, "Unable to find that user", 500)
		return
	}

	//get the activity they want to add out of the form
	actName := r.FormValue("addAct")

	//add an activity for the user using a db package func
	if err := database.AddUserAct(userID, actName); err != nil {
		http.Error(w, "Could not add activity for user", 500)
		return
	}

	http.Redirect(w, r, "/trackers", http.StatusFound)

}
