# safehodl

> Track your Bitcoin holdings value in a safe way from your terminal

```sh
  /$$$$$$             /$$$$$$          /$$   /$$  /$$$$$$  /$$$$$$$  /$$
 /$$__  $$           /$$__  $$        | $$  | $$ /$$__  $$| $$__  $$| $$
| $$  \__/  /$$$$$$ | $$  \__//$$$$$$ | $$  | $$| $$  \ $$| $$  \ $$| $$
|  $$$$$$  |____  $$| $$$$   /$$__  $$| $$$$$$$$| $$  | $$| $$  | $$| $$
 \____  $$  /$$$$$$$| $$_/  | $$$$$$$$| $$__  $$| $$  | $$| $$  | $$| $$
 /$$  \ $$ /$$__  $$| $$    | $$_____/| $$  | $$| $$  | $$| $$  | $$| $$
|  $$$$$$/|  $$$$$$$| $$    |  $$$$$$$| $$  | $$|  $$$$$$/| $$$$$$$/| $$$$$$$$
 \______/  \_______/|__/     \_______/|__/  |__/ \______/ |_______/ |________/

```

## Demo

[![asciicast](https://asciinema.org/a/nWronQ3nPXC3LdWt3DulWUCHt.svg)](https://asciinema.org/a/nWronQ3nPXC3LdWt3DulWUCHt)

## Features

- [x] Data on disk is encrypted using [AES encryption](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) + [GCM encryption algorithm](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
- [x] Pin code protection for secure access
- [x] Remove data if wrong access pin code is entered
- [x] Bitcoin holdings value calculated in USD and EUR ([Kraken.com public API](https://docs.kraken.com/rest/))

## Usage

```sh
$ safehodl --help

Track your Bitcoin holdings value in a safe way
https://github.com/alfonmga/safehodl

Usage:
  safehodl [flags]
  safehodl [command]

Available Commands:
  config      Configure SafeHODL
  help        Help about any command

Flags:
  -h, --help   help for safehodl

Use "safehodl [command] --help" for more information about a command.
```

## Build

Requirements:

- Golang

```sh
$ git clone https://github.com/alfonmga/safehodl
$ cd safehodl/
$ make build
```
