package roster

import (
	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

type TD struct {
	M *MagicCard
	C *Classroom
	N *Nav
}

type Nav struct {
	ClassList []*synergy.Stu415
	Path      string
	Title     string
	Status    string
	// 0 Student 1 Teacher 2 Admin
	Type Scope
}

// Classroom is teachers class info
type Classroom struct {
	// todo: remove stu415
	Stu415s []*synergy.Stu415
	// StarStrikes map students to StarStrike list for button data
	StarStrikes map[*synergy.Stu415][]*comment.StarStrike
	Teacher     string
	ClassName   string
}

// MagicCard holds the Stu415 and Starbars for a template
//  My and Strikes are mapped to their title for counting
type MagicCard struct {
	Name string
	ID   string

	S415 *synergy.Stu415
	// Star.Title to SS
	StarMap map[string][]*comment.StarStrike
	// StrikeMap
	StrikeMap map[string][]*comment.StarStrike
}

func (sv *View) MakeStudent(perm string) (*MagicCard, error) {
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
	return &MagicCard{
		S415:      stu,
		StarMap:   starsM,
		StrikeMap: strikesM}, nil
}

func (sv *View) MakeClassroom(teacher, section string) (*Classroom, error) {
	s415s, err := sv.store.ListStudents(section)
	if err != nil {
		return nil, err
	}
	myss, err := sv.store.GetMyStarStrikes(teacher)

	if err != nil {
		return nil, err
	}

	// starstrikes takes generic starstrikes and puts each student perm in for use with buttons
	ss := make(map[*synergy.Stu415][]*comment.StarStrike)
	for _, stu := range s415s {
		s := make([]*comment.StarStrike, 0)
		for _, mss := range myss {
			mss.PermID = stu.PermID
			mss.Teacher = teacher
			s = append(s, mss)
		}
		ss[stu] = s
	}
	return &Classroom{
		Stu415s:     s415s,
		StarStrikes: ss,
		Teacher:     teacher,
	}, nil

}

// MakeSchoolClassroom returns list of every unique stu415
func (sv *View) MakeSchoolClassroom(teacher string) (*Classroom, error) {
	s415s, err := sv.store.ListAllStudents()
	if err != nil {
		return nil, err
	}
	myss, err := sv.store.GetMyStarStrikes(teacher)
	if err != nil {
		return nil, err
	}

	// starstrikes takes generic starstrikes and puts each student perm in for use with buttons
	ss := make(map[*synergy.Stu415][]*comment.StarStrike)
	for _, stu := range s415s {
		s := make([]*comment.StarStrike, 0)
		for _, mss := range myss {
			mss.PermID = stu.PermID
			mss.Teacher = teacher
			s = append(s, mss)
		}
		ss[stu] = s
	}
	return &Classroom{
		Stu415s:     s415s,
		StarStrikes: ss,
		Teacher:     teacher,
	}, nil

}

func (sv *View) MakeTeacherMagicCard(name string) (*MagicCard, error) {

	stars, strikes, err := sv.store.SelectTeacherStarStrikes(name, 100)
	if err != nil {
		return nil, err
	}

	starsM := make(map[string][]*comment.StarStrike)
	for _, ss := range stars {
		match, exists := starsM[ss.Title]
		if exists {
			starsM[ss.Title] = append(match, ss)
			continue
		}
		starsM[ss.Title] = []*comment.StarStrike{ss}
	}

	stikesM := make(map[string][]*comment.StarStrike)
	for _, ss := range strikes {
		match, exists := stikesM[ss.Title]
		if exists {
			stikesM[ss.Title] = append(match, ss)
			continue
		}
		stikesM[ss.Title] = []*comment.StarStrike{ss}
	}

	return &MagicCard{

		Name:      name,
		ID:        name,
		StarMap:   starsM,
		StrikeMap: stikesM}, nil
}

func (sv *View) MakeStudentMagicCard(perm string) (*MagicCard, error) {

	stu, err := sv.store.SelectStu415(perm)
	if err != nil {
		return nil, err
	}
	ss, err := sv.store.GetStarStrikesByPerm(perm)
	if err != nil {
		return nil, err
	}

	stars := make(map[string][]*comment.StarStrike)
	strikes := make(map[string][]*comment.StarStrike)

	for _, s := range ss {
		// 0 star 1 minor 2 strik 3 major
		if s.Cat == comment.Star {
			tStars, exists := stars[s.Title]
			if !exists {
				stars[s.Title] = []*comment.StarStrike{s}
				continue
			}
			// title exists: add to list
			stars[s.Title] = append(tStars, s)
		} else {
			// ss.Cat is 1+: A strike
			tStrikes, exists := strikes[s.Title]
			if !exists {
				strikes[s.Title] = []*comment.StarStrike{s}
				continue
			}
			// strike exists: add to rest
			strikes[s.Title] = append(tStrikes, s)
		}
	}

	return &MagicCard{

		Name:      stu.StudentName,
		ID:        stu.PermID,
		S415:      stu,
		StarMap:   stars,
		StrikeMap: strikes}, nil

}

// MakeNav returns data struct for use in <head> / <nav>
func (sv *View) MakeNav(teacher, path, title string, stype Scope) (*Nav, error) {
	classlist, err := sv.store.ListClasses(teacher)
	if err != nil {
		return nil, err
	}
	n := &Nav{
		ClassList: classlist,
		Path:      path,
		Title:     title,
		Type:      stype,
	}
	return n, nil
}
