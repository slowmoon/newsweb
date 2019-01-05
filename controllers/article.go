package controllers

import (
	"github.com/astaxie/beego/orm"
	"newsweb/models"
	"github.com/astaxie/beego"
	"path/filepath"
	"time"
	"strings"
	"errors"
	"mime/multipart"
)

const IMAGE_DIR = "/static/img/"


type ArticleController struct{
	beego.Controller
}

func (a *ArticleController)ShowAdd(){
	o:= orm.NewOrm()
	var articles []models.ArticleType
	qs := o.QueryTable("article_type")
	i, _ := qs.All(&articles)

	a.Data["articleTypes"] = articles[:i]
	a.Layout = "layout.html"
	a.TplName = "add.html"
}

func (a *ArticleController)HandleAdd(){
	  title := a.GetString("articleName")
	  content := a.GetString("content")
	  articleTypeId, err:= a.GetInt64("select")
		if err != nil{
			beego.Error("传输数据失败", err)
			a.Data["error"] = "添加文件失败，请稍后重试"
			a.TplName = "add.html"
			return
		}

	  file, header , err := a.GetFile("uploadname")
	  if title=="" || content ==""  || err != nil {
		  beego.Error("传输数据失败",err)
		  a.Data["error"] = "添加文件失败，请稍后重试"
		  a.TplName = "add.html"
		  return
		}
	  defer file.Close()
	 
	  path, err := buildPath(file, header)
	  if err != nil {
		  beego.Error(err)
		  a.Data["error"] = "添加文件失败，请稍后重试"
		  a.TplName = "index.html"
		  return 
	  }
	  err = a.SaveToFile("uploadname", "."+ path)
	  if err != nil {
		  beego.Error(err)
		  a.Data["error"] = "添加文件失败，请稍后重试"
		  a.TplName = "index.html"
		  return 
	  }

	  o := orm.NewOrm()
	  var articleType models.ArticleType
	  articleType.Id = articleTypeId
	  err = o.Read(&articleType)
	  if err != nil{
		  beego.Error(err)
		  a.Data["error"] = "文件类型出错"
		  a.TplName = "index.html"
		  return 
	  }

	  var article models.Article
	  
	  article.Title = title
	  article.Content = content
	  article.Image = path
	  article.ArticleType = &articleType

	  beego.Info(article.Title, article.Content, article.Image)
	 _, err =  o.Insert(&article)
	
	if err != nil {
		beego.Error("保存数据出错", err)
		a.Data["error"] = "保存数据库异常"
		a.TplName = "add.html"
		return
	}
	  a.Redirect("/index", 302)
}

func (a *ArticleController)ShowDetail(){
		

		id, err := a.GetInt("id")
		beego.Info("id", id)
		if err!=nil{
			beego.Info("查询文章详情出错", err)
			a.Data["error"] = "查询文章详情出错"
			a.TplName = "index.html"
			return
		}
		o := orm.NewOrm()
		o.Begin()
		defer o.Commit()
		
		var article models.Article
		article.Id = int64(id)

		err = o.QueryTable("article").RelatedSel("ArticleType").Filter("Id", article.Id).One(&article)
		
		if err != nil{
			beego.Error("read article fail", err)
			a.Data["error"] = "查询出错"
			a.TplName = "index.html"
			return
		}
		article.Count++
		o.Update(&article)
		
		qm := o.QueryM2M(&article, "Users")
		userName := a.GetSession("userName")
		var user models.User
		var ok bool
		 if user.UserName, ok = userName.(string);!ok{
				beego.Error("error")
				o.Rollback()
				a.Redirect("/login", 302)
				return
		 }

		err= o.Read(&user, "UserName")
		if err != nil{
			o.Rollback()
			beego.Error(err)
			a.Redirect("/index", 302)
		}
		_, err = qm.Add(user)
		if err != nil{
			o.Rollback()
			beego.Error(err)
			a.Redirect("/index", 302)
		}
		//o.LoadRelated(&article, "Users")
		
		var users []*models.User
		_, err = o.QueryTable("User").Distinct().Filter("Articles__Article__Id", id).All(&users)

		if err != nil{
			beego.Error("error related select ....", err)
			o.Rollback()
			a.Redirect("/index", 302)
			return
		}
		a.Data["article"] = article
		a.Data["users"] = users

		a.Layout = "layout.html"
		a.TplName = "content.html"
}


