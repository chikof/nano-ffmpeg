package cmd

import "testing"

func TestParseThemeOverride(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty means no override", input: "", want: ""},
		{name: "dark", input: "dark", want: "dark"},
		{name: "light", input: "light", want: "light"},
		{name: "trim and normalize", input: " LIGHT ", want: "light"},
		{name: "invalid", input: "blue", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseThemeOverride(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got nil", tc.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
