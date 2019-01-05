package controllers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
	"fmt"
)

type RedisController struct{
	beego.Controller
}

func (r *RedisController)ShowSql(){
	 conn, err := redis.Dial("tcp", ":6379")
	 defer conn.Close()
	 if err != nil{
		beego.Error(err)
		return 
	 }
	 reply, err := conn.Do("mget", "m1", "m2")
	var s1 int 
	var s2 string
	 t ,err := redis.Values(reply, err)
	 if err != nil{
		 fmt.Println(err)
		 return
	 }
	 redis.Scan(t, &s1, &s2)
	 beego.Info(t)
	 r.Ctx.WriteString(fmt.Sprintf("%d%s", s1, s2))
}







func init(){
}