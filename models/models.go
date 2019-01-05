package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)


type User struct{
	Id int64  
	UserName string  `orm:"unique"`
	Password string
	Articles []*Article  `orm:"reverse(many)"`
}
 

type Article struct{
	Id int64 `orm:"pk;auto;"`
	Title string  `orm:"size(100)"`
	Time time.Time `orm:"type(datetime);auto_now"`
	Count int32 `orm:"default(0)"`
	Image string `orm:"null"`
	Content string `orm:size(1024)`
	ArticleType *ArticleType `orm:"rel(fk);set_null;null"`
	Users []*User `orm:"rel(m2m)"`
}

//price   orm:"digits(19);decimals(4)"
type ArticleType struct{
	Id int64
	Type string `orm:size(20)`
	Articles []*Article `orm:"reverse(many)"`
}


//host.docker.internal    docker 内部使用
func init(){
	orm.RegisterDataBase("default", "mysql", "root:244121@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}



