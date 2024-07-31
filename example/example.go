package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chinmayrelkar/dogpool"
	"github.com/chinmayrelkar/dogpool/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Execute(appContext context.Context, task dao.Task) error {
	fmt.Println("Executing ", task.Name, " with args ", task.Args, " at ", time.Now(), " id ", task.ID)
	return nil
}

func main() {
	gormDb, gormErr := gorm.Open(
		mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN:        os.Getenv("MYSQL_DSN"),
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	if gormErr != nil {
		panic(gormErr)
	}
	dbConfig, err := gormDb.DB()
	if err != nil {
		panic(err)
	}
	dbConfig.SetMaxIdleConns(10)
	dbConfig.SetMaxOpenConns(100)
	dbConfig.SetConnMaxLifetime(100)

	worker := dogpool.NewWorker(context.TODO(), gormDb)
	worker.Register("ExampleTask", Execute)
	go worker.Run()

	dogpool.NewScheduler(gormDb).ScheduleTask("ExampleTask", map[string]string{
		"hello": "world",
	})
	time.Sleep(10 * 6 * time.Second)
	worker.Exit()
}
