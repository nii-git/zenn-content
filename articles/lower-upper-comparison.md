---
title: "Pythonで大文字と小文字の変換と速度比較"
emoji: "🐍"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["Python"]
published: false
---

## 概要
Pythonの文字列で、"Python"のような文字列を"python”に変換したいことがあります。
本記事では、大文字や小文字を判定する方法や変換する方法をまとめました。
また、大量の文字列に対して変換する際の速度比較も実施しました。


## 変換方法
### 大文字から小文字に変換
str型のlower()を使用します。

https://docs.python.org/ja/3/library/stdtypes.html#str.lower

```python
upper_str = "THIS IS UPPERCASE"
to_lower_str = upper_str.lower()

print(to_lower_str)  # 'this is uppercase'
```


### 小文字から大文字に変換
str型のupper()を使用します。

https://docs.python.org/ja/3/library/stdtypes.html#str.upper

```python
lower_str = "this is lowercase"
to_upper_str = lower_str.upper()

print(to_upper_str)  # 'THIS IS LOWERCASE'
```


## 判別方法
### 大文字があるか確認
str型のislower()を使用します。
大文字と小文字が区別可能な文字がすべて小文字で、かつ小文字が1文字以上あればTrueを返します。
1文字でも大文字が入っていた場合や、区別不能な文字列のみ(数字等)はFalseとなります。

https://docs.python.org/ja/3/library/stdtypes.html#str.islower

```python
"python".islower()  # True
"Python".islower()  # False
"12345".islower()   # False
"a12345".islower()  # True
```


### 小文字があるか確認
str型のisupper()を使用します。
大文字と小文字が区別可能な文字がすべて大文字で、かつ大文字が1文字以上あればTrueを返します。
1文字でも小文字が入っていた場合や、区別不能な文字列のみ(数字等)はFalseとなります。

https://docs.python.org/ja/3/library/stdtypes.html#str.isupper

```python
"PYTHON".isupper()  # True
"Python".isupper()  # False
"12345".isupper()   # False
"A12345".isupper()  # True
```


## 配列内の大文字小文字変換
### 配列の文字列を全て大文字にしたい場合
リストの内包表記と、上記のupper()を使用します。

```python
words = ["apple","banana","orange"]
upper_words = [w.upper() for w in words]
print(upper_words) # ['APPLE', 'BANANA', 'ORANGE']
```

### 配列の文字列を全て小文字にしたい場合
リストの内包表記と、上記のlower()を使用します。

```python
words = ['APPLE', 'BANANA', 'ORANGE']
lower_words = [w.lower() for w in words]
print(lower_words) # ['apple', 'banana', 'orange']
```


## 特定の文字群を全て小文字にしたい時の速度比較

英文は基本的に文の先頭の単語1文字目を大文字としますが、全て小文字に変換したい場合があります。
この際2通りのアプローチが考えられますが、どちらの方が効率的なのでしょうか。色々なケースに対して試してみます。

```
A. 全ての文字列にlower()を実行する
B. islower()でFalseになった文字列のみlower()を実行する
```


### テスト1: ニュース記事に対して実行

一般的な英ニュース記事を全て小文字にする場合を検討してみます。
テストで使用するデータは[BBC](https://www.bbc.com/)から約750記事を抽出し、単語毎に分割したものを使用します。

```python:準備
wordList = ['I', 'did', 'it', 'adding', '"I', 'got', 'my', "mom's", 'gun', 'last', 'night"', ...]  # 中略

print(len(wordList)) # 564023
```


```python:実行例
# アプローチA: 全ての文字列にlower()を実行する
all_lower_start = time.time()
all_lower_list = [w.lower() for w in wordList]
all_lower_end = time.time()

# アプローチB: islower()でFalseになった文字列のみlower()を実行する
check_islower_start = time.time()
check_islower_list = [w if w.islower() else w.lower() for w in wordList]
check_islower_end = time.time()

print("Test1:ニュース記事で比較")
print("all_lower:" + str(all_lower_end-all_lower_start))
print("check_lower:" + str(check_islower_end-check_islower_start))
```

```text:実行結果
Test1:ニュース記事で比較
all_lower:0.024682044982910156
check_lower:0.03107905387878418
```

`アプローチA:全ての文字列にlower()を実行する方法`の方が約18%早いという結果になりました。



### テスト2: 文字列が全て大文字の場合

テスト1で使用した、BBCニュース750記事を一旦全て大文字に変換したデータではどうなるでしょうか。
アプローチBの方法では、全ての文字列に対してislower()とlower()を実行するので遅くなりそうだと予想できます。

```python:準備
# wordListを全て大文字化
wordList = [w.upper() for w in wordList]
```

実行例は出力以外テスト1と同様なので省略します。

```text:実行結果
Test2: 全て大文字だった場合の比較
all_lower:0.029858827590942383
check_lower:0.04039716720581055
```

`アプローチA:全ての文字列にlower()を実行する方法`の方が約26%早いという結果になりました。
テスト1より差が大きくなっています。

### テスト3: 文字列が全て小文字の場合

次は、テスト1で使用したBBCニュース750記事を一旦全て小文字に変換したデータでテストします。
アプローチBの方法では全ての文字列にislower()のみを実行するため、アプローチAより早くなるのではないでしょうか。

実行例は出力以外テスト1と同様なので省略します。

```text:実行結果
Test3: 全て小文字だった場合の比較
all_lower:0.030844926834106445
check_lower:0.032996177673339844
```

予想とは反して、`アプローチA:全ての文字列にlower()を実行する方法`の方が約6%早いという結果になりました。
どちらも全ての要素に対してlower()またはislower()を実行している為、islower()の方がlower()に比べて時間的な処理コストが大きいのかもしれません。


## 速度比較まとめ

実行結果をまとめたものが下記の表です。
数値の単位は**ミリ秒**です。小数点第3位を四捨五入しています。

| テストケース | 全て小文字にする | 大文字判定後小文字にする |
| ---- | ---- | ---- | 
| ニュース記事 | 24.68 | 31.07|
| 全て大文字 | 29.86 | 40.40 |
| 全て小文字 | 30.84 | 33.00 |


全てのパターンで`アプローチA:全ての文字列にlower()を実行する方法`の方が早いという結果になりました。
どのような文字列に対しても、判定せずに片っ端からlower()を実行した方が良いようです。

もし効率的な判定と変換方法をご存知の方がいらっしゃいましたらコメント頂けると幸いです。