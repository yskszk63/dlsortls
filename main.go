package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func runTsserver() error {
	binName := "typescript-language-server"
	bin, err := exec.LookPath(binName)
	if err != nil {
		return err
	}

	args := []string{binName, "--stdio"}
	env := os.Environ()
	return syscall.Exec(bin, args, env)
}

func runDenoLsp() error {
	binName := "deno"
	bin, err := exec.LookPath(binName)
	if err != nil {
		return err
	}

	args := []string{binName, "lsp"}
	env := os.Environ()
	return syscall.Exec(bin, args, env)
}

func main() {
	if exists, err := exists("package.json"); err != nil {
		log.Fatal(err)
	} else if exists {
		if err := runTsserver(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if exists, err := exists("deno.json"); err != nil {
		log.Fatal(err)
	} else if exists {
		runDenoLsp()
		return
	}

	log.Fatal("failed to detect launch `deno lsp` or `typescript-language-server`.")
}
