---
title: "Go1.20リリースお知らせ記事を翻訳する(後編)"
emoji: "😃"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["go","English","英語"]
published: false
---



## Outline

前編はこちらです。

https://zenn.dev/nii/articles/read-go-1-20-release-article

2023年2月1日にGo 1.20がリリースされ、The Go Blogにてリリース告知の記事が投稿されました。
本記事ではブログの翻訳を記載します。誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

https://go.dev/blog/go1.20

Go 1.20の詳細な機能説明に関しては、リリース告知ではなくリリースノートをご覧ください。
https://go.dev/doc/go1.20

## Index
### 前編([前回記事](https://zenn.dev/nii/articles/read-go-1-20-release-article))

- Go 1.20がリリースされました!
- 言語の変更
- ツールの改善

### 後編(本記事)

- 標準ライブラリの追加
- パフォーマンスの向上
- あとがき


## 標準ライブラリの追加
- 新しい`crypto/ecdh`パッケージは、NIST曲線とCurve25519を介した楕円曲線ディフィー・ヘルマン鍵共有を明示的にサポートします。

> - The new crypto/ecdh package provides explicit support for Elliptic Curve Diffie-Hellman key exchanges over NIST curves and Curve25519.

::: message
explicit: (形)明白な
Elliptic Curve Diffie-Hellman key exchanges: 楕円曲線ディフィー・ヘルマン鍵共有のこと。事前の秘密の共有無しに、盗聴の可能性のある通信路を使って、暗号鍵の共有を可能にする、公開鍵暗号方式の暗号プロトコルである。[wikipedia](https://ja.wikipedia.org/wiki/%E6%A5%95%E5%86%86%E6%9B%B2%E7%B7%9A%E3%83%87%E3%82%A3%E3%83%95%E3%82%A3%E3%83%BC%E3%83%BB%E3%83%98%E3%83%AB%E3%83%9E%E3%83%B3%E9%8D%B5%E5%85%B1%E6%9C%89)より引用
:::

---
> 新しい関数である `errors.Join` はエラーをラッピングしたリストを返します。これにより、もしエラータイプが `Unwrap() []error`メソッドを実装していれば再度エラーを取得できるようになるかもしれません。

> - The new function errors.Join returns an error wrapping a list of errors which may be obtained again if the error type implements the Unwrap() []error method.

::: message
obtain: (動)得る。今回は`be obtained`と受け身なので、エラーが取得されるという意味
:::

---
- 新しい `http.ResponseController` 型は、`http.ResponseWriter`インターフェースでは扱えなかった、HTTPリクエスト毎に制御できる拡張された機能を提供します。

> - The new http.ResponseController type provides access to extended per-request functionality not handled by the http.ResponseWriter interface.

::: message
functionality: (名)機能性
:::

---
- `httputil.ReverseProxy`のフォワードプロキシには、前の`Director`フックに取って代わって新しい`Rewrite`フック関数が含まれるようになりました。

> - The httputil.ReverseProxy forwarding proxy includes a new Rewrite hook function, superseding the previous Director hook.
> 
::: message
forwarding proxy: フォワードプロキシ
supersede: (動) 取って代わる
:::

---
- 新しい　`context.WithCancelCause` 関数は、与えられたエラーと共にコンテクストをキャンセルする手法を提供します。このエラーは、新関数である`context.Cause`を呼ぶことによって取得することができます。

> - The new context.WithCancelCause function provides a way to cancel a context with a given error. That error can be retrieved by calling the new context.Cause function.

::: message
retrieve: (動)取り戻す、回収する
:::

---
- 新しい`os/exec.Cmd`フィールドの`Cancel`と`WaitDelay`は、関連づけられたContextがキャンセルかプロセスが終了した際の `Cmd` の振る舞いを明示します。

> - The new os/exec.Cmd fields Cancel and WaitDelay specify the behavior of the Cmd when its associated Context is canceled or its process exits.

::: message
specify: (動)明示する、明記する
associate: (動)関連づける。ここではContextにかかる過去分詞の形容詞用法
:::


## パフォーマンスの向上
- コンパイラとガベージコレクターはメモリのオーバーヘッドを減らし、全体的なCPU性能を2%向上させました。

> - Compiler and garbage collector improvements have reduced memory overhead and improved overall CPU performance by up to 2%.

::: message
overhead: オーバーヘッドのこと。あるコンピューターの処理を実行するのに付随する作業を指すものである。たいていは、処理に時間がかかるようになるなど、システムの負荷になるものを指す。[IT用語辞典バイナリ](https://www.sophia-it.com/content/%E3%82%AA%E3%83%BC%E3%83%90%E3%83%BC%E3%83%98%E3%83%83%E3%83%89)より引用
:::

- 特にコンパイルをターゲットにする作業により、ビルドが10%程度早くなりました。これによって、Go1.17の時と同じ程度のスピードでビルドできます。

> - Work specifically targeting compilation times led to build improvements by up to 10%. This brings build speeds back in line with Go 1.17.

::: message
in line with~: ~と一致して
:::

---

::: message

:::








## あとがき
When building a Go release from source, Go 1.20 requires a Go 1.17.13 or newer release. In the future, we plan to move the bootstrap toolchain forward approximately once a year. Also, starting with Go 1.21, some older operating systems will no longer be supported: this includes Windows 7, 8, Server 2008 and Server 2012, macOS 10.13 High Sierra, and 10.14 Mojave. On the other hand, Go 1.20 adds experimental support for FreeBSD on RISC-V.

For a complete and more detailed list of all changes see the full release notes.

Thanks to everyone who contributed to this release by writing code, filing bugs, sharing feedback, and testing the release candidates. Your efforts helped to ensure that Go 1.20 is as stable as possible. As always, if you notice any problems, please file an issue.

Enjoy Go 1.20!
