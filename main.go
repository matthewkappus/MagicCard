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

	staffView, err := roster.NewView(s, "tmpl/*.tmpl.html", roster.Teacher)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/login", staffView.Login)

	if *devMode {
		http.HandleFunc("/devTeacher", roster.DevTeacherLogin)
		http.HandleFunc("/devStudent", roster.DevStudentLogin)
	}

	// list of students
	staffView.HF("/studentSearch", staffView.Search)

	staffView.HF("/card", staffView.MagicCard)
	staffView.HF("/profile", staffView.Profile)
	staffView.HF("/addComment", staffView.AddComment)

	staffView.HF("/class", staffView.ClassEdit)

	// Admin Tools
	// http.HandleFunc("/admin/addMyStarStrike", staffView.AddMyStarStrikeAll)
	// http.HandleFunc("/admin/myStarStrikeForm", staffView.MyStarStrikeForm)
	// http.HandleFunc("/admin", staffView.AdminHome)

	http.HandleFunc("/", staffView.Home)

	http.ListenAndServe(":8080", nil)

}
