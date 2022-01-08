package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type lspKind uint8

const (
	unknown lspKind = iota + 1
	deno
	typescript
)

func (v lspKind) cmd() []string {
	switch v {
	case deno:
		return []string{"deno", "lsp"}
	case typescript:
		return []string{"typescript-language-server", "--stdio"}
	default:
		return nil
	}
}

func exists(f fs.FS, fnames ...string) (bool, error) {
	for _, fname := range fnames {
		if _, err := fs.Stat(f, fname); err == nil {
			return true, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
	}

	return false, nil
}

func detect(f fs.FS) (lspKind, error) {
	if exists, err := exists(f, "package.json", "jsconfig.json"); err != nil {
		return unknown, err
	} else if exists {
		return typescript, nil
	}

	if exists, err := exists(f, "deno.json", "deno.jsonc"); err != nil {
		return unknown, err
	} else if exists {
		return deno, nil
	}

	return unknown, nil
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
	fs := os.DirFS(".")
	kind, err := detect(fs)
	if err != nil {
		log.Fatal(err)
	}

	if kind == unknown {
		log.Fatal("Couldn't decide whether to launch `deno lsp` or `typescript-language-server`.")
	}

	if err := execProg(kind.cmd()); err != nil {
		log.Fatal(err)
	}
}
