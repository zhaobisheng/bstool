package BMail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/smtp"
	"path/filepath"
	"strconv"

	"github.com/bstool/Logger"

	"github.com/go-gomail/gomail"
)

//不使用SSL
func SendToMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	err := smtp.SendMail(addr, auth, from, to, msg)
	return err
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := SmtpConnect(addr)
	if err != nil {
		Logger.Errorln("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				Logger.Errorln("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

//return a smtp client
func SmtpConnect(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		Logger.Errorln("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func GetAttachMsg(body, email, toEmail, Subject, attaFile string) []byte {
	_, attaFileName := filepath.Split(attaFile)
	boundary := "----=bison" //boundary 用于分割邮件内容，可自定义. 注意它的开始和结束格式
	mime := bytes.NewBuffer(nil)
	//设置邮件
	mime.WriteString(fmt.Sprintf("From: %s<%s>\r\nTo: %s\r\nCC: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\n", email, email, toEmail, toEmail, Subject))
	mime.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	bodyMsg := make([]byte, base64.StdEncoding.EncodedLen(len(body)))
	mime.WriteString("Content-Description: This is a multi-part message in MIME format.\r\n")
	//邮件普通Text正文
	mime.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	mime.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	mime.WriteString("This is a multipart message in MIME format.")
	mime.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	mime.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	mime.WriteString("Content-Transfer-Encoding: base64\r\n")
	base64.StdEncoding.Encode(bodyMsg, []byte(body))
	mime.Write(bodyMsg)
	mime.WriteString(fmt.Sprintf("\n--%s\r\n", boundary))
	mime.WriteString("Content-Type: application/octet-stream\r\n")
	mime.WriteString("Content-Description: 附一个文件\r\n")
	mime.WriteString("Content-Transfer-Encoding: base64\r\n")
	mime.WriteString("Content-Disposition: attachment; filename=\"" + attaFileName + "\"\r\n\r\n")
	//读取并编码文件内容
	attaData, err := ioutil.ReadFile(attaFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	b := make([]byte, base64.StdEncoding.EncodedLen(len(attaData)))
	base64.StdEncoding.Encode(b, attaData)
	mime.Write(b)
	mime.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	msg := mime.Bytes()
	return msg
}

func SendMail(server, portSTR, email, pwd, toEmail, Subject, Content, file string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", email, "搬砖人Bison")       // 发件人
	m.SetHeader("To", m.FormatAddress(toEmail, "收件人"))  // 收件人
	m.SetHeader("Cc", m.FormatAddress(toEmail, "收件人"))  //抄送
	m.SetHeader("Bcc", m.FormatAddress(toEmail, "收件人")) // 暗送
	m.SetHeader("Subject", Subject)                     // 主题
	m.SetBody("text/html", Content)                     // 正文
	if file != "" && len(file) > 0 {
		m.Attach(file) //添加附件
	}
	port, _ := strconv.Atoi(portSTR)
	d := gomail.NewPlainDialer(server, port, email, pwd) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
