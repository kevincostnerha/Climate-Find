Before you can run this program, you will need to:
    1. go get golang.org/x/crypto/bcrypt
    2. go get github.com/go-sql-driver/mysql

Then you will be able to run the program by typing: 
    go run signup.go

LOCAL URL: http://localhost:8080/

The database is correctly hooked up to Sarah's so you'll be able to see each new entry there in a table called user_data.
In it's current iteration the table can only hold the user_id(PK), username, and password. 

The HTML pages and CSS sheet were provided by Melissa. I removed some of the additional inputs that she had for the sake of simplicity, on my part, and am more than willing to add them back after we meet on Tuesday and make sure everything we have so far works the way they're supposed too.

EDITS MADE ON 11/19/2018
Pushed up an updated version of the file with a few changes:
    1. the file now has comments so that people new to go can understand the flow of things
    2. created more handlers to catch more errors
    3. the loginPage function no longer directs the user to user.html, but will instead greet the user with a message 
       that reads, "Hello, 'username'". This will allow us to directly see whether or not the user has entered into the 
       correct account or not.
       
       If you want the page to go to the user page when the login is successful, 
       replace the line:
            w.Write([]byte("Hello, " + databaseUsername))
       
       with the lines:
       switch {
	   default:
	       http.Redirect(w, r, "/user", 301)
       }
