package main
import (
	"net/http"
	// "io"
	"html/template"
	// "text/template"
	"database/sql"
	"html"
	// "os"
	// "time"
	"fmt"
	"strconv"
)

import _ "github.com/go-sql-driver/mysql"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "Welcome to the barn!")
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func rideHandler(w http.ResponseWriter, r *http.Request) {
	// get horses you can ride (currently all horses)

	rows, err := db.Query("SELECT * FROM Horses")
	if err != nil {
	        fmt.Printf("error: %s", err)
	}
	defer rows.Close()

	// struc for horse info
	type Horse struct {
		Id int
		Name string
	}
	
	// list of horses from DB
	type HorseList []Horse
	var myhorselist HorseList
	for rows.Next() {
	        var id int
	        var name string
	        var age int
	        if err := rows.Scan(&id, &name, &age); err != nil {
	                fmt.Printf("error: %s\n", err)
	        }
	        myhorselist = append(myhorselist, Horse{id, name})

	}

	// if err := horses.Err(); err != nil {
	//         fmt.Printf("error: %s", err)
	// }

	// fmt.Printf("horsesLIST: %v\n", myhorselist)
	t, _ := template.ParseFiles("templates/ride.html")
	t.Execute(w, myhorselist)
}

func ridingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("riding yoo...\n")

	r.ParseForm()
	d := r.FormValue("id")
	horse_id, _ := strconv.Atoi(d)

	result, err := db.Exec("INSERT INTO rides(starttime, horse_id) VALUES(NOW(), ?)", horse_id)

	// fmt.Printf("LACROIX: %v\n", result)

	rideId, err := result.LastInsertId()
	if err != nil {
	    println("Error:", err.Error())
	} else {
	    println("LastInsertId:", rideId)
	}

	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add starttime.", 500)
		return
	}

	http.Redirect(w, r, "/", 301)

	// return a ride ID

	// fmt.Printf("ID: %v\n", d)
}

var db *sql.DB
var err error

func newHorseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("party")
	if r.Method != "POST" {
		fmt.Printf("post")
		t, _ := template.ParseFiles("templates/newHorse.html")
		t.Execute(w, nil)
		return
	}

	horseName := html.EscapeString(r.FormValue("horseName"))
	fmt.Printf("%s\n\n\n", horseName)

	// var user string
	// _, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)

	_, err = db.Exec("INSERT INTO horses(Name) VALUES(?)", horseName)
										// INSERT INTO centaur (horse) VALUES ('$horse');
	fmt.Printf("guy")
	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add horse.", 500)
		return
	}

	http.Redirect(w, r, "/", 301)
}


func main() {
	// fmt.Printf("%s", os.Getenv("FOO"))
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/centaur")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/newHorse", newHorseHandler)
	http.HandleFunc("/ride", rideHandler)
	http.HandleFunc("/riding", ridingHandler)
	// http.HandleFunc("/ride", testRide)
	http.ListenAndServe(":8080", nil)
}