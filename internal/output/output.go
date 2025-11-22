package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type output struct {
	directory string
}

func (o *output) WriteFile(filename string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		return fmt.Errorf("Error marshalling data: %w", err)
	}

	target := filepath.Join(o.directory, filename)

	err = os.WriteFile(target, jsonData, 0644)

	if err != nil {
		return fmt.Errorf("Error writing to file: %w", err)
	}

	return nil
}

func NewOutput(directory string) *output {
	out := output{
		directory: directory,
	}

	return &out
}
