package roster

import "testing"

func Test_formatName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		
		{ name: "Normal Name", args: args{name: "Kappus, Matthew D."}, want: "Matthew Kappus"},
		{ name: "No MI", args: args{name: "Kappus, Matthew"}, want: "Matthew Kappus"},
		{ name: "Empty Name", args: args{name: ""}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatName(tt.args.name); got != tt.want {
				t.Errorf("formatName() = %v, want %v", got, tt.want)
			}
		})
	}
}
