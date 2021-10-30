package main

import "fmt"

type Vedio struct {
	DIanzan int
	Shoucang int
	Toubi int
	Yjsl int
}

func(v Vedio) SetInfo(dz int,sc int,tb int,yjsl int){
	v.DIanzan = dz
	v.Shoucang = sc
	v.Toubi = tb
	v.Yjsl = yjsl
}

func(v Vedio) PrinInfo(){
	fmt.Printf("点赞数:%v\n",v.DIanzan)
	fmt.Printf("收藏数:%v\n",v.Shoucang)
	fmt.Printf("投币数:%v\n",v.Toubi)
	fmt.Printf("一键三连数:%v\n",v.Yjsl)
}
func main()  {
	var p1 = Vedio{
		DIanzan: 10,
		Shoucang: 20,
		Toubi: 30,
		Yjsl: 40,
	}

	p1.SetInfo(p1.DIanzan,p1.Shoucang,p1.Toubi,p1.Yjsl)//点赞数，收藏数，投币数，一键三连数
	p1.PrinInfo()
}