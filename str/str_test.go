package str

import (
	"reflect"
	"testing"
)

func TestTrim(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"no whitespace", "hello", "hello"},
		{"leading spaces", "  hello", "hello"},
		{"trailing spaces", "hello  ", "hello"},
		{"both sides spaces", "  hello  ", "hello"},
		{"tabs", "\t\thello\t\t", "hello"},
		{"newlines", "\n\nhello\n\n", "hello"},
		{"mixed whitespace", "\r\n\t hello world \t\n\r", "hello world"},
		{"only whitespace", "\r\n\t  \t\n\r", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.input); got != tt.want {
				t.Errorf("Trim() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrimLeft(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"no whitespace", "hello", "hello"},
		{"leading spaces", "  hello", "hello"},
		{"trailing spaces", "hello  ", "hello  "},
		{"both sides spaces", "  hello  ", "hello  "},
		{"tabs", "\t\thello\t\t", "hello\t\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimLeft(tt.input); got != tt.want {
				t.Errorf("TrimLeft() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrimRight(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"no whitespace", "hello", "hello"},
		{"leading spaces", "  hello", "  hello"},
		{"trailing spaces", "hello  ", "hello"},
		{"both sides spaces", "  hello  ", "  hello"},
		{"tabs", "\t\thello\t\t", "\t\thello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimRight(tt.input); got != tt.want {
				t.Errorf("TrimRight() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSplitBy(t *testing.T) {
	tests := []struct {
		name  string
		input string
		sep   string
		want  []string
	}{
		{"empty string", "", ",", []string{}},
		{"single value", "hello", ",", []string{"hello"}},
		{"multiple values", "a,b,c", ",", []string{"a", "b", "c"}},
		{"with spaces", " a , b , c ", ",", []string{"a", "b", "c"}},
		{"empty parts", "a,,b", ",", []string{"a", "b"}},
		{"trailing separator", "a,b,", ",", []string{"a", "b"}},
		{"leading separator", ",a,b", ",", []string{"a", "b"}},
		{"different separator", "a:b:c", ":", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitBy(tt.input, tt.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComma(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty string", "", []string{}},
		{"single value", "hello", []string{"hello"}},
		{"multiple values", "a,b,c", []string{"a", "b", "c"}},
		{"with spaces", " a , b , c ", []string{"a", "b", "c"}},
		{"empty parts", "a,,b", []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Comma(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Comma() = %v, want %v", got, tt.want)
			}
		})
	}
}
