package clog

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func ReadLog(r *http.Request) {

	c := appengine.NewContext(r)
	query := &log.Query{
		AppLogs:  true,
		Versions: []string{"1"},
	}

	for results := query.Run(c); ; {
		record, err := results.Next()
		if err == log.Done {
			log.Infof(c, "Done processing results")
			break
		}
		if err != nil {
			log.Errorf(c, "Failed to retrieve next log: %v", err)
			break
		}
		log.Infof(c, "Saw record %v", record)
	}

}

// type LogType int

// // log type
// const (
// 	INFO LogType = iota
// 	WARN
// 	ERROR
// 	DEBUG
// )

// type BufLog struct {
// 	tabe  string
// 	buf   *bytes.Buffer
// 	l     *log.Logger
// 	store *db.Store
// }

// // NewBufLog
// func NewBufLog(table string, store *db.Store) (*BufLog, error) {
// 	b := bytes.NewBuffer([]byte{})
// 	l := &BufLog{
// 		tabe:  table,
// 		buf:   b,
// 		l:     log.New(b, "", log.LstdFlags),
// 		store: store,
// 	}
// 	return l, nil
// }

// func (bl  *BufLog) Print(t LogType, v ...interface{}) error  {
// 	bl.l.Print(logTypeToString(t), v)

// 	// store log
// 	// bl.store.Log(bl.tabe, logTypeToString(t), v)
// 	return nil
// }

// func (bl  *BufLog) bufToDB()error  {

// }

// func logTypeToString(t LogType) string {
// 	switch t {
// 	case INFO:
// 		return "INFO"
// 	case WARN:
// 		return "WARN"
// 	case ERROR:
// 		return "ERROR"
// 	case DEBUG:
// 		return "DEBUG"
// 	default:
// 		return "UNKNOWN"
// 	}
// }
