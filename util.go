package main

import (
	"flag"
	"strings"
)

type stringListValue []string

func newStringListValue(val []string, p *[]string) *stringListValue {
	*p = val
	return (*stringListValue)(p)
}

func (s *stringListValue) Get() interface{} {
	return []string(*s)
}

func (s *stringListValue) Set(val string) error {
	*s = strings.Split(val, ",")
	return nil
}

func (s *stringListValue) String() string {
	return strings.Join(*s, ",")
}

func StringListVar(p *[]string, name string, val []string, usage string) {
	flag.Var(newStringListValue(val, p), name, usage)
}
