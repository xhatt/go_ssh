package apps

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
)

const prev = "[↑] 返回上一页"

var (
	H          = flag.Bool("help", false, "显示帮助信息")
	S          = flag.Bool("s", false, "载入ssh配置 config '~/.ssh/config'")
	C          = flag.String("c", configName, "服务器配置文件名")
	configName = "go_ssh.yaml"
	detailLen  = 12
	logs       = GetLogger()

	cursor int
	keys   = &promptui.SelectKeys{
		Prev:     promptui.Key{Code: promptui.KeyPrev, Display: promptui.KeyPrevDisplay},
		Next:     promptui.Key{Code: promptui.KeyNext, Display: promptui.KeyNextDisplay},
		PageUp:   promptui.Key{Code: promptui.KeyBackward, Display: promptui.KeyBackwardDisplay},
		PageDown: promptui.Key{Code: promptui.KeyForward, Display: promptui.KeyForwardDisplay},
		//Search:   promptui.Key{Code: readline.CharEsc, Display: "Esc"},
		//Search:   promptui.Key{Code: promptui.KeyBackspace, Display: "Esc"},
	}
)

func Run() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		return
	}

	if *H {
		flag.Usage()
		return
	}

	if *S {
		err := LoadSshConfig()
		if err != nil {
			logs.Error("load ssh config error", err)
			os.Exit(1)
		}
	}
	if *C != "" {
		err := LoadConfig(*C)
		if err != nil {
			logs.Error("load config error", err)
			os.Exit(1)
		}
	}

	trees := GetConfig()
	if trees == nil {
		// 说明没有
		fmt.Println("没有任何服务器")
		os.Exit(0)
	}
	node := choose(nil, trees, 0)
	if node == nil {
		return
	}
	fmt.Println("正在连接。。。")
	client := NewClient(node)
	client.Login()
}

// 获取当前页的长度
func initLength(trees []*Node) {
	DomainLen = 0
	NameLen = 0
	MaxLen = 0
	IDLen = 0
	for _, item := range trees {
		if _nameLen := ZhLen(item.Name); _nameLen > NameLen {
			NameLen = _nameLen
			MaxLen = NameLen
		}
		if _domainLen := ZhLen(fmt.Sprintf("%s@%s", item.user(), item.Host)); _domainLen > DomainLen {
			DomainLen = _domainLen
		}
		if DomainLen > MaxLen {
			MaxLen = DomainLen
		}
		if _idLen := ZhLen(item.ID); _idLen > IDLen {
			IDLen = _idLen
		}
	}
	for _, item := range trees {
		if item.ID == "" {
			continue
		} else {
			item.ID = AppendLeft(item.ID, " ", IDLen)
		}
		if item.Host != "" {
			item.Host = AppendRight(item.Host, " ", DomainLen-ZhLen(item.Name))
		}
		if item.ChildrenCount != 0 {
			item.Name = AppendRight(item.Name, " ", NameLen)
		} else {
			item.Name = AppendRight(item.Name, " ", NameLen+4)
		}
	}
}

func getTemplates() *promptui.SelectTemplates {
	templates := &promptui.SelectTemplates{
		Label:    "✨ {{ . | green}}",
		Active:   "{{`➤ ` | yellow }}{{if .ID}}{{ .ID | yellow  }}{{`.`|yellow}} {{end}}{{if .ChildrenCount}}{{`[+] `| yellow  }}{{.Name | yellow}}{{` > `| yellow}}{{.ChildrenCount|yellow}}{{`个服务器`|yellow}}{{else}}{{ .Name | yellow  }}{{if .Host}}{{` > `|yellow}}{{if .User}}{{.User | yellow}}{{else}}{{`root` | yellow}}{{end}}{{`@` | yellow}}{{.Host | yellow}}{{end}}{{end}}",
		Inactive: "  {{if .ID}}{{ .ID | faint  }}{{`.`|faint}} {{end}}{{if .ChildrenCount}}{{`[+] `| faint  }}{{.Name | faint}}{{` | `}}{{.ChildrenCount|faint}}{{`个服务器`|faint}}{{else}}{{ .Name | faint  }}{{if .Host}}{{` | `}}{{if .User}}{{.User | faint}}{{else}}{{`root` | faint}}{{end}}{{`@` | faint}}{{.Host | faint}}{{end}}{{end}}",
		//Details: fmt.Sprintf("{{if .ID}}%s\n{{if not .ChildrenCount}}"+
		//	"{{`%s`|faint}}{{ .Name | yellow}}\n{{`%s`|faint}}"+
		//	"{{if .User}}{{.User | yellow}}\n{{else}}{{`root` | yellow}}\n"+
		//	"{{end}}{{`%s`|faint}}{{ .Host | yellow}}\n{{`%s`|faint}}{{ .Method | yellow}}\n{{`%s`|faint}}{{ .Port | yellow}}"+
		//	"{{else}}{{`%s` | faint}}{{ .Name|yellow}}\n{{range .Children}}    {{ .ID }}. {{.Name}}\n{{end}} {{end}}{{end}}",
		//	FormatSeparator("详细信息", "-", detailLen+MaxLen),
		//	AppendLeft("服务器名称:", " ", detailLen),
		//	AppendLeft("用户名:", " ", detailLen),
		//	AppendLeft("域名或IP:", " ", detailLen),
		//	AppendLeft("鉴权方式:", " ", detailLen),
		//	AppendLeft("端口:", " ", detailLen),
		//	AppendLeft("分组名称:", " ", detailLen),
		//),
	}
	return templates
}

func choose(parent, trees []*Node, i int) *Node {

	initLength(trees)
	prompt := promptui.Select{
		Label:        "请选择要连接的服务器：",
		Items:        trees,
		Templates:    getTemplates(),
		Size:         10,
		HideSelected: true,
		Keys:         keys,
		Searcher: func(input string, index int) bool {
			node := trees[index]
			input = strings.ToLower(input)
			content := strings.ToLower(fmt.Sprintf("%s %s %s %s", node.ID, node.Name, node.User, node.Host))
			if strings.Contains(input, " ") {
				for _, key := range strings.Split(input, " ") {
					key = strings.TrimSpace(key)
					if key != "" {
						if !strings.Contains(content, key) {
							return false
						}
					}
				}
				return true
			}
			if strings.Contains(content, input) {
				return true
			}
			return false
		},
	}
	index, _, err := prompt.RunCursorAt(i, 0)

	if err != nil {
		return nil
	}

	node := trees[index]
	if node.ID == "" {
		// 选择了返回上层，删掉这个节点
		if parent == nil {
			return choose(nil, GetConfig(), 0)
		}
		_node := parent[cursor]
		_node.Children = append(_node.Children[:index], _node.Children[index+1:]...)
		return choose(nil, parent, cursor)
	}
	if len(node.Children) > 0 {
		first := node.Children[0]
		if first.Name != prev {
			first = &Node{Name: prev}
			node.Children = append(node.Children[:0], append([]*Node{first}, node.Children...)...)
		}
		cursor = index
		return choose(trees, node.Children, 0)
	}

	return node
}
