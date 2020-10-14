package XStringUtil

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseInt(str string) int64 {
	intData, err := strconv.ParseInt(str, 10, 64) //strconv.Atoi(str)
	if err == nil {
		return intData
	} else {
		return 0
	}
}

func ToString(src interface{}) string {
	str := fmt.Sprintf("%v", src)
	return str
}

func ToStringFormat(format string, src interface{}) string {
	str := fmt.Sprintf(format, src)
	return str
}

func ParseFloat(str string) float64 {
	floatData, err := strconv.ParseFloat(str, 10)
	if err == nil {
		return floatData
	} else {
		return 0
	}
}

func Replace(str, oldStr, newStr string) string {
	return strings.Replace(str, oldStr, newStr, -1)
}

func GetSplitMap(soldier string) []*KVstruct { //[]map[string]int64
	//soldiers = "4,18,630154,19,138068,20,32252,21,32000,"
	//fmt.Println("GetSplitMap:", soldier)
	soldier = strings.Replace(soldier, " ", "", -1)
	if len(soldier) <= 2 {
		return nil
	}
	tagFlag := ","
	if soldier[len(soldier)-1:] == tagFlag {
		soldier = soldier[:len(soldier)-1]
	}
	firstIndex := strings.Index(soldier, tagFlag)
	//fmt.Println("strings.Count(soldier, tagFlag)/2:", strings.Count(soldier, tagFlag)/2, soldier[:firstIndex], firstIndex)
	if strconv.Itoa(strings.Count(soldier, tagFlag)/2) == soldier[:firstIndex] {
		soldier = soldier[firstIndex+1:]
		//fmt.Println("soldier:", soldier)
	}
	rs := strings.Split(soldier, tagFlag)
	//fmt.Println("GetSplitMap-rs:", rs)
	mSoldiers := make([]*KVstruct, len(rs)/2)
	//fmt.Println("GetSplitMap-rs:", soldier, " rs:", rs, len(rs))
	for index := 0; index < len(rs); index += 2 {
		//fmt.Println("GetSplitMap-index:", index)
		mSoldiers[index/2] = &KVstruct{Key: rs[index], Val: ParseInt(rs[index+1]), ValType: "int64"}
	}
	/*mSoldiers := make([]map[string]int64, len(rs)/2)
	for index := 0; index < len(rs); index += 2 {
		tempMap := make(map[string]int64)
		tempMap[rs[index]] = ParseInt(rs[index+1])
		mSoldiers[index/2] = tempMap
		//mSoldiers[rs[index]] = ParseInt(rs[index+1])
	}*/
	return mSoldiers
}

func Split(str string) []string {
	str = strings.Replace(str, " ", "", -1)
	if len(str) <= 2 {
		return nil
	}
	tagFlag := ","
	if str[len(str)-1:] == tagFlag {
		str = str[:len(str)-1]
	}
	return strings.Split(str, tagFlag)
}

func SplitTag(str, tagFlag string) []string {
	//tagFlag := ","
	str = strings.Replace(str, " ", "", -1)
	if str[len(str)-1:] == tagFlag {
		str = str[:len(str)-1]
	}
	return strings.Split(str, tagFlag)
}
