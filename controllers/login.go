package controllers

import (
	"github.com/astaxie/beego/orm"
	"newsweb/models"
	"github.com/astaxie/beego"
)


type LoginController struct{
	beego.Controller
}

func (l *LoginController)ShowLogin(){
	
	  l.Data["userName"] = l.Ctx.GetCookie("userName")
	  l.Data["remember"] = l.Ctx.GetCookie("remember")
      l.TplName = "login.html"		
}

func (l *LoginController)Logout(){
	l.DelSession("userName")
	l.Redirect("/login", 302)
}

//post 
func(l *LoginController)Login(){
	 var user models.User
	 userName := l.GetString("userName")
	 password := l.GetString("password")
	 remember := l.GetString("remember")

	if userName=="" || password =="" {
		l.Data["message"] = "用户名或者密码不能为空"
	    l.TplName = "login.html"
		return
	}
	o := orm.NewOrm()

	user.UserName = userName
	err := o.Read(&user, "userName")
	if err !=nil{
		l.Data["message"] = "用户名不存在"
		l.TplName = "login.html"
		return
	}
	if user.Password != password{
		l.Data["message"] = "密码出错"
		l.TplName = "login.html"
		return
	}
	l.Data["message"] = "登录成功"

	l.Ctx.SetCookie("userName", userName, 3600*24)
	if remember == ""{
	    beego.Info("not set remember me!")
		l.Ctx.SetCookie("userName", userName, -1)
		l.Ctx.SetCookie("remember", remember, -1)
	}else if remember == "on"{
		l.Ctx.SetCookie("userName", userName, 3600*24)
		l.Ctx.SetCookie("remember", remember, 3600*24)
	}
	l.SetSession("userName", userName)
	l.Redirect("/index", 302)
}