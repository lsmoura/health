package health

import "embed"

//go:embed schema.sql
var f embed.FS

func Schema() (string, error) {
	data, err := f.ReadFile("schema.sql")
	if err != nil {
		return "", err
	}

	return string(data), nil
}