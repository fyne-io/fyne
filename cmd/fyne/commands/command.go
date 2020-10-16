package commands

// Command defines the required functionality to provide a subcommand to the "fyne" tool.
type Command interface {
	AddFlags()
	PrintHelp(string)
	Run(args []string)
}
