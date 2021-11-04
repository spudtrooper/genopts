#!/bin/sh

set -e

for f in genopts/testdata/*.golden; do 
    gofmt $f > t
    mv t $f
done