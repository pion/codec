#!/bin/sh

go build && ./examples | ffplay -i - -vf format=yuv422p -vcodec h264