package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 路由重定向，请求转发，ANY ，NoRoute，路由组
func main() {
	r := gin.Default()

	// ---------------路由重定向(到其他网址)---------------------
	r.GET("/index", func(c *gin.Context) {
		//c.IndentedJSON(200,gin.H{
		// "status":"ok",
		//})
		//  重定向到另一个地址
		c.Redirect(http.StatusMovedPermanently, "https://dashen.tech")
	})

	// ---------------转发(到其他接口)---------------------
	r.GET("/a", func(c *gin.Context) {
		// 跳转到 /b对应的路由； 地址栏的URL不会变，是请求转发  (而上面是**重定向**到https://dashen.tech，地址栏的地址改变了)
		c.Request.URL.Path = "/b" //把请求的URL修改
		r.HandleContext(c)        //继续后续的处理
	})

	r.GET("/b", func(context *gin.Context) {
		context.IndentedJSON(200, gin.H{
			"status": "这是b路径，请求成功！",
		})
	})

	// ---------------Any: 任意http method(适用于restful，因为有同名的不同http method的方法)---------------------
	// 路径为/test的，不管是GET还是POST或者PUT等，均路由到这里
	r.Any("/test", func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			c.IndentedJSON(200, gin.H{
				"info": "api为/test，GET方式~",
			})

		case "POST":
			c.IndentedJSON(200, gin.H{
				"info": "api为/test，POST方式~",
			})

		}
	})

	// --------------------NoRoute: 当匹配不到路径时到这里(其实gin默认有个404，取代之)--------------------
	r.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"msg": "任何没有匹配到的路径都会走到这里...(可以重定向到 寻找走丢儿童的网站)",
		})
	})

	// // ----------------Group: 路由组--------------------
	//把公用的前缀提取出来，创建一个路由组
	userGroup := r.Group("/user")
	{ //通常加一个大括号 让代码看起来更有调理

		userGroup.GET("/getName", func(context *gin.Context) {
			//等价于 /user/getname
			context.IndentedJSON(http.StatusOK, gin.H{
				"name": "张三",
			})
		})
		userGroup.GET("/getGender", func(context *gin.Context) {
			context.IndentedJSON(200, gin.H{
				"gender": "男",
			})
		})

		userGroup.POST("/updateAvatar", func(context *gin.Context) {
			context.IndentedJSON(http.StatusOK, gin.H{
				"msg": "头像更新成功",
			})
		})

		// 可以嵌套路由组
		// 即 /user/address
		addressGroup := userGroup.Group("/address")

		// 即/user/address/city
		addressGroup.GET("/city", func(c *gin.Context) {
			c.IndentedJSON(200, gin.H{
				"city": "阿姆斯特丹",
			})
		})
	}

	fmt.Println("路由规则初始化完毕")
	r.Run()

	//go func() {
	// r.Run()
	//}()

	//select {}

}
