package roster

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Scopes: 0 Guest 1 Student 2 Teacher 3 Admin
type Scope int

const (
	Guest Scope = iota
	Student
	Teacher
	Admin
)

type Session struct {
	// student perm or teacher name
	User    string
	Expires time.Time
	Scope   int
}

func (v *View) GetSession(sid string) (s *Session, err error) {
	s = new(Session)
	s.User, _, s.Expires, s.Scope, err = v.store.GetSession(sid)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (v *View) StartSession(user string, name string, s Scope, w http.ResponseWriter, r *http.Request) {
	if _, err := v.GetSession(name); err == nil {
		http.Error(w, "Session already exists", http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// delte old sessions
	v.store.DeleteOldSessions(user)
	// create new session
	sid := uuid.NewString()
	exp := time.Now().Add(8 * time.Hour)
	var scope string
	var scopeInt int
	switch s {
	case Guest:
		scope = "0"
		scopeInt = 0
	case Student:
		scope = "1"
		scopeInt = 1
	case Teacher:
		scope = "2"
		scopeInt = 2
	case Admin:
		scope = "3"
		scopeInt = 3
	}

	if err := v.store.InsertSession(user, sid, exp, scopeInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "sid", Value: sid, Expires: exp})
	http.SetCookie(w, &http.Cookie{Name: "name", Value: name, Expires: exp})
	http.SetCookie(w, &http.Cookie{Name: "user", Value: user, Expires: exp})
	http.SetCookie(w, &http.Cookie{Name: "scope", Value: scope, Expires: exp})

}


func (v  *View) EndSession(w http.ResponseWriter, r *http.Request)  {
	http.SetCookie(w, &http.Cookie{Name: "sid", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "name", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "user", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "scope", MaxAge: -1})

	http.Redirect(w, r, "/", http.StatusFound)	
}