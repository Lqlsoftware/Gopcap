// +build !php

// default build with out php
package php

var UsePhp = false

type Plugin struct {
	engine	interface{}
	context interface{}
}


func GetThreadPhp() *Plugin {
	return nil
}

func (phpPlugin *Plugin)SetPhpWriter(a interface{}) {
	return
}

func (phpPlugin *Plugin)Exec(a interface{}) error {
	return nil
}

func (phpPlugin *Plugin)closeThreadPhp() {
	return
}