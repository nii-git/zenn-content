---
title: "【Golang】同一docker-composeファイル内にあるMySQLサーバーに接続する方法"
emoji: "🐬"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["docker","go","dockercompose","mysql"]
published: false
---

## TL;DR
同一`docker-compose.yml`ファイルにあるgoサーバからmysqlサーバーに接続させるには、2点気をつける必要がある
1. dbのaddressは`docker-compose.yml`内のmysqlのservice名にする必要がある(下記例だと`mysqlService`)
2. mysqlサーバーが立ち上がるまでリトライ処理を入れる


```yml:docker-compose.yml
services:
  mysqlService:
    image: mysql:8.2.0
```

## はじめに

docker-composeを用いて、golangのAPIサーバーとmysqlを接続する際に2点詰まったところがあるので備忘録として残しておきます。

## dbのaddressは`docker-compose.yml`内のmysqlのservice名にする必要がある

下記のdocker-compose.ymlを例に、go言語で動く`api`コンテナから`mysqlContainer`に接続するケースを考えます。

```yml:docker-compose.yml
version: '3'

services:
  mysqlContainer:
    image: mysql:8.2.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_DATABASE: sample_db
      MYSQL_PASSWORD: password
      MYSQL_ROOT_PASSWORD: password
  api:
    build:
      context: .
    ports:
      - "1323:1323"
```

`api`コンテナでは、`go-mysql-driver`を用いて下記のように接続します。

```go:api
db, err = sql.Open("mysql","[root[:password]@][tcp[(localhost:3306)]]/sample_db")
```

ここで`address`にlocalhost:3306と記載すると、`dial tcp 127.0.0.1:3306: connect: connection refused`エラーが発生します。

これは`mysql`コンテナがlocalhost(≒127.0.0.1)ではないことが原因です。
`docker inspect`コマンドで`mysql`コンテナを調べてみたところ、`172.19.0.2`であることがわかりました（ここの値は実行するたびに書き変わるようです）。

```txt
$ docker inspect {mysqlContainerID}
{
	...
	"Gateway": "172.19.0.1",
	"IPAddress": "172.19.0.2",
	"IPPrefixLen": 16,
	...
}

```

`addresss`にdocker-compose.ymlファイル内で記載したmysqlのservice名を指定すると、動的に`mysql`コンテナのIPに書き換えてくれます。

今回の例ですと、`address = mysqlContainer`となります。書き換えて実行すると、接続先が`172.19.0.3`と書き変わっています。

```go:api
db, err = sql.Open("mysql","[root[:password]@][tcp[(mysqlContainer:3306)]]/sample_db")
```

```txt:実行結果
dial tcp 172.19.0.3:3306: connect: connection refused
```

接続先IPは書き換わりましたが、mysqlサーバーの準備ができる前に接続を行なっているためにエラーが発生しています。その際の対応を次項で説明します。


## mysqlサーバーが立ち上がるまでリトライ処理を入れる
docker-composeには起動順序を指定できる[depends_on](https://docs.docker.jp/v1.11/compose/compose-file.html#compose-file-depends-on)オプションがありますが、あくまでも起動順序を指定できるだけであり、コンテナの準備ができているかは保証しません。

したがって、下記のようにmysqlコンテナが先に立ち上がったとしてもエラーが発生してしまう可能性があります。

![connectionエラー発生例](/images/others/docker-mysql-connections.png)


そこで、アプリケーション側でリトライ処理を実装するように修正しました。
最大試行回数を決めておき、接続に失敗する度にSleepをするシンプルな仕組みです。

```go: retry処理の例
	for r := 1; r <= MaxDBRetryCount; r++ {
		println("NewDB Connection Attept #" + strconv.Itoa(r))
		
		// 前項の例ではDSNをstringで表していましたが、本項ではFormatDSNを使用しています
		// https://pkg.go.dev/github.com/go-sql-driver/mysql#Config.FormatDSN
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


これでmysqlコンテナの実行が立ち上がるまでapiサーバーが待つようになりました。

![connection解消例](/images/others/docker-mysql-connections2.png)

## 実行結果
下記のコードを用いて実行してみます。
https://github.com/nii-git/zenn-content/tree/main/sample/docker-mysql-connections

 API接続1回目では失敗していますが、接続先が`172.21.0.2:3306`となっており動的にIPが書き換えられていることがわかります。

```txt:実行結果(1)
$ docker-compose up
[+] Building 4.2s (10/10) FINISHED       
...
Attaching to docker-mysql-connections-api-1, docker-mysql-connections-mysqlContainer-1
docker-mysql-connections-api-1             | NewDB Connection Attept #1
docker-mysql-connections-api-1             | NewDB Connection Error:dial tcp 172.21.0.2:3306: connect: connection refused
```

mysqlサーバーが起動しましたが、まだ準備は完了していません。
その間にAPIサーバーが２度目の接続を試みますが、当然拒否されてしまいます。

```txt:実行結果(2)
docker-mysql-connections-mysqlContainer-1  | 2023-11-19 05:09:37+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.2.0-1.el8 started.
...
docker-mysql-connections-api-1             | NewDB Connection Attept #2
docker-mysql-connections-api-1             | NewDB Connection Error:dial tcp 172.21.0.2:3306: connect: connection refused
```

mysqlサーバーがsample_dbデータベースを作成し、その後`ready for connections`メッセージを出力しました。
その直後にAPIサーバーが10回目の接続を試みた際、無事成功しています。

```txt:実行結果(3)
docker-mysql-connections-mysqlContainer-1  | 2023-11-19 05:09:43+00:00 [Note] [Entrypoint]: Creating database sample_db
...
docker-mysql-connections-api-1             | NewDB Connection Attept #9
docker-mysql-connections-api-1             | NewDB Connection Error:dial tcp 172.21.0.2:3306: connect: connection refused
...
docker-mysql-connections-mysqlContainer-1  | 2023-11-19T05:09:45.621827Z 0 [System] [MY-010931] [Server] /usr/sbin/mysqld: ready for connections. Version: '8.2.0'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server - GPL.
docker-mysql-connections-api-1             | NewDB Connection Attept #10
docker-mysql-connections-api-1             | NewDB Connection Succeeded!
docker-mysql-connections-api-1 exited with code 0
```

