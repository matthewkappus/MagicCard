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

	http.HandleFunc("/search", sv.TeacherLock(sv.Search))
	http.HandleFunc("/class", sv.TeacherLock(sv.ClassEdit))
	http.HandleFunc("/classComment", sv.TeacherLock(sv.ClassAddComment))

	http.HandleFunc("/profile", sv.TeacherLock(sv.Profile))
	http.HandleFunc("/starbaredit", sv.TeacherLock(sv.StarBarEdit))
	http.HandleFunc("/starbardelete", sv.TeacherLock(sv.StarBarDelete))
	http.HandleFunc("/starbarcreate", sv.TeacherLock(sv.StarBarCreate))

	http.HandleFunc("/addComment", sv.TeacherLock(sv.AddComment))
	http.HandleFunc("/card", sv.TeacherLock(sv.Card))

	http.HandleFunc("/", sv.Home)

	http.ListenAndServe(":8080", nil)

}
