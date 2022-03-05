#!/bin/sh

# set -e

base=$(dirname $0)

out=some/someopts.go
rm -f $out
mkdir -p $(dirname $out)

function test_update() {
	cp ~/go/bin/goimports ~/go/bin/goimports2
	go run main.go --update
	go run main.go --update --goimports ~/go/bin/goimports2
	echo "ok  	test_update"
	return 0
}

function test_filetoupdate() {
	local filetoupdate=some/filetoupdate.go
	cat << EOF > $filetoupdate
package some

func Func() {
	TakesOpts(Foo(true), Bar("bar"), Baz(1.0))
}

func TakesOpts(opts ...SomeOpts) {
	// nothing
}
EOF
	local optionsfiletoupdate=some/filetoupdateoptions.go
	cat << EOF > $optionsfiletoupdate
package some

//go:generate genopts --prefix=FileToUpdate --outfile=some/filetoupdate.go 'foo' 'bar:int' 'baz:string'

EOF
	go run main.go --update_file $filetoupdate
	local expectedfiletoupdate=$base/testdata/filetoupdate.go
	diff $expectedfiletoupdate $filetoupdate
	echo "ok  	test_filetoupdate"
	return 0
}


function test_filetonotupdate() {
	local filetonotupdate=some/filetonotupdate.go
	cat << EOF > $filetonotupdate
package some

func Func2() {}
EOF
	go run main.go --update_file $filetonotupdate
	local expectedfiletonotupdate=$base/testdata/filetonotupdate.go
	diff $expectedfiletonotupdate $filetonotupdate
	echo "ok  	test_filetonotupdate"
	return 0
}

function test_pipe() {
	local driver=testdata/some/driver.go
	touch $out
	echo "package some" >> $out
	echo >> $out
	go run main.go --opt_type SomeOpts foo bar:string baz:float64 >> $out
	go build $out $driver
	echo "ok  	test_pipe"
	return 0
}


test_update 			|| echo "no  	test_update"
test_filetoupdate 		|| echo "no  	test_filetoupdate"
test_filetonotupdate	|| echo "no  	test_filetonotupdate"
# test_pipe 				&& echo "no  	test_pipe"

echo "ok  	github.com/spudtrooper/genopts/integtest"