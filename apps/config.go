package apps

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/atrox/homedir"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type Node struct {
	Name           string           `yaml:"name"`
	Host           string           `yaml:"host"`
	User           string           `yaml:"user"`
	Port           int              `yaml:"port"`
	Key            string           `yaml:"key"`
	Passphrase     string           `yaml:"passphrase"`
	Password       string           `yaml:"password"`
	CallbackShells []*CallbackShell `yaml:"callback-shells"`
	Children       []*Node          `yaml:"children"`
	Jump           []*Node          `yaml:"jump"`
	ID             string
	ChildrenCount  int
	Method         string // 鉴权方式
}

type CallbackShell struct {
	Cmd   string        `yaml:"cmd"`
	Delay time.Duration `yaml:"delay"`
}

func (n *Node) String() string {
	return n.Name
}

func (n *Node) host() string {
	return strings.Replace(n.Host, " ", "", -1)
}

func (n *Node) user() string {
	if n.User == "" {
		return "root"
	}
	return strings.Replace(n.User, " ", "", -1)
}

func (n *Node) port() int {
	if n.Port <= 0 {
		return 22
	}
	return n.Port
}

func (n *Node) password() ssh.AuthMethod {
	if n.Password == "" {
		return nil
	}
	return ssh.Password(n.Password)
}

var (
	config    []*Node
	NameLen   int
	DomainLen int
	MaxLen    int
	IDLen     int
)

func GetConfig() []*Node {
	return config
}

// 添加编号 处理数据
func HandleNode(c []*Node) []*Node {
	var temp []*Node
	for _, item := range c {
		if item.Host == "" && item.Children == nil {
			// 说明是单个服务器，必须要host
			continue
		} else if item.Name == "" && item.Children != nil {
			// 说明是分组，必须有名字
			continue
		}
		temp = append(temp, item)
	}

	for index, item := range temp {
		item.ID = fmt.Sprintf("%d", index+1)
		if item.Password != "" {
			item.Method = "密码"
		}
		if item.Password == "" {
			item.Method += "秘钥"
		}
		if item.Key != "" && item.Password != "" {
			item.Method += "密码和秘钥"
		}
		if item.Port == 0 {
			item.Port = 22
		}
		if item.Name == "" {
			item.Name = item.Host
		}
		if item.User == "" {
			item.User = "root"
		}
		if item.Children != nil {
			item.ChildrenCount = len(item.Children)
			//item.Name = fmt.Sprintf("[+] %s", item.Name)
			HandleNode(item.Children)
		}
	}
	return temp
}

func LoadConfig(configName string) error {
	b, err := LoadConfigBytes(configName)
	if err != nil {
		return err
	}
	var c []*Node
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	config = append(config, HandleNode(c)...)

	return nil
}

func LoadSshConfig() error {
	u, err := user.Current()
	if err != nil {
		l.Error(err)
		return nil
	}
	f, _ := os.Open(path.Join(u.HomeDir, ".ssh/config"))
	cfg, _ := ssh_config.Decode(f)
	var nc []*Node
	for _, host := range cfg.Hosts {
		alias := fmt.Sprintf("%s", host.Patterns[0])
		hostName, err := cfg.Get(alias, "HostName")
		if err != nil {
			return err
		}
		if hostName != "" {
			port, _ := cfg.Get(alias, "Port")
			if port == "" {
				port = "22"
			}
			var c = new(Node)
			c.Name = alias
			c.Host = hostName
			c.User, _ = cfg.Get(alias, "User")
			c.Port, _ = strconv.Atoi(port)
			keyPath, _ := cfg.Get(alias, "IdentityFile")
			c.Key, _ = homedir.Expand(keyPath)
			nc = append(nc, c)
			// fmt.Println(c.Alias, c.Host, c.User, c.Port, c.Key)
		}
	}
	config = append(config, HandleNode(nc)...)
	return nil
}

func LoadConfigBytes(names string) ([]byte, error) {

	//currentPath, err := os.Getwd()
	currentPath, err := os.Executable()
	p := path.Join(filepath.Dir(currentPath), names)

	_, err = os.Stat(p)
	if os.IsNotExist(err) && names == configName {
		fmt.Println("未找到配置文件，已初始化了示例配置文件 ", p)
		InitConfig(p)
	} else if os.IsNotExist(err) {
		fmt.Println("未找到配置文件", p)
		os.Exit(1)
	}

	sshw, err := ioutil.ReadFile(p)
	if err == nil {
		return sshw, nil
	} else {

	}
	return nil, err
}
