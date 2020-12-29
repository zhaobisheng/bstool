package BSFtp

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SSHConnect(user, password, host, port string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 5 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%s", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func SftpUploadFile(sftpClient *sftp.Client, localFile, remotePath string) error {
	var err error
	// 用来测试的本地文件路径 和 远程机器上的文件夹
	index := strings.LastIndex(remotePath, "/")
	if index < len(remotePath)-1 {
		remotePath = remotePath + "/"
	}
	var localFilePath = localFile
	var remoteDir = remotePath
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	_, err1 := sftpClient.ReadDir(remoteDir)
	if err1 != nil {
		err2 := sftpClient.MkdirAll(remoteDir)
		if err2 != nil {
			return err2
		}
	}
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		len, err := dstFile.Write(buf)
		if len <= 0 && err != nil {
			return err
		}
	}
	fmt.Println("copy " + localFilePath + " to remote server finished!")
	return nil
}
