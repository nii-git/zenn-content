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

//todo:public ルートテーブル

## 2. ドメイン設定
### 2.1 ドメイン取得
- Route53 > ドメイン > ドメインを登録に移動
  - https://us-east-1.console.aws.amazon.com/route53/domains/home?region=ap-northeast-1#/
- ドメインを登録を選択
- 使用したいドメインを入力
- 期間、自動更新の有無を入力
- 連絡先情報を入力

### 2.2 証明書を作成
- AWS Certificate Manager > 証明書一覧に移動
  - https://ap-northeast-1.console.aws.amazon.com/acm/home?region=ap-northeast-1#/certificates/list
- ホストゾーンの作成を選択
- 
