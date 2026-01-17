# Go DevOps Tool - EXAMEN PROJET GO

Ce projet est un outil en console développé en Go pour manipuler des fichiers, extraire des données web, gérer des processus système et appliquer des mesures de sécurité.

## Information Étudiant
- **Nom** : Hafsa
- **Date** : 16/01/2026
- **Niveau visé** : 18/20

## Fonctionnalités implémentées

### Niveau 10 : FileOps
- **Menu interactif** : Boucle de menu complète.
- **Configuration** : Lecture initiale depuis `config.json` (avec flag `--config`).
- **Analyse de fichier** : Taille, lignes, stats mots (ignorant les numériques), filtres (mots-clés), Head / Tail.
- **Traitement par lot (Batch)** : Analyse de tous les `.txt`, génération d'un `index.txt`, `report.txt` et fusion dans `merged.txt`.

### Niveau 13 : WebOps
- **Wikipédia** : Extraction de paragraphes via `goquery`.
- **Analyse** : Application des statistiques de mots sur le contenu extrait.
- **Sortie** : Sauvegarde dans `out/wiki_<article>.txt`.

### Niveau 16 : ProcOps
- **Multi-plateforme** : Fonctionne sur Windows (tasklist/taskkill) et macOS (ps/kill).
- **Lister** : Liste les N premiers processus.
- **Filtrer** : Recherche par mot-clé dans les noms de processus.
- **Kill sécurisé** : Confirmation explicite avant de terminer un processus.

### Niveau 18 : SecureOps
- **Verrouillage** : Système de *Lockfile* (.lock) pour simuler un verrouillage de fichier.
- **Lecture seule** : Modification des attributs système (Windows via `attrib`, macOS via `chmod`).
- **Journalisation** : Audit de toutes les actions sensibles (Kill, Lock, RO) dans `out/audit.log`.

## Procédure d'exécution

1. **Prérequis** :
   - Go installé.
   - Dépendances : `go get github.com/PuerkitoBio/goquery`.

2. **Installation** :
   ```bash
   go mod tidy
   ```

3. **Exécution** :
   ```bash
   go run main.go
   # Ou avec une config spécifique
   go run main.go --config config.json
   ```

## Description du travail effectué
Le projet a été structuré en paquets (`fileops`, `procops`, `secureops`) pour une meilleure maintenabilité. La gestion des erreurs est centrale, assurant que les entrées utilisateurs invalides ou les problèmes système ne fassent pas planter le programme. L'utilisation de `runtime.GOOS` permet une portabilité réelle entre Windows et macOS pour les outils système.
