package main

import (
	"flag"
	"os"
	"github.com/xtls/xray-core/main/commands/base"
	_ "github.com/xtls/xray-core/main/distro/all"
	"github.com/kardianos/service"
	"log"
)

// program 结构体将实现 service.Interface
type program struct{}

// Start 方法将在服务启动时调用
func (p *program) Start(s service.Service) error {
	// Start 应该是非阻塞的，所以启动一个goroutine
	go p.run()
	return nil
}

// run 包含原始的 Xray 启动逻辑
func (p *program) run() {
	// 重新设置命令行参数，以确保与原有逻辑兼容
	os.Args = getArgsV4Compatible()
	
	base.RootCommand.Long = "Xray is a platform for building proxies."
	base.RootCommand.Commands = append(
		[]*base.Command{
			{
				Use:   "run",
				Short: "Run Xray",
				Long:  "Run Xray core service.",
			},
			{
				Use:   "version",
				Short: "Version info",
				Long:  "Show version information.",
			},
		},
		base.RootCommand.Commands...,
	)
	// 执行 Xray 的命令行解析和命令执行
	base.Execute()
}

// Stop 方法将在服务停止时调用
func (p *program) Stop(s service.Service) error {
	// 执行任何需要的清理操作
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

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		// 处理服务控制命令
		cmd := os.Args[1]
		if err = service.Control(s, cmd); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err = s.Run(); err != nil {
		logger.Error(err)
	}
}

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
