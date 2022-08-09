package apps

import (
	"flag"
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
	"strings"
)

var (
	H            = flag.Bool("help", false, "显示帮助信息")
	S            = flag.Bool("s", false, "载入ssh配置 config '~/.ssh/config'")
	C            = flag.String("c", configName, "服务器配置文件名")
	configName   = "go_ssh.yaml"
	logs         = GetLogger()
	ClearContent = "\033[K" // 清除从光标到行尾的内容
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
	node := choose(trees)
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

// getServers 将服务器信息打印出来
func getServers(trees []*Node, i int) []string {
	var content []string
	noResult := true
	for index, item := range trees {
		if item == nil {
			content = append(content, ClearContent)
		} else if index == i {
			noResult = false
			content = append(content, item.Str(true))
		} else {
			noResult = false
			content = append(content, item.Str(false))
		}
	}
	if noResult && len(trees) != 0 {
		// 说明搜索没搜到任何服务器
		content[1] = "  🍵 未找匹配到数据" + ClearContent
	}
	return content
}

type ServerInfo struct {
	CurrentIndex  int
	Nodes         []*Node
	nodes         []*Node
	SearchContent string
	searchContent string
	Length        int
	height        int // 内容的高度
}

// \033[0m 关闭所有属性
//\033[1m 设置高亮度
//\033[4m 下划线
// \033[5m 闪烁
//\033[7m 反显
//\033[8m 消隐
//\033[30m — \033[37m 设置前景色
//\033[40m — \033[47m 设置背景色
//\033[nA 光标上移n行
//\033[nB 光标下移n行
//\033[nC 光标右移n行
//\033[nD 光标左移n行
//\033[y;xH设置光标位置
//\033[2J 清屏
//\033[K 清除从光标到行尾的内容
//\033[s 保存光标位置
//\033[u 恢复光标位置
//\033[?25l 隐藏光标
//\033[?25h 显示光标
// HideCursor 隐藏光标
func HideCursor() {
	fmt.Printf("\033[?25l")
}

func ShowCursor() {
	fmt.Printf("\033[?25h")
}

func MoveCursorUP(y int) {
	// 相对位置移动，而不是按照整个屏幕定位
	// \033[nA 光标上移n行
	// \033[nC 光标右移n行
	fmt.Printf("\033[%dA", y)
}
func MoveCursorRight(x int) {
	// 往右移动光标
	fmt.Printf("\033[%dC", x)
}
func SaveCursor() {
	fmt.Printf("\033[s")
}

func RecoveryCursor() {
	fmt.Printf("\033[u")

}
func (s *ServerInfo) getTips() []string {
	// 根据搜索内容匹配服务器信息
	if len(s.SearchContent) != 0 && s.searchContent != s.SearchContent {
		var nodes []*Node
		for _, node := range s.nodes {
			if strings.Contains(node.Name, s.SearchContent) || strings.Contains(node.Host, s.SearchContent) || strings.Contains(node.User, s.SearchContent) {
				nodes = append(nodes, node)
			}
		}

		s.Length = len(nodes) - 1
		q := len(s.nodes) - len(nodes)
		for i := 0; i < q; i++ {
			nodes = append(nodes, nil)
		}
		s.Nodes = nodes
		s.CurrentIndex = 0
	} else if len(s.SearchContent) == 0 && s.searchContent != s.SearchContent {
		s.Nodes = s.nodes
		s.CurrentIndex = 0
		s.Length = len(s.Nodes) - 1

	}
	s.searchContent = s.SearchContent
	// 根据输入的内容计算光标移动的偏差

	return []string{fmt.Sprintf("🔍 输入自动搜索：%s"+ClearContent, s.SearchContent), Green("✨ 请选择要连接的服务器：")}
}

func (s *ServerInfo) getContent() []string {
	// 获取本次要打印的内容
	var content []string
	content = append(content, s.getTips()...)
	content = append(content, getServers(s.Nodes, s.CurrentIndex)...)
	return content
}

func (s *ServerInfo) Draw() {
	content := s.getContent()
	height := len(content)
	if height > s.height {
		s.height = height
	}
	RecoveryCursor()
	for _, s := range content {
		fmt.Println(s)
	}
	// 计算第一列输入的提示语句和已输入的内容的长度
	MoveCursorUP(s.height)
	MoveCursorRight(ZhLen(content[0]) - 1)
}

func NewServerInfo(trees []*Node) *ServerInfo {
	initLength(trees)
	return &ServerInfo{
		Nodes:  trees,
		nodes:  trees,
		Length: len(trees) - 1,
	}
}

func choose(trees []*Node) *Node {
	SaveCursor()
	serverInfo := NewServerInfo(trees)
	serverInfo.Draw()
	// 绘制之后，开始监听键盘
	node := serverInfo.HandleKeyboard()

	return node
}

// HandleKeyboard 处理键盘事件
func (s *ServerInfo) HandleKeyboard() *Node {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			panic(err)
		}
	}()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Rune != 0 {
			s.handleChar(event.Rune)
		} else if event.Key != 0 {
			ret := s.handleKey(event.Key)
			switch ret {
			case 1:
				return s.Nodes[s.CurrentIndex]
			case 2:
				return nil
			}
		}
		s.Draw()
	}
}

// 处理字母按键
func (s *ServerInfo) handleChar(char rune) {
	ch := string(char)
	s.SearchContent += ch
}

func (s *ServerInfo) deleteSearchContent() {
	searchContent := []rune(s.SearchContent)
	if len(s.SearchContent) != 0 {
		searchContent = searchContent[:len(searchContent)-1]
		s.SearchContent = string(searchContent)
		s.Draw()
	}
}
func (s *ServerInfo) clear() {
	for i := 0; i < s.height; i++ {
		fmt.Println(ClearContent)
	}
	MoveCursorUP(s.height)
	ShowCursor()
}

// 处理键盘除字母键以外的按键
func (s *ServerInfo) handleKey(key keyboard.Key) int {
	switch key {
	//case keyboard.KeyArrowRight, keyboard.KeyArrowLeft, keyboard.KeyArrowDown, keyboard.KeyArrowUp:
	case keyboard.KeyArrowUp:
		if s.CurrentIndex == 0 {
			s.CurrentIndex = s.Length
		} else {
			s.CurrentIndex--
		}
	case keyboard.KeyArrowDown:
		if s.CurrentIndex == s.Length {
			s.CurrentIndex = 0
		} else {
			s.CurrentIndex++
		}
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		s.deleteSearchContent()
	case keyboard.KeyEnter:
		s.clear()
		return 1
	case keyboard.KeyCtrlC:
		s.clear()
		return 2
	}
	return 0
}
