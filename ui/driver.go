package ui

var driver Driver

type Driver interface {
	CreateWindow(string) Window
	Run()
}
