package apps

import (
	"fmt"
	"os"
	"unicode"
)

// 计算字符宽度（中文）
func ZhLen(str string) int {
	length := 0
	for _, c := range str {
		if unicode.Is(unicode.Scripts["Han"], c) {
			length += 2
		} else {
			length += 1
		}
	}

	return length
}

//右填充
func AppendRight(body string, char string, maxlength int) string {
	length := ZhLen(body)
	if length >= maxlength {
		return body
	}

	for i := 0; i < maxlength-length; i++ {
		body = body + char
	}

	return body
}

//左填充
func AppendLeft(body string, char string, maxlength int) string {
	length := ZhLen(body)
	if length >= maxlength {
		return body
	}

	for i := 0; i < maxlength-length; i++ {
		body = char + body
	}

	return body
}

// 左右填充
// title 主体内容
// c 填充符号
// maxlength 总长度
// 如： title = 测试 c=* maxlength = 10 返回 ** 返回 **
func FormatSeparator(title string, c string, maxlength int) string {
	charslen := (maxlength - ZhLen(title)) / 2
	chars := ""
	for i := 0; i < charslen; i++ {
		chars += c
	}
	return chars + title + chars
}

// 没有配置文件时，初始化一份示例文件
func InitConfig(p string) {
	example := "# 欢迎使用Go SSH 工具\n# 字段说明\n# name         ： 自定义的服务器名字 可不填\n# user         ： 服务器名 不填默认 root\n# host         ： 服务器域名或ip  ！！！必填！！！  不填的话，列表不会出现这条配置\n# port         ： 端口号  不填默认 22\n# password     ： 密码    不填默认用秘钥\n# key          ： 私钥    不填默认  ~/.ssh/id_rsa\n# passphrase   ： 私钥的密码  可不填\n# jump         ： 跳板机 配置同上\n\n- { user: appuser, port: 22, password: 123456 }\n- { name: 使用秘钥, user: appuser, host: 192.168.8.35, port: 22}\n- { name: 使用带密码的秘钥, user: appuser, host: 192.168.8.35, port: 22, passphrase: abcdefghijklmn}\n- { name: 不配置端口, user: appuser, host: 192.168.8.35 }\n- { name: 不配置用户名, host: 192.168.8.35 }\n- { name: 不设置密码, host: 192.168.8.35 }\n- { name: ⚡️ 带emoji表情, host: 192.168.8.35 }\n- name: 通过跳板机连接\n  user: appuser\n  host: 192.168.8.35\n  port: 22\n  password: 123456\n  jump:\n    - user: appuser\n      host: 192.168.8.36\n      port: 2222\n"
	dstFile, err := os.Create(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(example)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
