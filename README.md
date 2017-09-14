# gouv

[![Build Status](https://travis-ci.org/linxGnu/gouv.svg?branch=master)](https://travis-ci.org/linxGnu/gouv)
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/gouv)](https://goreportcard.com/report/github.com/linxGnu/gouv)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/gouv/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/gouv?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/gouv)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jmoiron/sqlx/master/LICENSE)

Compatible with latest libuv (1.14.x) and Golang version >= 1.8. 

Inspired by [go-uv](https://github.com/mattn/go-uv) with bugs fixed and much more additional libuv handles.

# prerequisites

* libuv:

```
wget https://dist.libuv.org/dist/v1.14.1/libuv-v1.14.1.tar.gz
tar xzf libuv-v1.14.1.tar.gz
cd libuv-v1.14.1
sh autogen.sh
./configure
make
sudo make install
sudo ldconfig
```

* pkg-config:

```
sudo apt-get install -y pkg-config
```