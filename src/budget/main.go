package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	budget    *Budget
	httpPort  string
	httpAddr  string = "localhost"
	sessions         = make(SessionManager)
	debugMode bool   = false
	funcMap          = template.FuncMap{
		"startsWith": strings.HasPrefix,
	}
	templatePath           string
	layoutPattern          string
	partialTemplatePattern string
	layout                 *template.Template
	partialTmpl            *template.Template
)

func mustString(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}

func init() {
	flag.BoolVar(&debugMode, "debug", false, "-debug")
	flag.StringVar(&httpPort, "http.port", "8000", "-http.port=8000")
	defaultPath := filepath.Join(mustString(os.Getwd()), "templates/")
	flag.StringVar(&templatePath, "template.path", defaultPath, fmt.Sprintf("template.path=%s", defaultPath))
	flag.Parse()
	layoutPattern = filepath.Join(templatePath, "layout.html.tmpl")
	partialTemplatePattern = filepath.Join(templatePath, "_*.tmpl")
	layout = template.Must(template.ParseGlob(layoutPattern)).Funcs(funcMap)
	partialTmpl = template.Must(layout.ParseGlob(partialTemplatePattern))
}

func main() {
	googleClientId := os.Getenv("GOOGLE_CLIENTID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENTSECRET")
	hostAddress := strings.TrimLeft(strings.TrimSuffix(httpAddr, ":"), "https://") + ":" + strings.TrimPrefix(httpPort, ":")
	host := "http://" + hostAddress
	hostEnv := os.Getenv("HOST")
	if hostEnv != "" {
		host = hostEnv
	}
	googleOauth2Config := &oauth2.Config{
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  host + "/google_oauth2redirect",
	}
	budget = NewBudget()
	http.Handle("/", http.HandlerFunc(indexHandler()))
	http.Handle("/favicon.ico", http.HandlerFunc(fileHandler))
	http.Handle("/login", getHandler(loginHandler()))
	http.Handle("/logout", getHandler(authHandler(logoutHandler())))
	http.HandleFunc("/google_login", googleLoginHandler(googleOauth2Config))
	http.HandleFunc("/google_oauth2redirect", getHandler(googleRedirectHandler(googleOauth2Config)))
	log.Printf("Listening on port %s\n", httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}

func MustString(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}
func fileHandler(w http.ResponseWriter, r *http.Request) {}

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("budget")
		session := GetSessionFromCookie(cookie)
		if r.Method == "GET" {
			if r.URL.Path != "/" {
				if session != nil {
					session.URL = r.URL
					session.AddError(fmt.Errorf("Die eingegebene Seite existiert nicht: '%s'", r.URL.Path))
				}
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			var page *PageObject
			if session != nil {
				session.URL = r.URL
				page = PageForSession(session)
				page.Set("Budget", budget)
			} else {
				page = NewPageObject()
				page.Set("RequestPath", "")
			}
			renderTemplate(w, "index", page)
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
				budget.Ausgaben[key] = NewBudgetColumnEntry(value)
			} else {
				budget.Einnahmen[key] = NewBudgetColumnEntry(value)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
}

func loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderTemplate(w, "login", &PageObject{})
			return
		}
	}
}

type authHandlerFunc func(http.ResponseWriter, *http.Request, *Session)

func logoutHandler() authHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *Session) {
		sessions.Remove(s)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func authHandler(fn authHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("budget")
		if err != nil {
			log.Printf("No cookie found: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		session := GetSessionFromCookie(cookie)
		if session == nil {
			log.Printf("No session found\n")
			r.Header.Set("X-Referer", r.URL.String())
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		session.URL = r.URL
		fn(w, r, session)
	}
}

func getHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
}

var contentTemplateString = `{{define "content"}}{{template "%s" .}}{{end}}`

func renderTemplate(w http.ResponseWriter, tmpl string, page *PageObject) {
	formattedTemplateString := fmt.Sprintf(contentTemplateString, tmpl)
	contentTemplate := template.Must(template.Must(layout.Clone()).Parse(formattedTemplateString))
	var buf bytes.Buffer
	var err error
	// TODO(mw): this is a really dirty hack to use this function with no pageObject
	if page == nil {
		err = contentTemplate.Execute(&buf, nil)
	} else {
		err = contentTemplate.Execute(&buf, page)
	}
	if err != nil {
		log.Printf("Template error(%T): %v\n", err, err)
		debug.PrintStack()
		http.Error(w, fmt.Sprintf("%T: %v\n", err, err), http.StatusInternalServerError)
	} else {
		io.Copy(w, &buf)
	}
}
