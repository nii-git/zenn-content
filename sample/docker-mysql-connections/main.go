package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

func main() {

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		println("LoadLocation Error:" + err.Error())
		os.Exit(-1)
	}

	// サンプルなのでパスワード等をハードコードしているが、
	// 環境変数から取得するのが望ましい
	c := mysql.Config{
		DBName:    "sample_db",
		User:      "root",
		Passwd:    "password",
		Addr:      "mysqlContainer:3306",
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	// わかりやすいように1秒毎に再接続しているが、
	// 実際は接続頻度を減らしたほうが良い
	MaxDBRetryCount := 30
	SleepTime := 1 * time.Second

	var db *sql.DB

	for r := 1; r <= MaxDBRetryCount; r++ {
		println("NewDB Connection Attept #" + strconv.Itoa(r))
		db, err = sql.Open("mysql", c.FormatDSN())
		if err != nil {
			println("NewDB Connection Error:" + err.Error())
			time.Sleep(SleepTime)
			continue
		}

		err = db.Ping()
		if err != nil {
			println("NewDB Connection Error:" + err.Error())
			time.Sleep(SleepTime)
			continue
		} else {
			break
		}
	}

	if err != nil {
		println("NewDB Connection Failed")
		os.Exit(-1)
	}

	println("NewDB Connection Succeeded!")
	os.Exit(0)
}
