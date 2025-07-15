package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/darkseear/gophkeeper/client/internal/comand"
)

var (
	version   = "N/A"
	buildDate = "N/A"
)

func main() {
	if runtime.GOOS == "windows" && len(os.Args) == 1 {
		// Получаем путь к текущему исполняемому файлу
		exe, _ := os.Executable()
		exePath, _ := filepath.Abs(exe)

		// Создаем команду для запуска cmd.exe
		cmd := exec.Command("cmd.exe", "/C", "start", "cmd.exe", "/K", exePath, "--help")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}

		// Запускаем
		_ = cmd.Run()
		return
	}
	rootCmd := comand.NewRootCmd(version, buildDate)
	// Выполняем команду
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
