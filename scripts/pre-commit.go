package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var pathToGitLeaksWin = "C:\\gitleaks"

func main() {
	os := runtime.GOOS
	gitLeaskEnabled := gitConfig("gitleaks.enable")
	if gitLeaskEnabled == "false" || gitLeaskEnabled == "" {
		return
	}
	fullPathToGitleaks := ""
	if strings.Contains(strings.ToLower(os), "windows") {
		log.Printf("Operation system defined as: Windows\n")
		fullPathToGitleaks = installGitleaksForWindows()
	}
	if strings.Contains(strings.ToLower(os), "linux") {
		installGitleaksForLinux()
	}
	if strings.Contains(strings.ToLower(os), "darwin") {
		installGitleaksForMac()
	}

	runGitleaks(fullPathToGitleaks)
}
func gitConfig(key string) string {
	cmd := exec.Command("git", "config", "--get", key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func installGitleaksForMac() {
	panic("unimplemented")
}

func installGitleaksForLinux() {
	panic("unimplemented")
}

func installGitleaksForWindows() string {
	fullPathToGitLeaks := fmt.Sprintf("%s\\gitleaks.exe", pathToGitLeaksWin)
	_, err := os.Stat(fmt.Sprintf("%s\\gitleaks.exe", pathToGitLeaksWin))
	if err == nil {
		log.Printf("Found gitleaks: %s\n", fullPathToGitLeaks)
		return fullPathToGitLeaks
	}
	downloadURL := "https://github.com/gitleaks/gitleaks/releases/download/v8.18.1/gitleaks_8.18.1_windows_x64.zip"
	extractCmd := fmt.Sprintf("Expand-Archive %s\\gitleaks.zip -DestinationPath %s", pathToGitLeaksWin, pathToGitLeaksWin)
	mkdir := fmt.Sprintf("mkdir %s", pathToGitLeaksWin)
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("%s; Invoke-WebRequest -Uri %s -OutFile %s\\gitleaks.zip ; %s", mkdir, downloadURL, pathToGitLeaksWin, extractCmd))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Fatal("Error during gitleaks installation:", err)
	}
	log.Printf("Gitleaks installed: %s\n", fullPathToGitLeaks)
	return fullPathToGitLeaks
}

func runGitleaks(pathToGitleaks string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error duting current directory difinition:", err)
	}
	cmd := exec.Command(pathToGitleaks, "detect", fmt.Sprintf("-s=%s", currentDir), "--verbose")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Println("Gitleaks return validation error")
		os.Exit(1)
	}
}
