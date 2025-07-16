package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/darkseear/gophkeeper/client/internal/command"
)

var (
	version   = "N/A"
	buildDate = "N/A"
)

func main() {
	if runtime.GOOS == "windows" && len(os.Args) == 1 {
		// Получаем путь к текущему исполняемому файлу
		exe, err := os.Executable()
		if err != nil {
			log.Fatalf("Error getting executable path: %v\n", err)
		}

		exePath, err := filepath.Abs(exe)
		if err != nil {
			log.Fatalf("Error getting absolute path: %v\n", err)
		}

		// Создаем команду для запуска cmd.exe
		cmd := exec.Command("cmd.exe", "/C", "start", "cmd.exe", "/K", exePath, "--help")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}

		// Запускаем
		if err = cmd.Run(); err != nil {
			log.Fatalf("Error running cmd: %v\n", err)
		}

		return
	}
	rootCmd := command.NewRootCmd(version, buildDate)
	// Выполняем команду
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
