package pkg2

import (
	"fmt"
)

//包变量可见性

var a = 100 //首字母小写，外部包不可见，只能在当前包中使用

const Mode  = 1 //首字母大写，对外部包是可见的，可在其他包中使用

type person struct {   //首字母小写，外部包不可见，只能在当前包中使用
	name 	string
}

func Add(x,y int) (z int)  {  //首字母大写，外部包可见，可在其他包中使用
	return x+y
}

func  age()  {  //首字母小写，外部包不可见，只能在当前包中使用
	var Age = 18   //函数局部变量，外部包不可见，只能在当前函数内使用
	fmt.Println(Age)

}

//结构体中的字段和接口中的方法名如果首字母都是大写，外部包可以访问这些字段和方法。例如：

type  Student struct {
	Name 	string
	class 	string    //字段名必须大写外部包才能调用
}

func (s *Student) SetName()  {  //方法名必须大写，外部包才能调用
	s.Name = "tom"
	s.class = "6"

}

type Payer interface {
	init()   //仅限包内访问的方法
	Pay()    //可以在包外访问的方法
}

//包的导入
//包名是哦从$GIPATH/src后开始计算的，使用 / 进行路径分隔

//匿名导入包
//如果只是希望导入包，而不使用内部的数据时，可以使用匿名导入包，匿名导入的包与其他方式导入的包一样会被编译到可执行文件中

func init()  {
	fmt.Println(Mode)

}














