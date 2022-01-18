package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

type Student struct {
	ID	int
	Gender	string
	Name	string
}

type Class struct {
	Title	string
	Students	[]*Student
}

func jsonTest1()  {
	c := &Class{
		Title:    "101",
		Students: make([]*Student,0,200),
	}

	for i := 0; i < 10; i++ {
		stu := &Student{
			ID:     i,
			Gender: "男",
			Name:   fmt.Sprintf("stu%02d",i),
		}
		c.Students = append(c.Students,stu)
	}

	//json序列化：结构体--> json格式字符串
	data, err := json.Marshal(c)     //c是一个指针
	if err != nil {
		fmt.Println("json marshal failed")
		return
	}
	fmt.Printf("json:%s\n",data)

	fmt.Printf("%T\n",c)
	//json反序列化：json格式的字符串-->结构体
	str := `{"Title":"101","Students":[{"ID":0,"Gender":"男","Name":"stu00"},{"ID":1,"Gender":"男","Name":"stu01"},{"ID":2,"Gender":"男","Name":"stu02"},{"ID":3,"Gender":"男","Name":"stu03"},{"ID":4,"Gender":"男","Name":"stu04"},{"ID":5,"Gender":"男","Name":"stu05"},{"ID":6,"Gender":"男","Name":"stu06"},{"ID":7,"Gender":"男","Name":"stu07"},{"ID":8,"Gender":"男","Name":"stu08"},{"ID":9,"Gender":"男","Name":"stu09"}]}`
	c1 := &Class{}
	err = json.Unmarshal([]byte(str),c1)   //interface传指针跟值都可以吗
	if err != nil {
		fmt.Println("json unmarshal failed!")
		return
	}
	fmt.Printf("%#v\n",c1)
}




type student struct {
	id   int
	name string
	age  int
}

func demo(ce []student) {
	//切片是引用传递，是可以改变值的
	ce[1].age = 999
	// ce = append(ce, student{3, "xiaowang", 56})
	// return ce
}

func test1()  {
	var ce []student  //定义一个切片类型的结构体
	ce = []student{
		student{1, "xiaoming", 22},
		student{2, "xiaozhang", 33},
	}
	fmt.Println(ce)
	demo(ce)
	fmt.Println(ce)
}

func mapSort()  {
	map1 := make(map[int]string,5)
	map1[1] = "www"
	map1[2] = "http"
	map1[3] = "."
	map1[4] = "com"
	map1[5] = "cn"
	sli := []int{}

	for k,_ := range map1 {
		sli = append(sli,k)
	}
	sort.Ints(sli)
	for i := 0;i<len(map1);i++ {
		fmt.Println(map1[sli[i]])
	}
}


func length(s string) int {
	println("call length.")
	return len(s)
}


// 外部引用函数参数局部变量
func add(base int) func(int) int {
	return func(i int) int {
		base += i
		return base
	}
}

func factorial(i int) int  {
	if i <= 1 {
		return 1
	} else  {
		return i*factorial(i-1)
	}
}

type Test struct {
	name	string
}

func (t *Test) Close()  {
	fmt.Println(t.name,"closed")

}

//可变参数通常要作为函数的最后一个参数
//本质上函数的可变参数是通过切片来实现的
func intSum2(x ...int) int  {  //x可变参数实际是一个切片
	fmt.Println(x)
	sum := 0
	for _,v := range x {
		sum += v
	}
	return sum
}

//go语言中函数支持多返回值，函数如果有多个返回值必须用（）将所有返回值包裹起来。
func calu(x,y int) (sum,sub int)  {
	sum = x + y
	sub = x -y
	return       //此时可以return后可以不添加具体需要返回的参数，因为返回值我们在函数定义时已经命名
}

//高阶函数分为函数作为参数和函数作为返回值两部分

//函数作为参数
func adder(x,y int) int {
	return x + y
}
func suber(x,y int) int {
	return x - y
}

func calc(x,y int,op func(int,int) int) int  {
	return op(x,y)
}

