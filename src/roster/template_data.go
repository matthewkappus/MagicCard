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
	// todo: remove stu415
	Stu415s []*synergy.Stu415
	// StarStrikes map students to StarStrike list for button data
	MyStars   map[*synergy.Stu415][]*comment.StarStrike
	MyStrikes map[*synergy.Stu415][]*comment.StarStrike
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

	mystars := make(map[*synergy.Stu415][]*comment.StarStrike)
	mystrikes := make(map[*synergy.Stu415][]*comment.StarStrike)

	for _, ss := range starstrikes {
		if ss.Cat == comment.Star {
			for _, s415 := range s415s {
				list, exists := mystars[s415]
				if exists {
					list = append(list, ss)
					mystars[s415] = list
				} else {
					mystars[s415] = []*comment.StarStrike{ss}
				}
			}
			// append strike to student
		} else {
			for _, s415 := range s415s {
				list, exists := mystrikes[s415]
				if exists {
					list = append(list, ss)
					mystrikes[s415] = list
				} else {
					mystrikes[s415] = []*comment.StarStrike{ss}
				}
			}
		}
	}

	return &Classroom{
		Stu415s:   s415s,
		MyStars:   mystars,
		MyStrikes: mystrikes,
		Teacher:   teacher,
	}, nil

}

// MakeSchoolClassroom returns list of every unique stu415
func (sv *StaffView) MakeSchoolClassroom(teacher string) (*Classroom, error) {
	s415s, err := sv.store.ListAllStudents()
	if err != nil {
		return nil, err
	}
	starstrikes, err := sv.store.GetMyStarStrikes(teacher)
	if err != nil {
		return nil, err
	}

	mystars := make(map[*synergy.Stu415][]*comment.StarStrike)
	mystrikes := make(map[*synergy.Stu415][]*comment.StarStrike)

	for _, ss := range starstrikes {
		if ss.Cat == comment.Star {
			for _, s415 := range s415s {
				list, exists := mystars[s415]
				if exists {
					list = append(list, ss)
					mystars[s415] = list
				} else {
					mystars[s415] = []*comment.StarStrike{ss}
				}
			}
			// append strike to student
		} else {
			for _, s415 := range s415s {
				list, exists := mystrikes[s415]
				if exists {
					list = append(list, ss)
					mystrikes[s415] = list
				} else {
					mystrikes[s415] = []*comment.StarStrike{ss}
				}
			}
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
