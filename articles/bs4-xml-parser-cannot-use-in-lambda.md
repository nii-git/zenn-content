---
title: "BeautifuleSoupのxmlパーサーがLambdaで使えない"
emoji: "🤨"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["AWS","Lambda","Python3","BeautifulSoup","xml"]
published: false
---

# 事象
Lambdaにて、BeautifulSoupを用いてxmlをパースしようとしたところ、エラーが発生した

```python
    soup = BeautifulSoup(input_text,"xml")
```

```
"errorMessage": "Couldn't find a tree builder with the features you requested: xml. Do you need to install a parser library?",
  "errorType": "FeatureNotFound",
```

アップロードは下記を参考にzip形式で行った。
https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/python-package.html


# 試したこと

パーサーを`html.parser`に変更すると、xml形式をパースできなくなる。
対応策として、lxmlライブラリをzipに含めても解消せず。

```sh
$ pip3 install --target ./package lxml
```

# 解決法

xmltodictライブラリを使いましょう。
xmlを辞書型に変換してくれるツールです。

https://pypi.org/project/xmltodict/


```sh
$ pip3 install --target ./package xtodict
Collecting xmltodict
  # ローカルのpipでもインストールしたためcacheになっています
  Using cached xmltodict-0.13.0-py2.py3-none-any.whl (10.0 kB)
Installing collected packages: xmltodict
Successfully installed xmltodict-0.13.0
$ cd package
$ zip -r ../my_deployment_package.zip .
$ cd ../
$ auto_webscraping % zip my_deployment_package.zip lambda_function.py
```

## サンプルコード

下記はBBCニュースRSS.xmlから、記事一覧のurlを取得するコードの例です。

### 旧コード
```python
    resultLinks = []
    selectRule = 'link'

    # この行でエラーが発生
    soup = BeautifulSoup(req.text,"xml")

    links = soup.select(selectRule)

    for l in links:
        resultLinks.append(l.get_text())

    return resultLinks
```

### 新コード
```python
    resultLinks = []

    dict_xml_data = xmltodict.parse(req.text)
    items = dict_xml_data["rss"]["channel"]["item"]

    for i in items:
        resultLinks.append(i["link"])

    return resultLinks
```

