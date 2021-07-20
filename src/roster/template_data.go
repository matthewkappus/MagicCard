package roster

import (
	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

type TD struct {
	S *Student
	C *Classroom
	N *Nav
}
type Nav struct {
	ClassList []*synergy.Stu415
	Path      string
	Title     string
	Status    string
}

// Classroom is teachers class info
type Classroom struct {
	Stu415s   []*synergy.Stu415
	MyStars   []*comment.StarStrike
	MyStrikes []*comment.StarStrike
	ClassList []*synergy.Stu415
	Teacher   string
	ClassName string
}

// Student holds the Stu415 and Starbars for a template
//  My and Strikes are mapped to their title for counting
type Student struct {
	S415 *synergy.Stu415
	// Star.Title to SS
	StarMap map[string][]*comment.StarStrike
	// StrikeMap
	StrikeMap map[string][]*comment.StarStrike
}

func (sv *StaffView) MakeStudent(perm string) (*Student, error) {
	stu, err := sv.store.SelectStu415(perm)
	if err != nil {
		return nil, err
	}

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

func (sv *StaffView) MakeClassroom(teacher, section string) (*Classroom, error) {
	s415s, err := sv.store.ListStudents(section)
	if err != nil {
		return nil, err
	}
	starstrikes, err := sv.store.GetMyStarStrikes(teacher)
	if err != nil {
		return nil, err
	}

	mystars := make([]*comment.StarStrike, 0)
	mystrikes := make([]*comment.StarStrike, 0)

	for _, ss := range starstrikes {
		if ss.Cat == comment.Star {
			mystars = append(mystars, ss)
		} else {
			mystrikes = append(mystrikes, ss)
		}
	}
	return &Classroom{
		Stu415s:   s415s,
		MyStars:   mystars,
		MyStrikes: mystrikes,
		Teacher:   teacher,
	}, nil

}

//
func (sv *StaffView) MakeNav(teacher, path, title string) (*Nav, error) {
	classlist, err := sv.store.ListClasses(teacher)
	if err != nil {
		return nil, err
	}

	n := &Nav{
		ClassList: classlist,
		Path:      path,
		Title:     title,
	}
	return n, nil
}
