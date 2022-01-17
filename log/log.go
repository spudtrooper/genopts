package log

import reallog "log"

const prefix = "[genopts] "

func Printf(tmpl string, args ...interface{}) {
	reallog.Printf(prefix+tmpl, args...)
}

func Println(s string) {
	reallog.Println(prefix + s)
}
