# Concurrent files reading test

## How to run instructions

### golang
on a local machine from the current folder
```shell
$ go run main.go
```
in a container
```shell
$ docker build -f Dockerfile.golang -t gotest .
$ docker run --rm gotest
```

### javascript (node)
on a local machine from the current folder
```shell
$ node async.js
```
in a container
```shell
$ docker build -f Dockerfile.nodejs -t jstest .
$ docker run --rm jstest
```