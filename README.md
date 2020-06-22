[![Build Status](https://travis-ci.org/sitewhere/swctl.svg?branch=master)](https://travis-ci.org/sitewhere/swctl) [![Go Report Card](https://goreportcard.com/badge/github.com/sitewhere/swctl)](https://goreportcard.com/report/github.com/sitewhere/swctl) [![GoDoc](https://godoc.org/github.com/sitewhere/swctl?status.svg)](https://godoc.org/github.com/sitewhere/swctl)

![SiteWhere](https://s3.amazonaws.com/sitewhere-branding/SiteWhereLogo.svg)

---

# SiteWhere Control CLI

## Build

For building it requires go 1.11+.

```console
go build
```

## Install

```console
go install
```

### Install swctl binary with curl on Linux

1 - Download the latest release and install with the command:

```bash
curl -L https://github.com/sitewhere/swctl/releases/download/v0.0.2/swctl.linux.amd64 -o swctl && \
chmod +x ./swctl && sudo mv ./swctl /usr/local/bin/swctl
```

### Install swctl binary with curl on macOS

1 - Download the latest release and install with the command:

```bash
curl -L https://github.com/sitewhere/swctl/releases/download/v0.0.2/swctl.darwin.amd64 -o swctl && \
chmod +x ./swctl && sudo mv ./swctl /usr/local/bin/swctl
```

### Install swctl binary with curl on Windows

```bash
curl -L https://github.com/sitewhere/swctl/releases/download/v0.0.2/swctl.windows.amd64.exe -o swctl.exe
```
