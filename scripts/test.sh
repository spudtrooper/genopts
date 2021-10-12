#!/bin/sh

dir=$(dirname $0)
go test ./genopts
$dir/integtest.sh