package main

import (
	"io/fs"
	"testing"
	"time"
)

type dummyfileinfo struct{}

func (dummyfileinfo) Name() string {
	return "dummy"
}
func (dummyfileinfo) Size() int64 {
	return 0
}
func (dummyfileinfo) Mode() fs.FileMode {
	return fs.FileMode(0)
}
func (dummyfileinfo) ModTime() time.Time {
	return time.Now()
}
func (dummyfileinfo) IsDir() bool {
	return false
}
func (dummyfileinfo) Sys() interface{} {
	return nil
}

type dummyfile struct{}

func (dummyfile) Stat() (fs.FileInfo, error) {
	return dummyfileinfo{}, nil
}
func (dummyfile) Read(b []byte) (int, error) {
	return 0, fs.ErrInvalid
}
func (dummyfile) Close() error {
	return nil
}

type dummyfs struct {
	files map[string]error
}

func (f *dummyfs) Open(name string) (fs.File, error) {
	err, exists := f.files[name]
	if !exists {
		return nil, fs.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	return dummyfile{}, nil
}

func TestExists(t *testing.T) {
	fs := &dummyfs{
		files: map[string]error{
			"main.go":    nil,
			"secret.txt": fs.ErrPermission,
		},
	}

	tests := []struct {
		name string
		path string
		want bool
		err  bool
	}{
		{"exists main.go", "main.go", true, false},
		{"not exists", "notexists.txt", false, false},
		{"perm", "secret.txt", false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v, err := exists(fs, test.path)
			if !test.err && err != nil {
				t.Fail()
			}
			if test.want != v {
				t.Fail()
			}
		})
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		files map[string]error
		name  string
		want  lspKind
		err   bool
	}{
		{map[string]error{"package.json": nil}, "package.json", typescript, false},
		{map[string]error{"jsconfig.json": nil}, "jsconfig.json", typescript, false},
		{map[string]error{"deno.json": nil}, "deno.json", deno, false},
		{map[string]error{"deno.jsonc": nil}, "deno.jsonc", deno, false},
		{map[string]error{}, "unknown", unknown, false},
		{map[string]error{"package.json": fs.ErrPermission}, "package.json err", unknown, true},
		{map[string]error{"deno.json": fs.ErrPermission}, "deno.json err", unknown, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &dummyfs{
				files: test.files,
			}
			kind, err := detect(f)
			if err != nil && !test.err {
				t.Fail()
			}
			if kind != test.want {
				t.Fail()
			}
		})
	}
}

func TestLspKindCmd(t *testing.T) {
	tests := []struct {
		name  string
		input lspKind
		want  []string
	}{
		{"deno", deno, []string{"deno", "lsp"}},
		{"typescript", typescript, []string{"typescript-language-server", "--stdio"}},
		{"unknown", unknown, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := test.input.cmd()
			if cmd == nil && test.want != nil {
				t.Fail()
			}

			if len(cmd) != len(test.want) {
				t.Fail()
			}
			for i, a := range cmd {
				if a != test.want[i] {
					t.Fail()
				}
			}
		})
	}
}

func TestExecProg(t *testing.T) {
	tests := []struct {
		name string
		cmd  []string
		err  string
	}{
		{"directory", []string{"/"}, "exec: \"/\": permission denied"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := execProg(test.cmd)
			if err == nil {
				t.Fail()
			}
			if err.Error() != test.err {
				t.Fatal(err)
			}
		})
	}
}
