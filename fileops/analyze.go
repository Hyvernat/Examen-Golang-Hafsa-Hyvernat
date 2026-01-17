package fileops

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Vérifie que le chemin existe et que c’est un fichier
func CheckFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return os.ErrInvalid
	}
	return nil
}

// Infos sur le fichier : taille en octets et nombre de lignes
func FileInfo(path string) (int64, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}
	return info.Size(), lines, nil
}

// Statistiques mots : nombre de mots (ignore les nombres) et longueur moyenne
func WordStats(path string) (int, float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordCount := 0
	totalLength := 0

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			if _, err := strconv.Atoi(word); err == nil {
				continue
			}
			wordCount++
			totalLength += len(word)
		}
	}

	avg := 0.0
	if wordCount > 0 {
		avg = float64(totalLength) / float64(wordCount)
	}
	return wordCount, avg, nil
}

// Compte les lignes contenant un mot-clé
func CountLinesWithKeyword(path, keyword string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			count++
		}
	}
	return count, nil
}

// Filtre les lignes contenant ou ne contenant pas le mot-clé
func FilterLines(path, keyword, outFile string, include bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if (include && strings.Contains(line, keyword)) || (!include && !strings.Contains(line, keyword)) {
			out.WriteString(line + "\n")
		}
	}
	return nil
}

// N premières lignes → head.txt
func Head(path string, N int, outFile string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() && count < N {
		out.WriteString(scanner.Text() + "\n")
		count++
	}
	return nil
}

// N dernières lignes → tail.txt
func Tail(path string, N int, outFile string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	start := len(lines) - N
	if start < 0 {
		start = 0
	}

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, line := range lines[start:] {
		out.WriteString(line + "\n")
	}
	return nil
}
