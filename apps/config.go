package apps

import (
	"fmt"
	"github.com/atrox/homedir"
	"github.com/kevinburke/ssh_config"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	Name          string  `yaml:"name"`
	Host          string  `yaml:"host"`
	User          string  `yaml:"user"`
	Port          int     `yaml:"port"`
	Key           string  `yaml:"key"`
	Passphrase    string  `yaml:"passphrase"`
	Password      string  `yaml:"password"`
	Jump          []*Node `yaml:"jump"`
	ID            string
	ChildrenCount int
	Method        string // 鉴权方式
	pinyin        string // 拼音
	firstPinyin   string // 首字母拼音
}

type CallbackShell struct {
	Cmd   string        `yaml:"cmd"`
	Delay time.Duration `yaml:"delay"`
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

func getPinyin(style *pinyin.Args, hans string) (string, string) {
	p := pinyin.Pinyin(hans, *style)
	pinyinStr := ""
	pinyinFirst := ""
	for _, item := range p {
		pinyinStr += item[0]
		pinyinFirst += string(item[0][0])
	}
	return pinyinStr, pinyinFirst
}

// 添加编号 处理数据
func HandleNode(c []*Node) []*Node {
	var temp []*Node
	for _, item := range c {
		if item.Host == "" {
			// 说明是单个服务器，必须要host
			continue
		}
		temp = append(temp, item)
	}
	pinyinStyle := pinyin.NewArgs()
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
		// 服务器名转pinyin 为了搜索
		//fmt.Println()
		//item.pinyin = pinyin.Pinyin(item.Name, pinyinStyle)
		item.pinyin, item.firstPinyin = getPinyin(&pinyinStyle, item.Name)
	}
	return temp
}

func ColorStr(s string, color int) string {
	// 给字符串上色
	return fmt.Sprintf("%s%dm%v%s", "\033[", color, s, fmt.Sprintf("%s%dm", "\033[", 0))
}

func Faint(s string) string {
	return ColorStr(s, 2)
}
func Yellow(s string) string {
	return ColorStr(s, 33)
}

func Green(s string) string {
	return ColorStr(s, 32)
}

func (n *Node) Str(selected bool, showDetail bool) string {
	s := "   %s. %s %s %s@%s%s\033[K"

	if selected {
		format := []interface{}{n.ID, n.Name, ">", n.User, n.Host}
		if showDetail {

			if n.Method == "密码" {
				format = append(format, fmt.Sprintf("   密码：%s  端口：%d", n.Password, n.Port))
			} else {
				format = append(format, fmt.Sprintf("   端口：%d", n.Port))
			}

		} else {
			format = append(format, "")
		}
		temp := fmt.Sprintf(s, format...)
		s = Yellow(temp)
	} else {
		s = Faint(fmt.Sprintf(s, n.ID, n.Name, "|", n.User, n.Host, ""))
	}
	return s
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
