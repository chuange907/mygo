package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strings"
	"time"
	_ "github.com/go-redis/redis/v8"
)

type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	//创建一个默认的路由引擎
	r := gin.Default()
	//GET：请求方式；/hello：请求路径
	//当客户端以get方法请求/hello路径时，会执行后面的匿名函数
	r.GET("/hello", func(c *gin.Context) {
		//c.JSON：返回json格式的数据
		c.JSON(200, gin.H{
			"message": "Hello,World",
		})
	})

	r.GET("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GET",
		})

	})

	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POST",
		})
	})

	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PUT",
		})
	})

	r.DELETE("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DELETE",
		})

	})

	r.LoadHTMLGlob("templates/**/*")
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.html", gin.H{
			"title": "posts/index",
		})

	})

	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.html", gin.H{
			"title": "users/index",
		})

	})

	//获取querystring参数
	//querystring 指的是URL中 ? 后面携带的参数，例如：/user/search?username=小王子&address=沙河 。获取请求的querystring参数的方法如下：

	r.GET("user/search", func(c *gin.Context) {
		username := c.DefaultQuery("username", "小王子")
		//username := c.Query("username")
		address := c.Query("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})

	//获取form参数
	//当前端请求的数据通过form表单提交时，例如向 /user/search 发送一个post请求时，获取请求数据的方式如下：
	r.POST("/user/search", func(c *gin.Context) {
		// DefaultPostForm取不到值时会返回指定的默认值
		//username := c.DefaultPostForm("username", "小王子")
		username := c.PostForm("username")
		address := c.PostForm("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})

	//获取json参数
	//当向前端请求的数据通过json提交时，例如向 /json 发送一个 POST请求，则获取请求参数的方式如下：
	r.POST("/json", func(c *gin.Context) {
		b, _ := c.GetRawData() //从c.Request.Body获取请求数据
		//定义map或结构体
		var m map[string]interface{}
		//反序列化
		_ = json.Unmarshal(b, &m)
		c.JSON(http.StatusOK, m)

	})

	//获取path参数
	//请求的参数通过URL路径传递，例如： /user/search/小王子/沙河 .获取请求URL路径中的参数的方式如下。
	r.GET("/user/search/:username/:address", func(c *gin.Context) {
		username := c.Param("username")
		address := c.Param("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})

	})

	//参数绑定
	//为了能够更方便的获取请求相关参数，提高开发效率，我们可以基于请求的Content-Type识别请求数据类型并利用反射机制自动提取
	// QueryString、form表单、JSON、XML等参数到结构体中。下面的示例代码演示了 .ShouldBuild强大的功能，它能够基于请求自动提取JSON
	// form表单 和 QueryString类型的数据，并把值绑定到指定的结构体对象。
	//Binding from JSON

	//绑定json的示例  ({"user": "q1mi", "password": "123456"})
	r.POST("/loginJSON", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err == nil {
			fmt.Printf("login info:%#v\n", login)
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	//绑定form表单示例(user=q1mi&password=123456)
	r.POST("/loginForm", func(c *gin.Context) {
		var login Login
		//ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	//绑定QueryString示例 (/loginQuery?user=q1mi&password=123456)
	r.GET("/loginForm", func(c *gin.Context) {
		var login Login
		//ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

	})

	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.sogo.com")

	})
	/*
		ShouldBind会按照下面的顺序解析请求中的数据完成绑定：
		如果是 GET 请求，只使用 Form 绑定引擎（query）。
		如果是 POST 请求，首先检查 content-type 是否为 JSON 或 XML，然后再使用 Form（form-data）。
	*/
	/*
		r.POST("/upload", func(c *gin.Context) {
			//单个文件上传
			file,err := c.FormFile("f1")
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{
					"message":err.Error(),
				})
				return
			}
			log.Println(file.Filename)
			dst := fmt.Sprintf("C:\\Users\\10922\\go\\src\\mygo\\%s",file.Filename)
			//上传文件到指定目录
			c.SaveUploadedFile(file,dst)
			c.JSON(http.StatusOK,gin.H{
				"message":fmt.Sprintf("'%s' upload!",file.Filename),
			})

		})
	*/
	//多个文件上传
	r.MaxMultipartMemory = 8 << 20
	r.POST("/uploads", StaCost(), func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["f1"]
		for index, file := range files {
			log.Println(file.Filename)
			dst := fmt.Sprintf("C:\\Users\\10922\\go\\src\\mygo\\%s_%d", file.Filename, index)
			c.SaveUploadedFile(file, dst)
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("'%s' upload!", file.Filename),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%d' upload!", len(files)),
		})

	})

	r.POST("/auth",authHandler)

	r.GET("/home", JWTAuthMiddleware(), homeHandler)



	r.Run()

	//新建一个没有任何默认中间件的路由
	r1 := gin.New()
	//注册一个全局中间件
	r1.Use(StaCost())
	r1.GET("/test1", func(c *gin.Context) {
		name := c.MustGet("name").(string) //从上下文中取值
		log.Println(name)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})
	//r1.Run()

	server01 := &http.Server{
		Addr:              ":8080",
		Handler:           router01(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	server02 := &http.Server{
		Addr:              ":9090",
		Handler:           router02(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	//借助errgroup.Group或者自行开启两个 goroutine分别启动两个服务
	g.Go(func() error {
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	//ro := setupRouter()
	//if err := ro.Run(); err != nil {
	//	fmt.Println("startup service failed,err:%v\n", err)
	//}
}

//生成JWT和解析JWT
//我们在这里直接使用JWT-go 这个库来实现我们生成JWT和解析JWT的功能

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//定义JWT的过期时间，这里以2小时为例
const TokenExpireDuration = time.Hour * 2
//接下来还需要定义Secret：
var MySecret = []byte("夏天夏天悄悄过去")
//生成JWT
//GenToken生成JWT
func GenToken(username string) (string,error)  {
	//创建一个我们自己使用的声明
	c:=MyClaims{
		Username:       username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt:time.Now().Add(TokenExpireDuration).Unix(),  //过期时间
			Issuer: "my-project",
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodES256,c)
	//使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

//解析JWT
//ParseToken解析JWT
func ParseToken(tokenString string) (*MyClaims,error)  {
	//解析token
	token,err := jwt.ParseWithClaims(tokenString,&MyClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return MySecret,nil
	})

	if err != nil {
		return nil, err
	}
	if claims,ok := token.Claims.(*MyClaims);ok && token.Valid {
		return claims,nil
	}
	return nil, errors.New("invalid token")
	
}





//gin中间件
//gin框架允许开发者在处理请求的过程中，加入用户自己的钩子函数，这个钩子函数就叫做中间件，中间件适合处理一些公共的业务逻辑
//比如登录认证、权限校验、数据分页、记录日志、耗时统计等。
//定义中间件
//gin中的中间件必须是一个 gin.HandlerFunc 类型。例如我们像下面的代码一样统计请求耗时的中间件。

func StaCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		//可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
		c.Set("name", "小王子")
		//调用剩余该请求的剩余处理程序
		c.Next()
		cost := time.Since(start)
		log.Println(cost)
	}

}

//运行多个服务
//我们可以在多个端口启动服务，例如：

var (
	g errgroup.Group
)

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)

	})
	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "welcome server02",
			},
		)
	})
	return e

}

func authHandler(c *gin.Context)  {
	//用户发送用户名和密码过来
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"code":2001,
			"msg":"无效的参数",
		})
		return
	}
	//校验用户名和密码是否正确
	if user.Username == "qimi" && user.Password == "qimi123" {
		//生成token
		tokenString,_ := GenToken(user.Username)
		c.JSON(http.StatusOK,gin.H{
			"code":2000,
			"msg":"success",
			"data":gin.H{"token":tokenString},
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":2002,
		"msg":"鉴权失败",
	})
	return
}
//用户通过上面的接口获取Token之后，后续就会携带着Token再来请求我们的其他接口，这个时候就需要对这些请求的Token进行校验操作了，
// 很显然我们应该实现一个检验Token的中间件，具体实现如下：
// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}


func homeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}

type UserInfo struct {
	Username string
	Password string
}