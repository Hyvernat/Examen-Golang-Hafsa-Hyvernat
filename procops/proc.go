package procops

import (
	"bytes"
	"encoding/csv"
	"os/exec"
	"runtime"
	"strings"
)

// ProcessInfo représente les informations de base d'un processus
type ProcessInfo struct {
	PID  string
	Name string
}

// ListProcesses liste les N premiers processus selon l'OS
func ListProcesses(topN int) ([]ProcessInfo, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/FO", "CSV", "/NH")
	} else {
		cmd = exec.Command("ps", "-Ao", "pid,comm", "--no-headers")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	if runtime.GOOS == "windows" {
		reader := csv.NewReader(bytes.NewReader(output))
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		for i, record := range records {
			if i >= topN && topN > 0 {
				break
			}
			processes = append(processes, ProcessInfo{
				PID:  record[1],
				Name: record[0],
			})
		}
	} else {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}
			if i >= topN && topN > 0 {
				break
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				processes = append(processes, ProcessInfo{
					PID:  fields[0],
					Name: fields[1],
				})
			}
		}
	}

	return processes, nil
}

// FilterProcesses recherche des processus par nom
func FilterProcesses(keyword string) ([]ProcessInfo, error) {
	all, err := ListProcesses(0)
	if err != nil {
		return nil, err
	}

	var filtered []ProcessInfo
	keyword = strings.ToLower(keyword)
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Name), keyword) {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

// KillProcess termine un processus par son PID
func KillProcess(pid string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/PID", pid, "/T", "/F")
	} else {
		cmd = exec.Command("kill", "-9", pid)
	}

	return cmd.Run()
}
