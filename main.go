package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/MagicCard/src/roster"

	"flag"
)

var devMode = flag.Bool("dev", false, "Run in development mode")

func main() {

	flag.Parse()

	var s *db.Store
	var err error

	if *devMode {
		fmt.Println("starting dev mode")
		s, err = db.OpenStore("data/cards.db")
	} else {
		s, err = db.OpenCloudStore()
	}

	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	sv, err := roster.NewView(s)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/login", sv.Login)

	if *devMode {
		fmt.Println("getting temp access")
		http.HandleFunc("/temp", roster.Temp)
	}

	// list of students
	http.HandleFunc("/students", sv.TeacherLock(sv.Search))

	// new starststrike handlers
	http.HandleFunc("/card", sv.MagicCard)

	http.HandleFunc("/teacher", sv.TeacherLock(sv.Profile))

	http.HandleFunc("/", sv.Home)

	http.ListenAndServe(":8080", nil)

}
