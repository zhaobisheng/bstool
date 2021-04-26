package ProcessBar

import (
	"fmt"
	"os"
)

func ShowProcess(index, totalNum int64, show bool) string {
	fixedNum := float64(totalNum) * 0.0002
	if int(fixedNum) < 1 {
		fixedNum = 1.0
	}
	formatStr := "%.2f%%"
	if index >= totalNum {
		if show {
			fmt.Fprintf(os.Stdout, formatStr+"\r", 100.00)
		} else {
			return fmt.Sprintf(formatStr, 100.00)
		}
	} else if index%int64(fixedNum) == 0 {
		value := (float64(index) / float64(totalNum)) * 100
		if show {
			fmt.Fprintf(os.Stdout, formatStr+"\r", value)
		} else {
			return fmt.Sprintf(formatStr, value)
		}
	}
	return ""
}

func ShowProcessWithFormat(index, totalNum int64, show bool, formatStr string) string {
	fixedNum := float64(totalNum) * 0.0002
	if int(fixedNum) < 1 {
		fixedNum = 1.0
	}
	if index >= totalNum {
		if show {
			fmt.Fprintf(os.Stdout, formatStr+"\r", 100.00)
		} else {
			return fmt.Sprintf(formatStr, 100.00)
		}
	} else if index%int64(fixedNum) == 0 {
		value := (float64(index) / float64(totalNum)) * 100
		if show {
			fmt.Fprintf(os.Stdout, formatStr+"\r", value)
		} else {
			return fmt.Sprintf(formatStr, value)
		}
	}
	return ""
}
