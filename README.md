# stackit-nuke

[![license](https://img.shields.io/github/license/cgroschupp/stackit-nuke.svg)](https://github.com/cgroschupp/stackit-nuke/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/cgroschupp/stackit-nuke.svg)](https://github.com/cgroschupp/stackit-nuke/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgroschupp/stackit-nuke)](https://goreportcard.com/report/github.com/cgroschupp/stackit-nuke)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/cgroschupp/stackit-nuke/total)
![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/cgroschupp/stackit-nuke/latest/total)

## Overview

Remove all resources from an StackIT account.

## Example

### Example config

config.yaml:
```yaml
accounts:
  046344C1-B68B-425D-8386-F4789283E2FD: {}
```

### Build locally
```sh
go install .
# Make sure the go path is inside the PATH
# export PATH=$(go env GOPATH)/bin:$PATH
stackit-nuke run --project-id <project-id> -c config.yaml
```