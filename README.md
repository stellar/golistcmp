# golistcmp

[![Build Status](https://github.com/stellar/golistcmp/workflows/Go/badge.svg)](https://github.com/stellar/golistcmp/actions)

A tool for comparing the output of `go list -m -json` executions.

## Install

```
go get github.com/stellar/golistcmp
```

## Usage

```
Usage of golistcmp:
  golistcmp <go list before> <go list after>

Example:
  git checkout master
  go list -m -json all > go.list.master
  git checkout mybranchwithchanges
  golistcmp go.list.master <(go list -m -json all)

Flags:
  -help
        print this help
```
