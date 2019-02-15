package idGenerator

import (
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/core/errors"
)

func AddModule(mod string) error  {

	mysqldb,err := sqlx.Connect("mysql", "root:Windows2000!@(10.1.62.230:3306)/shihang?charset=utf8")

	if err != nil  {
		return err
	}
	defer func() {mysqldb.Close()}()



	var id int
	err = mysqldb.Get(&id, "SELECT count(*) FROM `id_worker` where bizmod =?", mod)
	if(id>0){
		return  errors.New("该模块已经存在")
	}


	ret:=mysqldb.MustExec(`insert id_worker (bizmod,stub) VALUES (?,1);`,mod)
	last_id,err:=ret.LastInsertId();
	if(last_id>0){
		return nil
	}

	return  errors.New("未知错误")

}
