// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"com.github.bider/idGenerator"
	"com.github.bider/util"
	"fmt"
	log2 "github.com/kataras/golog"
	"github.com/kataras/iris"
	"net"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log2.SetFormatter(&log2.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log2.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log2.SetLevel(log2.WarnLevel)
}
func main() {

	//c, _, err := zk.Connect([]string{"10.1.62.22"}, time.Second) //*10)

	//if err != nil {
	//	panic(err)
	//}

	//a,b,f:=c.Get("/service/productdata/1")
	//	fmt.Println(a,b,f)
	//	children, stat, ch, err := c.ChildrenW("/")
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("%+v %+v\n", children, stat)
	//	e := <-ch
	//	fmt.Printf("%+v\n", e)

	app := iris.Default()
	ip := util.GetIp()
	log2.Info("ip:", ip)

	app.Get("/idworker/{module:string max(20)}", func(ctx iris.Context) {

		module := ctx.Params().GetStringDefault("module", "")
		if (module == "") {
			ctx.JSON(iris.Map{"error": 1, "msg": "业务模块参数不正确",})
			return;
		}
		idWorker, _ := idGenerator.GetIdWorker(module)

		newid, err := idWorker.Id();
		if err != nil {
			ctx.JSON(iris.Map{"error": 1, "msg": err.Error()})
			return
		}
		ctx.JSON(iris.Map{"error": 0, "data": newid})

	})
	app.Get("/idworker/add/{module:string max(20)}", func(ctx iris.Context) {
		token := ctx.URLParamDefault("token", "")

		if token != "ec887868-2738-414D-a952-a418f31B9937" {
			ctx.JSON(iris.Map{"error": 1, "msg": "token不正确"})
			return
		}

		module := ctx.Params().GetStringDefault("module", "")
		if (module == "") {
			ctx.JSON(iris.Map{"error": 1, "msg": "业务模块参数不正确",})
			return;
		}

		err := idGenerator.AddModule(module)
		if err != nil {
			ctx.JSON(iris.Map{"error": 1, "msg": err.Error(),})
			return;
		}

		ctx.JSON(iris.Map{"error": 0, "data": "56"})
	})
	app.Run(iris.Addr(":8080"))
}

func mac() {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	for _, inter := range interfaces {
		fmt.Println(inter.Name)
		mac := inter.HardwareAddr //获取本机MAC地址
		fmt.Println("MAC = ", mac)
	}
}
