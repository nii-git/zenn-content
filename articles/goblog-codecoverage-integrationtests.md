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

<!-- ここのbutは何を表すのか？ -->

さて、前のスクリプトを呼び出してカバレッジのためのツールをビルドし、その結果のプロファイルを後処理する別のラッパースクリプトを実装していきましょう。

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
# カバレッジのためのビルドをするスクリプトに"-cover"を渡し、GOCOVERDIRをセットして実行します
# Pass in "-cover" to the script to build for coverage, then
# run with GOCOVERDIR set.
#
GOCOVERDIR=covdatafiles \
  /bin/sh integration_test.sh -cover $PKGARGS
#
# 結果の後処理を行います
# Post-process the resulting profiles.
#
go tool covdata percent -i=covdatafiles
$
```

::: message
invoke: 呼び出す。callとの違いは[この記事](https://qiita.com/Ted-HM/items/7dde25dcffae4cdc7923#%E3%82%A4%E3%83%99%E3%83%B3%E3%83%88)がわかりやすい
post-process: 後処理
:::

---

ラッパーに関して、いくつか重要な点があります。
- `integration_test.sh`実行時に"-cover"を渡すことで、カバレッジを測定できる`mdtool.exe`バイナリを作成した
- カバレッジデータが書かれるファイル名をGOCOVERDIR環境変数にセットした
- テストが完了した時、プログラムの網羅度合いについてのレポートを生み出すために`go tool covdata percent`を実行した 　

> Some key things to note about the wrapper above:

> - it passes in the “-cover” flag when running integration_test.sh, which gives us a coverage-instrumented “mdtool.exe” binary
> - it sets the GOCOVERDIR environment variable to a directory into which coverage data files will be written
> - when the test is complete, it runs “go tool covdata percent” to produce a report on percentage of statements covered


---

新しいラッパースクリプトを実行した結果はこの通りです。

> Here’s the output when we run this new wrapper script:

```sh
$ /bin/sh wrap_test_for_coverage.sh
...
    gitlab.com/golang-commonmark/mdtool coverage: 48.1% of statements
# 注: covdatafiles内には381ファイルあります
# Note: covdatafiles now contains 381 files.
```
---

じゃじゃーん！これで`mdtool`アプリケーションのソースコードが動作する結合試験がどの程度うまく動作しているか知ることができました。

もしテストハーネスの効率を高めるために変更するのであれば、2つ目のカバレッジコレクションを実行させることでカバレッジレポートに変化が反映されていることを確認することができるでしょう。例えば、`integration_test.sh`に2行追加するとします。

> Voila! We now have some idea of how well our integration tests work in exercising the “mdtool” application’s source code.
> If we make changes to enhance the test harness, then do a second coverage collection run, we’ll see the changes reflected in the coverage report. For example, suppose we improve our test by adding these two additional lines to integration_test.sh:

```sh
./mdtool.exe +ty testdata/README.md  > /dev/null
./mdtool.exe +ta < testdata/README.md  > /dev/null
```

::: message
Voila: 〈フランス語〉ほら、じゃじゃーん、（もう）出来上がり◆完成を披露するときなどに、相手の注意を引くために言う。本来aに重アクセント記号が付く。([英辞郎 on the Web](https://eow.alc.co.jp/search?q=voila#:~:text=%E3%80%88%E3%83%95%E3%83%A9%E3%83%B3%E3%82%B9%E8%AA%9E%E3%80%89%E3%81%BB%E3%82%89%E3%80%81%E3%81%98%E3%82%83%E3%81%98%E3%82%83%E3%83%BC,%E3%81%B0%E3%80%81%E3%81%BB%E3%82%89%E5%87%BA%E6%9D%A5%E4%B8%8A%E3%81%8C%E3%82%8A%E3%80%82)より引用)
test harness: ソフトウェアテストで用いられるテスト実行用のソフトウェアのこと。([ITmedia エンタープライズ](https://www.itmedia.co.jp/im/articles/1111/07/news185.html)から引用)
:::

---

再度カバレッジテストを実行します。

網羅率が48%から54%に増えていて、変更の効果が現れていることがわかります。

> Running the coverage testing wrapper again:
> We can see the effects of our change: statement coverage has increased from 48% to 54%.

```sh
$ /bin/sh wrap_test_for_coverage.sh
finished processing 380 files, no crashes
    gitlab.com/golang-commonmark/mdtool coverage: 54.6% of statements
$
```

### カバーするパッケージを選択する

デフォルトでは、

`go build -cover`はビルドするGoモジュールの一部のパッケージを測定します。今回だと `gitlab.com/golang-commonmark/mdtool`のみ測定されます。

しかし、他のパッケージもカバレッジを測定できた方が便利な場合もあるでしょう。その場合、`go build -cover`に`-coverpkg`を渡してあげることで実現することができます。

> By default, “go build -cover” will instrument just the packages that are part of the Go module being built, which in this case is the gitlab.com/golang-commonmark/mdtool package.
> In some cases however it is useful to extend coverage instrumentation to other packages; this can be accomplished by passing “-coverpkg” to “go build -cover”.

::: message
accomplish: (動)成し遂げる、果たす
:::


---

<!--ここも翻訳がイマイチ instrumentのいい訳し方教えてください-->

今回使用したサンプルプログラムでは、`mdtool` は実は単に`gitlab.com/golang-commonmark/markdown`パッケージをラッパーしていただけにすぎません。なので、実装されているパッケージ群にマークダウンを含んでみるのも面白いかもしれません。

下記は`mdtool`のための`go.mod`です。

> For our example program, “mdtool” is in fact largely just a wrapper around the package gitlab.com/golang-commonmark/markdown, so it is interesting to include markdown in the set of packages that are instrumented.
> Here’s the go.mod file for “mdtool”:

```md
$ head go.mod
module gitlab.com/golang-commonmark/mdtool

go 1.17

require (
    github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
    gitlab.com/golang-commonmark/markdown v0.0.0-20211110145824-bf3e522c626a
)
```

---

//todo

> We can use the “-coverpkg” flag to control which packages are selected for inclusion in the coverage analysis to include one of the deps above. Here’s an example:

```sh
$ /bin/sh wrap_test_for_coverage.sh -coverpkg=gitlab.com/golang-commonmark/markdown,gitlab.com/golang-commonmark/mdtool
...
    gitlab.com/golang-commonmark/markdown   coverage: 70.6% of statements
    gitlab.com/golang-commonmark/mdtool coverage: 54.6% of statements
$
```


::: message

:::

---

>

::: message

:::
