#!/bin/sh

go build && ./stream | mpv --no-cache --untimed --no-demuxer-thread --vd-lavc-threads=1 - 