package figure

import "testing"

func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"ID", "id"},
		{"i", "i"},
		{"SnakeCase", "snake_case"},
	}

	for _, tc := range cases {
		got := toSnakeCase(tc.in)
		if got != tc.out {
			t.Errorf("expected %s got %s", tc.out, got)
		}
	}
}