//函数作为返回值
func do(s string) (func(int ,int) int, error)  {
	switch s {
	case "+":
		return adder,nil
	case "-":
		return suber,nil
	default:
		err := errors.New("无法识别的操作符")
		return nil,err

	}
}


//匿名函数和闭包
//函数当然还可以作为返回值，但是在go语言中函数内部不能再像之前那样定义函数了，只能定义匿名函数，匿名函数就是没有函数名的函数，匿名函数的定义格式如下：
//匿名函数需要保存到某个变量或者作为立即执行的函数
//匿名函数多用于实现回调函数以及闭包
func anynous()  {
	//将匿名函数保存到变量
	add := func(x,y int ) {
		fmt.Println(x+y)
	}

	add(10,20)

	//自执行函数：匿名函数定义完加()直接执行
	func(x,y int) {
		fmt.Println(x+y)
	}(10,20)

}

//闭包
//闭包指的是一个函数和其相关的引用环境组合而成的实体。简单来说， 闭包=函数+引用环境。

func adders() func(int) int {
	var x int

	//实际上返回的是一个匿名函数，f变量保存匿名函数的
	return func(y int) int {
		x += y
		return x
	}
}

//闭包进阶

func makeSuffixFunc(suffix string) func(string) string  {
	return func(name string) string {
		if !strings.HasSuffix(name,suffix){
			return name + suffix
		}
		return name

	}

}

