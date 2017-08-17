package internal

import (
	"errors"
	"fmt"
	"os"
)

func RunCheckMain(checkFunc CheckFunc) error {
	err := RunCheck(os.Stdin, os.Stdout, checkFunc)
	if err != nil {
		return fmt.Errorf("error processing check request: %v", err)
	}
	return nil
}

func RunInMain(inFunc InFunc) error {
	if len(os.Args) < 2 {
		return errors.New("in script requires a target directory argument")
	}
	err := RunIn(os.Stdin, os.Stdout, os.Args[1], inFunc)
	if err != nil {
		return fmt.Errorf("error processing in request: %v", err)
	}
	return nil
}

func RunOutMain(outFunc OutFunc) error {
	if len(os.Args) < 2 {
		return errors.New("out script requires a target directory argument")
	}
	err := RunOut(os.Stdin, os.Stdout, os.Args[1], outFunc)
	if err != nil {
		return fmt.Errorf("error processing out request: %v", err)
	}
	return nil
}

type MainRunner struct {
	checkFunc CheckFunc
	inFunc    InFunc
	outFunc   OutFunc
}

func (r *MainRunner) SetCheckFunc(checkFunc CheckFunc) {
	r.checkFunc = checkFunc
}

func (r *MainRunner) SetInFunc(inFunc InFunc) {
	r.inFunc = inFunc
}

func (r *MainRunner) SetOutFunc(outFunc OutFunc) {
	r.outFunc = outFunc
}

func (r MainRunner) RunMain() error {
	switch os.Args[0] {
	case "check":
		if r.checkFunc == nil {
			return errors.New("no CheckFunc set")
		}
		return RunCheckMain(r.checkFunc)
	case "in":
		if r.checkFunc == nil {
			return errors.New("no InFunc set")
		}
		return RunInMain(r.inFunc)
	case "out":
		if r.checkFunc == nil {
			errors.New("no OutFunc set")
		}
		return RunOutMain(r.outFunc)
	default:
		return errors.New("RunMain: os.Args[0] must be one of 'check', 'in', 'out'")
	}
}

var defaultMainRunner = &MainRunner{}

func RegisterCheckFunc(checkFunc CheckFunc) {
	defaultMainRunner.SetCheckFunc(checkFunc)
}

func RegisterInFunc(inFunc InFunc) {
	defaultMainRunner.SetInFunc(inFunc)
}

func RegisterOutFunc(outFunc OutFunc) {
	defaultMainRunner.SetOutFunc(outFunc)
}

func RunMain() error {
	return defaultMainRunner.RunMain()
}

func RunMainExit() {
	err := RunMain()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
