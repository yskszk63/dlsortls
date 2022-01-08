# dlsortss

Launch `deno lsp` or `typescript-language-server`.

If `package.json` exists -> launch typescript-language-server.
If `deno.json` exists -> launch deno lsp.

## Example

nvim-lspconfig

```lua
if not require'lspconfig.configs'.dlsortss then
  require'lspconfig.configs'.dlsortss = {
    default_config = {
      init_options = {
        enable = true,
        lint = false,
        unstable = false,
        hostInfo = 'neovim',
      },
      cmd = { "dlsortss" },
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
nvim_lsp.dlsortss.setup { }
```

## License

[MIT](LICENSE)
