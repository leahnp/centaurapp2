package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// get horses you can ride (currently all horses)
	rows, err := db.Query("SELECT * FROM Horses")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer rows.Close()

	// struc for horse info
	type Horse struct {
		Id   int
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

func rideHandler(w http.ResponseWriter, r *http.Request) {
	// get horses you can ride (currently all horses)
	rows, err := db.Query("SELECT * FROM Horses")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer rows.Close()

	// struc for horse info
	type Horse struct {
		Id   int
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

func startRideHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	horse_id := r.FormValue("id")
	// horse_idstr, _ := strconv.Atoi(horse_id)

	// add horse id and starttime to new ride entry
	result, err := db.Exec("INSERT INTO rides(starttime, horse_id) VALUES(NOW(), ?)", horse_id)
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

	http.Redirect(w, r, fmt.Sprintf("/riding?ride_id=%d", rideId), 301)
}

func stopRideHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ride_id := r.FormValue("ride_id")
	ride_idstr, _ := strconv.Atoi(ride_id)

	// // add stoptime to ride_id
	_, err := db.Exec("UPDATE rides SET stoptime=NOW() WHERE id=?", ride_id)
	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add starttime.", 500)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/rideSummary?ride_id=%d", ride_idstr), 301)
}

func rideSummaryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ride_id := r.FormValue("ride_id")

	ride_idint, _ := strconv.Atoi(ride_id)
	// get total ride duration
	ride_duration := rideDuration(ride_idint)
	fmt.Printf("RIDE DURATION: %v\n", ride_duration)
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
	fmt.Printf("%s\n\n\n", horseName)

	_, err = db.Exec("INSERT INTO horses(Name) VALUES(?)", horseName)
	if err != nil {
		fmt.Print("Error: %v", err)
		http.Error(w, "Insert error, unable to add horse.", 500)
		return
	}

	http.Redirect(w, r, "/", 301)
}

func ridingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ride_id := r.FormValue("ride_id")

	t, _ := template.ParseFiles("templates/riding.html")
	t.Execute(w, ride_id)
}

func horseSummaryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	horse_id := r.FormValue("id")
	// horse_idstr, _ := strconv.Atoi(horse_id)

	//select all rides for horse and return date
	// SELECT starttime FROM rides WHERE horse_id=1;
	rows, err := db.Query("SELECT * FROM rides WHERE horse_id=?", horse_id)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer rows.Close()
	// struc for horse info
	type Ride struct {
		Id         int
		Prettydate string
		Duration   time.Duration
	}

	// list of horses from DB
	type RideList []Ride
	var myridelist RideList
	for rows.Next() {
		var id int
		var starttime time.Time
		var stoptime time.Time
		if err := rows.Scan(&id, &horse_id, &starttime, &stoptime); err != nil {
			fmt.Printf("error: %s\n", err)
		}

		// idstr, _ := strconv.Atoi(id)
		// get total ride duration
		ride_duration := rideDuration(id)
		fmt.Printf("RUGS: %v", ride_duration)
		myridelist = append(myridelist, Ride{id, fmt.Sprintf(starttime.Format("Mon Jan _2 15:04:05 2006")), ride_duration})
	}

	fmt.Printf("BLUE: %v\n", myridelist)

	t, _ := template.ParseFiles("templates/horseSummary.html")
	t.Execute(w, myridelist)

}

type ride struct {
	Motion     [][]float32
	Horse_id   int
	Start_time int
	Stop_time  int
}

// test pass
func uploadDataHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var my_ride ride
	err := decoder.Decode(&my_ride)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	fmt.Printf("YEPPERS: %s, %d, %d, %d\n", my_ride.Motion, my_ride.Horse_id, my_ride.Start_time, my_ride.Stop_time)
	// alright put this in the rides table right meow....

	t, _ := template.ParseFiles("templates/rideSummary.html")
	t.Execute(w, nil)
}

func main() {
	// fmt.Printf("%s", os.Getenv("FOO"))
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
	http.HandleFunc("/ride", rideHandler)
	http.HandleFunc("/startride", startRideHandler)
	http.HandleFunc("/riding", ridingHandler)
	http.HandleFunc("/stopRide", stopRideHandler)
	http.HandleFunc("/rideSummary", rideSummaryHandler)
	http.HandleFunc("/horseSummary", horseSummaryHandler)
	http.HandleFunc("/upload_data", uploadDataHandler)
	http.ListenAndServe(":8080", nil)
}
