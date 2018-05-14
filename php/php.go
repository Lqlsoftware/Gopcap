// +build php

// build with tags php
package php

import (
	"io"

	"github.com/deuill/go-php"
)

var UsePhp = false

type Plugin struct {
	engine	*php.Engine
	context *php.Context
}

func GetThreadPhp() *Plugin {
	engine, err := php.New()
	check(err)
	context, _ := engine.NewContext()
	check(err)
	return &Plugin{engine, context}
}

func (phpPlugin *Plugin)SetPhpWriter(writer io.Writer) {
	phpPlugin.context.Output = writer
}

func (phpPlugin *Plugin)Exec(url string) error {
	return phpPlugin.context.Exec(url)
}

func (phpPlugin *Plugin)closeThreadPhp() {
	phpPlugin.engine.Destroy()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}