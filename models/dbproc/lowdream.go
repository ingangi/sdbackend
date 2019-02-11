package dbproc

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var LowDreamORM orm.Ormer
var sqlCache SqlCache

func InitSDSql() {
	user := beego.AppConfig.String("dsdb::user")
	passwd := beego.AppConfig.String("dsdb::passwd")
	host := beego.AppConfig.String("dsdb::urls")
	port, err := beego.AppConfig.Int("dsdb::port")
	dbname := beego.AppConfig.String("dsdb::dbname")
	if nil != err {
		port = 3306
	}
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	beego.Info("init mysql ...", user, passwd, host, port, dbname)
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname))
	err = orm.RegisterDataBase("low_dream", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname))
	if err != nil {
		beego.Error("init mysql db error.")
		return
	}

	beego.Info("init mysql ok")

	LowDreamORM = orm.NewOrm()
	LowDreamORM.Using("low_dream")
	sqlCache.init()
}

func SelectObjListByMainType(mainType int) (rows []ObjRow) {

	beego.Info("SelectObjListByMainType ...", mainType)

	stable := beego.AppConfig.String("dsdb::tbname")
	sql := ""
	if mainType > 0 {
		sql = fmt.Sprintf(`SELECT * FROM %s where main_type=%d`, stable, mainType)
	} else {
		sql = fmt.Sprintf(`SELECT * FROM %s`, stable) //for all
	}

	cache, exist := sqlCache.getObjListCache(sql)
	if exist {
		rows = cache
		beego.Info("hit cache for sql:", sql)
		return
	}

	if nil == LowDreamORM {
		beego.Error("SelectObjListByMainType failed: db not connected")
		return nil
	}
	num, err := LowDreamORM.Raw(sql).QueryRows(&rows)
	if err == nil {
		beego.Info(sql, "get item nums:", num)
		sqlCache.setObjListCache(sql, rows)
	}
	return rows
}

func StartBabyRecord(t int8) {
	beego.Info("StartBabyRecord ...", t)

	stable := beego.AppConfig.String("dsdb::tbbaby")
	sql := fmt.Sprintf(`insert into %s (type) values(%d)`, stable, t)

	if nil == LowDreamORM {
		beego.Error("StartBabyRecord failed: db not connected")
		return
	}
	_, err := LowDreamORM.Raw(sql).Exec()
	if err != nil {
		beego.Info("StartBabyRecord ... err:", err.Error())
	}
}

func StopBabyRecord(t int8) {
	beego.Info("StopBabyRecord ...", t)

	stable := beego.AppConfig.String("dsdb::tbbaby")
	// UPDATE t_baby_rcd SET stop_time=NOW(), cost_seconds=UNIX_TIMESTAMP(NOW())-UNIX_TIMESTAMP(start_time), state=1 WHERE state=0 and type=1 ORDER BY id desc LIMIT 1;
	sql := fmt.Sprintf(`UPDATE %s SET stop_time=NOW(), cost_seconds=UNIX_TIMESTAMP(NOW())-UNIX_TIMESTAMP(start_time), state=1 WHERE state=0 and type=%d ORDER BY id desc LIMIT 1;`, stable, t)

	if nil == LowDreamORM {
		beego.Error("StopBabyRecord failed: db not connected")
		return
	}
	_, err := LowDreamORM.Raw(sql).Exec()
	if err != nil {
		beego.Info("StopBabyRecord ... err:", err.Error())
	}
}

func SelectBabyRecordOfToday(t int8) (rows []BabyRow) {
	beego.Info("SelectBabyRecordOfToday ...", t)

	stable := beego.AppConfig.String("dsdb::tbbaby")

	// return last 15 items is better
	//sql := fmt.Sprintf(`SELECT * FROM %s where type=%d and start_time>='%s' ORDER BY id desc`, stable, t, time.Now().Format("2006-01-02")+" 00:00:00")
	sql := fmt.Sprintf(`SELECT * FROM %s where type=%d ORDER BY id desc limit 15`, stable, t)

	if nil == LowDreamORM {
		beego.Error("SelectBabyRecordOfToday failed: db not connected")
		return nil
	}
	num, err := LowDreamORM.Raw(sql).QueryRows(&rows)
	beego.Info(sql, "SelectBabyRecordOfToday get item nums:", num)
	if err != nil {
		beego.Info("SelectBabyRecordOfToday ... err:", err.Error())
	}
	return rows
}