//defer语句
//go语言中的defer语句会将其后面跟随的语句进行延迟处理，在defer归属的函数即将返回时，将延迟处理的语句按defer定义的逆序进行执行，也就是说，先被defer的语句最后被执行
//最后被defer的语句先被执行
func calc1(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func funcA() {
	fmt.Println("func A")
}

func funcB() {
	defer func() {
		err := recover()
		//如果程序出出现了panic错误,可以通过recover恢复过来
		if err != nil {
			fmt.Println("recover in B")
		}
	}()
	panic("panic in B")
}

func funcC() {
	fmt.Println("func C")
}


var (
	coins = 50
	users = []string{
		"Matthew", "Sarah", "Augustus", "Heidi", "Emilie", "Peter", "Giana", "Adriano", "Aaron", "Elizabeth",
	}
	distribution = make(map[string]int,len(users))
)

func dispatchCoin()  int {

	var disptchCoins int = 0
	for _,v := range users{
		var personCoins int = 0
		for _,char := range v{
			switch string(char) {   //将整数强制转换为字符串，安装ascii表转换
				case "e" , "E":
					personCoins+=1
				case "i" , "I":
					personCoins+=2
				case "o" , "O":
					personCoins+=3
				case "u" , "U":
					personCoins+=4
				}
		distribution[v] = personCoins
		}
		disptchCoins += personCoins
	}
	fmt.Println(distribution)
	//fmt.Println(disptchCoins)

	return coins - disptchCoins

}

// 遍历字符串
func traversalString() {
	s := "hello沙河"
	for i := 0; i < len(s); i++ { //byte
		fmt.Printf("%v(%c) ", s[i], s[i])
	}
	fmt.Println()
	for _, r := range s { //rune
		fmt.Printf("%v(%c)(%T) ",r, r, r)
		fmt.Println(string(r))
	}
	fmt.Println()
}

//只有当结构体实例化时，才会真正的分配内存，也就是必须实例化后才能使用结构体的字段。
//结构体本身也是一种类型，我们可以像生命内置类型一样使用var 关键字生命结构体类型

//基本实例化
type person1 struct {
	name	string
	city	string
	age		int8
}

//匿名结构体
//在定义一些临时数据等场景下还可以使用匿名结构体。

//创建指针类型的结构体
//我们hai可以使用new关键字对结构体进行实例化，得到的式结构体的地址。格式如下：


//结构体初始化
//没有初始化的结构体，其成员变量都是对应其类型的零值

//空结构体
//空结构体是不占内存空间的
var v struct{}

//构造函数
//go语言的结构体没有构造函数，我们可以自己实现，例如，下面的代码就实现了一个person的构造函数。因为 struct 是值类型，如果结构体比较复杂的话，值拷贝性能开销是比较大的，所以该构造函数
//返回的是结构体指针类型
func newPerson(name,city string, age int8) *person1  {
	return &person1{
		name: name,
		city: city,
		age:  age,
	}
}

//go语言中的方法是一种作用于特定类型变量的函数。这种特定类型的函数叫做接收者，接受这的概念类似其他语言的this 或者 self
//方法与函数的区别是，函数不属于任何类型，方法属于特定的类型

//构造函数 与 方法一起使用
type Person4 struct {
	name 	string
	age 	int8
}

//NewPerson构造函数
func NewPerson(name string, age int8) *Person4  {
	return &Person4{
		name: name,
		age:  age,
	}
}

//Dream Person4做梦的方法

func (p Person4)  Dream()  {
	fmt.Printf("%s的梦想是学习好go语言！\n",p.name)
	fmt.Printf("%s的年龄是%d\n",p.name,p.age)
}

//指针类型的接收者
//指针类型的接受这由一个结构体的指针组成，由于指针的特性，调用方法时修改接收者指针的任意成员变量，在方法结束后，修改都是有效的。这种方式就十分接近于其他语言中面向对象中的this，self

func (p *Person4) SetAge(newAge int8)  {
	p.age = newAge
}

//值类型的接收者
//当方法作用于值类型接收者时，go语言会在代码运行时将接收者的值复制一份。在值类型接收者的方法中可以获取接收者的成员值，但修改操作只是针对副本，无法修改接收者变量本身
func (p Person4) SetAge1(newAge int8)  {
	p.age = 28
}
//总结 什么时候应该使用指针类型接收者
/*
1、需要修改接收者中的值
2、接收者时拷贝代价比较大的大对象
3、保证一致性，如果有某个方法使用了指针接收者，那么其他的方法也应该使用指针接收者。
*/

/*
因为slice和map这两种数据类型都包含了指向底层数据的指针，因此我们在需要复制他们时要特别注意
*/

type Person struct {
	name 	string
	age 	int8
	dreams 	[]string
}

func (p *Person) SetDreams(dreams []string)  {
	p.dreams = dreams
}


/*
使用面向对象的思维方式编写一个学生信息管理系统
	1、学生有id、姓名、年龄、分数等
	2、程序提供展示学生列表、添加学生、编辑学生信息、删除学生等功能
*/

type Stu struct {
	id 	int
	name 	string
	age 	int
	score 	int
}

type Cla struct {
	mp 	map[int]*Stu
}

func printLog()  {
	fmt.Println("欢迎来到学生信息管理系统")
	fmt.Println("展示学生信息列表：1")
	fmt.Println("添加学生信息：2")
	fmt.Println("编辑学生信息：3")
	fmt.Println("删除学生信息：4")
	fmt.Println("展示学生个人信息:5")
	fmt.Printf("请输入对应的数字:")
}

func (c *Cla)  displayOne()   {
	var id int
	fmt.Print("请输入需要查询学生的id:")
	_,err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("获取id失败")
	}
	for k,v := range c.mp {
		if k == id {
			fmt.Printf("id:%d\t,姓名:%s\t,年龄:%d\t,分数:%d\n",v.id,v.name,v.age,v.score)
		} else {
			fmt.Println("查无此人")
		}
	}
}

func (c *Cla)  displayAll() {
	fmt.Printf("\t%s\t%s\t%s\t%s\n","id","姓名","年龄","分数")
	sortId := make([]int,0)
	for k,_ := range c.mp {
		sortId = append(sortId,k)
	}

	sort.Ints(sortId)
	for _,v := range sortId {
		s := c.mp[v]
		fmt.Printf("\t%d\t%s\t%d\t%d\t",s.id,s.name,s.age,s.score)
	}
}

