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

cp ~/go/bin/goimports ~/go/bin/goimports2
go run main.go --update
go run main.go --update --goimports ~/go/bin/goimports2


touch $out
echo "package some" >> $out
echo >> $out
go run main.go --opt_type SomeOpts foo bar:string baz:float64 >> $out && \
go build $out $driver

echo "ok  	github.com/spudtrooper/genopts/integtest"