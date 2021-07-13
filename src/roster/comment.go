package roster

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Comment struct {
	ID int `json:"id,omitempty"`
	//
	PermID string `json:"perm_id,omitempty"`
	// staff(name)
	Author  string `json:"author,omitempty"`
	Comment string `json:"comment,omitempty"`

	Created time.Time `json:"created,omitempty"`
	// max comment: 280
	IsMerrit bool `json:"is_merrit,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
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
