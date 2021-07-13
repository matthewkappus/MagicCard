package db

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func TestStore_GetStarBars(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		teacher string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*comment.StarBar
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			got, err := s.GetStarBars(tt.args.teacher)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetStarBars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.GetStarBars() = %v, want %v", got, tt.want)
			}
		})
	}
}
