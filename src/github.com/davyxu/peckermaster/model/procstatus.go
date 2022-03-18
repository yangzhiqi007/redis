package model

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	parseSupervisorStatus = regexp.MustCompile(`([a-z0-9_\:]+)(\s+)([A-Z]+)(\s)*([\S\d\s]+)`)
)

type ProcStatus struct {
	Name   string
	Status string
	Desc   string
}

func ParseStatusText(text string) (ret []ProcStatus) {
	reader := bufio.NewScanner(strings.NewReader(text))

	reader.Split(bufio.ScanLines)
	for reader.Scan() {

		line := strings.TrimSpace(reader.Text())

		result := parseSupervisorStatus.FindStringSubmatch(line)

		var s ProcStatus
		if len(result) == 6 {
			s.Name = result[1]
			s.Status = result[3]
			s.Desc = result[5]
			ret = append(ret, s)
		} else {
			log.Errorf("regmatch error: %s", line)
		}

	}

	return
}
