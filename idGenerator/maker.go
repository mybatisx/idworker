package idGenerator

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log2 "github.com/kataras/golog"


	"log"
	"os"
	"sync"
	"time"
)

type IdWorker struct {
	key string
	queue chan uint64
	running bool
	msg string
}
type IdWorkerDB struct {
	Stub int `db:"stub"`
	Id int `db:"id"`
}
var (
	workerMap sync.Map

)


func assert(err error) {
	if err != nil {
		log.Panicln(err)
	}
}



func GetIdWorker(key string) (*IdWorker, error) {

	if v, ok := workerMap.Load(key); ok {
		return v.(*IdWorker),nil
	}



	//_, err = strconv.ParseInt(string(kv.Value), 10, 64)

	//assert(err)

	ai := &IdWorker{
		key:key,
		running: true,
		queue: make(chan uint64, 98),
	}

	go ai.process()

	ai2,_:=workerMap.LoadOrStore(key,ai)
	ai= ai2.(*IdWorker)
	return ai,nil
}

func (iw *IdWorker) process() {

//	defer func() {recover()}()


	var mysqldb *sqlx.DB
	var err error

	closeDB := func() {
		mysqldb.Close()
	}



	for  {


		log2.Infof("进来了0,%s,%t",iw.key,iw.running)

		if iw.running==false {
			return
		}
		mysqldb,err = sqlx.Connect("mysql", "root:Windows2000!@(10.1.62.230:3306)/shihang?charset=utf8")
		if err != nil {
			fmt.Printf("Disabling MySQL tests:\n    %v", err)
			closeDB()
			os.Exit(0)
		}

		//ret:=mysqldb.MustExec(`REPLACE INTO id_worker (stub) VALUES (?);`,iw.key)
		//last_id,err:=ret.LastInsertId();

		//row := mysqldb.QueryRow("SELECT id FROM `id_worker` where stub =?", "test17")
		//var id int64
		//err = row.Scan(&id)

		var id int
		err = mysqldb.Get(&id, "SELECT count(*) FROM `id_worker` where bizmod =?", iw.key)

		if err != nil {
			iw.running=false
			iw.msg=err.Error()
			closeDB()
			return
		}
		if id ==0 {
			iw.running=false
			iw.msg=fmt.Sprintf("%s 存根不存在%d",iw.key,id)
			closeDB()
			return
		}
		if id !=1 {
			//err=fmt.Errorf("%s 多条存根 %n",iw.key,id)
			iw.running=false
			iw.msg=fmt.Sprintf("%s 多条存根 %d",iw.key,id)
			closeDB()
			return
		}

		 idWorkerDB :=IdWorkerDB{}
		err = mysqldb.Get(&idWorkerDB, "SELECT id,stub FROM `id_worker` where bizmod =? LIMIT 1", iw.key)

		tx := mysqldb.MustBegin()

		ret:=tx.MustExec("update `id_worker` set stub=stub+1 where id=? and stub=? and bizmod=?",idWorkerDB.Id,idWorkerDB.Stub,iw.key)

		affected,_:= ret.RowsAffected()
		var last_id uint64
		if affected==1 {
			err = tx.Commit()
			err = mysqldb.Get(&last_id, "SELECT stub FROM `id_worker` where bizmod =? and id=? and stub=?", iw.key,idWorkerDB.Id,idWorkerDB.Stub+1)
		} else{

			tx.Rollback()
		}

		closeDB()

		if last_id== 0{
			log2.Errorf("进来了5,不应该为0,%s,%d",iw.key,last_id)
			continue
		}

		 startId :=last_id*100
		//startId, err:= strconv.ParseInt(string(nextV), 10, 64)
		endId := startId+99
		log2.Infof("进来了1,开始塞,%d,%d",startId,endId)

		for i := startId; i <= endId; i++ {
			iw.queue <- i
		}
		log2.Infof("进来了2，塞完了,%d,%d",startId,endId)
	}





}

func (iw *IdWorker) Id() (uint64,error) {
	//return <-iw.queue

	select {
	case newid := <-iw.queue:
		return newid,nil
	case <-time.After(5 * time.Second):
		return 0,errors.New(iw.msg)
	}

}

func (iw *IdWorker) Close() {
	iw.running = false
	close(iw.queue)
}