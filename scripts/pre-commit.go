package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var pathToGitLeaksWin = "C:\\gitleaks"
var urlForWin = "https://github.com/gitleaks/gitleaks/releases/download/v8.18.1/gitleaks_8.18.1_windows_x64.zip"
var urlForLinux = "https://github.com/gitleaks/gitleaks/releases/download/v8.18.1/gitleaks_8.18.1_linux_x64.tar.gz"
var urlForMac = "https://github.com/gitleaks/gitleaks/releases/download/v8.18.1/gitleaks_8.18.1_darwin_x64.tar.gz"

func main() {
	os := runtime.GOOS
	gitLeaskEnabled := gitConfig("gitleaks.enable")
	if gitLeaskEnabled == "false" || gitLeaskEnabled == "" {
		enableGitLeaks()
	}
	fullPathToGitleaks := ""
	if strings.Contains(strings.ToLower(os), "windows") {
		log.Printf("Operation system defined as: Windows\n")
		fullPathToGitleaks = installGitleaksForWindows(urlForWin)
	} else if strings.Contains(strings.ToLower(os), "linux") {
		log.Printf("Operation system defined as: Linux\n")
		fullPathToGitleaks = installGitleaksForNix(urlForLinux)
	} else if strings.Contains(strings.ToLower(os), "darwin") {
		log.Printf("Operation system defined as: MacOS\n")
		fullPathToGitleaks = installGitleaksForNix(urlForMac)
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

func installGitleaksForNix(downloadURL string) string {
	if isCommandAvailableNix("gitleaks") {
		return "gitleaks"
	}

	extractCmd := "tar xvz -C ./gitleaks && sudo mv gitleaks/* /usr/local/bin/ && rm -d ./gitleaks"
	cmd := exec.Command("sh", "-c", fmt.Sprintf("mkdir -p gitleaks | curl -sSfL %s | %s", downloadURL, extractCmd))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal("Error during gitleaks installation:", err)
	}
	log.Println("Gitleaks installed")
	return "gitleaks"
}

func installGitleaksForWindows(downloadURL string) string {
	fullPathToGitLeaks := fmt.Sprintf("%s\\gitleaks.exe", pathToGitLeaksWin)
	_, err := os.Stat(fmt.Sprintf("%s\\gitleaks.exe", pathToGitLeaksWin))
	if err == nil {
		log.Printf("Found gitleaks: %s\n", fullPathToGitLeaks)
		return fullPathToGitLeaks
	}
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

func enableGitLeaks() {
	cmd := exec.Command("git", "config", "gitleaks.enable", "true")
	cmd.Output()
}

func isCommandAvailableNix(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runGitleaks(pathToGitleaks string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error during current directory definition:", err)
		log.Println(currentDir)

	}
	currentDir = filepath.Dir(currentDir)
	cmd := exec.Command(pathToGitleaks, "protect", "--staged", fmt.Sprintf("-s=%s", currentDir), "--verbose")
	log.Println("cmd = ", cmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Fatal("Error during gitleaks execution:", err)
	}
}