func (c *Cla)  add() {
	var id int
	var name string
	var age int
	var score int
	fmt.Print("输入id:")
	_,err := fmt.Scan(&id)
	fmt.Print("输入名字:")
	_,err = fmt.Scan(&name)
	fmt.Print("输入年龄:")
	_,err = fmt.Scan(&age)
	fmt.Print("输入分数:")
	_,err = fmt.Scan(&score)
	if err != nil {
		fmt.Println("获取信息失败！")
	} else {
		_,ok := c.mp[id]
		if ok {
			fmt.Println("学生已经存在")
			return
		}
		student := &Stu{
			id:    id,
			name:  name,
			age:   age,
			score: score,
		}
		c.mp[id] = student
		fmt.Println("保存成功")
		}
	}

func (c *Cla)  modify()  {
	var id int
	fmt.Println("请输入修改学生的id")
	_,err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("获取id失败")
	}
	_,ok := c.mp[id]
	if !ok {
		fmt.Println("要修改的学生id不存在")
		return
	}
	var name string
	var age int
	var score int
	fmt.Print("输入名字:")
	_,err = fmt.Scan(&name)
	fmt.Print("输入年龄:")
	_,err = fmt.Scan(&age)
	fmt.Print("输入分数:")
	_,err = fmt.Scan(&score)
	if err != nil {
		fmt.Println("获取信息失败！")
	}
	student := &Stu{
		id:    id,
		name:  name,
		age:   age,
		score: score,
	}
	c.mp[id] = student
	fmt.Println("修改成功")
}

func (c *Cla) del()  {
	var id int
	fmt.Println("请输入删除学生的id")
	_,err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("获取id失败")
	}
	_,ok := c.mp[id]
	if !ok {
		fmt.Println("要删除的id不存在")
		return
	}
	delete(c.mp,id)
	fmt.Println("删除成功！")
}


//init()函数介绍
//在go语言中执行时导入包语句会自动触发包内部的init()函数的调用。需要注意的是：init()函数没有参数也没有返回值。init()函数在程序运行时自动被调用执行，不能在代码中主动调用它

var h int = 8

const pi  = 3.14

func init()  {
	fmt.Println(h)
}

//接口就是一个需要实现的方法列表，只要实现了接口中的所有方法，就实现了这个接口

//接口类型变量
//接口类型变量能够存储所有实现了该接口的实例

type Sayer interface {
	say()
}

type dog struct {
}

type cat struct {
}

func (d dog) say()  {
	fmt.Println("汪汪汪")
}

func (c cat) say()  {
	fmt.Println("喵喵喵")
}


//值接收者和指针接收者实现接口的区别
type Mover interface {
	move()
}

func (d dog) move()  {
	fmt.Println("狗会动")
}

type People interface {
	Speak(string) string
}

type Student2 struct{}

func (stu *Student2) Speak(think string) (talk string) {
	if think == "sb" {
		talk = "你是个大帅比"
	} else {
		talk = "您好"
	}
	return
}

func show(a interface{})  {
	fmt.Printf("type:%T value:%#v\n",a,a)
}

//空接口是指没有定义任何方法的接口。因此任何类型都实现了空接口
//空接口类型的变量可以存储任意类型的变量


//使用接口的方式实现一个既可以往终端写日志也可以往文件写日志的简易日志库
type Logger interface {
	Info(string)
}

type FileLogger struct {
	filename 	string
}

func (fl *FileLogger) Info(msg string)  {
	var f *os.File
	var err1 error
	if checkFileIsExit(fl.filename) {
		f,err1 = os.OpenFile(fl.filename,os.O_APPEND|os.O_WRONLY,0666)   //打开文件
		fmt.Println("文件存在")
	} else {
		f,err1 = os.Create(fl.filename)  //创建文件
		fmt.Println("文件不存在")
	}
	defer  f.Close()
	n,err1 := io.WriteString(f,msg+"\n")  //写入文件(字符串)
	if err1 != nil {
		panic(err1)
	}
	fmt.Printf("写入%d个字节\n",n)
}

