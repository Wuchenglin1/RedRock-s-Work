package main

import (
	"fmt"
)
type Details struct {
	Author
	More

}
type More struct {
	Recommended []string
}
type Author struct {
	Name string
	VIP bool
	Signature string
	Focus int
}
func main() {
	var a Author
	a.Name = "逆风笑"
	a.VIP = false
	a.Signature = "哇这里原来能写签名"
	a.Focus = 1783000
	var b More
	b.Recommended = make([]string, 10)

	fmt.Printf("up主的名字: %v\n", a.Name)
	fmt.Printf("是否为vip%v\n", a.VIP)
	fmt.Printf("up主的签名%v\n", a.Signature)
	fmt.Printf("up主的关注数%v\n", a.Focus)
	fmt.Printf("up主视频的相关推荐%v", b.Recommended[0])
}