---
title: "The Go Blog - ガイド付き最適化(PGO)のプレビュー"
emoji: "😃"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["go","go1.20","翻訳","PGO","最適化"]
published: false
---

# Outline
本記事はThe Go Blogの"Profile-guided optimization preview"を翻訳したものになります。
誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

https://go.dev/blog/pgo-preview

# Profile-guided optimization preview
2023年2月8日 Michael Pratt

---

Goバイナリをビルドする際、コンパイラはできるだけパフォーマンスが高いバイナリを生成しようと最適化を行います。下記は最適化の一例です。
- 定数伝播: 定数式をコンパイル時に評価し、実行時の評価時間を削減
- エスケープ解析: ローカルに割り振られたオブジェクトへのheap割り当てを回避し、ガーベジコレクションのオーバーヘッドを削減
- インライン展開: 単純な関数を呼び出し側にコピーすることで、呼び出し側の最適化(さらなる定数伝播やエスケープ解析など)を可能に

Go言語はリリースを重ねるごとに最適化を改善していますが、これは簡単なことではありません。いくつかの最適化手法は調整可能ですが、コンパイラがすべての関数の性能を飛躍的に向上させることはできません。過度な最適化はパフォーマンスを損なったり、ビルド時間が著しく長くなってしまうかもしれないからです。他の最適化手法ではコンパイラに一般的なパスとそうでないパスの判断を委ねるというものもあります。しかしこの手法では、実行時までどちらのケースが一般的なパスなのかを知ることができないため、コンパイラは最善の推測をしなくてはなりません。

もしそれができるとしたら...？

