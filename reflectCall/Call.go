package reflectCall

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrParamsNotAdapted = errors.New("Function parameters mismatch.")
)

//定义控制器函数Map类型，便于后续快捷使用
type FuncMap map[string]reflect.Value

func InitFunc(size int) FuncMap {
	return make(FuncMap, size)
}

func (fMap FuncMap) FuncBand(key string, vpt interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("FuncBand fail:" + key + " is not callable.")
		}
	}()
	v := reflect.ValueOf(vpt)
	v.Type().NumIn()
	fMap[key] = v
	return
}

func (fMap FuncMap) FunCall(key string, args []interface{}) (result []reflect.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%s call fail:%s", key, e))
		}
	}()
	
	if _, ok := fMap[key]; !ok {
		err = errors.New(key + " does not exist.")
		return
	}
	fun := fMap[key] //reflect.ValueOf(fMap[key])
	fmt.Printf("%T:%+v\n", fun.Type(), fun.Type())
	//fmt.Println("NumField:", fun.Type().Key().Kind())
	if len(args) != fun.Type().NumIn() {
		if !strings.Contains(fun.Type().String(), "interface {}") {
			err = ErrParamsNotAdapted
			return
		}
	}
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	result = fun.Call(in)
	return
}

func GetFunList(enter interface{}) FuncMap {
	vf := reflect.ValueOf(enter)
	fmt.Printf("%T-%T:%+v\n", vf, vf.Type(), vf.Type())
	vft := vf.Type()
	//读取方法数量
	mNum := vf.NumMethod()
	fmt.Println("vft.NumMethod():", vft.NumMethod(), "mNum:", mNum)
	if mNum > 0 {
		crMap := make(FuncMap)
		fmt.Println("NumMethod:", mNum)
		//遍历路由器的方法，并将其存入控制器映射变量中
		for i := 0; i < mNum; i++ {
			mName := vft.Method(i).Name
			fmt.Println("index:", i, " MethodName:", mName)
			crMap[mName] = vf.Method(i) //<<<
		}
	}
	return nil
}
