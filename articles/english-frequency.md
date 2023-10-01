---
title: "【個人開発】英字ニュースから頻繁に使用されている単語を調査した話"
emoji: "💻"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["個人開発","英語","AWS","python","lambda","stepfunctions"]
published: false
---

## 1.はじめに
皆さん、英語勉強していますか？
最近では様々なドキュメントが翻訳されたり、自動翻訳ツールが充実してきていますが、まだまだ英語が必要になる場面も多いです。

私もTOEICの勉強を始め、なんとかレベルB(730~855)のスコアに到達できました。
基礎はそこそこ完成できたと考え、少し背伸びして英字新聞を買ってみたりしました。
しかし1ヶ月くらい新聞を読んだところ、あることに気づきました。

「単語、全然わからん」と。

TOEICと実際に使用されている単語は少なからずズレているようです。
今だにTOEIC頻出単語であるwheelbarrow(手押し車)やcanopy(天蓋)を英字新聞で見たことがありません。
そこで、英字新聞に使用されている単語の出現頻度を解析してデータベースに入れるツールを作成しました。

:::message alert
英字新聞の頻出単語とTOEICの頻出単語は少しズレがあるように感じたという意味です。
TOEIC等の試験を批判する意図は一切ありません。
:::



## 2.作ったもの
自動で英字ニュースサイトをスクレイピングし、単語毎の出現回数をデータベースに入れるバッチを作成しました。
日毎やサイト毎に単語の出現数を格納しているため、様々な用途の分析に使用することができます。

