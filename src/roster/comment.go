package roster

import (
	"fmt"
	"net/http"
	"time"
)

type Comment struct {
	ID int
	//
	PermID string
	// Teacher Email
	Email   string
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

	isMerrit := false
	if r.FormValue("isMerrit") == "true" {
		isMerrit = true
	}

	c := &Comment{
		PermID: r.FormValue("permID"),
		// can this get from the session?
		Email:    "samplesetinAdd@gmail.com",
		Comment:  r.FormValue("comment"),
		IsMerrit: isMerrit,
	}

	sv.store.InsertComment(c.PermID, c.Email, c.Comment, c.IsMerrit)
	return c, nil

}
