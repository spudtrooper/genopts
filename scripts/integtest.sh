#!/bin/sh

set -e

out=some/someopts.go
driver=some/driver.go
rm -f $out
mkdir -p $(dirname $out)

cat << EOF > $driver
package some

func Func() {
	TakesOpts(Foo(true), Bar("bar"), Baz(1.0))
}

func TakesOpts(opts ...SomeOpts) {
	// nothing
}
EOF


touch $out
echo "package some" >> $out
echo >> $out
go run main.go --opts_type SomeOpts foo bar:string baz:float64 >> $out && \
go build $out $driver

echo "ok  	github.com/spudtrooper/genopts/integtest"