# string-checker
A simple Github Action to check if certain strings are included in source code.

## How to use
#### Put a config file in your repository

e.g.
```yml
rules:
- name: "Hoge"
  message: "Hoge is not allowed to use"
  level: error
  caseSensitive: false
  targets:
    - hoge
```

**name**: name of a rule
**message**: message that you want to show when target strings are included
**level (error | warning)**: whether you want to just warn or make the workflow failed
**caseSensitive (bool)**: whether you want to check target strings with case-sensitive
**targets**: target strings

#### Add a workflow

e.g.
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
    confFilePath: 'string-checker/config.yml'
```

**filePattern**: file pattern that you explore (required, [pattern syntax](https://pkg.go.dev/path/filepath#Match))
**confFilePath**: path to config file (optional, default is `string-checker/config.yml`)
