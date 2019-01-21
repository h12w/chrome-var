package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Command struct {
	URI string `long:"uri"
		description:"URI for the browser to visit to grab the variable value"`
	Var string `long:"var"
		description:"JS variable name"`
}

func main() {
	cmd := &Command{}
	rootParser := flags.NewParser(cmd, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)
	args, err := rootParser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Println("chrome-var: cli to a get JS variable value from Chrome browser\n")
			rootParser.WriteHelp(os.Stderr)
			return
		}
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if err := cmd.Execute(args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func (c *Command) Execute([]string) error {
	value, err := getVarFromChrome(c.URI, c.Var)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}
