package main

type test string

func (*test) Hello() string {
	return "hello"
}

var Test test