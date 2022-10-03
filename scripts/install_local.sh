#!/bin/sh

set -e

go build main.go
# cp main ~/go/bin/genopts
rm -f ~/go/bin/genopts
ln -fns /Users/jeff/Projects/genopts/main ~/go/bin/genopts