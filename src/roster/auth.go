package roster

import (
	"fmt"
	"net/http"
	"net/url"
)

func (sv *StudentView) Login(w http.ResponseWriter, r *http.Request) {
	email := "Lindsey.Berg@aps.edu"

	// todo: check if student number
	teacher, name, err := sv.store.TeacherNameFromEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// todo: set to appengine url
	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: set a session id
	cookies := []*http.Cookie{
		{
			Name:  "name",
			Value: name,
			// cookie lasts a day
			MaxAge: 86400,
		},
		{
			Name:  "teacher",
			Value: teacher,
			// cookie lasts a day
			MaxAge: 86400,
		},

		{
			Name:  "email",
			Value: email,
			// cookie lasts a day
			MaxAge: 86400,
		},
	}

	http.SetCookie(w, cookies[0])
	http.SetCookie(w, cookies[1])
	http.SetCookie(w, cookies[2])
	fmt.Printf("set cookies for %v:\n%v\n", u, cookies)

	// redirect to /classes or /card (student)

	http.Redirect(w, r, "/classes", http.StatusTemporaryRedirect)
}
