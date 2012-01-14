package ircclient

// This plugin manages a common config-file pointer
// and the locks on it.

import (
	"sync"
	"github.com/kless/goconfig/config"
	"log"
)

type ConfigPlugin struct {
	ic       *IRCClient
	filename string
	Conf     *config.Config
	// Operations to the Config structure should be atomic
	lock sync.Mutex
}

func NewConfigPlugin(filename string) *ConfigPlugin {
	c, ok := config.ReadDefault(filename)
	if ok != nil {
		c = config.NewDefault()
		c.AddSection("Server")
		c.AddOption("Server", "host", "dpaulus.dyndns.org:6667")
		c.AddOption("Server", "nick", "testbot")
		c.AddOption("Server", "ident", "ident")
		c.AddOption("Server", "realname", "TestBot Client")
		c.AddOption("Server", "trigger", ".")
		c.AddSection("Auth")
		c.WriteFile(filename, 0644, "go-faui2k11 default config file")
		log.Println("Note: A new default configuration file has been generated in go-faui2k11.cfg. Please edit it to suit your needs and restart go-faui2k11 then")
		return nil
	}
	return &ConfigPlugin{filename: filename, Conf: c}
}
func (cp *ConfigPlugin) Register(cl *IRCClient) {
	cp.ic = cl
	cl.RegisterCommandHandler("version", 0, 0, cp)
	cl.RegisterCommandHandler("source", 0, 0, cp)
	cl.RegisterCommandHandler("writeconf", 0, 400, cp)
}
func (cp *ConfigPlugin) String() string {
	return "conf"
}
func (cp *ConfigPlugin) ProcessLine(msg *IRCMessage) {
	// Empty
}
func (cp *ConfigPlugin) Unregister() {
	cp.lock.Lock()
	cp.Conf.WriteFile("go-faui2k11.cfg", 0644, "go-faui2k11 config")
	cp.lock.Unlock()
}
func (cp *ConfigPlugin) Info() string {
	return "run-time configuration manager plugin"
}
func (cp *ConfigPlugin) ProcessCommand(cmd *IRCCommand) {
	switch cmd.Command {
	case "version":
		cp.ic.Reply(cmd, "This is go-faui2k11, version 0.01a")
	case "writeconf":
		cp.lock.Lock()
		cp.Conf.WriteFile("go-faui2k11.cfg", 0644, "go-faui2k11 config")
		cp.Conf, _ = config.ReadDefault(cp.filename)
		cp.lock.Unlock()
		cp.ic.Reply(cmd, "Successfully flushed cached config entries")
	case "source":
		cp.ic.Reply(cmd, "source is (atm) at http://bitbucket.org/dpaulus/go-faui2k11")
	}
}

func (cp *ConfigPlugin) Lock() {
	cp.lock.Lock()
}
func (cp *ConfigPlugin) Unlock() {
	cp.lock.Unlock()
}