func(a *ArticleController)ShowEdit(){
	 id , err := a.GetInt("id")
	 if err != nil {
		beego.Error("展示编辑出错!")
		a.Data["error"] = "展示编辑出错"
		a.Redirect("/index", 302)
		return
	 }
	 o := orm.NewOrm()
	 var article models.Article
	 article.Id = int64(id)
	 err = o.Read(&article)
	 if err != nil{
		 beego.Error("读取数据出错")
		 a.Data["error"] = "读取数据出错"
		 a.Redirect("/index", 302)
		 return
	 }
	 a.Data["article"] = article
	 a.Layout = "layout.html"
	 a.TplName = "update.html"

}

func buildPath(file multipart.File , header *multipart.FileHeader)(string, error){
	  curTime := time.Now().Format("20060102150304")
	  extName := strings.ToLower(filepath.Ext(header.Filename))
	  beego.Info("extName is :", extName)
	 	
	  if extName != ".jpg" && extName!= ".jpeg" && extName != ".png"{
		  beego.Error("upload is not a picture!")
		  return "", errors.New("update data is not a picture")
	  }

	  size := header.Size
	  if size > 1024*1024*10{
		   beego.Error("picture is too large")
		   return "", errors.New("picture is too big")
	  }
	  path := IMAGE_DIR + curTime + extName
	  return path , nil
}


func (a *ArticleController)EditArticle(){
	 id, err := a.GetInt("id")
	 if err != nil {
		  beego.Error("传输数据失败", err)
		  a.Data["error"] = "传输数据失败"
		  a.Redirect("/index", 302)
		  return
	 }
	 title := a.GetString("title")
	 content := a.GetString("content")

	 file, header , err := a.GetFile("uploadname")
	if  title=="" || content =="" || err != nil {
		  beego.Error("传输数据失败", err)
		  a.Data["error"] = "添加文件失败，请稍后重试"
		  a.Redirect("/index", 302)
		  return
	}
	defer file.Close()
	 
	  path, err := buildPath(file, header)
	  if err != nil {
		  beego.Error(err)
		  a.Data["error"] = "添加文件失败，请稍后重试"
		  a.Redirect("/index", 302)
		  return 
	  }
	  err = a.SaveToFile("uploadname", "."+ path)
	  if err != nil {
		  beego.Error(err)
		  a.Data["error"] = "添加文件失败，请稍后重试"
		  a.Redirect("/index", 302)
		  return 
	  }
	  o := orm.NewOrm()
	  var article models.Article
	  article.Id = int64(id)
	  article.Title = title
	  article.Content = content
	  article.Image = path
	  o.Update(&article)
	  a.Redirect("/index", 302)
	 
}


func(a *ArticleController)DelArticle(){
	id, err := a.GetInt("id")
	if err != nil{
		beego.Info("删除数据失败",err)
		a.Redirect("/index", 302)
		return	
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = int64(id)
	i, err := o.Delete(&article)
	if err != nil{
		beego.Info("删除数据失败",err)
		a.Redirect("/index", 302)
		return	
	}
	beego.Info("删除成功, 删除了:",i)
	a.Redirect("/index", 302)
}

func(a *ArticleController)ShowTypeIndex(){

	o := orm.NewOrm()
	var articles []models.ArticleType
	qs := o.QueryTable("article_type")
	i, err := qs.All(&articles)
	if err != nil{
		beego.Error("查询类型出错")
		a.Redirect("/index", 302)
	}
	a.Data["articleTypes"] = articles[:i]
	a.Layout = "layout.html"
	a.TplName = "addType.html"
}


func(a *ArticleController)AddType(){
    name :=  a.GetString("typeName")
	if name == ""{
		beego.Error("articleType name should not be empty!")
		a.Redirect("/article/addType", 302)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Type = name
	err := o.Read(&articleType, "type")
	if err == nil{
		//not empty already exists
		beego.Error("read data error", err)
		a.Data["error"] = "类型已存在"
		a.TplName = "addType.html"
		return 
	}
	_, err = o.Insert(&articleType)

	if err !=nil{
		beego.Error("类型插入失败", err)
	}
	a.Redirect("/article/addType", 302 )
}

func(a *ArticleController)DelArticleType(){
	 id ,err := a.GetInt("id")
	 if err != nil{
		 beego.Error("missing parameters ","id")
		 a.Redirect("/article/addType",  302)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = int64(id)
	_, err = o.Delete(&articleType)
	if err !=nil{
		 beego.Error("read articleType error ")
		 a.Redirect("/article/addType",  302)
		return
	}
	a.Redirect("/article/addType",  302)
}