下記は2023/9/18の[BBCトップニュース記事](https://feeds.bbci.co.uk/news/rss.xml)の単語出現頻度top3を取得した時の例です。
![DB](/images/english-frequency/1.png)


## 3. 全体図
全体の構成図は下記の通りです。

![全体図](/images/english-frequency/english_frequency-aws_latest.jpg)

外部英語ニュースサイトのrssから最新記事を取得するバッチ、記事を解析するバッチはそれぞれLambda(Python3.9)で作成しました。それぞれの実行順序は必ず決まっているため、StepFunctionで1つのステートマシンにまとめました。
1日1回実行して欲しいので、EventBridgeで定期実行化しています。

英字ニュースサイトから取得したデータは、記事の本文だけを抽出したファイルをS3に一時的に配置しています。バックアップ等で再度実行する際などに外部サイトへの負荷を減らす目的があります。
集計データはRDS(MySQL)に格納され、半永久的に保管されます。記事データは定期的に削除しています。

現在、集計データの確認API等を用意していません。管理者のみEC2からRDSに接続を行ってデータを確認できるようにしています。

## 4. 使用技術
全体図で説明できなかった箇所やつまづいた部分を解説します。

### 4.1 Lambda(Python3.9,3.11)
外部英語ニュースサイトのrssから最新記事を取得するバッチ、記事を解析するバッチをLambdaで作成しました。
最新記事を取得するバッチはPython3.11、記事を解析するバッチはPython3.9で動作しています。

#### 採用理由
それぞれ1日1回しか起動しないので、実行時のみ課金される点がコストを抑えられるからです。
また、責任範囲がコードのみで、ユーザーは実行サーバーの意識をしなくて良いという点もスピーディーに開発をしたい点にマッチしていると感じました。

Lambdaを採用する際の注意点として、実行時間は15分以内やメモリは10GBまでという制約がある点に注意してください。
今回の２つのバッチはメモリ128MBで十分に収まりますし、実行時間も90秒以内で完了しているので問題ありませんでした。

今後、データの収集範囲を広げた場合の影響範囲についても考える必要があります。
現在の作りでは引数にデータ取得元を指定しているため、さらに上位のStepFunctionを作成してもいいかもしれません。

![今後の構成案](/images/english-frequency/english_frequency-aws_latest.jpg)


#### Python 3.9を使用する理由
Lambda ランタイムでは、Pythonのバージョンは3.7~3.11までサポートしています(2023/10/1現在)。

:::message alert
Python 3.7は2023/11/27に非推奨となります。
https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/lambda-runtimes.html#runtime-support-policy
:::

記事を解析するバッチをランタイムをPython3.11で実行したところ、依存関係エラーが発生してしまいました。
エラーログを確認したところ、NumPyモジュールで発生しているようです。


:::details エラー文
[ERROR] Runtime.ImportModuleError: Unable to import module 'lambda_function': Unable to import required dependencies:
numpy: 

IMPORTANT: PLEASE READ THIS FOR ADVICE ON HOW TO SOLVE THIS ISSUE!

Importing the numpy C-extensions failed. This error can happen for
many reasons, often due to issues with your setup or how NumPy was
installed.

We have compiled some common reasons and troubleshooting tips at:

    https://numpy.org/devdocs/user/troubleshooting-importerror.html

Please note and check the following:

  * The Python version is: Python3.11 from "/var/lang/bin/python3.11"
  * The NumPy version is: "1.25.2"

and make sure that they are the versions you expect.
Please carefully study the documentation linked above for further help.

Original error was: No module named 'numpy.core._multiarray_umath'

Traceback (most recent call last):
[ERROR] Runtime.ImportModuleError: Unable to import module 'lambda_function': Unable to import required dependencies: numpy: IMPORTANT: PLEASE READ THIS FOR ADVICE ON HOW TO SOLVE THIS ISSUE! Importing the numpy C-extensions failed. This error can happen for many reasons, often due to issues with your setup or how NumPy was installed. We have compiled some common reasons and troubleshooting tips at: https://numpy.org/devdocs/user/troubleshooting-importerror.html Please note and check the following: * The Python version is: Python3.11 from "/var/lang/bin/python3.11" * The NumPy version is: "1.25.2" and make sure that they are the versions you expect. Please carefully study the documentation linked above for further help. Original error was: No module named 'numpy.core._multiarray_umath' Traceback (most recent call last):
:::

Numpyのバージョンは1.25.2を使用しているようですが、公式ドキュメントではPython3.11は対応していると記載されています。
バージョンの問題では無さそうです。

> NumPy 1.25.0 リリース
> 2023年1月17日 – Numpy 1.25.0 がリリースされました。 今回のリリースの目玉機能は次のとおりです。
> ...(中略)
> このリリースでサポートされている Python のバージョンは3.3.9 - 3.11 です。
> https://numpy.org/ja/news/ より引用

pomblue’s blogさんの記事によると、どうやらローカルでLambdaのzipファイルを作成したのが悪さしているようです。

[Lambda用のDocker image](https://hub.docker.com/r/amazon/aws-sam-cli-build-image-python3.9/tags)をpullし、その環境内でzipを作成したところ動作しました。
あまりに初見殺しすぎる。

> バージョンがどうのこうの、依存関係がどうのこうの言ってるぽい。
> でもバージョンは正しい。
> ローカルでモジュールzip作成やったのがいけないっぽい！
> LambdaはAmazon Linux(2)で動いているので、
> ・EC2でAmazon Linux(2)のインスタンスを立てる
> もしくは
> ・Lambdaのdocker imageを利用する
> pomblue’s blog - Lambdaで外部モジュールを使う方法 より引用 
> https://pomblue.hatenablog.com/entry/2021/06/08/230146

現在、Lambda用のDocker imageはPython3.9までしか無いようなので、仕方なくランタイムもPython3.9にしているという状況です。


#### Beautifulsoup4でxmlパーサーが使えない
またまた外部ライブラリ周りで詰まりました。

最新記事を取得するバッチでは、スクレイピングを用いて記事の内容を取得する必要があります。
その際にBeautifulsoup4というライブラリを使用したのですが、デフォルトではhtmlファイルのパーサーだけしか含まれていません。
今回はxmlファイルのパースが必要だったため、lxmlパッケージをインストールして実行しました。

```txt:エラー
Couldn't find a tree builder with the features you requested: lxml. Do you need to install a parser library?
```

[BeautifulSoup のエラー "Couldn't find a tree builder" の原因と対処法](https://qiita.com/PND/items/06e1053eeed69ec4f418) を参考にさせていただきましたが、こちらでも解決しませんでした。

どうしようもなかったので、根本的な解決にはなりませんが別ライブラリを使用しました。
同じ部分で詰まる方がいるかもしれないので、この件に関して記事を作成しました。

https://zenn.dev/nii/articles/bs4-xml-parser-cannot-use-in-lambda


#### 大文字小文字変換の速度比較
記事を解析するバッチでは、単語ごとに集計して出現回数をカウントします。
ここで、文頭にある大文字アルファベットが含まれる単語を全て小文字に変換しなくては上手く集計できなくなってしまいます。

変換方法として、下記の２種類のアプローチが考えられます。
1. 全ての文字列に小文字化を適用する
2. 小文字チェックを行い、違ったものの文字列のみ小文字化を実行する

検証した結果、手法1の片っ端から小文字化を適用する方が早いことがわかりました。
検証内容に関しては下記の記事にまとめています。
https://zenn.dev/nii/articles/lower-upper-comparison



細かいメモ等はスクラップを参照してください。
https://zenn.dev/nii/scraps/ae5358c8db16e0


## X.今後やりたいこと
現在は使用する英語のデータ収集ができた段階
データを利用するために下記をやりたい...

## memo
概要、なぜ作ろうと思ったか
  - toeicの単語でねえ、wheelbarrow,canopy

作ったもの

技術スタック（全体）
  - 全体図とかあるといいかも

技術スタック（詳細）
  - つまづいたところ
  - rss xmlパーサー

今後やりたいこと
  - これも図で説明したい

その他
  - scrap紹介
  - 参考文献

