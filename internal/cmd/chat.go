package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"toolbox/internal/chat"
)

const welcome = `使用API Key与chatGPT聊天，按两次回车发送消息。
为防止帐号封禁，【请确认网络连接方式】，输入Yes后回车开始...`

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat with chatGPT",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(welcome)

		var str string
		reader := bufio.NewReader(os.Stdin)
		str, _ = reader.ReadString('\n')
		str = strings.TrimSpace(str)

		if strings.EqualFold(str, "Yes") {
			chat.Do()
		}
	},
}

func init() {
	initFuncList = append(initFuncList, initChat)
}

func initChat() {
	rootCmd.AddCommand(chatCmd)
}
