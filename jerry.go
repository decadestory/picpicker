package main

import (
	"flag"
	"fmt"
	"jerry/jandan"
)

var nv bool
var ng bool
var c int
var s int

func init() {
	flag.BoolVar(&nv, "nv", false, "获取美女图片")
	flag.BoolVar(&ng, "ng", false, "获取GIF图片")
	flag.IntVar(&c, "c", 100, "获取图片的数量")
	flag.IntVar(&s, "s", 0, "跳过的数量")
}

func main() {

	flag.Parse()

	if flag.NFlag() == 0 {
		showUsage()
		return
	}

	if (nv && ng) || (!nv && !ng) {
		jandan.GetPic(c, s)
	} else if nv {
		jandan.GetNv(c, s)
	} else if ng {
		jandan.GetPic(c, s)
	}

}

func showUsage() {
	fmt.Println(`Usage of jerry:
	-c int 获取图片的数量 (default 100)
	-ng    获取无聊图片
	-nv    获取美女图片
	例如: jerry -nv -c 50
	
	`)
}
