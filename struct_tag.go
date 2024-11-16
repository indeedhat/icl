package icl

import (
	"errors"
	"strconv"
	"strings"
)

type tags struct {
	key       string
	env       string
	precision int
	isParam   bool
}

func parseTags(s string) (*tags, error) {
	if s == ".param" {
		return &tags{isParam: true}, nil
	}

	parts := strings.Split(s, ",")
	t := tags{
		key:       parts[0],
		precision: -1,
	}

	if strings.Contains(t.key, ".") {
		parts := strings.Split(t.key, ".")
		if len(parts) != 2 {
			return nil, errors.New("invalid icl key: " + t.key)
		}

		precision, err := strconv.Atoi(parts[1])
		if err != nil || precision < 1 {
			return nil, errors.New("invalid icl key: " + t.key)
		}

		t.key = parts[0]
		t.precision = precision
	}

	for _, part := range parts {
		if !strings.HasPrefix(part, "env(") && part[len(part)-1] != ')' {
			continue
		}
		t.env = part[4 : len(part)-1]
	}

	return &t, nil
}
