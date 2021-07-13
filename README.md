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

Example (built dependency graph comparison):
  git checkout master
  go list -json -deps -test ./... | jq -s 'map(select(.Module != null) | .Module) | unique | .[]' > go.list.json.master
  git checkout mybranchwithchanges
  go list -json -deps -test ./... | jq -s 'map(select(.Module != null) | .Module) | unique | .[]' > go.list.json.mybranchwithchanges
  golistcmp go.list.json.master go.list.json.mybranchwithchanges 

Example (full dependency graph comparison):
  git checkout master
  go list -m -json all > go.list.json.master
  git checkout mybranchwithchanges
  go list -m -json all > go.list.json.mybranchwithchanges
  golistcmp go.list.json.master go.list.json.mybranchwithchanges

Flags:
  -help
        print this help
```
