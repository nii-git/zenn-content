---
title: "【Golang】同一docker-composeファイル内にあるMySQLサーバーに接続する方法"
emoji: "🐬"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["docker","go","dockercompose","mysql"]
published: false
---

## TL;DR
同一`docker-compose.yml`ファイルにあるgoサーバからmysqlサーバーに接続させるには、2点気をつける必要がある
1. mysqlサーバーが立ち上がるまでリトライ処理を入れる
2. db接続先は`docker-compose.yml`内のmysqlのservice名にする必要がある(下記例だと`mysqlService`)


```yml:docker-compose.yml
services:
  mysqlService:
    image: mysql:8.2.0
```

## mysqlサーバーが立ち上がるまでリトライ処理を入れる
docker-composeには起動順序を指定できる[depends_on](https://docs.docker.jp/v1.11/compose/compose-file.html#compose-file-depends-on)オプションがありますが、あくまでも起動順序を指定できるだけであり、コンテナの準備ができているかは保証しません。

したがって、下記のようにmysqlコンテナが先に立ち上がったとしてもエラーが発生してしまう可能性があります。

![connectionエラー発生例](/images/others/docker-mysql-connections.png)


そこで、アプリケーション側でリトライ処理を実装するように修正しました。
最大試行回数を決めておき、接続に失敗する度にSleepをするシンプルな仕組みです。

```go: retry処理の例
	for r := 1; r < MaxDBRetryCount; r++ {
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
```


## サンプルコード