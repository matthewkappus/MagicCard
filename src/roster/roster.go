package roster

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/matthewkappus/MagicCard/src/db"
	"github.com/matthewkappus/Roster/src/synergy"
)

type StudentView struct {
	tmpls *template.Template
	store *db.Store
	// store sql.DB
}

// NewView takes a roster db and tmpl path and returns handler object
func NewView(store *db.Store) (*StudentView, error) {
	tmpls, err := template.ParseGlob("tmpl/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	return &StudentView{
		tmpls: tmpls,
		store: store,
	}, nil

}

func (sv *StudentView) Search(w http.ResponseWriter, r *http.Request) {
	stu401s, err := UpdateRoster()
	if err != nil {
		http.Error(w, "Could not get students\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	sv.tmpls.Lookup("search.tmpl.html").Execute(w, stu401s)
}

func (sv *StudentView) Home(w http.ResponseWriter, r *http.Request) {
	sv.tmpls.Lookup("login.tmpl.html").Execute(w, nil)
}

func (sv *StudentView) Card(w http.ResponseWriter, r *http.Request) {
	// todo: look up student in db: join with comments
	permid := r.FormValue("id")
	if len(permid) != 9 {
		http.NotFound(w, r)
		return
	}

	stu415, err := sv.store.SelectStu415(permid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card.tmpl.html").Execute(w, stu415)
}
func UpdateRoster() ([]*synergy.Stu415, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu415.csv")
	if err != nil {
		log.Fatal(err)
	}

	return synergy.ReadStu415sFromCSV(f)

}
