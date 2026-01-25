package utils

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"xytz/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

func StartDownload(program *tea.Program, url, formatID string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		go doDownload(program, url, formatID)
		return nil
	})
}

func doDownload(program *tea.Program, url, formatID string) {
	args := []string{"-f", formatID, "--no-playlist", "--newline", "-R", "infinite", url}
	cmd := exec.Command("yt-dlp", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("pipe error: %v", err)
		errMsg := fmt.Sprintf("pipe error: %v", err)
		program.Send(types.DownloadResultMsg{Err: errMsg})
		return
	}

	stderr, err2 := cmd.StderrPipe()
	if err2 != nil {
		log.Printf("stderr pipe error: %v", err2)
		errMsg := fmt.Sprintf("stderr pipe error: %v", err2)
		program.Send(types.DownloadResultMsg{Err: errMsg})
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("start error: %v", err)
		errMsg := fmt.Sprintf("start error: %v", err)
		program.Send(types.DownloadResultMsg{Err: errMsg})
		return
	}

	parser := NewProgressParser()
	readPipe := func(pipe io.Reader, source string) {
		parser.ReadPipe(pipe, func(percent float64, speed, eta string) {
			log.Printf("Progress from %s: %.2f%%, speed: %s, eta: %s", source, percent, speed, eta)
			program.Send(types.ProgressMsg{Percent: percent, Speed: speed, Eta: eta})
		})
	}

	go readPipe(stdout, "stdout")
	go readPipe(stderr, "stderr")
	err = cmd.Wait()

	if stdout != nil {
		stdout.Close()
	}
	if stderr != nil {
		stderr.Close()
	}
	if err != nil {
		errMsg := fmt.Sprintf("Download error: %v", err)
		program.Send(types.DownloadResultMsg{Err: errMsg})
	} else {
		program.Send(types.DownloadResultMsg{Output: "Download complete"})
	}
}
