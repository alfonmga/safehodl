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

[![safehodl demo](https://asciinema.org/a/430538.svg)](https://asciinema.org/a/430538)

_A random Bitcoin holdings amount was used for the demo._

## Features

- [x] Critical data (`$HOME/.safehodl`) is protected thanks to [Argon2id](https://en.wikipedia.org/wiki/Argon2) + [Argon2id KDF](https://pkg.go.dev/golang.org/x/crypto/argon2#IDKey) + [AES-256 bits encryption](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) + [Galois/Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode
- [x] Passphrase protection for secure SafeHODL access and usage
- [x] Remove data if incorrect passphrase is entered
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
