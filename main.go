package main

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os/exec"
	"regexp"
	"strings"
)

const MessagePattern = `Process (?P<process>\w+)\.(?P<pid>\d+) (?P<action>\w+) (?P<assertion>\w+) "(?P<description>[^"]+)"`

const Assertion = "NoDisplaySleepAssertion"

type Match struct {
	Process     string
	PID         string
	Action      string
	Assertion   string
	Description string
}

type State struct {
	Depth   int
	TurnOn  func() error
	TurnOff func() error
}

func (s *State) Create() (err error) {
	s.Depth++

	if s.Depth == 1 {
		return s.TurnOn()
	}

	return
}

func (s *State) Release() (err error) {
	if s.Depth == 0 {
		return
	}

	s.Depth--

	if s.Depth == 0 {
		return s.TurnOff()
	}

	return
}

type Event struct {
	Message string `json:"eventMessage"`
}

func turnOnSign() (err error) {
	slog.Info("Sign turned ON")
	return
}

func turnOffSign() (err error) {
	slog.Info("Sign turned OFF")
	return
}

func main() {
	cmd := exec.Command(
		"log", "stream",
		"--process", "powerd",
		"--style", "ndjson",
		"--level", "default",
	)

	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()
	scanner := bufio.NewScanner(stdout)

	cmd.Start()

	// ignore first line of output
	if scanner.Scan() {
		_ = scanner.Text()
	}

	state := &State{
		TurnOn:  turnOnSign,
		TurnOff: turnOffSign,
	}

	re := regexp.MustCompile(MessagePattern)

	var event Event
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.NewDecoder(strings.NewReader(line)).Decode(&event); err != nil {
			slog.Error("failed to decode log line", "line", line, "error", err)
			continue
		}

		matches := re.FindStringSubmatch(event.Message)
		if matches == nil {
			slog.Debug("failed to parse log line", "line", line)
			continue
		}

		match := Match{
			Process:     matches[re.SubexpIndex("process")],
			PID:         matches[re.SubexpIndex("pid")],
			Action:      matches[re.SubexpIndex("action")],
			Assertion:   matches[re.SubexpIndex("assertion")],
			Description: matches[re.SubexpIndex("description")],
		}

		slog.Debug("message received", "process", match.Process, "action", match.Action, "assertion", match.Assertion, "description", match.Description)

		if match.Assertion != Assertion {
			continue
		}

		if match.Description == "video-playing" {
			slog.Info("ignoring video-playing assertion")
			continue
		}

		switch match.Action {
		case "Created":
			if err := state.Create(); err != nil {
				slog.Error("failed to set state", "action", "created", "error", err)
			}
		case "Released":
			slog.Debug("NoDisplaySleepAssertion", "action", "released")
			if err := state.Release(); err != nil {
				slog.Error("failed to set state", "action", "released", "error", err)
			}
		}
	}
}
