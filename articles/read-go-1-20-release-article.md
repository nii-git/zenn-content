---
title: "Go1.20リリースお知らせ記事を読んでいく"
emoji: "😃"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["go","English","英語"]
published: false
---

## Outline
2023年2月1日にGo 1.20がリリースされ、The Go Blogにてリリース告知の記事が投稿されました。
本記事ではブログの翻訳を記載して行きます。誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

// TODO: 埋め込み確認
https://go.dev/blog/go1.20
[url](https://go.dev/blog/go1.20)

Go 1.20の詳細な機能説明に関しては、リリース告知ではなくリリースノートをご覧ください。
https://go.dev/doc/go1.20



# Go 1.20がリリースされました！
2023年 2月 1日 
Go開発チーム代表 Robert Griesemer

`
Go 1.20 is released!
Robert Griesemer, on behalf of the Go team
1 February 2023
`

本日、Go開発チームは Go 1.20をリリースすることを嬉しく思います。
ダウンロードは[こちら](https://go.dev/dl/)からすることが可能です。

開発期間の延長のおかげで、本バージョンではより早期段階での幅広いテストとコードベースの全体的な安定性の向上を実現しました。

`Today the Go team is thrilled to release Go 1.20, which you can get by visiting the download page.

Go 1.20 benefited from an extended development phase, made possible by earlier broad testing and improved overall stability of the code base.`

::: message
thrilled: (形)わくわくした、興奮した
benefit: (自動)利益を得る ここでは Go 1.20 にかかる 過去分詞の形容詞的用法
made possible by ~ : ~によって実現可能となる
:::

特に、PGOのプレビューサポート機能をリリースできたことを嬉しく思います。
PGOとは実行時間のプロファイル情報に基づき、コンパイラがアプリケーションや作業負荷に特化した最適化を可能とするものです。

`
We’re particularly excited to launch a preview of profile-guided optimization (PGO), which enables the compiler to perform application- and workload-specific optimizations based on run-time profile information.
`

::: message
PGO: Profile-Guided Optimization コンパイラ最適化の手法 詳細は[こちら](https://go.dev/doc/pgo)
application-specific: (形) アプリケーションに特化した ここではapplcation-とworkload- が specificにかかっている
:::

// ここまで


 Providing a profile to go build enables the compiler to speed up typical applications by around 3–4%, and we expect future releases to benefit even more from PGO. Since this is a preview release of PGO support, we encourage folks to try it out, but there are still rough edges which may preclude production use.

Go 1.20 also includes a handful of language changes, many improvements to tooling and the library, and better overall performance.

Language changes
The predeclared comparable constraint is now also satisfied by ordinary comparable types, such as interfaces, which will simplify generic code.
The functions SliceData, String, and StringData have been added to package unsafe. They complete the set of functions for implementation-independent slice and string manipulation.
Go’s type conversion rules have been extended to permit direct conversion from a slice to an array.
The language specification now defines the exact order in which array elements and struct fields are compared. This clarifies what happens in case of panics during comparisons.
Tool improvements
The cover tool now can collect coverage profiles of whole programs, not just of unit tests.
The go tool no longer relies on pre-compiled standard library package archives in the $GOROOT/pkg directory, and they are no longer shipped with the distribution, resulting in smaller downloads. Instead, packages in the standard library are built as needed and cached in the build cache, like other packages.
The implementation of go test -json has been improved to make it more robust in the presence of stray writes to stdout.
The go build, go install, and other build-related commands now accept a -pgo flag enabling profile-guided optimizations as well as a -cover flag for whole-program coverage analysis.
The go command now disables cgo by default on systems without a C toolchain. Consequently, when Go is installed on a system without a C compiler, it will now use pure Go builds for packages in the standard library that optionally use cgo, instead of using pre-distributed package archives (which have been removed, as noted above).
The vet tool reports more loop variable reference mistakes that may occur in tests running in parallel.
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

Previous article: Share your feedback about developing with Go
Blog Index

Why Go
Use Cases
Case Studies
Get Started
Playground
Tour
Stack Overflow
Help
Packages
Standard Library
About Go Packages
About
Download
Blog
Issue Tracker
Release Notes
Brand Guidelines
Code of Conduct
Connect
Twitter
GitHub
Slack
r/golang
Meetup
Golang Weekly
The Go Gopher
Copyright
Terms of Service
Privacy Policy
Report an Issue
System theme
Google logo