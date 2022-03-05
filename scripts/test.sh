#!/bin/sh

dir=$(dirname $0)
go test gen gitversion
$dir/integtest.sh