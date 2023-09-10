---
title: "Pythonで大文字と小文字の変換と速度比較"
emoji: "💻"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["Python"]
published: false
---


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
"12345".islower()  # False
"a12345".islower() # True
```


### 小文字があるか確認
str型のisupper()を使用します。
大文字と小文字が区別可能な文字がすべて大文字で、かつ大文字が1文字以上あればTrueを返します。
1文字でも小文字が入っていた場合や、区別不能な文字列のみ(数字等)はFalseとなります。

https://docs.python.org/ja/3/library/stdtypes.html#str.isupper

```python
"PYTHON".isupper()  # True
"Python".isupper()  # False
"12345".isupper()  # False
"A12345".isupper() # True
```


## 配列内の大文字小文字変換
### 配列の文字列を全て大文字にしたい場合
リストの内包表記と、上記のupper()を使用します。

```python3
words = ["apple","banana","orange"]
upper_words = [w.upper() for w in words]
print(upper_words) # ['APPLE', 'BANANA', 'ORANGE']
```

### 配列の文字列を全て小文字にしたい場合
リストの内包表記と、上記のlower()を使用します。

```python3
words = ['APPLE', 'BANANA', 'ORANGE']
lower_words = [w.lower() for w in words]
print(lower_words) # ['apple', 'banana', 'orange']
```


## 特定の文字群を全て小文字にしたい時の速度比較

英文は基本的に文の先頭の単語1文字目を大文字とするが、全て小文字に変換したい場合があります。
この際2通りのアプローチが考えられますが、どちらの方が効率的なのでしょうか。実際に試してみます。

`
1. 全ての文字列にlower()を実行する
2. islower()でFalseになった文字列のみlower()を実行する
`

### 前提条件
テストで使用するデータは[BBC](https://www.bbc.com/)から約750記事を抽出し、単語毎に分割したものを使用します。

```python3
wordList = ['I', 'did', 'it', 'adding', '"I', 'got', 'my', "mom's", 'gun', 'last', 'night"', '', 'shortly', 'after', 'the', 'shooting', '', 'The', 'boy', 'has', 'not', 'been', 'charged', 'Taylor', 'had', 'pleaded', 'guilty', 'in', 'June', 'to', 'a', 'federal', 'charge', 'of', 'using', 'marijuana', 'while', 'possessing', 'a', 'firearm', '', 'The', 'sentencing', 'in', 'that', 'case', 'is', 'also', '', 'scheduled', 'for', 'October', 'Sign', 'up', 'for', 'our', 'morning', 'newsletter', 'and', 'get', 'BBC', 'News', 'in', 'your', 'inbox', '', ...]  # 中略

print(len(wordList)) # 564023
```

### テスト1: ニュース記事に対して実行

```python3
    # 全体に対して小文字化
    all_lower_start = time.time()
    all_lower_list = [w.lower() for w in wordList]
    all_lower_end = time.time()

    # print(all_lower_list)

    # チェックして小文字化
    check_islower_start = time.time()
    check_islower_list = [w if w.islower() else w.lower() for w in wordList]
    check_islower_end = time.time()

    # print(check_islower_list)
    print("Test1:ニュース記事で比較")
    print("all_lower:" + str(all_lower_end-all_lower_start))
    print("check_lower:" + str(check_islower_end-check_islower_start))
```





もし効率的な判定と変換方法をご存知の方がいらっしゃいましたらコメント頂けると幸いです。