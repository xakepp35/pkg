package utils

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// identity parser for strings
func stringParser(s string) (string, error) {
	return s, nil
}

func TestParseRepeatedStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
		wantErr  bool
	}{
		{"a,b,c", []string{"a", "b", "c"}, false},
		{"single", []string{"single"}, false},
		{"", []string{""}, false}, // empty yields one empty element
		{"x,,y", []string{"x", "", "y"}, false},
	}

	for _, tc := range tests {
		out, err := ParseRepeated[string](tc.input, stringParser)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseRepeated[string](%q) expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRepeated[string](%q) unexpected error: %v", tc.input, err)
			continue
		}
		if !reflect.DeepEqual(out, tc.expected) {
			t.Errorf("ParseRepeated[string](%q) = %v; want %v", tc.input, out, tc.expected)
		}
	}
}

func TestParseRepeatedInts(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"1,2,3", []int{1, 2, 3}, false},
		{"42", []int{42}, false},
		{"", nil, true}, // strconv.Atoi("") -> error
		{"7,8,notanint,9", nil, true},
	}

	for _, tc := range tests {
		out, err := ParseRepeated[int](tc.input, func(s string) (int, error) {
			return strconv.Atoi(strings.TrimSpace(s))
		})
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseRepeated[int](%q) expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRepeated[int](%q) unexpected error: %v", tc.input, err)
			continue
		}
		if !reflect.DeepEqual(out, tc.expected) {
			t.Errorf("ParseRepeated[int](%q) = %v; want %v", tc.input, out, tc.expected)
		}
	}
}

func TestParseRepeatedCustomType(t *testing.T) {
	type Foo struct{ V string }
	parser := func(s string) (Foo, error) {
		return Foo{V: strings.ToUpper(s)}, nil
	}

	out, err := ParseRepeated[Foo]("foo,bar", parser)
	if err != nil {
		t.Fatalf("ParseRepeated[Foo] unexpected error: %v", err)
	}
	want := []Foo{{"FOO"}, {"BAR"}}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("ParseRepeated[Foo] = %v; want %v", out, want)
	}
}
