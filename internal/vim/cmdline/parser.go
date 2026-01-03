package cmdline

import (
	"errors"
	"strings"
	"unicode"
)

var ErrUnterminatedQuote = errors.New("unterminated quote")

func ParseInvocation(input string) (Invocation, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return Invocation{}, nil
	}
	if strings.HasPrefix(raw, ":") {
		raw = strings.TrimSpace(raw[1:])
	}
	parts, err := splitArgs(raw)
	if err != nil {
		return Invocation{}, err
	}
	if len(parts) == 0 {
		return Invocation{}, nil
	}

	name := parts[0]
	force := false
	if strings.HasSuffix(name, "!") {
		force = true
		name = strings.TrimSuffix(name, "!")
	}

	args := parts[1:]
	if len(args) > 0 && args[0] == "!" {
		force = true
		args = args[1:]
	}

	return Invocation{
		Name:  name,
		Args:  args,
		Force: force,
		Raw:   raw,
	}, nil
}

func splitArgs(input string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuotes := false

	flush := func() {
		if current.Len() == 0 {
			return
		}
		args = append(args, current.String())
		current.Reset()
	}

	for _, r := range input {
		switch {
		case r == '"':
			inQuotes = !inQuotes
		case unicode.IsSpace(r) && !inQuotes:
			flush()
		default:
			current.WriteRune(r)
		}
	}

	if inQuotes {
		return nil, ErrUnterminatedQuote
	}
	flush()
	return args, nil
}
