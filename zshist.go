package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var cmdCt = 0                   // command count
var fileCt = 0                  // file count
var nomatch = 1                 // command that doesn't match zsh format
var keys = []string{}           // dedup keys
var lines = map[string]string{} // the zsh history lines

var zshline = regexp.MustCompile(`^(?sm)(\: \d{10,}:\d{1,};)(.*)`) // zsh history line pattern

var v = flag.Bool("v", false, "Shows version")
var home = flag.String("home", os.Getenv("HOME"), "Home dir, defaults to $HOME env")

var Version = "0.0.1" // injected

func main() {
	flag.Parse()

	if *v {
		fmt.Println(Version)
		return
	}

	if *home == "" {
		log.Fatal("Please set $HOME env, or pass in as -home")
	}

	src := "/.zsh_history"
	bak := "/.zsh_history.bak"

	handle(src, bak)

	fmt.Printf("Parsed and merged %d files with %d commands\n", fileCt, cmdCt)
	fmt.Printf("Saved into ~%s with %d commands\n", src, len(keys))
	fmt.Printf("Backed up into ~%s\n", bak)
}

// handler does all the operation
func handle(src, bak string) {
	zsh, err := os.OpenFile(*home+src, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	parse(zsh)
	defer zsh.Close()

	// zsh pre history
	pre := *home + "/.zsh_history.pre-oh-my-zsh"
	if hist, err := os.Open(pre); err == nil {
		parse(hist)
		hist.Close()
		os.Rename(pre, pre+".bak")
	}

	// bash history
	if hist, err := os.Open(*home + "/.bash_history"); err == nil {
		defer hist.Close()
		parse(hist)
	}

	// Backup+save
	if err = backup(zsh, *home+bak); err == nil {
		err = save(zsh, keys, lines)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// parse parses given file handler
func parse(file *os.File) {
	scan := bufio.NewScanner(file)

	for scan.Scan() {
		line := strings.Trim(scan.Text(), " \t")
		if line == "" || line[0] == '#' {
			continue
		}

		// Multiline command
		for line[len(line)-1:] == `\` && scan.Scan() {
			line += "\n" + scan.Text()
		}

		cmd, ts := "", ""
		sub := zshline.FindAllStringSubmatch(line, 1)
		if len(sub) > 0 {
			cmd, ts = strings.Trim(sub[0][2], " \t"), sub[0][1]
		} else {
			cmd, ts, nomatch = line, fmt.Sprintf("%d", nomatch), nomatch+1
		}

		cmdCt++
		if cmd == "" || cmd[0] == '#' {
			continue
		}

		ots, ok := lines[cmd]
		if !ok {
			keys = append(keys, cmd)
		}
		if !ok || ts > ots {
			lines[cmd] = ts
		}
	}

	fileCt++
}

// backup takes a backup with .bak ext
func backup(zsh *os.File, bak string) error {
	zshb, err := os.Create(bak)
	if err != nil {
		return err
	}
	defer zshb.Close()

	// Copy and sync
	zsh.Seek(0, io.SeekStart)
	if _, err = io.Copy(zshb, zsh); err != nil {
		return err
	}
	return zshb.Sync()
}

// save saves the zsh file
func save(zsh *os.File, keys []string, lines map[string]string) (err error) {
	if len(keys) == 0 || len(lines) == 0 {
		return nil
	}

	// Empty
	if err = zsh.Truncate(0); err == nil {
		_, err = zsh.Seek(0, io.SeekStart)
	}
	if err != nil {
		return
	}

	// Sort
	sort.SliceStable(keys, func(i, j int) bool { return lines[keys[i]] < lines[keys[j]] })

	// Write by line
	for _, val := range keys {
		key := lines[val]
		if key[0] != ':' {
			key = ""
		}
		if _, err = zsh.WriteString(key + val + "\n"); err != nil {
			return err
		}
	}

	// Sync
	return zsh.Sync()
}
