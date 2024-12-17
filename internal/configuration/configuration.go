package configuration

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	e "github.com/liioan/faek/internal/errors"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

const settingsFilePath = "/faek_settings.json"
const settingsDirectoryPath = "/.config/faek"

type Settings struct {
	OutputStyle string `json:"outputStyle"`
	FileName    string `json:"fileName"`
	Language    string `json:"lang"`
	Indent      string `json:"indent"`
}

func getConfigDirectory() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dirname + settingsDirectoryPath, nil
}
func getConfigFilePath() (string, error) {
	dirname, err := getConfigDirectory()
	if err != nil {
		return "", err
	}
	return dirname + settingsFilePath, nil
}

func SaveUserSettings(settings *Settings) error {
	utils.LogToDebug(settings.Language)

	if settings.FileName == "" {
		settings.FileName = "faekOutput.ts"
	}

	if settings.OutputStyle == "" {
		settings.OutputStyle = string(v.Terminal)
	}

	if settings.Language == "" {
		settings.Language = string(v.TypeScript)
	}

	if settings.Indent == "" {
		settings.Indent = "2"
	}

	settings.FileName = strings.Split(settings.FileName, ".")[0]
	if settings.Language == string(v.TypeScript) {
		settings.FileName += ".ts"
	} else {
		settings.FileName += ".js"
	}

	bytes, err := json.Marshal(settings)
	if err != nil {
		return errors.New(e.CantMarshalJson)
	}

	filePath, err := getConfigDirectory()
	if err != nil {
		return errors.New(e.CantCreateConfigDirectory)
	}

	os.MkdirAll(filePath, 0755)
	filePath, err = getConfigFilePath()
	if err != nil {
		return errors.New(e.CantCreateConfigDirectory)
	}

	file, _ := os.Create(filePath)
	_, err = file.Write(bytes)

	if err != nil {
		return errors.New(e.CanSaveToFile)
	}
	return nil
}

func GetUserSettings() (Settings, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Settings{}, errors.New(e.FileDoesNotExists)
	}
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return Settings{}, errors.New(e.FileDoesNotExists)
	}
	s := Settings{}
	err = json.Unmarshal(fileBytes, &s)

	if err != nil {
		return Settings{}, errors.New(e.CantUnmarshalJson)
	}

	return s, nil
}
