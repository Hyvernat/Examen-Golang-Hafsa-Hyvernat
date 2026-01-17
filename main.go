package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go-devops-tool/fileops"   // FileOps
	"go-devops-tool/procops"   // ProcOps
	"go-devops-tool/secureops" // SecureOps

	"github.com/PuerkitoBio/goquery" // WebOps
)

// Structure de configuration
type Config struct {
	DefaultFile string `json:"default_file"`
	BaseDir     string `json:"base_dir"`
	OutDir      string `json:"out_dir"`
	DefaultExt  string `json:"default_ext"`
	ProcessTopN int    `json:"process_top_n"`
}

// Chargement de la configuration JSON
func loadConfig(path string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Menu principal
func showMenu() {
	fmt.Println("====== MENU PRINCIPAL ======")
	fmt.Println("1. FileOps - Analyse fichier")
	fmt.Println("2. FileOps - Batch (à implémenter)")
	fmt.Println("3. WebOps - Wikipédia")
	fmt.Println("4. ProcOps - Processus (à implémenter)")
	fmt.Println("5. SecureOps - Sécurité (à implémenter)")
	fmt.Println("0. Quitter")
	fmt.Print("Votre choix : ")
}

// Sous-menu FileOps
func showFileOpsMenu() {
	fmt.Println("------ FileOps ------")
	fmt.Println("1. Infos fichier")
	fmt.Println("2. Statistiques mots")
	fmt.Println("3. Head (N premières lignes)")
	fmt.Println("4. Tail (N dernières lignes)")
	fmt.Println("5. Compter lignes avec mot-clé")
	fmt.Println("6. Filtrer lignes avec/sans mot-clé")
	fmt.Println("0. Retour")
	fmt.Print("Votre choix : ")
}

// Sous-menu ProcOps
func showProcOpsMenu() {
	fmt.Println("------ ProcOps ------")
	fmt.Println("1. Lister les processus (Top N)")
	fmt.Println("2. Rechercher un processus")
	fmt.Println("3. Tuer un processus (Kill)")
	fmt.Println("0. Retour")
	fmt.Print("Votre choix : ")
}

// Sous-menu SecureOps
func showSecureOpsMenu() {
	fmt.Println("------ SecureOps ------")
	fmt.Println("1. Verrouiller un fichier (.lock)")
	fmt.Println("2. Déverrouiller un fichier")
	fmt.Println("3. Basculer Lecture Seule (Windows)")
	fmt.Println("0. Retour")
	fmt.Print("Votre choix : ")
}

func main() {
	// Gestion du flag --config
	configPath := flag.String("config", "config.json", "Chemin vers le fichier de configuration")
	flag.Parse()

	// Chargement de la config
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Println("Erreur chargement config:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration chargée avec succès")
	fmt.Println("Fichier par défaut:", config.DefaultFile)
	fmt.Println()

	// Création du dossier out si inexistant
	err = os.MkdirAll(config.OutDir, os.ModePerm)
	if err != nil {
		fmt.Println("Erreur création dossier out :", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// Boucle du menu principal
	for {
		showMenu()
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1: // FileOps
			for {
				showFileOpsMenu()
				var fchoice int
				fmt.Scanln(&fchoice)

				switch fchoice {
				case 1:
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					if err := fileops.CheckFile(path); err != nil {
						fmt.Println("Fichier invalide :", err)
						break
					}

					size, lines, err := fileops.FileInfo(path)
					if err != nil {
						fmt.Println("Erreur lecture fichier :", err)
						break
					}

					fmt.Println("Taille :", size, "octets")
					fmt.Println("Nombre de lignes :", lines)

				case 2:
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					words, avg, err := fileops.WordStats(path)
					if err != nil {
						fmt.Println("Erreur lecture fichier :", err)
						break
					}

					fmt.Println("Nombre de mots :", words)
					fmt.Printf("Longueur moyenne des mots : %.2f\n", avg)

				case 3:
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					fmt.Print("Nombre de lignes N : ")
					var N int
					fmt.Scanln(&N)

					outFile := filepath.Join(config.OutDir, "head.txt")
					if err := fileops.Head(path, N, outFile); err != nil {
						fmt.Println("Erreur :", err)
						break
					}
					fmt.Println("Head sauvegardé dans :", outFile)

				case 4:
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					fmt.Print("Nombre de lignes N : ")
					var N int
					fmt.Scanln(&N)

					outFile := filepath.Join(config.OutDir, "tail.txt")
					if err := fileops.Tail(path, N, outFile); err != nil {
						fmt.Println("Erreur :", err)
						break
					}
					fmt.Println("Tail sauvegardé dans :", outFile)

				case 5:
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					fmt.Print("Mot-clé : ")
					keyword, _ := reader.ReadString('\n')
					keyword = strings.TrimSpace(keyword)

					count, err := fileops.CountLinesWithKeyword(path, keyword)
					if err != nil {
						fmt.Println("Erreur :", err)
						break
					}
					fmt.Printf("Nombre de lignes contenant '%s' : %d\n", keyword, count)

				case 6:
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					fmt.Print("Mot-clé : ")
					keyword, _ := reader.ReadString('\n')
					keyword = strings.TrimSpace(keyword)

					fmt.Print("Inclure lignes contenant le mot ? (y/n) : ")
					resp, _ := reader.ReadString('\n')
					resp = strings.TrimSpace(strings.ToLower(resp))
					include := resp == "y"

					outFile := filepath.Join(config.OutDir, "filtered.txt")
					if err := fileops.FilterLines(path, keyword, outFile, include); err != nil {
						fmt.Println("Erreur :", err)
						break
					}
					fmt.Println("Filtre sauvegardé dans :", outFile)

				case 0:
					break

				default:
					fmt.Println("Choix invalide")
				}

				if fchoice == 0 {
					break
				}
				fmt.Println()
			}

		case 2: // FileOps Batch
			fmt.Print("Chemin du répertoire (défaut: data) : ")
			dir, _ := reader.ReadString('\n')
			dir = strings.TrimSpace(dir)
			if dir == "" {
				dir = config.BaseDir
			}

			// Report Stats
			if err := fileops.BatchWordStats(dir, config.OutDir); err != nil {
				fmt.Println("Erreur Rapport :", err)
			} else {
				fmt.Println("Rapport généré dans :", filepath.Join(config.OutDir, "report.txt"))
			}

			// Index
			indexFile := filepath.Join(config.OutDir, "index.txt")
			if err := fileops.BatchIndex(dir, indexFile); err != nil {
				fmt.Println("Erreur Index :", err)
			} else {
				fmt.Println("Index généré dans :", indexFile)
			}

			// Merge
			mergeFile := filepath.Join(config.OutDir, "merged.txt")
			if err := fileops.BatchMerge(dir, mergeFile); err != nil {
				fmt.Println("Erreur Fusion :", err)
			} else {
				fmt.Println("Fusion terminée dans :", mergeFile)
			}

		case 3: // WebOps Wikipédia
			fmt.Print("Nom de l'article Wikipédia (ex: Go_(langage)) : ")
			article, _ := reader.ReadString('\n')
			article = strings.TrimSpace(article)

			if article == "" {
				fmt.Println("Article vide !")
				break
			}

			url := "https://fr.wikipedia.org/wiki/" + article
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
				"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Erreur HTTP :", err)
				break
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Println("Page non trouvée, code :", resp.StatusCode)
				break
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				fmt.Println("Erreur parsing :", err)
				break
			}

			text := ""
			doc.Find("div#mw-content-text div.mw-parser-output p").Each(func(i int, s *goquery.Selection) {
				p := strings.TrimSpace(s.Text())
				if p != "" {
					text += p + "\n"
				}
			})

			if text == "" {
				fmt.Println("Aucun texte trouvé dans l'article !")
				break
			}

			// Stats mots
			words := strings.Fields(text)
			totalWords, totalLength := 0, 0
			for _, w := range words {
				if _, err := strconv.Atoi(w); err == nil {
					continue
				}
				totalWords++
				totalLength += len(w)
			}
			avgLength := 0.0
			if totalWords > 0 {
				avgLength = float64(totalLength) / float64(totalWords)
			}

			fmt.Println("Stats de l'article :")
			fmt.Println("Nombre de mots :", totalWords)
			fmt.Printf("Longueur moyenne des mots : %.2f\n", avgLength)

			outFile := filepath.Join(config.OutDir, "wiki_"+article+".txt")
			os.WriteFile(outFile, []byte(text), 0644)
			fmt.Println("Article sauvegardé dans :", outFile)

		case 4: // ProcOps
			for {
				showProcOpsMenu()
				var pchoice int
				fmt.Scanln(&pchoice)

				switch pchoice {
				case 1: // List
					topN := config.ProcessTopN
					if topN == 0 {
						topN = 10
					}
					procs, err := procops.ListProcesses(topN)
					if err != nil {
						fmt.Println("Erreur liste :", err)
						break
					}
					fmt.Printf("%-10s | %-30s\n", "PID", "NOM")
					fmt.Println("-------------------------------------------")
					for _, p := range procs {
						fmt.Printf("%-10s | %-30s\n", p.PID, p.Name)
					}

				case 2: // Filter
					fmt.Print("Rechercher (nom) : ")
					keyword, _ := reader.ReadString('\n')
					keyword = strings.TrimSpace(keyword)

					procs, err := procops.FilterProcesses(keyword)
					if err != nil {
						fmt.Println("Erreur recherche :", err)
						break
					}
					if len(procs) == 0 {
						fmt.Println("Aucun processus trouvé.")
					} else {
						fmt.Printf("%-10s | %-30s\n", "PID", "NOM")
						fmt.Println("-------------------------------------------")
						for _, p := range procs {
							fmt.Printf("%-10s | %-30s\n", p.PID, p.Name)
						}
					}

				case 3: // Kill
					fmt.Print("PID du processus à tuer : ")
					pid, _ := reader.ReadString('\n')
					pid = strings.TrimSpace(pid)

					if pid == "" {
						fmt.Println("PID vide !")
						break
					}

					fmt.Printf("Êtes-vous sûr de vouloir tuer le PID %s ? (yes/no) : ", pid)
					confirm, _ := reader.ReadString('\n')
					confirm = strings.TrimSpace(strings.ToLower(confirm))

					if confirm == "yes" {
						err := procops.KillProcess(pid)
						if err != nil {
							fmt.Println("Erreur lors du kill :", err)
						} else {
							fmt.Println("Processus", pid, "tué avec succès.")
							secureops.LogAction(config.OutDir, "Kill processus: "+pid)
						}
					} else {
						fmt.Println("Action annulée.")
					}

				case 0:
					break
				default:
					fmt.Println("Choix invalide")
				}

				if pchoice == 0 {
					break
				}
				fmt.Println()
			}
		case 5: // SecureOps
			for {
				showSecureOpsMenu()
				var schoice int
				fmt.Scanln(&schoice)

				switch schoice {
				case 1: // Lock
					fmt.Print("Chemin du fichier à verrouiller : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					if err := secureops.LockFile(path, config.OutDir); err != nil {
						fmt.Println("Erreur :", err)
					} else {
						fmt.Println("Fichier verrouillé avec succès.")
					}

				case 2: // Unlock
					fmt.Print("Chemin du fichier à déverrouiller : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					if err := secureops.UnlockFile(path, config.OutDir); err != nil {
						fmt.Println("Erreur :", err)
					} else {
						fmt.Println("Fichier déverrouillé avec succès.")
					}

				case 3: // Read-only
					fmt.Print("Chemin du fichier : ")
					path, _ := reader.ReadString('\n')
					path = strings.TrimSpace(path)

					fmt.Print("Activer Lecture Seule ? (y/n) : ")
					resp, _ := reader.ReadString('\n')
					resp = strings.TrimSpace(strings.ToLower(resp))
					ro := resp == "y"

					if err := secureops.SetReadOnly(path, ro); err != nil {
						fmt.Println("Erreur :", err)
					} else {
						status := "désactivée"
						if ro {
							status = "activée"
						}
						fmt.Println("Attribut Lecture Seule", status)
						secureops.LogAction(config.OutDir, fmt.Sprintf("SetReadOnly (%t): %s", ro, path))
					}

				case 0:
					break
				default:
					fmt.Println("Choix invalide")
				}

				if schoice == 0 {
					break
				}
				fmt.Println()
			}
		case 0:
			fmt.Println("Au revoir")
			return
		default:
			fmt.Println("Choix invalide")
		}
		fmt.Println()
	}
}
