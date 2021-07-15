package db

import (
	"github.com/matthewkappus/MagicCard/src/comment"
)

// starbar_db is a database of starbars.
// starbar_db_version = "0.1"
const (
	createStarBar          = `CREATE TABLE IF NOT EXISTS starbar(id INTEGER PRIMARY KEY AUTOINCREMENT, teacher, title, comment TEXT, isStar BOOLEAN, FOREIGN KEY(teacher) REFERENCES staff(name))`
	createStarBarSave      = `INSERT OR REPLACE INTO starbar(id, teacher, title, comment, isStar) VALUES(?, ?, ?, ?, ?)`
	selectStarBarByTeacher = `SELECT * FROM starbar WHERE teacher = ?`
	selectStarBarByID      = `SELECT * FROM starbar WHERE id = ?`
	insertStarBar          = `INSERT INTO starbar(teacher, title, comment, isStar) VALUES(?, ?, ?, ?)`
	updateStarBar          = `UPDATE starbar SET teacher = ?, title = ?, comment = ?, isStar = ? WHERE id = ?`
)

// const ErrInvalidStarBar = "Invalid StarBar"

func (s *Store) GetStarBarByID(id string) (*comment.StarBar, error) {
	sb := &comment.StarBar{}
	err := s.db.QueryRow(selectStarBarByID, id).Scan(&sb.ID, &sb.Teacher, &sb.Title, &sb.Comment, &sb.IsStar)
	return sb, err
}

func (s *Store) UpdateStarBar(id int, teacher, title, comments string, isStar bool) error {
	_, err := s.db.Exec(updateStarBar, teacher, title, comments, isStar, id)
	return err
}

// GetStarBars takes a staff.name and returns their StarBars or an error
func (s *Store) GetStarBars(teacher string) (stars, bars []*comment.StarBar, err error) {
	rows, err := s.db.Query(selectStarBarByTeacher, teacher)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	stars = make([]*comment.StarBar, 0)
	bars = make([]*comment.StarBar, 0)
	for rows.Next() {
		sb := new(comment.StarBar)
		if err := rows.Scan(&sb.ID, &sb.Teacher, &sb.Title, &sb.Comment, &sb.IsStar); err != nil {
			continue
		}
		if sb.IsStar {
			stars = append(stars, sb)

		} else {

			bars = append(bars, sb)
		}
	}
	return stars, bars, nil
}

func (s *Store) AddStarBar(teacher, title, comments string, isStar bool) (*comment.StarBar, error) {
	sb := &comment.StarBar{
		Teacher: teacher,
		Title:   title,
		Comment: comments,
		IsStar:  isStar,
	}

	// sb.IsValid()
	res, err := s.db.Exec(insertStarBar, teacher, title, comments, isStar)
	if err != nil {
		return nil, err
	}
	i64, err := res.LastInsertId()
	sb.ID = int(i64)
	if err != nil {
		return nil, err
	}
	return sb, nil

}
