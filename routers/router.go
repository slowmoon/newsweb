package routers

import (
	"github.com/astaxie/beego/context"
	"newsweb/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.InsertFilter("/article/*", beego.BeforeExec, p)
	beego.InsertFilter("/index", beego.BeforeExec, p)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:SaveRegister")
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:Login")
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")
	beego.Router("/index", &controllers.IndexController{}, "get:ShowIndex")
	beego.Router("/article/add", &controllers.ArticleController{}, "get:ShowAdd;post:HandleAdd")
	beego.Router("/article/view", &controllers.ArticleController{}, "get:ShowDetail")
	beego.Router("/article/edit", &controllers.ArticleController{}, "get:ShowEdit;post:EditArticle")
	beego.Router("/article/del", &controllers.ArticleController{}, "get:DelArticle")
	beego.Router("/article/addType", &controllers.ArticleController{}, "get:ShowTypeIndex;post:AddType")
	beego.Router("/article/delType", &controllers.ArticleController{} , "get:DelArticleType")

	beego.Router("/redis", &controllers.RedisController{}, "get:ShowSql")
}

var p = func(ctx *context.Context){
	userName := ctx.Input.Session("userName")
	name ,ok := userName.(string)  
	if name=="" || !ok{
		//未登录，跳转登录
		ctx.Redirect( 302, "/login")
		return
	}

}
