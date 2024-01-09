---
title: "DockerコンテナをECS Fargateにデプロイする"
emoji: "💻"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["aws","ecs","fargate","vpc"]
published: false
---

## 概要
ローカルで動かしていたコンテナをECSにデプロイし、外部からアクセスできるようにする方法をまとめた記事です。
既にローカルで動作するdockerファイルがあることを想定しているため、本記事ではdockerの作成方法は記載しません。

<!-- 混ぜてもいいかも注意点と概要　眠いんだ　朝4時から書いてるんだ -->

## 注意点
今回の構成ではECSとRDSをpublic subnetに配置している点に注意してください。

よりセキュアな構成にする場合は、ECSやRDSをprivate subnetに移動させることを検討してください。Natゲートウェイや[エンドポイント](https://dev.classmethod.jp/articles/privatesubnet_ecs/)の作成が必要です。

記事中のurlは基本的にregionが`ap-northeast-1`のものになっています。
他のリージョンで作成する方はリンクを踏まずに実施してください。

## 構成図
![現在の状態](/images/go-deploy-ecs/go-deploy-ecs5.jpg)

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


### 3.6 Aレコードの作成
- Route53 > ホストゾーン > 2.1で作成したホストゾーンを選択
- レコードを作成を選択
- 下記のように設定
  - レコード名: 空白
  - レコードタイプ: A - IPv4アドレスと一部のAWSリソースにトラフィックを...
  - エイリアス: チェックを入れる
  - トラフィックのルーティング先: Application Load BalancerとClassic Load Balancerへのエイリアス、東京リージョン、3.5で作成したELBを選択


### 3.7 現在の状態図
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


## 5. ECS設定
### 5.1 ECSロールの作成
- IAM > アクセス管理 > ロールに移動
  - https://us-east-1.console.aws.amazon.com/iam/home?region=ap-northeast-1#/roles
- ロールを作成を選択
- 下記のように設定
  - 信頼されたエンティティタイプ: AWSのサービス
  - サービスまたはユースケース: Elastic Container Service
  - ユースケース: Elastic Container Service Task
  - 許可ポリシー: AmazonECSTaskExecutionRolePolicy
  - ロール名: 好きな名前


### 5.2 クラスターの作成
- Amazon Elastic Container Service > クラスターを選択
  - https://ap-northeast-1.console.aws.amazon.com/ecs/v2/clusters?region=ap-northeast-1
- クラスターの作成を選択
- 下記のように設定
  - クラスター名: 好きな名前
  - インフラストラクチャ: AWS Fargate


### 5.3 タスク定義の作成
- Amazon Elastic Container Service > タスク定義を選択
  - https://ap-northeast-1.console.aws.amazon.com/ecs/v2/task-definitions?region=ap-northeast-1
- 新しいタスク定義の作成 > 新しいタスク定義の作成を選択
- 下記のように設定
  - タスク定義ファミリー: 好きな名前
  - 起動タイプ: AWS Fargate
  - オペレーティングシステム/アーキテクチャ: Linux/X86_64
  - ネットワークモード: awsvpc
  - CPU,メモリ: 好きなサイズ
  - タスクロール,タスク実行ロール: 5.1で作成したロール
  - コンテナの詳細_名前: 好きな名前
  - コンテナの詳細_イメージURI: 4.3でプッシュしたECRのURI
  - コンテナの詳細_必須コンテナ: はい
  - コンテナポート,プロトコル: 80,tcp
  - その他はデフォルトのまま

### 5.4 サービスの作成
- Amazon Elastic Container Service > クラスター > 5.2で作成したクラスターを選択
- サービス > 作成を選択
- 下記のように設定
  - コンピューティングオプション: 起動タイプ
  - 起動タイプ: FARGATE
  - プラットフォームのバージョン: 1.4.0
  - アプリケーションタイプ: サービス
  - ファミリー: 5.3で作成したタスク
  - リビジョン: 最新のもの
  - サービス名: 好きな名前
  - サービスタイプ: レプリカ
  - 必要なタスク: 好きな数
- ネットワーキングを選択し、下記のように設定
  - VPC: 1.1で作成したVPC
  - サブネット: 1.2で作成した2個のサブネット
  - セキュリティグループ: 既存のセキュリティグループを使用、3.3で作成したセキュリティグループを選択
  - パブリックIP: オンになっています
- ロードバランシングを選択し、下記のように設定
  - ロードバランサーの種類: Application Load Balancer
  - Application Load Balancer: 既存のロードバランサーを使用、3.5で作成したロードバランサーを選択
  - リスナー: 既存のリスナーを使用、443:HTTPS
  - ターゲットグループ: 既存のターゲットグループを使用、3.4で作成したターゲットグループを選択


### 5.5 現在の状態図
![現在の状態](/images/go-deploy-ecs/go-deploy-ecs5.jpg)


## 6. さいごに
登録したドメインにアクセスし、正しく表示されることを確認してください。
ECSタスクが起動失敗した場合、下記のランブックで診断してください。

https://docs.aws.amazon.com/ja_jp/systems-manager-automation-runbooks/latest/userguide/automation-aws-troubleshootecstaskfailedtostart.html

