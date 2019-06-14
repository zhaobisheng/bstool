package Cache

import (
	"nodeServer/Logger"
)

func CheckKey(Key string) bool {
	rs, err := Exists(Key)
	if err != nil {
		Logger.Errorln("func-CheckKey-error:", err)
		return false
	}
	if rs > 0 {
		return true
	}
	return false
}
