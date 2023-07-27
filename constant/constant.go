package constant

type SshConfig struct {
	SshHost     string         `json:"sshHost"`
	SshUser     string         `json:"sshUser"`
	SshStatus   int            `json:"sshStatus"` // 0 表示不开启、1表示开启监控
	SshPassword string         `json:"sshPassword"`
	SshType     string         `json:"sshType"` //password 或者 key,local 本地执行不需要登录
	SshKeyPath  string         `json:"sshKeyPath"`
	SshPort     int            `json:"sshPort"`
	Monitors    []ServerConfig `json:"monitors"`
}

type ServerConfig struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  int    `json:"status"`  // 0 表示不开启、1表示开启监控
	FindCmd string `json:"findCmd"` // 本地运行java暂时至支持一个连接
	Keyword string `json:"keyword"`
	StarCmd string `json:"starCmd"`
	FileDir string `json:"fileDir"`
}

type CmdInfo struct {
	Cmd     string
	Arg     string
	FileDir string
}

var Nohup = "nohup"
