package main

type command interface {
	addFlags()
	printHelp(string)
	run(args []string)
}
