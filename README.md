# dlsortls

Launch `deno lsp` or `typescript-language-server`.

- If `package.json` exists -> launch typescript-language-server.
- If `deno.json` exists -> launch deno lsp.

## Install

```
go install github.com/yskszk63/dlsortls@latest
```

## Example

nvim-lspconfig

```lua
if not require'lspconfig.configs'.dlsortls then
  require'lspconfig.configs'.dlsortls = {
    default_config = {
      init_options = {
        enable = true,
        lint = false,
        unstable = false,
        hostInfo = 'neovim',
      },
      cmd = { "dlsortls" },
      filetypes = {
        'javascript',
        'javascriptreact',
        'javascript.jsx',
        'typescript',
        'typescriptreact',
        'typescript.tsx',
      },
      root_dir = function(fname)
        local util = require'lspconfig.util'
        return util.root_pattern 'tsconfig.json'(fname)
          or util.root_pattern('package.json', 'jsconfig.json', '.git')(fname)
          or util.root_pattern('deno.json', 'deno.jsonc', 'tsconfig.json', '.git')
      end,
    }
  }
end
nvim_lsp.dlsortls.setup { }
```

## License

[MIT](LICENSE)
