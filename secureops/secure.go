package secureops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// LogAction enregistre une action dans le fichier audit.log
func LogAction(outDir, action string) error {
	logFile := filepath.Join(outDir, "audit.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, action))
	return err
}

// LockFile crée un fichier de verrouillage (.lock)
func LockFile(path, outDir string) error {
	lockFile := filepath.Join(outDir, filepath.Base(path)+".lock")
	if _, err := os.Stat(lockFile); err == nil {
		return fmt.Errorf("le fichier est déjà verrouillé")
	}

	err := os.WriteFile(lockFile, []byte("LOCKED"), 0644)
	if err != nil {
		return err
	}

	return LogAction(outDir, "Verrouillage de : "+path)
}

// UnlockFile supprime le fichier de verrouillage
func UnlockFile(path, outDir string) error {
	lockFile := filepath.Join(outDir, filepath.Base(path)+".lock")
	if _, err := os.Stat(lockFile); os.IsNotExist(err) {
		return fmt.Errorf("le fichier n'est pas verrouillé")
	}

	err := os.Remove(lockFile)
	if err != nil {
		return err
	}

	return LogAction(outDir, "Déverrouillage de : "+path)
}

// IsLocked vérifie si un fichier est verrouillé
func IsLocked(path, outDir string) bool {
	lockFile := filepath.Join(outDir, filepath.Base(path)+".lock")
	_, err := os.Stat(lockFile)
	return err == nil
}

// SetReadOnly tente de rendre un fichier en lecture seule selon l'OS
func SetReadOnly(path string, readOnly bool) error {
	if runtime.GOOS == "windows" {
		attr := "+R"
		if !readOnly {
			attr = "-R"
		}
		cmd := exec.Command("attrib", attr, path)
		return cmd.Run()
	} else {
		mode := os.FileMode(0644)
		if readOnly {
			mode = os.FileMode(0444)
		}
		return os.Chmod(path, mode)
	}
}
