package roster

import (
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

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
		S415:      stu,
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

	c := &Classroom{
		ClassList: classes,
		MyStars:   stars,
		MyStrikes: strikes,
		Teacher:   teacher,
	}
	stu, _ := sv.MakeStudent(stu415)

	cd := CardData{
		Student: stu,
		Class:   c,
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, cd)
}
