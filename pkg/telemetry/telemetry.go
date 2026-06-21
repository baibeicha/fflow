package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Event describes a single command execution event for telemetry.
type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Command    string    `json:"command"`
	FilesCount int64     `json:"files_count"`
	BytesTotal int64     `json:"bytes_total"`
	DurationMs int64     `json:"duration_ms"`
	IsSuccess  bool      `json:"is_success"`
	ErrorMsg   string    `json:"error_message,omitempty"`
}

// State stores telemetry settings and push history.
type State struct {
	Endpoint    string    `json:"endpoint"`
	LastPush    time.Time `json:"last_push"`
	Enabled     bool      `json:"enabled"`
	HasPrompted bool      `json:"has_prompted"`
}

var mu sync.Mutex

// AutoPushInterval defines how often background push is triggered.
const AutoPushInterval = 1 * time.Hour

func getPaths() (metricsPath, statePath string, err error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(cacheDir, "fflow")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "metrics.jsonl"), filepath.Join(dir, "state.json"), nil
}

// LoadState retrieves the current telemetry state.
func LoadState() State {
	var s State
	_, statePath, err := getPaths()
	if err != nil {
		return s
	}

	data, err := os.ReadFile(statePath)
	if err == nil {
		_ = json.Unmarshal(data, &s)
	}
	return s
}

// SaveState persists the telemetry state.
func SaveState(s State) {
	_, statePath, err := getPaths()
	if err != nil {
		return
	}
	data, _ := json.Marshal(s)
	_ = os.WriteFile(statePath, data, 0644)
}

// Record saves a telemetry event to a file and triggers auto-push if needed.
func Record(cmdName string, files, bytes int64, start time.Time, execErr error) {
	state := LoadState()
	if !state.Enabled {
		return
	}

	ev := Event{
		Timestamp:  time.Now().UTC(),
		Command:    cmdName,
		FilesCount: files,
		BytesTotal: bytes,
		DurationMs: time.Since(start).Milliseconds(),
		IsSuccess:  execErr == nil,
	}
	if execErr != nil {
		ev.ErrorMsg = execErr.Error()
	}

	metricsPath, _, err := getPaths()
	if err != nil {
		return
	}

	data, _ := json.Marshal(ev)
	data = append(data, '\n')

	mu.Lock()
	f, err := os.OpenFile(metricsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		_, _ = f.Write(data)
		f.Close()
	}
	mu.Unlock()

	if state.Endpoint != "" && time.Since(state.LastPush) > AutoPushInterval {
		triggerBackgroundPush()
	}
}

func triggerBackgroundPush() {
	exe, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.Command(exe, "telemetry", "push", "--silent")
	_ = cmd.Start()
}

// Push sends collected metrics to the specified endpoint.
func Push(endpoint string) error {
	metricsPath, _, err := getPaths()
	if err != nil {
		return fmt.Errorf("error getting paths: %w", err)
	}

	if _, err := os.Stat(metricsPath); os.IsNotExist(err) {
		return fmt.Errorf("metrics file not exists: %s", metricsPath)
	}

	sendingPath := metricsPath + ".sending"

	_ = os.Remove(sendingPath)

	if err := os.Rename(metricsPath, sendingPath); err != nil {
		return fmt.Errorf("can not open blocked metrics file: %w", err)
	}

	data, err := os.ReadFile(sendingPath)
	if err != nil {
		_ = os.Rename(sendingPath, metricsPath)
		return fmt.Errorf("error reading metrics file: %w", err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		_ = os.Remove(sendingPath)
		return fmt.Errorf("metrics file is empty")
	}

	jsonArray := "[" + string(bytes.ReplaceAll(bytes.TrimSpace(data), []byte("\n"), []byte(","))) + "]"

	resp, err := http.Post(endpoint, "application/json", bytes.NewBufferString(jsonArray))
	if err != nil {
		_ = os.Rename(sendingPath, metricsPath)
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		_ = os.Remove(sendingPath)

		state := LoadState()
		state.Endpoint = endpoint
		state.LastPush = time.Now()
		SaveState(state)

		return nil
	}

	_ = os.Rename(sendingPath, metricsPath)
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("server error %d: %s", resp.StatusCode, string(body))
}
