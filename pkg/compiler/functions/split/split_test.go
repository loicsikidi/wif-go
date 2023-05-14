package split

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		sep     string
		want    []string
		wantErr bool
	}{
		{
			name:    "splitting a string with a single separator",
			str:     "hello,world",
			sep:     ",",
			want:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "splitting a string with multiple separators",
			str:     "hello,,world",
			sep:     ",",
			want:    []string{"hello", "", "world"},
			wantErr: false,
		},
		{
			name:    "splitting a string with no separators",
			str:     "helloworld",
			sep:     ",",
			want:    []string{"helloworld"},
			wantErr: false,
		},
		{
			name:    "splitting an empty string",
			str:     "",
			sep:     ",",
			want:    []string{""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := split(tt.str, tt.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitN(t *testing.T) {
	var emptyList []string
	tests := []struct {
		name    string
		str     string
		sep     string
		n       int64
		want    []string
		wantErr bool
	}{
		{
			name:    "splitting a string with a single separator",
			str:     "hello,world",
			sep:     ",",
			n:       -1,
			want:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "splitting a string with multiple separators",
			str:     "hello,,world",
			sep:     ",",
			n:       -1,
			want:    []string{"hello", "", "world"},
			wantErr: false,
		},
		{
			name:    "splitting a string with no separators",
			str:     "helloworld",
			sep:     ",",
			n:       -1,
			want:    []string{"helloworld"},
			wantErr: false,
		},
		{
			name:    "splitting an empty string",
			str:     "",
			sep:     ",",
			n:       -1,
			want:    []string{""},
			wantErr: false,
		},
		{
			name:    "splitting a string with a limit of 0",
			str:     "hello hello hello",
			sep:     " ",
			n:       0,
			want:    emptyList,
			wantErr: false,
		},
		{
			name:    "splitting a string with a limit of 1",
			str:     "hello hello hello",
			sep:     " ",
			n:       1,
			want:    []string{"hello hello hello"},
			wantErr: false,
		},
		{
			name:    "splitting a string with a limit of 2",
			str:     "hello hello hello",
			sep:     " ",
			n:       2,
			want:    []string{"hello", "hello hello"},
			wantErr: false,
		},
		{
			name:    "splitting a string with a limit of 3",
			str:     "hello hello hello",
			sep:     " ",
			n:       3,
			want:    []string{"hello", "hello", "hello"},
			wantErr: false,
		},
		{
			name:    "splitting a string with a limit greater than the number of separators",
			str:     "hello hello hello",
			sep:     " ",
			n:       4,
			want:    []string{"hello", "hello", "hello"},
			wantErr: false,
		},
		{
			name:    "splitting a string with a limit less than the number of separators",
			str:     "hello hello hello",
			sep:     " ",
			n:       2,
			want:    []string{"hello", "hello hello"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitN(tt.str, tt.sep, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitN() = %v, want %v", got, tt.want)
			}
		})
	}
}
