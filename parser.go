package main

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

// Matches: users:(("nginx",pid=1234,fd=6)) or users:(("go",1234,6))
var procRegex = regexp.MustCompile(`"([^"]+)",(?:pid=)?(\d+)`)

func scanLines(stdout []byte, skipHeader bool, parser func(line string) (*Process, error)) ([]Process, error) {
	var processes []Process
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	isHeader := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if skipHeader && isHeader {
			isHeader = false
			continue
		}

		proc, err := parser(line)
		if err != nil {
			continue
		}
		if proc != nil {
			processes = append(processes, *proc)
		}
	}
	return processes, scanner.Err()
}

func ParseSSCmd(stdout []byte) ([]Process, error) {
	return scanLines(stdout, true, func(line string) (*Process, error) {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return nil, nil
		}

		p := &Process{Protocol: fields[0]}
		parseAddrPort(p, fields[4])

		for i := 5; i < len(fields); i++ {
			if strings.HasPrefix(fields[i], "users:") {
				fullUserField := strings.Join(fields[i:], " ")
				matches := procRegex.FindStringSubmatch(fullUserField)
				if len(matches) == 3 {
					p.Name = matches[1]
					p.PID, _ = strconv.Atoi(matches[2])
					break
				}
			}
		}

		if p.PID == 0 {
			return nil, nil
		}
		return p, nil
	})
}

func ParseLsofCmd(stdout []byte) ([]Process, error) {
	return scanLines(stdout, true, func(line string) (*Process, error) {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			return nil, nil
		}

		p := &Process{
			Name:     fields[0],
			PID:      toInt(fields[1]),
			Protocol: fields[7],
		}
		parseAddrPort(p, fields[8])
		return p, nil
	})
}

func ParseNetstatCmd(stdout []byte) ([]Process, error) {
	return scanLines(stdout, false, func(line string) (*Process, error) {
		fields := strings.Fields(line)

		if len(fields) < 4 || !strings.Contains(strings.ToUpper(line), "LISTENING") {
			return nil, nil
		}

		p := &Process{
			Protocol: fields[0],
		}

		for _, f := range fields {
			if strings.Contains(f, ":") || strings.Contains(f, ".") {
				parseAddrPort(p, f)
				break
			}
		}

		lastField := fields[len(fields)-1]
		if strings.Contains(lastField, "/") {
			parts := strings.Split(lastField, "/")
			p.PID = toInt(parts[0])
			p.Name = parts[1]
		} else {
			p.PID = toInt(lastField)
		}

		if p.PID == 0 {
			return nil, nil
		}

		return p, nil
	})
}

func parseAddrPort(p *Process, raw string) {
	lastColon := strings.LastIndex(raw, ":")
	if lastColon != -1 {
		p.Address = strings.Trim(raw[:lastColon], "[]")
		p.Port, _ = strconv.Atoi(raw[lastColon+1:])
	}
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
