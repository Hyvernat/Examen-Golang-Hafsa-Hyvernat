package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BatchWordStats : analyse les mots de tous les .txt et retourne un rapport
func BatchWordStats(dir string, outDir string) error {
	var report strings.Builder
	report.WriteString("=== Batch FileOps Report ===\n\n")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			words, avg, err := WordStats(path)
			if err != nil {
				return err
			}
			report.WriteString(fmt.Sprintf(
				"Fichier: %s\nMots: %d\nLongueur moyenne: %.2f\n\n",
				info.Name(),
				words,
				avg,
			))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outDir, "report.txt"), []byte(report.String()), 0644)
}

// BatchIndex : génère un index des fichiers .txt (chemin, taille, date)
func BatchIndex(dir string, outFile string) error {
	var index strings.Builder
	index.WriteString("Chemin | Taille | Date Modif\n")
	index.WriteString("--- | --- | ---\n")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			index.WriteString(fmt.Sprintf("%s | %d | %s\n", path, info.Size(), info.ModTime().Format("2006-01-02 15:04:05")))
		}
		return nil
	})

	if err != nil {
		return err
	}
	return os.WriteFile(outFile, []byte(index.String()), 0644)
}

// BatchMerge : fusionne tous les fichiers .txt dans un seul fichier
func BatchMerge(dir string, outFile string) error {
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			out.WriteString(fmt.Sprintf("--- FICHIER: %s ---\n", info.Name()))
			out.Write(content)
			out.WriteString("\n\n")
		}
		return nil
	})
	return err
}
