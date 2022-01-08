package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"syscall"
	"fmt"
	"path"
)

type lspKind uint8

const (
	deno lspKind = iota + 1
	typescript
)

func exists(cwd string, fnames ...string) (bool, error) {
	for _, f := range fnames {
		p := path.Join(cwd, f)
		if _, err := os.Stat(p); err == nil {
			return true, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
	}

	return false, nil
}

func detect(cwd string) (*lspKind, error) {
	if exists, err := exists(cwd, "package.json", "jsconfig.json"); err != nil {
		return nil, err
	} else if exists {
		r := typescript
		return &r, nil
	}

	if exists, err := exists(cwd, "deno.json", "deno.jsonc"); err != nil {
		return nil, err
	} else if exists {
		r := deno
		return &r, nil
	}

	return nil, fmt.Errorf("failed to detect to launch `deno lsp` or `typescript-language-server`.")
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

func execProg(args []string) error {
	bin := args[0]
	bin, err := exec.LookPath(bin)
	if err != nil {
		return err
	}

	env := os.Environ()
	return syscall.Exec(bin, args, env)
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	kind, err := detect(cwd)
	if err != nil {
		log.Fatal(err)
	}

	var args []string
	switch *kind {
	case deno:
		args = []string{"deno", "lsp"}
	case typescript:
		args = []string{"typescript-language-server", "--stdio"}
	}

	if err := execProg(args); err != nil {
		log.Fatal(err)
	}
}
