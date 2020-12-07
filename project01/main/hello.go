package main

import (
	"fmt"
	"sort"
	"time"
	"os"
	"go_code/my-study-go-code/project01/utils"
)


func display(str string) { 
	for w := 0; w < 6; w++ { 
		fmt.Println(str)
	} 
}

func main() {

	fmt.Println(utils.BoilingF)
	fmt.Println(utils.BoilingF2)
	
	fmt.Println(os.Args[0])
	fmt.Println(os.Args[0:len(os.Args)])
	
	 // 调用Goroutine 
	 go display("Welcome")
	  
	 // 正常调用函数
	 // display("nhooo.com") 

	 intValue := []int{10,20,5,8}
	 sort.Ints(intValue)
	 
	 fmt.Println("Ints:", intValue)

	 // 创建一个数组
	 arr := [7]string{"这", "是", "Golang", "基础", "教程", "在线", "www.nhooo.com"}

	 fmt.Println(len(arr))

	 for i := 0; i < len(arr); i++ {
		fmt.Println("数组:", arr[i])
	 }

	 // 显示数组
	 fmt.Println("数组:", arr)
 
	 // 创建切片
	 myslice := arr[1:6]
 
	 // 显示切片
	 fmt.Println("切片:", myslice)
 
	 // 显示切片的长度
	 fmt.Printf("切片长度: %d", len(myslice))
 
	 // 显示切片的容量
	 fmt.Printf("\n切片容量: %d\n", cap(myslice))
	
	 fmt.Println(time.Now())


}