---
title: "新しいEC2インスタンスタイプm7g,r7gについて/AWS News Blog"
emoji: "🥸"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["翻訳","AWS","EC2"]
published: false
---

# Outline
本記事はAWS New blogの"New Graviton3-Based General Purpose (m7g) and Memory-Optimized (r7g) Amazon EC2 Instances"を翻訳したものになります。
誤訳等や読みづらい点がありましたらコメント/Githubのissueにご連絡いただけると幸いです。

https://aws.amazon.com/jp/blogs/aws/new-graviton3-based-general-purpose-m7g-and-memory-optimized-r7g-amazon-ec2-instances/

# New Graviton3-Based General Purpose (m7g) and Memory-Optimized (r7g) Amazon EC2 Instances
<!-- リンク追加,スライド翻訳版作成 -->
2023年2月13日 Jeff Barr

2006年にm1.smallインスタンスをローンチして以来、メモリやCPU性能、CPUメーカーの選択肢を追加してきました。思えば長い道のりです。元々の汎用的なフリーサイズインスタンスは特定のケースに最適化された6つのタイプに進化し、全体で600以上のインスタンスが使用可能になりました。

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
本日、m7gとr7gという最新のAmazon EC2インスタンスタイプに関するお話をできることを嬉しく思います。両方とも最新のAWS Gravition3プロセッサーで動作し、第6世代(m6gとr6g)インスタンスに比べると25%ほどのパフォーマンス向上されており、EC2のさらなるパフォーマンスを引き出せるように設計されています。

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


Here are the specs for the M7g instances:

# Jeff Barrさんについて
