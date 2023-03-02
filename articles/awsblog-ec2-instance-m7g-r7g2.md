---
title: "新しいEC2インスタンスタイプm7g、r7gについて/AWS News Blog"
emoji: "🥸"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["翻訳","AWS","EC2"]
published: true
---

# Outline
本記事はAWS New blogの"New Graviton3-Based General Purpose (m7g) and Memory-Optimized (r7g) Amazon EC2 Instances"を翻訳したものになります。
誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

https://aws.amazon.com/jp/blogs/aws/new-graviton3-based-general-purpose-m7g-and-memory-optimized-r7g-amazon-ec2-instances/

# New Graviton3-Based General Purpose (m7g) and Memory-Optimized (r7g) Amazon EC2 Instances
<!-- リンク追加,スライド翻訳版作成 -->
2023年2月13日 Jeff Barr

2006年にm1.smallインスタンスをローンチして以来、[インスタンス](https://aws.amazon.com/ec2/instance-types/)のメモリやCPU性能、CPUメーカーの選択肢を追加してきました。思えば長い道のりです。元々の汎用的なフリーサイズインスタンスは特定のケースに最適化された6つのタイプに進化し、全体で600以上のインスタンスが使用可能になりました。

::: message
訳注: EC2インスタンスには下記の６タイプが存在しています。
- 一般用途向け　
- コンピューティング最適化
- メモリ最適化
- 高速コンピューティング
- ストレージ最適化
- HPC最適化
https://aws.amazon.com/jp/ec2/instance-types/
:::

## 新しいm7gとr7g
本日、m7gとr7gという最新のAmazon EC2インスタンスタイプに関するお話ができることを嬉しく思います。両方とも最新のAWS [Gravition3](https://aws.amazon.com/ec2/graviton/)プロセッサーで動作し、第6世代(m6gとr6g)インスタンスに比べると25%ほどのパフォーマンスの向上がされており、EC2のさらなるパフォーマンスを引き出せるように設計されています。

m7gインスタンスは、アプリケーションサーバー、マイクロサービス、ゲームサーバー、中程度のデータサーバー、フリートのキャッシュ等の汎用的な使用に向いています。r7gインスタンスは、オープンソースのデータベースやインメモリ型のキャッシュ、リアルタイムでのビッグデータ分析等のメモリに負荷がかかるような作業にぴったりでしょう。

m7gインスタンスのスペック表はこちらです。

| インスタンス名 | vCPU数 | メモリ | ネットワーク帯域幅 | EBS帯域幅 |
| ---- | ---- | ---- | ---- | ---- |
| m7g.medium | 1 | 4 GiB | up to 12.5 Gbps | up to 10 Gbps |
| m7g.large	| 2 | 8 GiB | up to 12.5 Gbps | up to 10 Gbps |
| m7g.xlarge | 4 | 16 GiB | up to 12.5 Gbps | up to 10 Gbps |
| m7g.2xlarge |	8 |	32 GiB | up to 15 Gbps | up to 10 Gbps |
| m7g.4xlarge |	16 | 64 GiB | up to 15 Gbps | up to 10 Gbps |
| m7g.8xlarge |	32 | 128 GiB | 15 Gbps | 10 Gbps |
| m7g.12xlarge | 48 | 192 GiB |	22.5 Gbps |	15 Gbps |
| m7g.16xlarge | 64 | 256 GiB |	30 Gbps | 20 Gbps |
| m7g.metal | 64 | 256 GiB | 30 Gbps | 20 Gbps |


r7gインスタンスのスペック表はこちらです。

| インスタンス名 | vCPU数 | メモリ | ネットワーク帯域幅 | EBS帯域幅 |
| ---- | ---- | ---- | ---- | ---- |
| r7g.medium | 1 | 8 GiB | up to 12.5 Gbps | up to 10 Gbps |
| r7g.large | 2 | 16 GiB | up to 12.5 Gbps | up to 10 Gbps |
| r7g.xlarge | 4 | 32 GiB | up to 12.5 Gbps | up to 10 Gbps |
| r7g.2xlarge | 8 | 64 GiB | up to 15 Gbps | up to 10 Gbps |
| r7g.4xlarge | 16 | 128 GiB | up to 15 Gbps | up to 10 Gbps |
| r7g.8xlarge | 32 | 256 GiB | 15 Gbps | 10 Gbps |
| r7g.12xlarge | 48 | 384 GiB | 22.5 Gbps | 15 Gbps |
| r7g.16xlarge | 64 | 512 GiB | 30 Gbps | 20 Gbps |
| r7g.metal | 64 | 512 GiB | 30 Gbps | 20 Gbps |

どちらのインスタンスタイプも[DDR5](https://en.wikipedia.org/wiki/DDR5_SDRAM)メモリを用いており、前世代のDDR4よりもメモリ帯域幅が最大50%上昇しています。下記の画像は新しいインスタンスで利用可能な主なパフォーマンスや容量の向上をまとめたものです。

![パフォーマンス向上まとめ](/images/others/awsblog-ec2-instance-m7g-r7g2.png)

もしまだアプリケーションをGravitonインスタンスで実行したことがない場合、ぜひ[AWS Graviton Ready Program](https://aws.amazon.com/ec2/graviton/partners/?partner-solutions-cards.sort-by=item.additionalFields.partnerNameLower&partner-solutions-cards.sort-order=asc&awsf.partner-solutions-filter-partner-type-storage=*all&awsf.partner-solutions-filter-partner-location-storage=*all)をご利用ください。このプログラムのパートナーは、アプリケーションの移行だけでなく、Gravitonインスタンスが提供する全ての機能を最大限に活用できるように助けてくれるでしょう。他にも、[Porting Advisor for Graviton](https://aws.amazon.com/about-aws/whats-new/2023/01/porting-advisor-graviton/)や[raviton Fast Start program](https://aws.amazon.com/ec2/graviton/fast-start/)に有益な情報があるのでチェックしてください。

これらのインスタンスは[AWS Nitro System](https://aws.amazon.com/ec2/nitro/)でビルドされているので、様々なセキュリティの恩恵を受けることができます。例として、常時メモリ暗号化、vCPUに特化したキャッシュ、ポインター認証のサポート等が挙げられます。また、保存データやインスタンスとボリューム間のデータ移行、スナップショット作成やスナップショットからのボリューム作成時にデータを守ってくれる[EBSボリュームの暗号化](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html)もサポートしています。これらのNitroによるセキュリティ特徴の詳細は[The Security Design of the AWS Nitro System](https://docs.aws.amazon.com/whitepapers/latest/security-design-of-aws-nitro-system/security-design-of-aws-nitro-system.html)をお読みください。

インスタンスのネットワークでは、[インスタンス-EBSボリューム間の通信用に最適化](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-optimized.html)されており、拡張ネットワークもサポートしています。詳細は[How do I enable and configure enhanced networking on my EC2 instances?]()をお読みください。16xlargeとmetalインスタンスでは、ノード間での高速な通信が必要なアプリケーション用に[Elastic Fabric Adapte](https://aws.amazon.com/hpc/efa/)(EFA)もサポートしています。


## 価格とリージョン
m7g、r7gインスタンスは下記のリージョンでオンデマンド、スポット、予約インスタンス、セービングプランタイプを提供開始しました。
- 米国東部(バージニア州)
- 英国東部(オハイオ州)
- 米国西部(オレゴン州)
- ヨーロッパ(アイルランド)

早速使ってみて、感想を教えてください！
[Jeff Barr](https://twitter.com/jeffbarr)
