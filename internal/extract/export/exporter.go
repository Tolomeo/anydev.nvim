package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Exporter struct {
	directory string
}

func (e *Exporter) Exported(subdirectory string, name string) bool {
	file := filepath.Join(e.directory, subdirectory, fmt.Sprintf("%s.json", name))

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (e *Exporter) ExportJson(subdirectory string, name string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		return fmt.Errorf("Error marshalling data: %w", err)
	}

	targetDirectory := filepath.Join(e.directory, subdirectory)
	err = os.MkdirAll(targetDirectory, 0777)

	if err != nil {
		return fmt.Errorf("Error creating directory '%s': %w", targetDirectory, err)
	}

	targetFile := filepath.Join(targetDirectory, fmt.Sprintf("%s.json", name))
	err = os.WriteFile(targetFile, jsonData, 0666)

	if err != nil {
		return fmt.Errorf("Error writing file '%s': %w", targetFile, err)
	}

	return nil
}

func NewExporter(directory string) *Exporter {
	exporter := Exporter{
		directory: directory,
	}

	return &exporter
}
