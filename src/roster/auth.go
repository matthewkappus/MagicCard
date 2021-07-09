package roster

import (
	"fmt"
	"net/http"
)

func (sv *StudentView) TeacherLock(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sv.isTeacher(r) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

// looks for teacher and key cookie, matches with staff db
func (sv *StudentView) isTeacher(r *http.Request) bool {
	teacherCookie, err := r.Cookie("teacher")
	if err != nil {
		return false
	}

	keyCookie, err := r.Cookie("key")
	if err != nil {
		return false
	}

	dbKey, err := sv.store.GetKeyByTeacher(teacherCookie.Value)
	if err != nil {
		fmt.Printf("error looking up key %s: %v", teacherCookie.Value, err)
		return false
	}
	fmt.Printf("isTeacher comparing cooKey: %s to staff key %s", keyCookie.Value, dbKey)
	return keyCookie.Value == dbKey
}
func (sv *StudentView) Login(w http.ResponseWriter, r *http.Request) {
	email := "Lindsey.Berg@aps.edu"

	// todo: check if student number
	teacher, name, key, err := sv.store.TeacherNameFromEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// todo: set a session id
	cookies := []*http.Cookie{
		{
			Name:  "name",
			Value: name,
			// cookie lasts a day
			MaxAge: 40000,
		},
		{
			Name:  "teacher",
			Value: teacher,
			// cookie lasts a day
			MaxAge: 40000,
		},

		{
			Name:  "key",
			Value: key,
			// cookie lasts a day
			MaxAge: 40000,
		},
	}

	http.SetCookie(w, cookies[0])
	http.SetCookie(w, cookies[1])
	http.SetCookie(w, cookies[2])

	// redirect to /classes or /card (student)

	http.Redirect(w, r, "/classes", http.StatusTemporaryRedirect)
}
