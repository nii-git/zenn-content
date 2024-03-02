---
title: "【型アノテーション】Pythonの矢印は何者か"
emoji: "🤔"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["python","python3","annotations"]
published: false
---


## この記事は何
Pythonのコードを読んでいたら急に矢印(->)が現れてびっくりした人向けの記事です。

つまりPython初学者向けです。
備忘録でもあります。

## TL;DR
Pythonでの -> 表記は型アノテーションを表しています。
変数や関数の引数や戻り値の**型のヒント**を定義できます。

型を定義するのではなくあくまで型のヒントを記述しているだけ、というところがポイントです。
これに違反していてもコンパイルエラーにならないことに注意してください。

## 実際に触ってみる
2つの値を足すだけの簡単な関数を定義してみましょう。

```python
def Addition(x, y):
    return x + y
```

これらを呼び出す際は、Addition関数に2つのint型を渡してあげればうまく動きます。

```python
result = Addition(2, 3)
print(result) # 5
```

実装した自分自身であれば、Addition関数にはint型を渡さなくてはいけないこともわかります。
ただ、他の人が誤って使ってしまったらどうなるでしょうか？

```python
invalid_result = Addtion(“2”,”3”)
print(invalid_result) # 23
```

Addition関数にはint型以外も入れることができてしまい、意図しない結果になってしまっています。
これらはコンパイルエラーとならないため、バグに繋がってしまう可能性があります。

:::message
stringは+演算子を用いることで文字列を結合することができます。
詳細: https://docs.python.org/ja/3.13/tutorial/introduction.html
:::

関数や変数の型制約を示すことができたら上記のようなバグは減る[1]と思いませんか？
そこで -> (型アノテーション)の登場です。

``` python
def Addition(x:int,y:int)->int:
    return x+y
```

上記のように書くことで、Additionの使い方に関して2点わかります。
・引数x,yは2つともint型である
・返り値もint型である

しかし、このままではpythonファイルは実行できてしまいます。
どのように型エラーを検知するのでしょうか。

## 型定義エラーを検知する方法
### Visual studio
拡張機能のPylanceを使います。
![Pylanceをインストール](/images/others/what-is-python-annotations.png)

インストール後、設定から`Python › Analysis: Type Checking Mode` をbasic以上に設定してください。
![Pylance設定](/images/others/what-is-python-annotations2.png)

VSCode再起動後、pythonファイルを確認すると適応されています。
![Pylance設定](/images/others/what-is-python-annotations3.png)



### mypy
Pythonの型チェックを行ってくれるツールがあります。
https://github.com/python/mypy

`インストール`
```
$ pip3 install mypy
$ mypy --version     
mypy 1.8.0 (compiled: yes)
```

`実行`
```
$ mypy sample/annotations-test.py 
sample/annotations-test.py:9: error: Argument 1 to "Addition" has incompatible type "str"; expected "int"  [arg-type]
sample/annotations-test.py:9: error: Argument 2 to "Addition" has incompatible type "str"; expected "int"  [arg-type]
```

## どのような型が使えるか
importを使用しないのであれば、Pythonの組み込み型が使用できます。

詳細は[Pythonドキュメント公式](https://docs.python.org/ja/3/library/stdtypes.html#bitwise-operations-on-integer-types) を参考にしてください。

| 型名 | 表記 | 表示例 | 
| ---- | ---- | ---- |
| 真偽型 | bool | True,False |
| 文字列型 | str | "Hello","world" |
| 整数型 | int | 0,1 |
| 浮動小数点数型 | float | 1.0,3.14 |
| 複素数型 | complex | (1+2j) |
| バイト型 | bytes | b'Hello' |
| バイト配列型 | bytearray | b'\x00\x00' |
| メモリビュー型 | memoryview | `<memory at 0x10491c580>` |
| リスト型 | list | [1,2] |
| タプル型 | tuple | (1,2) |
| レンジ型 | range | range(0,10) | 
| セット型 | set | {'o', 'g', 'e', 'h'} |
| フローズンセット型 | frozenset | frozenset({'o', 'g', 'e', 'h'}) |
| 辞書型 | dict | {"x":1, "y":2} |


他にも、全ての型を許可する typing.Any型もあります。(`import typing`が必要です)


## おわりに
型アノテーションを適切に利用することで、コードの品質を向上させ、信頼性の高いプログラムを構築することができます。
是非、Pythonの型アノテーションを活用して、より効率的で安全なコードを書くための一歩を踏み出してみてください。

[1] バグが0になるとは書けなかった

