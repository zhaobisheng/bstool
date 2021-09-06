package Hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func Sha256sum(tFile string) string {
	ha := sha256.New()
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
	sha256 := fmt.Sprintf("%x", ha.Sum(nil))
	return sha256
}
