package main

import (
	"flag"
	"os"
	"log"
	"github.com/xtls/xray-core/main/commands/base"
	_ "github.com/xtls/xray-core/main/distro/all"
	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	os.Args = getArgsV4Compatible()
	base.RootCommand.Long = "Xray is a platform for building proxies."
	base.RootCommand.Commands = append(
		[]*base.Command{
			{
				UsageLine: "run",
				Short:     "run the service",
				Long:      "Run the Xray proxy service.",
			},
			{
				UsageLine: "version",
				Short:     "show version",
				Long:      "Show the version of the Xray service.",
			},
		},
		base.RootCommand.Commands...,
	)
	base.Execute()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "XrayService",
		DisplayName: "Xray Proxy Service",
		Description: "This service manages the Xray proxy platform.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "install" || cmd == "uninstall" || cmd == "start" || cmd == "stop" {
			err = service.Control

func getArgsV4Compatible() []string {
	if len(os.Args) == 1 {
		return []string{os.Args[0], "run"}
	}
	if os.Args[1][0] != '-' {
		return os.Args
	}
	version := false
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&version, "version", false, "")
	// parse silently, no usage, no error output
	fs.Usage = func() {}
	fs.SetOutput(&null{})
	err := fs.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		// fmt.Println("DEPRECATED: -h, WILL BE REMOVED IN V5.")
		// fmt.Println("PLEASE USE: xray help")
		// fmt.Println()
		return []string{os.Args[0], "help"}
	}
	if version {
		// fmt.Println("DEPRECATED: -version, WILL BE REMOVED IN V5.")
		// fmt.Println("PLEASE USE: xray version")
		// fmt.Println()
		return []string{os.Args[0], "version"}
	}
	// fmt.Println("COMPATIBLE MODE, DEPRECATED.")
	// fmt.Println("PLEASE USE: xray run [arguments] INSTEAD.")
	// fmt.Println()
	return append([]string{os.Args[0], "run"}, os.Args[1:]...)
}

type null struct{}

func (n *null) Write(p []byte) (int, error) {
	return len(p), nil
}
