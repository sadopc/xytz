package utils

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const UnfinishedFileName = ".xytz_unfinished.json"

type UnfinishedDownload struct {
	URL       string    `json:"url"`
	FormatID  string    `json:"format_id"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
}

func GetUnfinishedFilePath() string {
	user, err := user.Current()
	if err != nil {
		return UnfinishedFileName
	}

	localDir := filepath.Join(user.HomeDir, ".local", "share", "xytz")

	if err := os.MkdirAll(localDir, 0755); err != nil {
		return filepath.Join(user.HomeDir, UnfinishedFileName)
	}

	return filepath.Join(localDir, UnfinishedFileName)
}

func LoadUnfinished() ([]UnfinishedDownload, error) {
	path := GetUnfinishedFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []UnfinishedDownload{}, nil
		}

		return nil, err
	}

	var downloads []UnfinishedDownload
	if err := json.Unmarshal(data, &downloads); err != nil {
		return nil, err
	}

	return downloads, nil
}

func SaveUnfinished(downloads []UnfinishedDownload) error {
	path := GetUnfinishedFilePath()
	data, err := json.MarshalIndent(downloads, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func AddUnfinished(download UnfinishedDownload) error {
	downloads, err := LoadUnfinished()
	if err != nil {
		return err
	}

	for i, d := range downloads {
		if d.URL == download.URL {
			downloads[i] = download
			return SaveUnfinished(downloads)
		}
	}

	downloads = append(downloads, download)
	return SaveUnfinished(downloads)
}

func RemoveUnfinished(url string) error {
	downloads, err := LoadUnfinished()
	if err != nil {
		return err
	}

	var newDownloads []UnfinishedDownload
	for _, d := range downloads {
		if d.URL != url {
			newDownloads = append(newDownloads, d)
		}
	}

	return SaveUnfinished(newDownloads)
}

func GetUnfinishedByURL(url string) *UnfinishedDownload {
	downloads, err := LoadUnfinished()
	if err != nil {
		return nil
	}

	for _, d := range downloads {
		if d.URL == url {
			return &d
		}
	}

	return nil
}
