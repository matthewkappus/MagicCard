package roster

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type StudentView struct {
	tmpls *template.Template
	// store sql.DB
}

// New view takes a roster db and tmpl path and returns handler object
func NewView() (*StudentView, error) {
	tmpls, err := template.ParseGlob("tmpl/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// todo: load roster from sql

	return &StudentView{
		tmpls: tmpls,
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

}

func (sv *StudentView) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// validate input: session, msg

	// insert db & log change

	// display confirmation
	sv.tmpls.Lookup("add.tmpl.html").Execute(w, fmt.Sprintf("%s awarded merrit with message:\n%s", r.FormValue("firstName"), r.FormValue("comment")))
}

func (sv *StudentView) Card(w http.ResponseWriter, r *http.Request) {
	// todo: look up student in db: join with comments
	id := r.FormValue("id")
	if len(id) != 9 {
		http.NotFound(w, r)
		return
	}

	stu401s, _ := UpdateRoster()
	var found *Stu401

	for _, stu := range stu401s {
		if stu.PermID == id {
			found = stu
			break
		}

	}

	// todo: add email and session info
	sv.tmpls.Lookup("card.tmpl.html").Execute(w, found)
}
func UpdateRoster() ([]*Stu401, error) {
	// todo: get from synergy
	f, err := os.Open("data/stu401.csv")
	if err != nil {
		log.Fatal(err)
	}

	return ReadStu401sFromCSV(f)

}
