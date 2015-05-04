package shell

/**

Provides the ability to list users

Copyright 2015 - Joseph Lewis <joseph@josephlewis.net>
                 Daniel Kumor <rdkumor@gmail.com>

All Rights Reserved

**/

import (
	"fmt"
)

// The clear command
type Su struct {
}

func (h Su) Help() string {
	return "Changes priviliges to a particular user: 'su username'"
}

func (h Su) Usage() string {
	return ""
}

func (h Su) Execute(shell *Shell, args []string) {
	if len(args) < 2 {
		fmt.Println(Red + "Must supply a name" + Reset)
		return
	}

	user, err := shell.sdb.ReadByNameOrEmail(args[1], args[1])
	if shell.PrintError(err) {
		return
	}

	userdevice, err := shell.sdb.ReadUserOperatingDevice(user)
	if shell.PrintError(err) {
		return
	}

	suOperator, err := shell.sdb.GetOperatorForDevice(userdevice)
	if shell.PrintError(err) {
		return
	}

	sushell := CreateShell(shell.sdb)
	sushell.operator = suOperator
	sushell.operatorName = user.Name

	sushell.Repl()
}

func (h Su) Name() string {
	return "su"
}
