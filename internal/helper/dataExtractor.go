package helper

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DataGame struct {
	Name       string
	PathToExec string
	PathToIcon string
}

func ExtractGamesData(pathToFolder string) ([]DataGame, error) {
	allPaths, err := GetAllLNKFiles(pathToFolder, ".lnk")
	if err != nil {
		return nil, err
	}
	result := make([]DataGame, 0, len(allPaths))

	for _, val := range allPaths {
		result = append(result, ExtractGameData(val))
	}
	return result, nil
}

func ExtractGameData(pathToFolder string) DataGame {
	var result DataGame
	// Инициализация COM
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatalf("Ошибка инициализации COM: %v", err)
	}
	defer ole.CoUninitialize()

	// Создание объекта WScript.Shell
	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		log.Fatalf("Ошибка создания объекта WScript.Shell: %v", err)
	}
	defer unknown.Release()

	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatalf("Ошибка получения интерфейса IDispatch: %v", err)
	}
	defer shell.Release()

	// Открытие ярлыка
	shortcut, err := oleutil.CallMethod(shell, "CreateShortcut", pathToFolder)
	if err != nil {
		log.Fatalf("Ошибка открытия ярлыка: %v", err)
	}
	defer shortcut.ToIDispatch().Release()

	// Извлечение пути к целевому файлу
	targetPath, err := oleutil.GetProperty(shortcut.ToIDispatch(), "TargetPath")
	if err != nil {
		log.Fatalf("Ошибка получения TargetPath: %v", err)
	}
	result.PathToExec = targetPath.ToString()

	// Извлечение последней рабочей директории(имени игры)
	workingDir, err := oleutil.GetProperty(shortcut.ToIDispatch(), "WorkingDirectory")
	if err != nil {
		log.Fatalf("Ошибка получения WorkingDirectory: %v", err)
	}
	result.Name = filepath.Base(workingDir.ToString())

	// Извлечение местоположения иконки
	iconLocation, err := oleutil.GetProperty(shortcut.ToIDispatch(), "IconLocation")
	if err != nil {
		log.Fatalf("Ошибка получения IconLocation: %v", err)
	}
	result.PathToIcon = iconLocation.ToString()

	return result
}

func GetAllLNKFiles(pathToDir string, format string) ([]string, error) {
	// Массив для хранения найденных файлов
	var files []string

	// Функция для поиска файлов с нужным расширением
	err := filepath.Walk(pathToDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Проверяем, является ли объект файлом и соответствует ли его расширение
		if !info.IsDir() && strings.HasSuffix(info.Name(), format) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка при обходе директории: %v\n", err)
	}

	return files, nil
}
