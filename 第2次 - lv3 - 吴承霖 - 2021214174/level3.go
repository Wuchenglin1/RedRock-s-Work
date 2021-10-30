package main

import "fmt"

func Receiver(x interface{}){
	switch x.(type){
	case string:
		fmt.Println("这是string类型")
	case int:
		fmt.Println("这是int类型")
	case bool:
		fmt.Println("这是bool类型")
	default:
		fmt.Println("判断失败")
	}

}
func main(){
	var a interface{}
	a = "请输入判断类型的内容"//请输入判断类型的内容
	Receiver(a)
}
