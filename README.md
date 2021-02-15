[![Go Report Card](https://goreportcard.com/badge/github.com/barandemirbas/ambystoma)](https://goreportcard.com/report/github.com/barandemirbas/ambystoma)
[![Go Reference](https://pkg.go.dev/badge/github.com/barandemirbas/ambystoma.svg)](https://pkg.go.dev/github.com/barandemirbas/ambystoma)

## What is Ambystoma?
**Ambystoma** is an HTTP server for **Frontend** Developers. 
You can serve your **HTML**, **JavaScript** and **CSS** files with live reload.

## Usage

```sh
ambystoma
```
Run this command in your work directory to serve files after installation.

## Install

### Source
```sh
git clone https://github.com/barandemirbas/ambystoma.git && cd ambystoma
make build && ls ./bin/
```
Go required for install from source.

### Linux 64-BIT

```sh
sudo wget https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-linux-x64 -O /usr/bin/ambystoma && sudo chmod +x /usr/bin/ambystoma

```

### Linux 32-BIT 
```sh
sudo wget https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-linux-x32 -O /usr/bin/ambystoma && sudo chmod +x /usr/bin/ambystoma

```

### Linux 64-BIT ARM

```sh
sudo wget https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-linux-arm64 -O /usr/bin/ambystoma && sudo chmod +x /usr/bin/ambystoma

```

### Linux 32-BIT ARM
```sh
sudo wget https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-linux-arm32 -O /usr/bin/ambystoma && sudo chmod +x /usr/bin/ambystoma

```

### MacOS
```sh
curl -LJO https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-mac-x64 && mv ambystoma-mac-x64 /usr/local/bin/ambystoma && chmod +x /usr/local/bin/ambystoma
```

### MacOS M1
```sh
git clone https://github.com/barandemirbas/ambystoma.git && cd ambystoma
make m1 && mv ./bin/ambystoma /usr/local/bin/ambystoma && chmod +x /usr/local/bin/ambystoma
```
Go required for install MacOS M1 release.

### Windows 64-BIT
[Download Here.](https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-windows-x64.exe)

### Windows 32-BIT
[Download Here.](https://github.com/barandemirbas/ambystoma/releases/download/v0.0.1/ambystoma-windows-x64.exe)
