package tools

import (
	"testing"
)

func Test_stringInSlice(t *testing.T) {
	type args struct {
		str  string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"stringInSlice Basic 1", args{"0", []string{"1", "2"}}, false},
		{"stringInSlice Basic 2", args{"1", []string{"1", "2"}}, true},
		{"stringInSlice Basic 3", args{"2", []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.args.str, tt.args.list); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_regexpMatch(t *testing.T) {
	type args struct {
		regex string
		word  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"regexpMatch false match", args{"(?i)(Story)[ ][a-zA-Z0-9]", "apple asdf123"}, false},
		{"regexpMatch true match 1", args{"(?i)(Story)[ ][a-zA-Z0-9]", "Story asdf123"}, true},
		{"regexpMatch true match 2", args{"(?i)(Story)[ ][a-zA-Z0-9]", "story asdf123"}, true},
		{"regexpMatch false match 2", args{"(?i)(Story)[ ][a-zA-Z0-9]", ""}, false},
		{"regexpMatch true match 2", args{"(?i)(Story)[ ][a-zA-Z0-9]", "stOrY asdf123"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := regexpMatch(tt.args.regex, tt.args.word); got != tt.want {
				t.Errorf("regexpMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
