---
title: "BeautifuleSoupのxmlパーサーがLambdaで使えない"
emoji: "🤨"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["翻訳","AWS","EC2"]
published: false
---

# 事象
Lambdaにて、BeautifulSoupを用いてxmlをパースしようとしたところ、エラーが発生した

```python3
    soup = BeautifulSoup(input_text,"xml")
```

```
"errorMessage": "Couldn't find a tree builder with the features you requested: xml. Do you need to install a parser library?",
  "errorType": "FeatureNotFound",
```

アップロードは下記を参考にzip形式で行った。
https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/python-package.html


# 対策

パーサーを`html.parser`に変更すると、xml形式をパースできなくなる。
対応策として、lxmlライブラリをzipに含めても解消せず。

```sh
$ pip3 install --target ./package lxml
```

