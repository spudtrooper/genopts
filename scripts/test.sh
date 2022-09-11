#!/bin/sh

dir=$(dirname $0)
rm -f gen/options.go  gen/someoptionoptions.go  gen/prefixoptions.go
go test gen/*.go
rm -f gen/options.go  gen/someoptionoptions.go  gen/prefixoptions.go
go test gitversion/*.go
$dir/integtest.sh
rm -f gen/options.go  gen/someoptionoptions.go  gen/prefixoptions.go