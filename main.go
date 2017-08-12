package main
import (
	"net/http"
	"io"
	"html/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "Welcome to the barn!")
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func rideHandler(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "Let's ride!")
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}


func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ride", rideHandler)
	http.ListenAndServe(":8080", nil)
}