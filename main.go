package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "newsweb/routers"
	"github.com/astaxie/beego"
	_ "newsweb/models"
)

func main() {
	beego.AddFuncMap("pre", getPre)
	beego.AddFuncMap("next", next)
	beego.Run()
}


func getPre(index int)int{
	 if index <=1 {
		return 1
	 }else{
		 return index-1
	 }
}

func next(index, max int)int{
	if index >= max{
		return max
	}else{
		return index+1
	}
}