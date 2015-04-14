package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
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
		params := r.PostForm
		outParam := params.Get("out")
		if outParam == "" {
			log.Printf("Error while parsing form: 'out' can't be blank")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		out, err := strconv.ParseBool(outParam)
		if err != nil {
			log.Printf("Error while parsing form: %T: %v\n")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		key := params.Get("name")
		valueParam := params.Get("value")
		value, err := strconv.ParseFloat(valueParam, 64)
		if err != nil {
			log.Printf("Error while parsing form: %T: %v\n")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		if out {
			budget.Ausgaben[key] = value
		} else {
			budget.Einnahmen[key] = value
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
