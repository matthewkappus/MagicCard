package main

import (
	"log"

	"github.com/matthewkappus/MagicCard/src/db"
)

// //go:embed tmpl/*tmpl.html
// var tmpls embed.FS

func main() {

	store, err := db.OpenStore("data/cards.db")
	if err != nil {
		log.Fatal(err)
	}

	sbs, err := store.GetStarBars("Kappus, Matthew D.")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(sbs)

	// sv, err := roster.NewView(store)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// http.HandleFunc("/login", sv.Login)

	// http.HandleFunc("/matty", roster.Matty)
	// http.HandleFunc("/search", sv.TeacherLock(sv.Search))
	// http.HandleFunc("/classes", sv.TeacherLock(sv.ListClasses))
	// http.HandleFunc("/class", sv.TeacherLock(sv.Class))
	// http.HandleFunc("/addComment", sv.TeacherLock(sv.Add))
	// http.HandleFunc("/card", sv.TeacherLock(sv.Card))

	// http.HandleFunc("/", sv.Home)

	// http.ListenAndServe(":8080", nil)
}
