package storage

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iliamikado/UrlShortener/internal/logger"
)

type FileSaver struct {
	file *os.File
}

type SavedURL struct {
	ID 			string	`json:"uuid"`
	ShortURL 	string	`json:"short_url"`
	OriginURL 	string	`json:"origin_url"`
	UserID		uint	`json:"user_id"`
}

func NewFileSaver(pathToFile string) *FileSaver {
	logger.Log.Info("Create file with path " + pathToFile)
	os.Mkdir(filepath.Dir(pathToFile), 0777)
	file, _ := os.OpenFile(pathToFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	return &FileSaver{file}
}

func (fs *FileSaver) GetAllData() (map[string]string, map[uint][]string) {
	m := make(map[string]string)
	usersURLs := make(map[uint][]string)
	b, _ := io.ReadAll(fs.file)
	for _, s := range strings.Split(string(b), "\n") {
		var savedURL SavedURL
		err := json.Unmarshal([]byte(s), &savedURL)
		if err == nil {
			m[savedURL.ID] = savedURL.OriginURL
			usersURLs[savedURL.UserID] = append(usersURLs[savedURL.UserID], savedURL.ID)
		}
	}
	return m, usersURLs
}

func (fs *FileSaver) AddURL(savedURL SavedURL) {
	str, _ := json.Marshal(savedURL)
	fs.file.WriteString(string(str) + "\n")
}