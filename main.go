package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
)

var (
	httpPort string
	httpAddr string
)

func init() {
	flag.StringVar(&httpPort, "http.port", "8000", "-http.port=8000")
	flag.StringVar(&httpAddr, "http.addr", "", "-http.addr=localhost")
}

var budget *Budget

func main() {
	budget = newBudget()
	http.Handle("/", http.HandlerFunc(budgetHandler))
	log.Fatal(http.ListenAndServe(httpAddr+":"+httpPort, nil))
}

func newBudget() *Budget {
	return &Budget{
		Einnahmen: make(map[string]float64),
		Ausgaben:  make(map[string]float64),
	}
}

type Budget struct {
	Einnahmen map[string]float64
	Ausgaben  map[string]float64
}

func MustString(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}

var index = template.Must(template.New("index").ParseFiles("./index.tmpl.html"))

func budgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := index.ExecuteTemplate(w, "index.tmpl.html", budget)
		if err != nil {
			log.Printf("Error while executing template: %T: %v\n", err, err)
		}
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error while parsing form: %T: %v\n")
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
