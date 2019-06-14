package Utils

import (
	"strings"

	"os"
	"path"

	"github.com/dutchcoders/goftp"
)

func FtpConnect(user, password, host, port, path string) (*goftp.FTP, error) {
	ftp, err := goftp.ConnectTimeout(host + ":" + port)
	if err == nil {
		//fmt.Println("path:", path)
		//defer ftp.Close()
		err = ftp.Login(user, password)
		if err == nil {
			err = ftp.Cwd(path)
			if err != nil {
				if strings.Contains(path, "/") {
					temp := strings.Split(path, "/")
					for index := 0; index < len(temp); index++ {
						var realPath string
						for index1 := 0; index1 <= index; index1++ {
							if temp[index1] != "" {
								realPath = realPath + temp[index1] + "/"
							}
						}
						realPath = realPath[:len(realPath)-1]
						err1 := ftp.Cwd(realPath)
						if err1 != nil {
							err = ftp.Mkd(temp[index])
							if err != nil {
								return nil, err
							} else {
								err2 := ftp.Cwd(temp[index])
								if err2 != nil {
									return nil, err2
								} else {
									return ftp, nil
								}
							}
						}
					}
				} else {
					err = ftp.Mkd(path)
					if err != nil {
						return nil, err
					}
				}
				err = ftp.Cwd(path)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	return ftp, nil
}

func FtpUploadFile(ftp *goftp.FTP, localFile string) error {
	// Upload a file
	var file *os.File
	var err error
	if file, err = os.Open(localFile); err != nil {
		return err
	}
	var remoteFileName = path.Base(localFile)

	if err := ftp.Stor(remoteFileName, file); err != nil {
		return err
	}
	return nil
}
