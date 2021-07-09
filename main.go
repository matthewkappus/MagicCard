package main

import (
	"log"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/MagicCard/src/roster"
)

func main() {

	store, err := db.OpenStore("data/cards.db")
	if err != nil {
		log.Fatal(err)
	}

	// if err = store.UpdateStu415("data/stu415.csv"); err != nil {
	// 	log.Fatal(err)
	// }

	sv, err := roster.NewView(store)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/search", sv.Search)
	http.HandleFunc("/classes", sv.ListClasses)
	http.HandleFunc("/class", sv.Class)
	http.HandleFunc("/addComment", sv.Add)
	http.HandleFunc("/card", sv.Card)
	// login
	http.HandleFunc("/", sv.Home)

	http.ListenAndServe(":8080", nil)
}
