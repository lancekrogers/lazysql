package cmdline

import "testing"

func TestParseInvocation(t *testing.T) {
	cases := []struct {
		input string
		name  string
		force bool
		args  []string
	}{
		{input: ":w file.sql", name: "w", force: false, args: []string{"file.sql"}},
		{input: "w!", name: "w", force: true, args: nil},
		{input: "wq! file.sql", name: "wq", force: true, args: []string{"file.sql"}},
		{input: "w ! file.sql", name: "w", force: true, args: []string{"file.sql"}},
		{input: "e \"file name.sql\"", name: "e", force: false, args: []string{"file name.sql"}},
	}

	for _, tc := range cases {
		inv, err := ParseInvocation(tc.input)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.input, err)
		}
		if inv.Name != tc.name || inv.Force != tc.force {
			t.Fatalf("parse %q: expected %s force=%v, got %s force=%v", tc.input, tc.name, tc.force, inv.Name, inv.Force)
		}
		if len(inv.Args) != len(tc.args) {
			t.Fatalf("parse %q: expected %d args, got %d", tc.input, len(tc.args), len(inv.Args))
		}
		for i, arg := range tc.args {
			if inv.Args[i] != arg {
				t.Fatalf("parse %q: expected arg %q, got %q", tc.input, arg, inv.Args[i])
			}
		}
	}
}

func TestParseInvocationEmpty(t *testing.T) {
	inv, err := ParseInvocation("   ")
	if err != nil {
		t.Fatalf("parse empty: %v", err)
	}
	if inv.Name != "" || inv.Raw != "" || inv.Force || len(inv.Args) != 0 {
		t.Fatalf("expected empty invocation, got %+v", inv)
	}
}

func TestParseInvocationUnterminatedQuote(t *testing.T) {
	_, err := ParseInvocation(`e "unterminated`)
	if err == nil {
		t.Fatal("expected unterminated quote error")
	}
	if err != ErrUnterminatedQuote {
		t.Fatalf("unexpected error: %v", err)
	}
}