func checkFileIsExit(filename string) bool  {
	if _,err := os.Stat(filename);os.IsNotExist(err) {
		return false
	}
	return true
}

type ConsoleLogger struct {

}

func (cl *ConsoleLogger) Info(msg string)  {
	fmt.Println(msg)
}

func homework()  {
	var logger Logger
	FileLogger := &FileLogger{"log.txt"}
	logger = FileLogger
	logger.Info("Hello")
	logger.Info("how are you")

	ConsoleLogger := &ConsoleLogger{}
	logger = ConsoleLogger
	logger.Info("Hello")
	logger.Info("how are you")

}


//在go语言中，使用reflect.TypeOf()函数可以获得任意值的类型对象（reflect.Type），程序通过类型对象可以访问任意值的类型信息

func reflectType(x interface{})  {
	v := reflect.TypeOf(x)
	fmt.Printf("Type:%v\n",v)
}

// type name 和 type kind
//在反射中关于类型还划分为两种：类型(Type) 和 种类(Kind)。因为在go语言中我们可以使用type关键字构造很多自定义类型，而种类(Kind)就是指底层的类型，但在反射中，
//当需要区分指针、结构体等大品种的类型时，就会用到种类(Kind).

type myint int64

func reflectType1(x interface{})  {
	t := reflect.TypeOf(x)
	fmt.Printf("type:%v,kind:%v\n",t.Name(),t.Kind())
}
//在go语言反射中像数组、切片、Map、指针等类型的变量，它们的 .Name()都是返回 空。


//VauleOf
//reflect.ValueOf()返回的是reflect.Value类型，其中包含了原始值的值信息。reflect.Value与原始值直之间可以互相转换
//通过反射获取值

func reflectValue( x interface{})  {
	v := reflect.ValueOf(x)
	fmt.Printf("%v\n",v)
	k := v.Kind()
	fmt.Printf("%v\n",k)
	switch k {
	case reflect.Int64:
		//v.Int()从反射中获取整形的原始值，然后通过int64()强制类型转换
		fmt.Printf("type is int64,value is %d\n",int64(v.Int()))
	case reflect.Float32:
		//v.Float()从反射中获取浮点型的原始值，然后通过float32()强制转换
		fmt.Printf("type is float32,value is %f\n",float32(v.Float()))
	case reflect.Float64:
		//v.Float()从反射中获取浮点型的原始值，然后通过float64()强制转换
		fmt.Printf("type is float64,value is %f\n",float64(v.Float()))
	}
}

type student2 struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// 给student添加两个方法 Study和Sleep(注意首字母大写)
func (s student2) Study() string {
	msg := "好好学习，天天向上。"
	fmt.Println(msg)
	return msg
}

func (s student2) Sleep() string {
	msg := "好好睡觉，快快长大。"
	fmt.Println(msg)
	return msg
}

func printMethod(x interface{}) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	fmt.Println(t.NumMethod())
	for i := 0; i < v.NumMethod(); i++ {
		methodType := v.Method(i).Type()
		fmt.Printf("method name:%s\n", t.Method(i).Name)
		fmt.Printf("method:%s\n", methodType)
		// 通过反射调用方法传递的参数必须是 []reflect.Value 类型
		var args = []reflect.Value{}
		v.Method(i).Call(args)
	}
}



func hello(i int)  {
	defer wg.Done()    //gorouting结束就登记-1
	fmt.Println("Hello Gorouting!",i)
}

func recv(c chan int)  {
	ret := <- c
	fmt.Println("接收成功",ret)
}

//worker pool (goroutine池)
//在工作中我们通常会使用可以指定启动的goroutine数量 - worker pool 模式，控制 goroutine 的数量，防止 goroutine泄露和暴涨
func worker(id int, jobs <- chan int, results chan <- int)  {
	for j := range jobs {
		fmt.Printf("worker:%d start job:%d\n",id,j)
		time.Sleep(time.Second)
		fmt.Printf("worker:%d end job:%d\n",id,j)
		results <- j *2
	}

}


