package Hash

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func Sha1sum(tFile string) string {
	ha := sha1.New()
	f, err := os.Open(tFile)
	if err != nil {
		fmt.Println("error1!", err)
		return ""
	}
	defer f.Close()
	if _, err := io.Copy(ha, f); err != nil {
		fmt.Println("error2:", err)
		return ""
	}
	sha1 := fmt.Sprintf("%x", ha.Sum(nil))
	return sha1
}
