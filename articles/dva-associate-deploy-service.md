---
title: "AWSデベロッパーアソシエイト(DVA) デプロイツールを触ってみる"
emoji: "👨‍💻"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["aws","aws認定試験","dva","dva-02"]
published: false
---

## Outline
AWS Certified Developer - Associate認定試験の対象のうち、デベロッパーツールを実際に設定する記事になります。 //todo
参考書だけの暗記ではなく、実際に画面を見てイメージしやすくするための助けになれば幸いです。 //todo

対象となるサービスは下記です。
- AWS CodeCommit
- AWS CodeBuild
- AWS CodePipeline
- AWS CodeDeploy
- AWS Elastic Beanstalk

なお、2023年2月28日よりDVA試験のバージョンが新しくなりますが、本記事で紹介するサービスは`試験の対象となる主要なツール、テクノロジー、概念`リストに含まれています。
https://d1.awsstatic.com/ja_JP/training-and-certification/docs-dev-associate/AWS-Certified-Developer-Associate_Exam-Guide_C02.pdf


## 事前準備
IAM Userの作成
一旦AdministratorAccess権限

ログインしてユーザー名確認。

## AWS CodeCommit
### 概要
> AWS CodeCommit は、クラウド内のアセット(ドキュメント、ソースコード、バイナリファイルなど) を非公開で保存および管理するために使用できるアマゾン ウェブ サービスによってホストされるバージョン管理サービスです。([公式ドキュメント](https://docs.aws.amazon.com/ja_jp/codecommit/latest/userguide/welcome.html))より引用

Githubのようなサービスだと思っておきましょう。
リポジトリを作成し、基本的な操作を行なってみます。

### 使ってみる
下記ドキュメントを参考に進めていきます。
https://docs.aws.amazon.com/ja_jp/codecommit/latest/userguide/getting-started-cc.html

---
リポジトリを作成

作成画面
dva-sample、dva試験用のサンプルです。タグなし

--- 
接続
3パターンあるが、今回はHTTPS接続によるGit認証情報の設定の方法実施

ローカルgitバージョン確認
2.28以降であることを確認
> Git バージョン 1.7.9 以降をサポートしています。Git バージョン 2.28 は、初期コミットのブランチ名の構成をサポートしています。
```terminal
$ git --version
git version 2.37.1 (Apple Git-137.1)
```

IAM > User > セキュリティ認証情報 > AWS CodeCommit の HTTPS Git 認証情報 > 作成 > csv DL

`git clone https://git-codecommit.{region}.amazonaws.com/v1/repos/{repoName}` 実施
username,passはcsvでDLしたもの

```terminal
$ git clone https://git-codecommit.ap-northeast-1.amazonaws.com/v1/repos/dva-sample
Cloning into 'dva-sample'...
Username for 'https://git-codecommit.ap-northeast-1.amazonaws.com': ******
Password for 'https://dva-user-at-162265692617@git-codecommit.ap-northeast-1.amazonaws.com':  ******
warning: You appear to have cloned an empty repository.
$ ls
dva-sample aws_new_user_credentials.csv
```

---
適当にファイル作ってpushしてみる

gitの操作でも可能だが、せっかくなのでaws-cliを使ってみる


### 試験に関して
Gitの基本操作がわかっていればなんとかなりそう //todo
CodeCommit APIの名前も