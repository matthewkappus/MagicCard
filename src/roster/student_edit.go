package roster

import (
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

// Student holds the Stu415 and Starbars for a template
//  Stars and Strikes are mapped to their title for counting
type Student struct {
	S *synergy.Stu415
	// Star.Title to SS
	StarMap map[string][]*comment.StarStrike
	// StrikeMap
	StrikeMap map[string][]*comment.StarStrike
}

func (sv *StaffView) MakeStudent(stu *synergy.Stu415) (*Student, error) {
	sss, err := sv.store.GetStarStrikesByPerm(stu.PermID)
	if err != nil {
		return nil, err
	}

	starsM := make(map[string][]*comment.StarStrike)
	strikesM := make(map[string][]*comment.StarStrike)
	for _, ss := range sss {
		// 0 star 1 minor 2 strik 3 major
		if ss.Cat == comment.Star {
			// todo: make titles lowercase, space-trimmed
			stars, found := starsM[ss.Title]
			if !found {
				starsM[ss.Title] = []*comment.StarStrike{ss}
				continue
			}
			// title exists: add to list

			stars = append(stars, ss)
			starsM[ss.Title] = stars
		} else {
			// ss.Cat is 1+: A strike
			strikes, found := strikesM[ss.Title]
			if !found {
				strikesM[ss.Title] = []*comment.StarStrike{ss}
				continue
			}
			// strike exists: add to rest
			strikes = append(strikes, ss)
			strikesM[ss.Title] = strikes
		}
	}
	return &Student{
		S:         stu,
		StarMap:   starsM,
		StrikeMap: strikesM}, nil
}

// MagicCard shows students StarStrikes and allows you to edit
func (sv *StaffView) MagicCard(w http.ResponseWriter, r *http.Request) {
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

	classes, _ := sv.ClassList(r)
	teacher := sv.GetTeacher(r)
	stars, strikes, _ := sv.store.GetTeacherStarStrikes(teacher)

	ci := &ClassInfo{
		ClassList: classes,
		Stars:     stars,
		Strikes:   strikes,
		Teacher:   teacher,
		Title:     stu415.StudentName + " Magic Card",
		Path:      "classes",
	}
	stu, _ := sv.MakeStudent(stu415)

	data := struct {
		Info *ClassInfo
		S *Student
	}{
		Info: ci,
		S: stu,
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, data)
}
