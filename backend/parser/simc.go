package parser

import (
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	Summary   Summary `json:"summary"`
	Stats     Stats   `json:"stats"`
	RawOutput string  `json:"raw_output"`
}

type Summary struct {
	Player       string  `json:"player"`
	Class        string  `json:"class"`
	Spec         string  `json:"spec"`
	DPS          float64 `json:"dps"`
	Iterations   int     `json:"iterations"`
	FightLength  string  `json:"fight_length"`
}

type Stats struct {
	Strength     int     `json:"strength"`
	Agility      int     `json:"agility"`
	Intellect    int     `json:"intellect"`
	Crit         float64 `json:"crit"`
	Haste        float64 `json:"haste"`
	Mastery     float64 `json:"mastery"`
	Versatility  float64 `json:"versatility"`
}

func ParseSimC(output string) Result {
	res := Result{
		RawOutput: output,
	}

	lines := strings.Split(output, "\n")

	dpsRe := regexp.MustCompile(`DPS:\s+([\d\.]+)`)
	playerRe := regexp.MustCompile(`Player:\s+(\S+)`)
	specRe := regexp.MustCompile(`Spec:\s+(.+)`)
	iterRe := regexp.MustCompile(`Iterations:\s+(\d+)`)

	for _, line := range lines {
		if m := dpsRe.FindStringSubmatch(line); m != nil {
			res.Summary.DPS, _ = strconv.ParseFloat(m[1], 64)
		}
		if m := playerRe.FindStringSubmatch(line); m != nil {
			res.Summary.Player = m[1]
		}
		if m := specRe.FindStringSubmatch(line); m != nil {
			res.Summary.Spec = strings.TrimSpace(m[1])
		}
		if m := iterRe.FindStringSubmatch(line); m != nil {
			res.Summary.Iterations, _ = strconv.Atoi(m[1])
		}
	}

	return res
}
