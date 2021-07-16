package roster

import (
	"fmt"
	"net/http"
	"time"
)

type Comment struct {
	ID int `json:"id,omitempty"`
	//
	PermID string `json:"perm_id,omitempty"`
	// staff(name)
	Teacher string `json:"teacher,omitempty"`
	Comment string `json:"comment,omitempty"`
	// Title is a catagory of the comment
	Title string `json:"title,omitempty"`

	Created time.Time `json:"created,omitempty"`
	// max comment: 280
	IsStar   bool `json:"is_star,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
}

// IsValid checks for valid field sizes before committing to db
func (c *Comment) IsValid() bool {
	return c.PermID != "" && c.Teacher != "" && c.Comment != ""
}

func (sv *StaffView) AddComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	c := &Comment{
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
	http.Redirect(w, r, "/", http.StatusAccepted)
}
