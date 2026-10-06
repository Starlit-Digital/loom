package loom

import "github.com/cshaiku/loom/internal/structured"

func readStructuredFile(path string) ([]byte, error) {
	data, err := structured.ReadFile(path, int(MaxInputBytes))
	if err != nil {
		return nil, err
	}
	return structured.JSON(data, int(MaxInputBytes))
}
