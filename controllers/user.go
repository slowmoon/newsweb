package controllers

import (
	"newsweb/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)


type UserController struct{
	beego.Controller
}

func (u *UserController)ShowRegister(){
	 u.TplName = "register.html"
}

func (u *UserController)SaveRegister(){
	  var  user  models.User
	  user.UserName = u.GetString("userName")
	  user.Password = u.GetString("password")
	 if user.UserName != "" && user.Password!= "" {
 		beego.Info("receive user register: ", user)
		  o := orm.NewOrm()
	 	 i, _ := o.Insert(&user)
	 	 beego.Info(i)
		 u.Redirect("/login", 302)
		 return	
	}
	u.TplName = "register.html"
}