package roster

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Comment struct {
	ID int
	//
	PermID string
	// staff(name)
	Author  string
	Comment string

	Created time.Time
	// max comment: 280
	IsMerrit bool
	IsActive bool
}

func (sv *StudentView) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	c, err := sv.insertComment(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// display confirmation
	sv.tmpls.Lookup("add.tmpl.html").Execute(w, fmt.Sprintf("awarded merrit:\n%v", c))
}

func (sv *StudentView) insertComment(r *http.Request) (*Comment, error) {
	// validate input: session, msg
	// insert db & log change
	// INSERT INTO comment(perm_id, email, comment, is_merrit, is_active) VALUES(?,?,?,?,?)
	r.ParseForm()
	isMerrit := false
	if r.FormValue("isMerrit") == "true" {
		isMerrit = true
	}

	name, err := r.Cookie("name")
	if err != nil {
		log.Printf("insertComment without name. error: %v", err)
	}
	c := &Comment{
		PermID: r.PostFormValue("permID"),
		// can this get from the session?
		Author:   name.Value,
		Comment:  r.PostFormValue("comment"),
		IsMerrit: isMerrit,
	}

	sv.store.InsertComment(c.PermID, c.Author, c.Comment, c.IsMerrit)
	return c, nil

}
