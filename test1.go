package main

import (
	"fmt"
)

func myTest(a [5]int, target int) {
	// 遍历数组
	for i := 0; i < len(a); i++ {
		other := target - a[i]
		// 继续遍历
		for j := i + 1; j < len(a); j++ {
			if a[j] == other {
				fmt.Printf("(%d,%d)\n", i, j)
			}
		}
	}
}

type person struct {
	name string
	city string
	age  int8
}

type student1 struct {
	name string
	age  int
}

func main1() {
	var p2 = new(person)
	p2.name = "测试"
	p2.age = 18
	p2.city = "北京"
	fmt.Printf("p2=%#v\n", p2)
	fmt.Printf("p2=%v\n", *p2)

	a := [5]int{1, 2, 3, 4, 5}
	b := &a
	c := *b
	fmt.Printf("%#v,%p\n", b, b)
	fmt.Printf("%p\n", &a)
	fmt.Println((*b)[1])
	fmt.Println(c)

	m := make(map[string]*student)
	stus := []student{
		{name: "pprof.cn", age: 18},
		{name: "测试", age: 23},
		{name: "博客", age: 28},
	}

	fmt.Printf("%#v\n",stus)

	for _, stu := range stus {
		stu1 := stu
		m[stu.name] = &stu1   //	其实&stu是一个地址并不会改变
		fmt.Printf("%p\n",&stu1)
	}

	for k, v := range m {
		fmt.Println(k, "=>", v.name)
	}



}
