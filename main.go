package main
import (
	"net/http"
	"html/template"
	"database/sql"
	"html"
	"time"
	"fmt"
	"strconv"
	"io/ioutil"
)

import _ "github.com/go-sql-driver/mysql"

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, myhorselist)
}

func uploadDataHandler(w http.ResponseWriter, r *http.Request) {
  ride_id := r.URL.RawQuery;
  body_bytes, _ := ioutil.ReadAll(r.Body)
  body_str := string(body_bytes)
  defer r.Body.Close()

  fmt.Printf("YEPPERS (%d): %s\n", ride_id, body_str)
  // alright put this in the rides table right meow....

	_, err2 := db.Exec("UPDATE rides SET motion=JSON_MERGE(motion, ?) WHERE id=?", body_str, ride_id)
	if err2 != nil {
		fmt.Print("Errorhere: %v", err2)
		http.Error(w, "Insert error, unable to add starttime.", 500)
		return
	}
	// t, _ := template.ParseFiles("templates/rideSummary.html")
	// t.Execute(w, nil)
}

func startRideHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	horse_id := r.FormValue("id")

	// add horse id and starttime to new ride entry
	result, err := db.Exec("INSERT INTO rides(horse_id, starttime, motion) VALUES(?, NOW(), '[]')", horse_id)
	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add starttime.", 500)
		return
	}

	// get id of ride
	rideId, err := result.LastInsertId()
	if err != nil {
	    println("Error: unable to get last inserted id from db", err.Error())
	} else {
	    println("LastInsertId:", rideId)
	}

	t, _ := template.ParseFiles("templates/riding.html")
	t.Execute(w, rideId)
}

// func testHandler(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("templates/test.html")
// 	t.Execute(w, r)	
// }

func stopRideHandler(w http.ResponseWriter, r *http.Request) {
	ride_id := r.URL.RawQuery;

	// add stoptime to ride_id
	_, err := db.Exec("UPDATE rides SET stoptime=NOW() WHERE id=?", ride_id)
	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add stoptime.", 500)
		return
	}
	// http.Redirect(w, r, fmt.Sprintf("/rideSummary?ride_id=%d", ride_id), 301)
}

func rideSummaryHandler(w http.ResponseWriter, r *http.Request) {
	ride_id := r.URL.RawQuery;
	ride_idint, _ := strconv.Atoi(ride_id)
	// get total ride duration
	ride_duration := rideDuration(ride_idint)
	t, _ := template.ParseFiles("templates/rideSummary.html")
	t.Execute(w, ride_duration)
}

func rideDuration(ride_id int) time.Duration {
	// query db and get starttime and stoptime for ride
	rows, err := db.Query("SELECT starttime, stoptime FROM rides WHERE id=?", ride_id)
	if err != nil {
	  fmt.Printf("error: %s", err)
	}
	defer rows.Close()
	var delta time.Duration

	for rows.Next() {
    var starttime time.Time
    var stoptime time.Time

    if err := rows.Scan(&starttime, &stoptime); err != nil {
      fmt.Printf("error: %s\n", err)
    }
    // get ride duration
    fmt.Printf("Start %v, stop %v", starttime, stoptime)
    delta = stoptime.Sub(starttime)
	}
  return delta
}

var db *sql.DB
var err error

func newHorseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Printf("post")
		t, _ := template.ParseFiles("templates/newHorse.html")
		t.Execute(w, nil)
		return
	}
	horseName := html.EscapeString(r.FormValue("horseName"))

	_, err = db.Exec("INSERT INTO horses(Name) VALUES(?)", horseName)
	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add horse.", 500)
		return
	}

	http.Redirect(w, r, "/", 301)
}


func horseSummaryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	horse_id := r.FormValue("id")
	//select all rides for horse and return date
	rows, err := db.Query("SELECT * FROM rides WHERE horse_id=?", horse_id)
	if err != nil {
	  fmt.Printf("error: %s", err)
	}
	defer rows.Close()
	// struc for horse info
	type Ride struct {
		Id int
		Prettydate string
		Duration time.Duration
	}
	
	// list of horses from DB
	type RideList []Ride
	var myridelist RideList
	for rows.Next() {
		var id int
		var starttime time.Time
		var stoptime time.Time
		var motion string
    if err := rows.Scan(&id, &horse_id, &starttime, &stoptime, &motion); err != nil {
      fmt.Printf("error: %s\n", err)
    }

		// get total ride duration
		ride_duration := rideDuration(id)
    myridelist = append(myridelist, Ride{id, fmt.Sprintf(starttime.Format("Mon Jan _2 15:04:05 2006")), ride_duration})
	}
	t, _ := template.ParseFiles("templates/horseSummary.html")
	t.Execute(w, myridelist)	
}


func main() {
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/centaur?parseTime=true")
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
	// http.HandleFunc("/test", testHandler)
	http.HandleFunc("/startride", startRideHandler)
	http.HandleFunc("/stopRide", stopRideHandler)
	http.HandleFunc("/rideSummary", rideSummaryHandler)
	http.HandleFunc("/horseSummary", horseSummaryHandler)
	http.HandleFunc("/upload_data", uploadDataHandler)
	http.ListenAndServe(":8080", nil)
}