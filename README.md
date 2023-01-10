# string-checker
A simple Github Action to check if certain strings are included in the source code.

## Inputs
### filePattern
*optional* - `string` - file pattern to search for. (refer to [pattern syntax](https://pkg.go.dev/path/filepath#Match).)

### filePaths
*optional* - `string` - file paths to search for.
### confFilePath
*optional* - `string` - path to config file. default is `string-checker/config.yml`.

## Config
### name
*required* - `string` - name of the rule.

### message
*required* - `string` - message to show when target strings are detected.

### level
*required* - `(error | warning)` - whether to just warn or to make the workflow failed.

### caseSensitive
*required* - `bool` - whether to check target strings with case-sensitive.

### targets
*required* - `[]string` - target strings.

## How to use
Put a config file in your repository and set strings that you want search for, notification level and so on. Configure a GitHub Action with inputs. (your repository must be checked out before running this action) `filePattern` or `filePaths` must be specified. If you set both `filePattern` and `filePaths`, then only files matching both will be targeted.


## Example
`string-checker/config.yml`
```yml
rules:
  - name: "Hoge"
    message: "Hoge is not allowed to use"
    level: error
    caseSensitive: false
    targets:
      - hoge
```

`.github/workflows/string-checker.yml`
```yml
name: string checker
on: [push]
jobs:
build:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@master
    - name: test
      id: test
      uses: sy-tencho/string-checker@main
      with:
        filePattern: '**/*.sql'
```

## Tips
If only `filePattern` is set, GitHub Action may show "unchanged files with check annotations". To prevent this, set changed files in `filePaths`.
