---
title: "test"
emoji: "😃"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["test"]
published: false
---


## Outline
本記事はThe Go Blogの"Code coverage for Go integration tests"を翻訳したものになります。
誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

https://go.dev/blog/integration-test-coverage

## Code coverage for Go integration tests

Than McIntosh
2023年3月8日

コードカバレッジツールは、テストスイートが実行された際にソースコードの割合がどれだけカバーされ、実行されているか、開発者の判断の支えになります。

Go言語では、Go1.2リリースの際にパッケージレベルのコードカバレッジ測定機能を導入を初め、何度かのアップデートを行なっています。これらは`go test`コマンドを実行する際、`-cover`フラグをつけることで実行できます。


> Code coverage tools help developers determine what fraction of a source code base is executed (covered) when a given test suite is executed.
> Go has for some time provided support (introduced in the Go 1.2 release) to measure code coverage at the package level, using the "-cover" flag of the “go test” command.

::: message
fraction of ~: ~の割合
test suite: テストスイートはテストの目的や条件が似ている複数のテストケースを一括りにしたもの( [Qmedia!](https://q-media.jp/test-suite)より引用)
:::

---

これらのツールはほとんどのケースで上手く動作するのですが、大きなGoアプリのカバレッジを測定するのは苦手です。このようなアプリケーションでは、技術者はパッケージレベルの結合テストに加え、統合テストを記載することで全体のプログラムの振る舞いを検証していたことでしょう。

> This tooling works well in most cases, but has some weaknesses for larger Go applications. For such applications, developers often write “integration” tests that verify the behavior of an entire program (in addition to package-level unit tests).



::: message
integration test: 総合テスト、統合テスト
:::

---

このようなテストは、一般的にはアプリケーション全体のバイナリを必要とします。独立したパッケージを単独にテストするのとは反対に、構成する全てのパッケージが正しく動作しているかを保証するために全体のバイナリを代表的な一入力値（サーバーなら負荷テスト等）を用いて実行します。

> This type of test typically involves building a complete application binary, then running the binary on a set of representative inputs (or under production load, if it is a server) to ensure that all of the component packages are working correctly together, as opposed to testing individual packages in isolation.

::: message
involve: (動)必要とする、巻き込む
representative: (形)代表的な、典型的な
isolation: (名)孤独、分離
:::

---

統合テストのバイナリは`go test`ではなく`go build`でビルドされているので、今まではGoツールを用いてこれらのテストのカバレッジプロファイルを集めるのは簡単ではありませんでした。

> Because the integration test binaries are built with “go build” and not “go test”, Go’s tooling didn’t provide any easy way to collect a coverage profile for these tests, up until now.

::: message
up until now: 今までは
:::

---

Go 1.20では、`go build -cover`を用いることでカバレッジ測定プログラムを使用が可能になりました。プログラム全体のバイナリを統合テストに渡すことによって、カバレッジテストの範囲が広がります。

この記事では新機能がどのようにして動作するかの例と、統合テストからカバレッジプロファイルを収集する方法の概要を説明していきます。


> With Go 1.20, you can now build coverage-instrumented programs using “go build -cover”, then feed these instrumented binaries into an integration test to extend the scope of coverage testing.
> In this blog post we’ll give an example of how these new features work, and outline some of the use cases and workflow for collecting coverage profiles from integration tests.

### 例

それでは、非常に小さなサンプルプログラムを使い、シンプルな総合テストを書き、それのカバレッジプロファイルを収集してみましょう。
マークダウン処理ツールである`mdtool`を用いて演習を進めていきます(URLは下記です)。

これはマークダウンからHTMLに変換するライブラリである`gitlab.com/golang-commonmark/markdown`のデモプログラムです。

https://pkg.go.dev/gitlab.com/golang-commonmark/mdtool

> Let’s take a very small example program, write a simple integration test for it, and then collect a coverage profile from the integration test.
> For this exercise we’ll use the “mdtool” markdown processing tool from gitlab.com/golang-commonmark/mdtool. This is a demo program designed to show how clients can use the package gitlab.com/golang-commonmark/markdown, a markdown-to-HTML conversion library.

### mdtoolのセットアップ

まず、`metool`のコピーをダウンロードしてみましょう（再現性を高めるために特定のバージョンを使用しています）。

> First let’s download a copy of “mdtool” itself (we’re picking a specific version just to make these steps reproducible):

```sh
$ git clone https://gitlab.com/golang-commonmark/mdtool.git
...
$ cd mdtool
$ git tag example e210a4502a825ef7205691395804eefce536a02f
$ git checkout example
...
```


::: message
reproducible: (形)再現可能な
:::

### シンプルな結合テスト

さて、これから`mdtool`のシンプルな結合テストを書いていきます。これから書くテストは`mdtool`バイナリをビルドした後、一連のマークダウンファイルを入力値として取り実行します。

この非常にシンプルなスクリプトはテストデータディレクトリにあるそれぞれのファイルに対し`mdtool`バイナリを実行し、outputの検証とプログラムがクラッシュしていないことを確認します。


> Now we’ll write a simple integration test for “mdtool”; our test will build the “mdtool” binary, then run it on a set of input markdown files. 
> This very simple script runs the “mdtool” binary on each file from a test data directory, checking to make sure that it produces some output and doesn’t crash.

```sh
$ cat integration_test.sh
#!/bin/sh
BUILDARGS="$*"
#
# 下記コマンドが正常に終了しない場合はテストを終了させます
# Terminate the test if any command below does not complete successfully.
#
set -e
#
# テスト用入力をダウンロードします(`website`リポジトリには複数のmarkdownファイルが含まれます)
# Download some test inputs (the 'website' repo contains various *.md files).
#
if [ ! -d testdata ]; then
  git clone https://go.googlesource.com/website testdata
  git -C testdata tag example 8bb4a56901ae3b427039d490207a99b48245de2c
  git -C testdata checkout example
fi
#
# テスト用にmdtoolバイナリをビルドします
# Build mdtool binary for testing purposes.
#
rm -f mdtool.exe
go build $BUILDARGS -o mdtool.exe .
#
# `testdata`を入力値としてツールを実行します
# Run the tool on a set of input files from 'testdata'.
#
FILES=$(find testdata -name "*.md" -print)
N=$(echo $FILES | wc -w)
for F in $FILES
do
  ./mdtool.exe +x +a $F > /dev/null
done
echo "finished processing $N files, no crashes"
$
```

実行結果の例は下記の通りです。

```sh
$ /bin/sh integration_test.sh
...
finished processing 380 files, no crashes
$

```

---

成功です。`mdtool`バイナリが正しくインプットファイルを処理したことを検証できました。しかし、これによってどのくらいのソースが実行されたのでしょうか？

次のセッションでは、ソースの実行率を解明するためにカバレッジプロファイルを集めてみましょう。

> Success: we’ve verified that the “mdtool” binary successfully digested a set of input files… but how much of the tool’s source code have we actually exercised? 
> In the next section we’ll collect a coverage profile to find out.

::: message
exercise: (動)働かせる
find out: 解く、見破る
:::

### カバレッジデータを集めるために統合テストを使用する



> Let’s write another wrapper script that invokes the previous script, but builds the tool for coverage and then post-processes the resulting profiles:

```sh
$ cat wrap_test_for_coverage.sh
#!/bin/sh
set -e
PKGARGS="$*"
#
# Setup
#
rm -rf covdatafiles
mkdir covdatafiles
#
# Pass in "-cover" to the script to build for coverage, then
# run with GOCOVERDIR set.
#
GOCOVERDIR=covdatafiles \
  /bin/sh integration_test.sh -cover $PKGARGS
#
# Post-process the resulting profiles.
#
go tool covdata percent -i=covdatafiles
$
```

::: message
invoke: 呼び出す。callとの違いは[この記事](https://qiita.com/Ted-HM/items/7dde25dcffae4cdc7923#%E3%82%A4%E3%83%99%E3%83%B3%E3%83%88)がわかりやすい
:::

---
>

::: message

:::

