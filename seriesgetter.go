package main

import (
	"encoding/json"
	"errors"
	"strings"
)

type intermediateSeriesObject struct {
	Episodes []Episode
}

const JsonPrefixLine = `<script id="series-object" type="application/json">`

func ExtractEpisodesList(body string) ([]Episode, error) {
	lines := strings.Split(body, "\n")
	found := false
	var i int
	var line string
	for i, line = range lines {
		if strings.TrimSpace(line) == JsonPrefixLine {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("series-object not found in series page")
	}
	seriesObjectBody := strings.TrimSpace(lines[i+1])
	var s intermediateSeriesObject
	err := json.Unmarshal([]byte(seriesObjectBody), &s)
	if err != nil {
		return nil, err
	}
	return s.Episodes, nil
}
