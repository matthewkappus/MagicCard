package main

import (
	"log"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/MagicCard/src/roster"
)

// //go:embed tmpl/*tmpl.html
// var tmpls embed.FS

func main() {

	store, err := db.OpenStore("data/cards.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	sv, err := roster.NewView(store)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/login", sv.Login)

	http.HandleFunc("/search", sv.TeacherLock(sv.Search))
	// http.HandleFunc("/classes", sv.TeacherLock(sv.ListClasses))
	http.HandleFunc("/matty", roster.Matty)
	http.HandleFunc("/class", sv.TeacherLock(sv.ClassEdit))
	// todo: omni lock (teacher/student can access)

	http.HandleFunc("/profile", sv.TeacherLock(sv.Profile))
	http.HandleFunc("/starbaredit", sv.TeacherLock(sv.StarBarEdit))
	http.HandleFunc("/starbardelete", sv.TeacherLock(sv.StarBarDelete))
	http.HandleFunc("/starbarcreate", sv.TeacherLock(sv.StarBarCreate))

	http.HandleFunc("/addComment", sv.TeacherLock(sv.Add))
	http.HandleFunc("/card", sv.TeacherLock(sv.Card))

	http.HandleFunc("/", sv.Home)

	http.ListenAndServe(":8080", nil)
}
