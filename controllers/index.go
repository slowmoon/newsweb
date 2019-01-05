package controllers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"newsweb/models"
	"math"
	"encoding/gob"
	"bytes"
)

type IndexController struct{
	beego.Controller
}

type PageInterface interface{
	 StartPage()
}

type PageController struct{
	pageSize int  `page:"pagesize";default:10`
	pageIndex int  `page:"pageIndex";default:1`
	totalCount int  `page:"totalSize"`
}

func (i *IndexController)ShowIndex(){
	/* userName := i.GetSession("userName")
	if userName==nil{
		beego.Error("user not login")
		i.Redirect("/login", 302)
	}
 */
	 articleId, err1 := i.GetInt64("articleType")
	
	 o := orm.NewOrm()
	 qs := o.QueryTable("article")

	 var articles []models.Article
	 
	 if err1 == nil{
		 qs =  qs.RelatedSel("ArticleType").Filter("ArticleType__Id", articleId)
	 }else{
		 qs = qs.RelatedSel("ArticleType")
	 }	
	
	count, err := qs.Count()

	if  err != nil {
		beego.Error("查询数据条目数出错", err)
		i.Data["error"] = "查询条目出错"
		i.TplName = "login.html"
		return
	}
	beego.Info("总记录数", count)

	pageSize := 2
	pageCount := int(math.Ceil(float64(count)/ float64(pageSize)))
	beego.Info("count=", pageCount, "size=", pageSize)
	index , err:= i.GetInt("index")   //获取页码
	
	if err!= nil {
		index = 1
	}

	t, err := qs.Limit( pageSize,(index-1)*pageSize).All(&articles)
	
	if err != nil {
		beego.Error("查询数据出错", err)
		i.Data["error"] = "查询出错"
		i.TplName = "login.html"
		return
	} 
	
	//查询缓存
	
	var articleTypes []models.ArticleType
	conn,err := redis.Dial("tcp", ":6379")
	if err != nil{
		beego.Error("redis 获取失败:", err)
	}else{
		//判断存不存在
		defer conn.Close()
		reply, err := conn.Do("get", "articleTypes")
		if err != nil || reply == nil{
			beego.Error("缓存读取失败。。。", err)
			//从mysql数据库中取数据
			qs2 := o.QueryTable("article_type")
			qs2.All(&articleTypes)
			var buf bytes.Buffer
			gob.NewEncoder(&buf).Encode(&articleTypes)
			conn.Do("set", "articleTypes", buf.Bytes())   //保存
		}else{
			ret ,err := redis.Bytes(reply, err)
			if err!=nil{
				beego.Error(err)
			}
			err = gob.NewDecoder(bytes.NewReader(ret)).Decode(&articleTypes)
			beego.Info(articleTypes)
		}
	}
	
	i.Data["count"] = count
	i.Data["pageCount"] = pageCount
	i.Data["index"] = index   
	i.Data["articles"] = articles[:t]
	i.Data["articleTypes"] = articleTypes[:]
	if err1 == nil{
		i.Data["articleTypeId"] = articleId
	}
	i.Layout = "layout.html"
	i.TplName = "index.html"
}