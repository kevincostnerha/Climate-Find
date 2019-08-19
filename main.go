package main

import (
	"climatefind/database"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl = template.Must(template.ParseGlob(filepath.Join("tmpl", "*.html")))

// function to signup for an account
func homePage(w http.ResponseWriter, r *http.Request) {
	if !validateLogin(r) {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	tmpl.ExecuteTemplate(w, "homepage.html", "")
}

func redirectHome(w http.ResponseWriter, r *http.Request) {
	if !validateLogin(r) {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func main() {

	database.Connect()
	defer database.DestroyConnection()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", redirectHome)
	http.HandleFunc("/index.html", homePage)

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/user", userPage)
	http.HandleFunc("/survey", surveyPage)
	http.HandleFunc("/video", videoPage)
	http.HandleFunc("/trackers", trackersPage)
	http.HandleFunc("/recommendation", recommendationPage)
	http.HandleFunc("/addActivity", addActivity)
	http.HandleFunc("/doActivity", doActivity)

	fmt.Println("Launching on 0.0.0.0:8110, Process Id", os.Getpid())
	log.Fatal(http.ListenAndServe("0.0.0.0:8110", nil))
}
