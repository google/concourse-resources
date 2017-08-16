package internal

import (
	"log"
	"os"
)

func RunCheckMain(checkFunc CheckFunc) {
	err := RunCheck(os.Stdin, os.Stdout, checkFunc)
	if err != nil {
		log.Fatalf("error processing check request: %v", err)
	}
}

func RunInMain(inFunc InFunc) {
	if len(os.Args) < 2 {
		log.Fatalln("in script requires a target directory argument")
	}
	err := RunIn(os.Stdin, os.Stdout, os.Args[1], inFunc)
	if err != nil {
		log.Fatalf("error processing in request: %v", err)
	}
}

func RunOutMain(outFunc OutFunc) {
	if len(os.Args) < 2 {
		log.Fatalln("out script requires a target directory argument")
	}
	err := RunOut(os.Stdin, os.Stdout, os.Args[1], outFunc)
	if err != nil {
		log.Fatalf("error processing out request: %v", err)
	}
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

func (r MainRunner) RunMain() {
	switch os.Args[0] {
	case "check":
		if r.checkFunc == nil {
			log.Fatalln("no CheckFunc set")
		}
		RunCheckMain(r.checkFunc)
	case "in":
		if r.checkFunc == nil {
			log.Fatalln("no InFunc set")
		}
		RunInMain(r.inFunc)
	case "out":
		if r.checkFunc == nil {
			log.Fatalln("no OutFunc set")
		}
		RunOutMain(r.outFunc)
	default:
		log.Fatalln("RunMain: os.Args[0] must be one of 'check', 'in', 'out'")
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

func RunMain() {
	defaultMainRunner.RunMain()
}
