package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"sdbackend/models/dbproc"
	//"strings"

	//"encoding/json"
	//"strings"
	//"prototcp/typedefs"
)

type MainController struct {
	beego.Controller
}

//func (c *MainController) Get() {
//	c.TplName = "index.tpl"
//
//	beego.Info(c.Data)
//}

func (c *MainController) ObjList() {
	//c.TplName = "Test.tpl"
	beego.Info(c.Data)
	//sqlstr := c.GetString("sql")
	//date := strings.Replace(c.GetString("date"), "-", "", -1)
	//beego.Info(date, sqlstr)

	result := dbproc.SelectObjListByMainType(dbproc.MyAtoi(c.GetString("maintype")))

	data, err := json.Marshal(&result)
	if err != nil {
		beego.Info(err)
	}
	//beego.Info("data...", string(data))
	c.Ctx.WriteString(string(data))
}
