package roster

import (
	"fmt"
	"net/http"
	"time"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func (sv *StaffView) AddComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	c := &comment.Card{
		PermID: r.PostFormValue("permID"),
		// can this get from the session?
		Teacher:  sv.GetTeacher(r),
		Title:    r.PostFormValue("title"),
		Comment:  r.PostFormValue("comment"),
		IsStar:   r.PostFormValue("isStar") == "true",
		Created:  time.Now(),
		IsActive: true,
	}

	err := sv.store.InsertComment(c.PermID, c.Teacher, c.Comment, c.Title, c.IsStar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Created comment: %#v\n", *c)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
