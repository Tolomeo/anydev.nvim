package scripts

import (
	"embed"
	"fmt"
)

//go:embed *.lua
var scripts embed.FS

func Read(name string) (string, error) {
	script, err := scripts.ReadFile(fmt.Sprintf("%s.lua", name))

	if err != nil {
		return "", err
	}

	return string(script), nil

}