var (
	x      int64
	wg     sync.WaitGroup
	lock   sync.Mutex
	rwlock sync.RWMutex
)

func write() {
	// lock.Lock()   // 加互斥锁
	rwlock.Lock() // 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
	rwlock.Unlock()                   // 解写锁
	// lock.Unlock()                     // 解互斥锁
	wg.Done()
}

func read() {
	// lock.Lock()                  // 加互斥锁
	rwlock.RLock()               // 加读锁
	time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
	rwlock.RUnlock()             // 解读锁
	// lock.Unlock()                // 解互斥锁
	wg.Done()
}

func main() {
	log.Println("这时一条很普通的日志")
	v := "很普通的"
	log.Printf("这时一条%s日志",v)
	//log.Fatalln("这时一条会触发fatal的日志")
	//log.Panicln("这时一条会触发panic的日志")
	//logger会打印每条日志信息的日期，时间，默认输出到系统的标准错误。Fatal系列函数会在写入日志信息后调用os.Exit(1).Panic系列函数会在写入日志信息后panic

	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	log.Println("这是一条很普通的日志")

	logFile ,err := os.OpenFile("./xxx.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644)
	if err !=nil {
		fmt.Println("open file failed,err:",err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile|log.Lmicroseconds|log.Ldate)
	log.Println("这是一条很普通的日志")
	log.SetPrefix("[小王子]")
	log.Println("这是一条很普通的日志")

	logger := log.New(logFile, "<New>", log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println("这是自定义的logger记录的日志。")


	/*
	var name string
	var age int
	var married bool
	var delay time.Duration
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 18, "年龄")
	flag.BoolVar(&married, "married", false, "婚否")
	flag.DurationVar(&delay, "d", 0, "时间间隔")
	fmt.Println(name,age,married,delay)
*/

/*
	//flag.Type()
	//基本格式 flag.Type(flag名,默认值,帮助信息)*Type
	name := flag.String("name","张三","姓名")
	age := flag.Int("age",18,"年龄")
	married := flag.Bool("married",false,"婚否")
	delay := flag.Duration("d",0,"时间间隔")
	fmt.Println(name,age,married,delay)
	//需要注意的是，此时 name 、age、married、delay均为对应类型的指针
*/


/*
	fmt.Println(time.Now())
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano())

	now := time.Now()
	later := now.Add(time.Hour)
	fmt.Println(later)
*/


/*
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go read()
	}

	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
*/
/*
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x)
		case ch <- i:
		}
	}
*/


/*
   	jobs := make(chan int,100)
   	results := make(chan int, 100)
   	//开启3个goroutine
   	for w := 1; w <= 3; w++ {
   		go	worker(w,jobs,results)
   	}
   	//5个任务
   	for j:=1;j<=5;j++ {
   		jobs <- j
   	}
   	close(jobs)
   	//输出结果
   	for a := 1; a<=5; a++ {
   		i := <-results
   		fmt.Println(i)
   	}
*/
	//for i := range results{
	//	fmt.Println(i)
	//}

/*
	//for range从通道中循环取值
	//当向通道中发送完数据时，我们可以通过 close 函数来关闭通道
	//当通道别关闭时，再往该通道中发送值就会引发panic,从该通道取值的操作会先取完通道中的值，再然后取到的值一直都时对应类型的零值。
	//channel联系
	ch1 := make(chan int)
	ch2 := make(chan int)
	//开启goroutine将0~100的数发送到ch1中
	go func() {
		for i := 0; i < 100; i++{
			ch1 <- i
		}
		close(ch1)
	}()

	//开启goroutine从ch1中接收值，并将该值的平方发送到ch2中
	go func() {
		for {
			i, ok := <- ch1   //通道关闭后再取值 ok=false
			if !ok {
				break
			}
			ch2 <- i*i
		}
		close(ch2)
	}()

	//在主goroutine中从ch2中接收值打印
	for i := range ch2  {  //通道关闭以后就会退出for range循环
 		fmt.Println(i)
	}
	//从上面的例子中我们可以看到有两种方式在接收值的时候判断该通道是否被关闭，不过我们通常使用的是 for range的方式。使用 for range遍历通道，当通道被关闭的时候就会退出 for range
*/

/*
	//无缓存通道上的发送操作会阻塞，直到另一个goroutine在该通道上执行接收操作，这时值才能发送成功，两个goroutine将继续执行。相反，如果接收操作先执行，接收方
	//的goroutine将阻塞，直到另一个goroutine在该通道上发送一个值。
	ch := make(chan int)
	go recv(ch)
	ch <- 10
	fmt.Println("发送成功")
*/

/*
	//创建多个gorouting
	for i:=0;i<10;i++ {
		wg.Add(1)   //启动一个gorouting就登记+1
		go hello(i)
		time.Sleep(time.Second)
	}
	//time.Sleep(time.Second)
	fmt.Println("main end")
	wg.Wait()   //等待所有登记的gorouting都结束

*/
	//printMethod(student2{})

/*
	stu1 := student2{
		Name:  "小王子",
		Score: 90,
	}

	t := reflect.TypeOf(stu1)
	fmt.Println(t.Name(), t.Kind()) // student struct
	// 通过for循环遍历结构体的所有字段信息
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("name:%s index:%d type:%v json tag:%v\n", field.Name, field.Index, field.Type, field.Tag.Get("json"))
	}

	// 通过字段名获取指定结构体字段信息
	if scoreField, ok := t.FieldByName("Score"); ok {
		fmt.Printf("name:%s index:%d type:%v json tag:%v\n", scoreField.Name, scoreField.Index, scoreField.Type, scoreField.Tag.Get("json"))
*/
	}



/*
	var a float32 = 3.14
	var b int64 = 100
	reflectValue(a)
	reflectValue(b)
	//将int类型的原始值转换为reflect.Value类型
	c := reflect.ValueOf(10)
	fmt.Printf("type c:%T\n",c)
*/

/*
	// reflect typeof name kind
	var a *float32
	var b myint
	var c rune
	reflectType1(a)
	reflectType1(b)
	reflectType1(c)
	type person struct {
		name 	string
		age 	int
	}
	type  book struct {
		title 	string
	}
	var d = person{
		name: "沙河小王子",
		age:  18,
	}

	var e = book{title:"《跟小王子学习go语言》"}
	reflectType1(d)
	reflectType1(e)
*/

	//reflect.TypeOf()
	//var a float32 = 3.14
	//reflectType(a)
	//var b int32 = 100
	//reflectType(b)


	//实现日志库
	//homework()

	//空接口可以存储任意类型的数据
	//var studentInfo = make(map[string]interface{})
	//studentInfo["name"] = "沙河娜扎"
	//studentInfo["age"] = 18
	//studentInfo["married"] = false
	//fmt.Println(studentInfo)

	// 定义一个空接口x
	//var x interface{}
	//s := "Hello 沙河"
	//x = s
	//fmt.Printf("type:%T value:%v\n", x, x)
	//i := 100
	//x = i
	//fmt.Printf("type:%T value:%v\n", x, x)
	//b := true
	//x = b
	//fmt.Printf("type:%T value:%v\n", x, x)


/*
var x Mover
	var wangcai = dog{}  //wangcai是一个dog类型
	x = wangcai    //x可以接收dog类型

	var fugui = &dog{}   //富贵是*dog类型
	x = fugui   //x可以接收*dog类型
	x.move()

	var peo People = &Student2{}
	think := "bitch"
	fmt.Println(peo.Speak(think))
*/

/*
	var x Sayer   //声明一个Sayer类型的变量x
	a := dog{}   //实例化一个dog
	b := cat{}    //实例化一个cat
	x = a
	x.say()
	x = b
	x.say()
*/

/*
	//包的导入
	test2 := mm.Student{}
	test2.SetName()
	fmt.Println(test2.Name)

	//init函数的调用顺序
	fmt.Println("Hello 川")

	//外部包中函数的调用
	 k:= g.Add(4,2)
	 l:= g.Del(4,2)
	 z := g.Mul(4,2)
	 x := g.Div(4,2)
	 fmt.Println(k,l,z,x)
*/


/*
	//学生信息管理系统
		c := &Cla{}
		c.mp = make(map[int]*Stu, 200)
		for {
			continueTag:
				printLog()
				var do int
				_, err := fmt.Scan(&do)
				if err != nil {
					fmt.Println("输入有误")
				}

				switch do {
					case 1:
						c.displayAll()
					case 2:
						c.add()
					case 3:
						c.modify()
					case 4:
						c.del()
					case 5:
						c.displayOne()
					default:
						fmt.Println("输入有误")
				}
				fmt.Println("继续：1   退出：2")
				var insertCode int
				_,err =  fmt.Scan(&insertCode)
				if err != nil {
					fmt.Println("输入有误！")
				}
				if insertCode == 1 {
					//goto跳转实现操作结束后的下一步动作
					goto continueTag
				} else {
					break
				}
		}
*/

/*
	//结构体实例化
	var p person1
	p.name = "沙河娜扎"
	p.city = "北京"
	p.age = 18
	fmt.Printf("%v\n",p)
	fmt.Printf("%#v\n",p)

	//匿名结构体
	var user struct{Name string; Age int}
	user.Name = "小王子"
	user.Age = 28
	fmt.Printf("%#v\n",user)
	fmt.Printf("%v\n",user)

	//创建指针类型的结构体
	var p2 = new(person1)
	p2.name = "指针类型结构体"
	p2.city = "北京"
	fmt.Printf("%#v\n",p2)
	fmt.Printf("%#v\n",*p2)

	//取结构体的地址实例化
	//使用&对结构体进行取地址操作相当于对该结构体进行了一次new实例化操作。
	p3 := &person1{}
	fmt.Printf("%T\n",p3)
	fmt.Printf("p3=%#v\n",p3)
	p3.name = "川"
	p3.city = "深圳"
	p3.age = 28
	fmt.Printf("%v\n",p3)

	//空结构体是不占内存空间的
	fmt.Println(unsafe.Sizeof(p3))

	//构造函数
	p9 := newPerson("川","深圳",28)
	fmt.Printf("%#v\n",p9)

	p4 := NewPerson("小王子", 25)
	p4.Dream()
	p4.SetAge(28)

	//结构体和方法补充知识点
	p1 := Person{name:"小王子",age:24}
	data := []string{"吃饭","睡觉","打豆豆"}
	p1.SetDreams(data)
	fmt.Printf("%v",p1.dreams)

	data[1] = "不睡觉"
	fmt.Println(p1.dreams)
*/
	//left := dispatchCoin()
	//fmt.Println("剩下：",left)

	//traversalString()

	//mapSort()
	//ts := []Test{{"a"},{"b"},{"c"}}
	//for _,t := range  ts {
	//	fmt.Printf("%p\n",&t)
	//	defer t.Close()
	//}
	//

	//sum := intSum2(1,2,3,4)
	//fmt.Println(sum)

	//ret2 := calc(10,20,adder)
	//fmt.Println(ret2)

	//anynous()

	//var f = adders()   //变量f是一个函数并且它引用了其外部作用域中的x变量，此时f就是一个闭包。在f的生命周期内，变量x也一直有效
	//fmt.Println(f(10))
	//fmt.Println(f(20))
	//fmt.Println(f(30))

	//x := 1
	//y := 2
	//defer calc1("AA", x, calc1("A", x, y))
	//x = 10
	//defer calc1("BB", x, calc1("B", x, y))
	//y = 20


	//funcA()
	//funcB()
	//funcC()

	//aaa := []string{"a","b"}
	//bbb := aaa
	//
	//bbb[1] =  "c"
	//
	//fmt.Println(aaa)






