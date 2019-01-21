package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-vgo/robotgo"
)

func getVarFromChrome(uri, jsVar string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
	default:
		panic("only Mac supported for now")
	}

	chromeProcessName := "Google Chrome"
	// activate or open chrome browser
	fpid, err := findProcess(chromeProcessName)
	if err != nil {
		if err := openChrome(); err != nil {
			return "", err
		}
		time.Sleep(time.Second)
		fpid, err = findProcess(chromeProcessName)
		if err != nil {
			return "", err
		}
	}
	activeWin := robotgo.GetActive()
	defer robotgo.SetActive(activeWin)
	if err := robotgo.ActivePID(fpid); err != nil {
		return "", err
	}
	time.Sleep(time.Second)

	robotgo.KeyTap("t", "command")       // create new tab
	defer robotgo.KeyTap("w", "command") // defer close the tab
	robotgo.PasteStr(uri)                // open dev site
	robotgo.KeyTap("enter")
	time.Sleep(time.Second)
	robotgo.KeyTap("j", "alt", "command") // open dev console
	time.Sleep(2 * time.Second)
	copyStatement := fmt.Sprintf("copy(%s)", jsVar)
	robotgo.PasteStr(copyStatement) // copy jsVar
	robotgo.KeyTap("enter")
	time.Sleep(2 * time.Second)
	jsValue, err := robotgo.ReadAll() // read service token
	if err != nil {
		return "", err
	}
	switch jsValue {
	case "", "undefined", copyStatement:
		return "", fmt.Errorf("empty %s from %s", jsVar, uri)
	}
	return jsValue, nil
}

func openChrome() error {
	return exec.Command("open", "/Applications/Google Chrome.app").Start()
}

func findProcess(name string) (int32, error) {
	processes, err := robotgo.Process()
	if err != nil {
		return 0, err
	}
	for _, process := range processes {
		if process.Name == name {
			return process.Pid, nil
		}
	}
	return 0, errors.New("process not found")
}
