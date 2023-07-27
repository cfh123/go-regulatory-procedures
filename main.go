package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"go-regulatory-procedures/constant"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var fileDir = "/host.json"
	// 打开文件
	resultInit, _ := ReadFile(fileDir)
	var configArray []constant.SshConfig
	// json 结构体数组
	if err := json.Unmarshal([]byte(resultInit), &configArray); err != nil {
		log.Printf("json 转换失败 :%v\n", err)
	}

	//每1Min一次刷新状态
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 执行定时任务
			fmt.Printf("\n-----------Begin--------------\n")

			// 循环配置
			for _, sshConfig := range configArray {
				if sshConfig.SshStatus == 1 {
					fmt.Printf("开始执行服务器监控配置 : %v  \n", sshConfig.SshHost)
					switch sshConfig.SshType {
					case "password":
						// 账号密码登录远程ssh
						PasswordSsh(&sshConfig)
						break
					case "key":
						// 密钥对登录远程ssh
						PasswordSsh(&sshConfig)
						break
					case "local":
						// 本地执行
						LocalExec(&sshConfig)
						break
					}
				}
			}

			// top 命令
			//CmdTop()

			// todo 其他服务

			fmt.Printf("-----------end--------------\n")
		}
	}
}

// PasswordSsh 用于登录ssh
func PasswordSsh(hostConfig *constant.SshConfig) {
	//创建ssh登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            hostConfig.SshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if hostConfig.SshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(hostConfig.SshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(hostConfig.SshKeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", hostConfig.SshHost, hostConfig.SshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	// 循环监控具体程序命令
	for _, ServerItem := range hostConfig.Monitors {
		if ServerItem.Status == 1 {
			//创建ssh-session
			session, err := sshClient.NewSession()
			if err != nil {
				log.Fatal("创建ssh session 失败", err)
			}
			defer session.Close()

			//执行远程命令
			combo, err := session.CombinedOutput(ServerItem.FindCmd)
			if err != nil {
				log.Fatal("远程执行cmd 失败", err)
			}
			fmt.Println(string(combo))
			isStart := strToArr(string(combo), ServerItem.Keyword)
			fmt.Printf("this %s need restart :%v\n", ServerItem.Name, isStart)
			if isStart {
				//创建ssh-session
				session, err := sshClient.NewSession()
				if err != nil {
					log.Fatal("创建ssh session 失败", err)
				}
				defer session.Close()

				session.CombinedOutput(ServerItem.StarCmd)
				log.Println("重启命令输出:", string(combo))

				// 发送邮箱
				constant.SendMailSelf(ServerItem.StarCmd)
			}
		}
	}
}

// LocalExec 本地执行
func LocalExec(hostConfig *constant.SshConfig) {
	// 循环监控具体程序命令
	for _, ServerItem := range hostConfig.Monitors {
		if ServerItem.Status == 1 {
			split := strings.Split(ServerItem.FindCmd, "|")
			var javaItem []constant.CmdInfo
			for _, itemStr := range split {
				if strings.TrimSpace(itemStr) != "" {
					itemItemArr := strings.Split(itemStr, " ")
					javaItem = append(javaItem, constant.CmdInfo{Cmd: itemItemArr[0], Arg: itemItemArr[1]})
				}
			}
			// 组装命令服务
			execFor := CmdExecFor(&javaItem[0], &javaItem[1], ServerItem.Keyword)
			fmt.Printf("this %s need restart :%v\n", ServerItem.Name, execFor)
			if execFor {
				if ServerItem.Type == "java" {
					CmdJava(ServerItem.FileDir)
				} else if ServerItem.Type == "docker" {
					all := strings.ReplaceAll(ServerItem.StarCmd, "docker restart ", "")
					cmd := exec.Command("docker", "restart", strings.TrimSpace(all))
					if err := cmd.Start(); err != nil {
						fmt.Println("Failed to start docker server:", err)
					}
				}
			}
		}
	}
}

// ReadFile 文件读取
func ReadFile(fileName string) (string, error) {

	var result = ""
	pwd, _ := os.Getwd()
	//fmt.Println("pwd:", pwd)

	file, err := os.Open(pwd + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return result, err
	}
	defer file.Close()

	// 创建 scanner 读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 逐行输出文件内容
		//fmt.Println(scanner.Text())
		result = result + scanner.Text()
	}

	// 检查错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return result, err
	}
	return result, nil
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func CmdTop() {
	// top 命令
	cmd := exec.Command("/bin/bash", "-c", "top -bn1")
	out, e := cmd.Output()
	if e != nil {
		log.Printf("failed due to :%v\n", e)
	}
	fmt.Println(string(out))

	strToArr(string(out), "")
}

func CmdJava(javaFile string) {
	_, err := os.Stat(javaFile)
	if err == nil {
		// 文件存在
	} else if os.IsNotExist(err) {
		// 文件不存在
		fmt.Println(javaFile+"file IsNotExist:", err)
		return
	} else {
		// 其他错误
	}
	// 确定执行
	cmd := exec.Command(constant.Nohup, "java", "-jar", javaFile, ">/dev/null", "2>&1", "&")
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start jar file:", err)
	}
}

// CmdExecFor false ，标识不需要重启 true 标识需要重启
func CmdExecFor(thisExec *constant.CmdInfo, nextExec *constant.CmdInfo, runStr string) bool {
	var output bytes.Buffer

	cmd1 := exec.Command(thisExec.Cmd, thisExec.Arg)
	cmd2 := exec.Command(nextExec.Cmd, nextExec.Arg)

	//cmd1 := exec.Command("netstat", "-ano")
	//cmd2 := exec.Command("findstr", "9099")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r
	cmd2.Stdout = &output

	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	w.Close()
	cmd2.Wait()

	fmt.Println(output.String())

	return strToArr(output.String(), runStr)
}

// strToArr false ，标识不需要重启 true 标识需要重启
func strToArr(outStr string, runStr string) bool {
	arr := strings.Split(outStr, "\n")
	for _, item := range arr {
		// 判断字符串是否非空 只有一种情况下返回false ，标识不需要重启
		if strings.TrimSpace(item) != "" && strings.Contains(item, runStr) {
			//fmt.Println("find process is : " + item)
			return false
		}
	}
	return true
}
