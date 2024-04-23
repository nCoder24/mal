package reader

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should tokenize an expression",
			args: args{str: "(+ 1 2)"},
			want: []string{"(", "+", "1", "2", ")"},
		},
		{
			name: "should ignore whitespaces",
			args: args{str: " (+ 	1 2 )"},
			want: []string{"(", "+", "1", "2", ")"},
		},
		{
			name: "should tokenize nested expression",
			args: args{str: "(+ 1 (* 3 4))"},
			want: []string{"(", "+", "1", "(", "*", "3", "4", ")", ")"},
		},
		{
			name: "should tokenize quote expression",
			args: args{str: "'(+ 1 2)"},
			want: []string{"'", "(", "+", "1", "2", ")"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
