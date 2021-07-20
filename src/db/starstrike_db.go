package db

import "github.com/matthewkappus/MagicCard/src/comment"

// starstrike table
const (
	createStarStrike                 = `CREATE TABLE IF NOT EXISTS starstrike (id INTEGER PRIMARY KEY, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	insertStarStrike                 = `INSERT INTO starstrike(perm_id, teacher, comment, title, cat, isActive) VALUES(?,?,?,?,?,?);`
	selectStarStrikeByPermID         = `SELECT * FROM starstrike WHERE perm_id = ?`
	selectNewestStarStrikesByTeacher = `SELECT * FROM starstrike WHERE teacher=? LIMIT ?`

	// select count(cat) from starstrike where cat=1 and perm_id="980016917"
)

// mystarstrike table holds teacher's saved starstrikes
const (
	createMystarStrike  = `CREATE TABLE IF NOT EXISTS mystarstrike (id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true, FOREIGN KEY(teacher) REFERENCES stu415(teacher))`
	selectMyStarStrikes = `SELECT * FROM mystarstrike WHERE teacher = ?`
)

func (s *Store) GetMyStarStrikes(teacher string) ([]*comment.StarStrike, error) {
	rows, err := s.db.Query(selectMyStarStrikes, teacher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := make([]*comment.StarStrike, 0)
	for rows.Next() {
		// appends a new struct to ss, then scan its values into the struct
		ss = append(ss, new(comment.StarStrike))
		if err := rows.Scan(&ss[len(ss)-1].ID, &ss[len(ss)-1].Teacher, &ss[len(ss)-1].Title, &ss[len(ss)-1].Comment, &ss[len(ss)-1].Cat); err != nil {
			continue
		}
	}
	return ss, nil
}

func (s *Store) GetStarStrikesByPerm(id string) ([]*comment.StarStrike, error) {
	rows, err := s.db.Query(selectStarStrikeByPermID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := make([]*comment.StarStrike, 0)
	for rows.Next() {
		sb := new(comment.StarStrike)
		if err = rows.Scan(&sb.ID, &sb.PermID, &sb.Teacher, &sb.Comment, &sb.Title, &sb.Created, &sb.Cat, &sb.IsActive); err != nil {
			continue
		}
		ss = append(ss, sb)
	}

	return ss, nil

}

// depricated: todo: remove
func (s *Store) GetTeacherStarStrikes(teacher string) (stars, strikes []*comment.StarStrike, err error) {
	rows, err := s.db.Query(selectMyStarStrikes, teacher)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	stars = make([]*comment.StarStrike, 0)
	strikes = make([]*comment.StarStrike, 0)
	for rows.Next() {
		ss := new(comment.StarStrike)
		if err := rows.Scan(&ss.ID, &ss.Teacher, &ss.Title, &ss.Comment, &ss.Cat); err != nil {
			continue
		}
		if ss.Cat == comment.Star {
			stars = append(stars, ss)

		} else {

			strikes = append(strikes, ss)
		}
	}
	return stars, strikes, nil
}