どのコードが本番環境で使われているかの決定的な情報がないので、コンパイラはパッケージのソースコードだけしか操作することができません。しかし我々には本番環境の振る舞いを評価するツールがあります。[プロファイリング](https://go.dev/doc/diagnostics#profiling)です。コンパイルにプロファイルを渡すことができたら、より賢い選択をすることができるでしょう。例えば、より使用される関数を積極的に最適化できたり、より正確に一般的なケースを選択したり、などです。

このように、アプリケーションの振る舞いのプロファイルをコンパイラで使用する手法を`Profile-Guided Optimization`,PGOと呼びます(`Feedback-Directed Optimization`,FDOとも呼ばれます)。

Go 1.20では、初めてPGOをプレビュー版としてサポートするバージョンです。完全なガイドを閲覧するためには[profile-guided optimization user guide](https://go.dev/doc/pgo)を参照してください。これらはまだ荒削りな部分が残っているので商用環境では使用しない方が良いでしょう。しかしあなたが試してみて、私たちに様々なフィードバックや[見つけたバグ](https://github.com/golang/go/issues/new/choose)を送ってくれることを心よりお待ちしております。

## 例
マークダウンをHTMLに変換してくれるサービスをビルドしてみましょう。ユーザーが`/render`にマークダウンファイルをアップロードし、それをHTMLに変換したものを返します。これらは`gitlab.com/golang-commonmark/markdown`を使用することで簡単に実装することができます。

### セットアップ
```
$ go mod init example.com/markdown
$ go get gitlab.com/golang-commonmark/markdown@bf3e522c626a
```

```go:main.go
package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    _ "net/http/pprof"

    "gitlab.com/golang-commonmark/markdown"
)

func render(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }

    src, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("error reading body: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    md := markdown.New(
        markdown.XHTMLOutput(true),
        markdown.Typographer(true),
        markdown.Linkify(true),
        markdown.Tables(true),
    )

    var buf bytes.Buffer
    if err := md.Render(&buf, src); err != nil {
        log.Printf("error converting markdown: %v", err)
        http.Error(w, "Malformed markdown", http.StatusBadRequest)
        return
    }

    if _, err := io.Copy(w, &buf); err != nil {
        log.Printf("error writing response: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func main() {
    http.HandleFunc("/render", render)
    log.Printf("Serving on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

```

ビルドし、サーバーを立ち上げます。
```
$ go build -o markdown.nopgo.exe
$ ./markdown.nopgo.exe
2023/01/19 14:26:24 Serving on port 8080...
```

他のターミナルからマークダウンファイルを送ってみましょう。下記はGoプロジェクトのREADMEのサンプルです。
```
$ curl -o README.md -L "https://raw.githubusercontent.com/golang/go/c16c2c49e2fa98ae551fc6335215fadd62d33542/README.md"
$ curl --data-binary @README.md http://localhost:8080/render
<h1>The Go Programming Language</h1>
<p>Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.</p>
...
```

### プロファイリング
これで動作するサービスができたので、プロファイルを収集してPGOを用いて再ビルドしてみましょう。パフォーマンスが向上するはずです。

`main.go`内で、CPUプロファイルを収集するため、サーバーに`/debug/pprof/profile`エンドポイントを自動的に追加する[`net/http/pprof`](https://pkg.go.dev/net/http/pprof)をインポートしました。

コンパイラが本番環境での動作を得られるために、一般的には本番環境のプロファイルを集めたいと思うはずです。この例では本番環境を持っていないので、今回はプロファイルを取得する間に負荷をかけるシンプルな動作のプログラムを作りたいと思います。下記のプログラムを`load/main.go`にコピーし、負荷をかけてください(先程実行したサーバーがまだ実行中であることを確認してください)。

https://go.dev/play/p/yYH0kfsZcpL

```
$ go run example.com/markdown/load
```

実行中にプロファイルをサーバーからダウンロードしてみましょう。
```
$ curl -o cpu.pprof "http://localhost:8080/debug/pprof/profile?seconds=30"
```

これらが完了したら、サーバーと負荷をかけるプログラムを停止させてください。

### プロファイルを使用してみる
GoツールチェインがPGOを用いてビルドするには、`-pgo`フラグを用います。これらはプロファイルのパスを指定するか、`auto`を指定することでメインパッケージのディレクトリにある`default.pgo`を使用させることができます。

リポジトリに`default.pgo`プロファイルをコミットしておくことをお勧めします。ソースコードと一緒にプロファイルを保存することで、ユーザーはリポジトリを(バージョン管理システムや、`go get`を用いることで)取得するだけで自動的にプロファイルにアクセスすることができます。また、ビルドを再現可能な状態に保つこともできます。

Go 1.20では、`-pgo`はoffがデフォルトとなっているため、ユーザーは引き続き`-pgo=auto`と指定する必要があります。しかし、将来のアップデートでは`-pgo`はautoがデフォルトとなる予定なので、誰でもPGOの恩恵を自動的に得られるようになります。

それでは、PGO有りでビルドしてみましょう。
```
$ mv cpu.pprof default.pgo
$ go build -pgo=auto -o markdown.withpgo.exe
```

### 評価
さて、PGOにおけるパフォーマンスの効果がどの程度なのかを確認するために、負荷生成器のベンチマークバージョンを使用してみましょう。

まずは、PGOなしでサーバーを測定してみましょう。サーバーを起動します。
```
$ ./markdown.nopgo.exe
```

サーバーが実行している間、何回かベンチマークを実行します。
```
$ go test example.com/markdown/load -bench=. -count=20 -source ../README.md > nopgo.txt
```

測定が完了したら、pgo無しのサーバーを止めてpgo有りのサーバーを起動します。
```
$ ./markdown.withpgo.exe
```

同様に、サーバーが実行している間何回かベンチマークを実行します。
```
$ go test example.com/markdown/load -bench=. -count=20 -source ../README.md > withpgo.txt
```

ベンチマークが完了したら、早速結果を比較してみましょう、
```
$ go install golang.org/x/perf/cmd/benchstat@latest
$ benchstat nopgo.txt withpgo.txt
goos: linux
goarch: amd64
pkg: example.com/markdown/load
cpu: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz
        │  nopgo.txt  │            withpgo.txt             │
        │   sec/op    │   sec/op     vs base               │
Load-12   393.8µ ± 1%   383.6µ ± 1%  -2.59% (p=0.000 n=20)
```

PGO有りの場合、2.6%程度早くなっています！Go1.20では、PGOを有効にすることでCPU使用率が2~4%改善します。プロファイルはアプリケーションの振る舞いにおける様々な有用な情報を提供してくれますが、本バージョンではこの情報の一部を用いてインライン化を行っているだけに過ぎません。将来のアップデートでは、より多くのPGOのメリットを享受できるように改善し、パフォーマンスの改善を行い続けていきます。


## 次のステップ
この例では、プロファイルを収集した後全く同じソースコードを用いて際ビルドしています。現実世界では開発は常に進んでいるため、先週のコードの本番環境からプロファイルを集め、それを本日のソースコードに使用してビルドするかもしれません。このようなケースでも、全く問題ありません！PGOは小さい変更点を問題なく扱うことができるのです。

PGOのより詳細な情報やベストプラクティス、注意点などは(profile-guided optimization user guide)[https://go.dev/doc/pgo]をご覧ください。

フィードバックをお待ちしています！PGOはまだプレビュー版であり、使用しづらい点や正しく動作しない点等のご意見を聞けることを楽しみにしています。こちらのリンクから起票を行なってください。

https://go.dev/issue/new

# 関連記事
Go 1.20リリース記事を翻訳した記事もあります。こちらもどうぞ(宣伝)

https://zenn.dev/nii/articles/read-go-1-20-release-article
https://zenn.dev/nii/articles/read-go-1-20-release-article2