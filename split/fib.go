package split


//Fib是一个计算第n个斐波那契数列的函数
func Fib(n int) int  {
	if n < 2 {
		return n
	}
	return Fib(n-1)+Fib(n-2)
}