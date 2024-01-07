---
title: "GoのEchoフレームワークで作成したコンテナをECS Fargateにデプロイする"
emoji: "💻"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["go","aws","echo","ecs","fargate","vpc"]
published: false
---

ねむい。ちゃんと推敲すること。

## 概要
Go言語のEchoフレームワークを用いて開発をした際、どのようにECS Fargateにデプロイするかまとめた記事です。
バックエンドはわかるけど、ネットワークが全くわからない...という方を対象にしています。

ローカルで動かしていたEchoフレームワークをデプロイし、外部からアクセスできるようにする

前提で、ローカルにdockerファイルがあることを想定

// 混ぜてもいいかも注意点と概要　眠いんだ　朝4時から書いてるんだ

## 注意点
今回の構成ではECSとRDSをpublic subnetに配置している点に注意してください。

よりセキュアな構成にする場合は、ECSやRDSをprivate subnetに移動させることを検討してください。
Natゲートウェイや[エンドポイント](https://dev.classmethod.jp/articles/privatesubnet_ecs/)の作成が必要です。

また、コンソール状でぽちぽちやっているため、IaCで作成したい方には　//TODO

記事中のurlは基本的にregionが`ap-northeast-1`のものになっています。
他のリージョンで作成する方はリンクを踏まずに実施してください。

## 構成図


## 1. VPC設定

### 1.1 VPC作成
- VPC > 仮想プライベートクラウド > お使いのVPCに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#vpcs
- VPCを作成を選択
- 下記のように設定
  - 作成するリソース: VPCのみ
  - 名前タグ: わかりやすい名前
  - IPv4 CICDブロック: IPv4CIDRの手動入力
  - IPv4 CIDR: 大規模であれば`10.0.0.0/16`, 中規模であれば`172.16.0.0/16`
    - 今回は中規模なので`172.16.0.0/16`を使用する
    - 参考: [VPC CIDR ブロック](https://docs.aws.amazon.com/ja_jp/vpc/latest/userguide/vpc-cidr-blocks.html)
  - Ipv6CIDRブロック: IPv6 CIDR ブロックなし
  - テナンシー: デフォルト


### 1.2 パブリックサブネット作成
- VPC > 仮想プライベートクラウド > サブネットに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#subnets
- サブネットを作成を選択
- 新しいサブネットを追加を選択
- VPC ID: 1.1で作成したVPCを選択
- サブネット 1 (2 個中)を下記のように設定
  - サブネット名: 好きな名前
  - アベイラビリティーゾーン: ap-northeast-1a
  - IPv4 VPC CIDR block: `172.16.0.0/16`
  - IPv4 subnet CIDR block: `172.16.10.0/24`
- サブネット 2 (2 個中)を下記のように設定
  - サブネット名: 好きな名前
  - アベイラビリティーゾーン: ap-northeast-1c
  - IPv4 VPC CIDR block: `172.16.0.0/16`
  - IPv4 subnet CIDR block: `172.16.20.0/24`


### 1.3 パブリックルートテーブル作成
- VPC > 仮想プライベートクラウド > ルートテーブルに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#RouteTables
- ルートテーブルを作成を選択
- 下記のように設定
  - 名前: 好きな名前
  - VPC: 1.1で作成したVPCを選択
- 作成後、作成したルートテーブルを選択 > アクション > サブネットの関連付けを編集 を選択
- 1.2で作成したサブネットを選択し、保存

### 1.4 プライベートサブネット作成
- VPC > 仮想プライベートクラウド > サブネットに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#subnets
- サブネットを作成を選択
- 新しいサブネットを追加を選択
- VPC ID: 1.1で作成したVPCを選択
- サブネット 1 (2 個中)を下記のように設定
  - サブネット名: 好きな名前
  - アベイラビリティーゾーン: ap-northeast-1a
  - IPv4 VPC CIDR block: `172.16.0.0/16`
  - IPv4 subnet CIDR block: `172.16.30.0/24`
- サブネット 2 (2 個中)を下記のように設定
  - サブネット名: 好きな名前
  - アベイラビリティーゾーン: ap-northeast-1c
  - IPv4 VPC CIDR block: `172.16.0.0/16`
  - IPv4 subnet CIDR block: `172.16.40.0/24`


### 1.5 プライベートルートテーブル作成
- VPC > 仮想プライベートクラウド > ルートテーブルに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#RouteTables
- ルートテーブルを作成を選択
- 下記のように設定
  - 名前: 好きな名前
  - VPC: 1.1で作成したVPCを選択
- 作成後、作成したルートテーブルを選択 > アクション > サブネットの関連付けを編集 を選択
- 1.4で作成したサブネットを選択し、保存

### 1.6 インターネットゲートウェイを作成
- VPC > 仮想プライベートクラウド > インターネットゲートウェイに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#igws
- インターネットゲートウェイの作成を選択
- 下記のように設定
  -  名前: 好きな名前
-  作成後、作成したインターネットゲートウェイを選択 > アクション > VPCにアタッチを選択
-  VPC: 1.1で選択したVPC を選択


### 1.7 現在の状態図
![現在の状態](/images/go-deploy-ecs/go-deploy-ecs1.jpg)


## 2. ドメイン設定
### 2.1 ドメイン取得
- Route53 > ドメイン > ドメインを登録に移動
  - https://us-east-1.console.aws.amazon.com/route53/domains/home?region=ap-northeast-1#/
- ドメインを登録を選択
- 使用したいドメイン名を入力
- 期間、自動更新の有無を入力
- 連絡先情報を入力

### 2.2 証明書を作成
- AWS Certificate Manager > 証明書一覧に移動
  - https://ap-northeast-1.console.aws.amazon.com/acm/home?region=ap-northeast-1#/certificates/list
- リクエストを選択
- 下記のように設定
  - 証明書タイプ: パブリック証明書をリクエスト
  - 完全修飾ドメイン名: 2.1で登録したドメイン名
  - 検証方法: DNS検証
  - キーアルゴリズム: RSA2048
- 該当の証明書ステータスが`検証保留中`になっていることを確認

### 2.3 証明書を検証
- 10分程度で2.1で登録したメールアドレス宛にRoute 53からメールが届く
  - 送信元: Amazon Route 53 (noreply@registrar.amazon.com)
  - 件名: Verify your email address for {ドメイン名}
- 本文中にあるverify linkを踏めばOK

### 2.4 レコードを作成
- 2.3で作成した証明書の詳細 > ドメイン > Route53でレコードを作成を選択
- 2.1で登録したドメインを選択し、レコードを作成を選択
- 5分程度で該当の証明書ステータスが`発行済み`になっていることを確認

### 2.5 現在の状態図
![現在の状態](/images/go-deploy-ecs/go-deploy-ecs2.jpg)


## 3. ルーティング設定
### 3.1 パブリックルートテーブルの設定
- VPC > 仮想プライベートクラウド > ルートテーブルに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#RouteTables:
- 1.3で作成したパブリックルートテーブルを選択 > アクション > ルートを編集
- 下記のように設定
  - 送信先: `0.0.0.0/0`
  - ターゲット: インターネットゲートウェイ/{1.6で作成したもの}

### 3.2 ELB用セキュリティグループの作成
- VPC > セキュリティ > セキュリティグループに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#SecurityGroups
- セキュリティグループを作成
- 下記のように設定
  - セキュリティグループ名: 好きな名前
  - VPC: 1.1で作成したVPC
  - インバウンドルール1
    - タイプ: HTTP
    - ソース: Anywhere-IPv4
  - インバウンドルール2
    - タイプ: HTTPS
    - ソース: Anywhere-IPv4
  - アウトバウンドルール
    - タイプ: すべてのトラフィック
    - 送信先: Anywhere-IPv4


### 3.3 ECS用セキュリティグループの作成
- VPC > セキュリティ > セキュリティグループに移動
  - https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#SecurityGroups
- セキュリティグループを作成
- 下記のように設定
  - セキュリティグループ名: 好きな名前
  - VPC: 1.1で作成したVPC
  - インバウンドルール
    - タイプ: HTTP
    - ソース: カスタム、3.2で作成したセキュリティグループ
  - アウトバウンドルール
    - タイプ: すべてのトラフィック
    - 送信先: Anywhere-IPv4


### 3.4 ターゲットグループの作成
- EC2 > ロードバランシング > ターゲットグループに移動
  - https://ap-northeast-1.console.aws.amazon.com/ec2/home?region=ap-northeast-1#TargetGroups
- ターゲットグループの作成を選択
- 下記のように設定
  - ターゲットタイプの選択: IPアドレス
  - ターゲットグループ名: 好きな名前
  - プロトコル : ポート: HTTP:80
  - IP アドレスタイプ: IPv4
  - VPC: 1.1で作成したVPC
  - プロトコルバージョン: HTTP1
  - ヘルスチェックプロトコル: HTTP
  - ヘルスチェックパス: /
  - ターゲット登録は何もせずにターゲットグループの作成を選択


### 3.5 ELBの作成
- EC2 > ロードバランシング > ロードバランサーに移動
  - https://ap-northeast-1.console.aws.amazon.com/ec2/home?region=ap-northeast-1#LoadBalancers:
- ロードバランサーの作成を選択
- 下記のように設定
  - ロードバランサータイプ:Application Load Balancer
  - ロードバランサー名: 好きな名前
  - スキーム: インターネット向け
  - IPアドレスタイプ: IPv4
  - VPC: 1.1で作成したVPC
  - マッピング: ap-northeast-1a, ap-northeast-1c(1.2で使用した2個のパブリックサブネットを選択)
  - セキュリティグループ: 3.3で作成したセキュリティグループ
  - リスナープロトコル,ポート: HTTPS/443
  - デフォルトアクション: 転送先 3.4で作成したターゲットグループ
  - セキュリティカテゴリ: 全てのセキュリティポリシー
  - ポリシー名: 推奨のもの
  - 証明書の取得先: ACMから
  - 証明書(ACMから): 2.2で作成した証明書


### 3.6 現在の状態図
![現在の状態](/images/go-deploy-ecs/go-deploy-ecs3.jpg)


## 4. ECR設定
### 4.1 リポジトリ作成
- Amazon Elastic Container Registry > Private regisstry > Repositoriesに移動
  - https://ap-northeast-1.console.aws.amazon.com/ecr/private-registry/repositories?region=ap-northeast-1
- リポジトリを作成を選択
- 下記の通りに設定
  - 可視化設定: プライベート
  - リポジトリ名: 好きな名前


### 4.2 イメージをプッシュ
- ローカル環境にて、dockerファイルがあるディレクトリで実行
  - 4.1で作成したリポジトリ > プッシュコマンドの表示で確認した方が確実だが、dockerをビルドする際に`--platform linux/x86_64`を忘れないようにする
- 下記が完了した際、4.1で作成したリポジトリにlatestイメージタグのイメージがあることを確認する

```
docker build -t {APP_NAME} --platform linux/x86_64 .
docker tag {APP_NAME} {ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/{APP_NAME}
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin {ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com
docker push {ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/{APP_NAME}
```

### 4.3 現在の状態図
![現在の状態](/images/go-deploy-ecs/go-deploy-ecs4.jpg)