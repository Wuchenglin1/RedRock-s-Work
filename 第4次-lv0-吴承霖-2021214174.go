package main

import (
	"fmt"
	"strings"
)

/*
预习作业lv4-发金币
你要发金币了，需要分配给以下几个人：Matthew,Sarah,Augustus,Heidi,Emilie,Peter,Giana,Adriano,Aaron,Elizabeth。
分配规则如下：
a. 名字中每包含1个'e'或'E'分1枚金币
b. 名字中每包含1个'i'或'I'分2枚金币
c. 名字中每包含1个'o'或'O'分3枚金币
d: 名字中每包含1个'u'或'U'分4枚金币
写一个程序，计算每个用户分到多少金币，以及最后发出去了多少金币？
程序结构如下，请实现 ‘dispatchCoin’ 函数
*/
	var userName = []string{"Matthew", "Sarah", "Augustus", "Heidi", "Emilie", "Peter", "Giana", "Adriano", "Aaron", "Elizabeth"}
	var	distribution = make(map[string]int, len(userName))
	var sum int//总共分发的金币数

	func dispatchCoin ()int{
		for v := range distribution {
			sum += distribution[v]
		}
		return sum
	}


	func main() {
		for i,v := range userName{
				c1 := strings.Count(userName[i], "e")
				c2 := strings.Count(userName[i], "E")
				c3 := strings.Count(userName[i], "i")*2
				c4 := strings.Count(userName[i], "I")*2
				c5 := strings.Count(userName[i], "o")*3
				c6 := strings.Count(userName[i], "O")*3
				c7 := strings.Count(userName[i], "u")*4
				c8 := strings.Count(userName[i], "U")*4
				distribution[v] += c1+c2+c3+c4+c5+c6+c7+c8
		}

		for i,v := range distribution{
			fmt.Printf("用户%v分发到的金币数是%v\n",i,v)
		}
		left := dispatchCoin()
		fmt.Println("剩下：", left)
	}
