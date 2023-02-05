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
本記事ではブログの翻訳を記載して行きます。誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

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


Standard library additions
The new crypto/ecdh package provides explicit support for Elliptic Curve Diffie-Hellman key exchanges over NIST curves and Curve25519.
The new function errors.Join returns an error wrapping a list of errors which may be obtained again if the error type implements the Unwrap() []error method.
The new http.ResponseController type provides access to extended per-request functionality not handled by the http.ResponseWriter interface.
The httputil.ReverseProxy forwarding proxy includes a new Rewrite hook function, superseding the previous Director hook.
The new context.WithCancelCause function provides a way to cancel a context with a given error. That error can be retrieved by calling the new context.Cause function.
The new os/exec.Cmd fields Cancel and WaitDelay specify the behavior of the Cmd when its associated Context is canceled or its process exits.
Improved performance
Compiler and garbage collector improvements have reduced memory overhead and improved overall CPU performance by up to 2%.
Work specifically targeting compilation times led to build improvements by up to 10%. This brings build speeds back in line with Go 1.17.
When building a Go release from source, Go 1.20 requires a Go 1.17.13 or newer release. In the future, we plan to move the bootstrap toolchain forward approximately once a year. Also, starting with Go 1.21, some older operating systems will no longer be supported: this includes Windows 7, 8, Server 2008 and Server 2012, macOS 10.13 High Sierra, and 10.14 Mojave. On the other hand, Go 1.20 adds experimental support for FreeBSD on RISC-V.

For a complete and more detailed list of all changes see the full release notes.

Thanks to everyone who contributed to this release by writing code, filing bugs, sharing feedback, and testing the release candidates. Your efforts helped to ensure that Go 1.20 is as stable as possible. As always, if you notice any problems, please file an issue.

Enjoy Go 1.20!
