---
title: "ジャルジャルの煩悩ネタの人気ネタを解析する奴"
emoji: "😎"
type: "idea" # tech: 技術記事 / idea: アイデア
topics: ["ジャルジャル","GCP","Youtube","YoutubeAPI","Go"]
published: false
---

皆様、[ジャルジャル](https://profile.yoshimoto.co.jp/talent/detail?id=179)というコンビ芸人はご存知でしょうか？

![ジャルジャル](/images/others/jarujaru.jpg)

独特なキャラクターと無限のネタで笑わせてくれるお二人のネタが大好きで、よく開発の作業用BGMとして利用させていただいております。

そんな中でも私のお気に入りなのが、2019年の大晦日に開催した「大晦日に108回もジャルってんじゃねえよ」。
除夜の鐘が`108回`鳴ることにかけて、披露するネタの数が**108回**、開演時間も**10時間8分**と正気の沙汰でないライブとなっています（料金も10800円というこだわりっぷりです）。

なんと、この108本のコントがYoutube上で無料で見ることができるのです。
なんていい時代なのでしょうか。ジャルジャルさんには頭が上がりませんね。

(この109件の動画を一気にアップしたせいでチャンネル登録欄がジャルジャルで埋め尽くされ、登録者が1.8万人減ったニュースを見た時は爆笑しました)

https://www.youtube.com/playlist?app=desktop&list=PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz


---

この煩悩プレーリストを周回している時にふと思いました。
どのコントが人気あるのだろうか、と。

というわけで、解析するためにツールを作成しました。
再生リスト内の動画の再生回数や高評価、コメント数を取得するツールです。

Youtube Playlist Analyze
https://github.com/nii-git/youtube-app

使用技術や使い方に関してはランキングの下に記載しています。

# ランキング
[ジャルジャルタワー JARUJARU TOWER](https://www.youtube.com/@jarujarutower365)チャンネルにアップされている[１０８本！煩悩ネタ！](https://www.youtube.com/playlist?app=desktop&list=PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz)再生リストの中から、下記の項目で5位までを記載します。
- 再生回数
- 高評価数
- コメント数
- 高評価率(高評価数/再生回数)
- コメント率(コメント数/再生回数)

データは2023年3月15日のものです。

## 再生回数ランキング
1位: [煩悩ネタ！『嘘つき通す奴』](https://www.youtube.com/watch?v=dKB8Z_KJG3I) 1386028回再生
2位: [煩悩ネタ！『第一声、間違えて、友達失う奴』](https://www.youtube.com/watch?v=f_r4F6QTZx4) 966121回再生
3位: [煩悩ネタ！『オキンタマデカ男って奴』](https://www.youtube.com/watch?v=hFnhrRqDrA0)963429回再生
4位: [煩悩ネタ！『あだ名コロコロ変わる奴』](https://www.youtube.com/watch?v=TJDTMyVX4fw)504954回再生
5位: [煩悩ネタ！『大学の入学を辞めただけの奴』](https://www.youtube.com/watch?v=wg-hRa2CFEs)455808回再生

## 高評価数ランキング
1位: [煩悩ネタ！『オキンタマデカ男って奴』](https://www.youtube.com/watch?v=hFnhrRqDrA0)16830高評価
2位: [煩悩ネタ！『嘘つき通す奴』](https://www.youtube.com/watch?v=dKB8Z_KJG3I)14872高評価
3位: [煩悩ネタ！『第一声、間違えて、友達失う奴』](https://www.youtube.com/watch?v=f_r4F6QTZx4)12499高評価
4位: [煩悩ネタ！『声の周波数のせいで離れたら英語に聞こえる奴』](https://www.youtube.com/watch?v=Au-mjjkqe7g)7572高評価
5位: [煩悩ネタ！『あだ名コロコロ変わる奴』](https://www.youtube.com/watch?v=TJDTMyVX4fw)6467高評価

## コメント数ランキング
1位: [煩悩ネタ！『107本を通して感じたことをコントにする奴』](https://www.youtube.com/watch?v=xZDBDH8eseY)861コメント
2位: [煩悩ネタ！『大学の入学を辞めただけの奴』](https://www.youtube.com/watch?v=wg-hRa2CFEs)761コメント
3位: [煩悩ネタ！『第一声、間違えて、友達失う奴』](https://www.youtube.com/watch?v=f_r4F6QTZx4)746コメント
4位: [煩悩ネタ！『オキンタマデカ男って奴』](https://www.youtube.com/watch?v=hFnhrRqDrA0)721コメント
5位: [煩悩ネタ！『嘘つき通す奴』](https://www.youtube.com/watch?v=dKB8Z_KJG3I)605コメント

## 高評価率ランキング
1位: [煩悩ネタ！『笛が落ちずにすんだ奴』](https://www.youtube.com/watch?v=N24IQTFXHds)3.41%
2位: [煩悩ネタ！『上司の実家のインターホン押しただけの奴』](https://www.youtube.com/watch?v=_CGIhM4Rd_k)2.21%
3位: [煩悩ネタ！『超能力あった奴』](https://www.youtube.com/watch?v=4Smk4Y3ZcNQ)2.07%
4位: [煩悩ネタ！『声の周波数のせいで離れたら英語に聞こえる奴』](https://www.youtube.com/watch?v=Au-mjjkqe7g)2.05%
5位: [煩悩ネタ！『声でっか、服小さっ！な奴』](https://www.youtube.com/watch?v=FrOQDo7jgdg)2.03%

## コメント率ランキング
1位: [煩悩ネタ！『107本を通して感じたことをコントにする奴』](https://www.youtube.com/watch?v=xZDBDH8eseY)0.762%
2位: [煩悩ネタ！『キャラ濃いくせに取材NGな奴』](https://www.youtube.com/watch?v=8562Gw23l7w)0.184%
3位: [煩悩ネタ！『チャラ男番長っていう奴』](https://www.youtube.com/watch?v=hzpp0Ryn6rw)0.179%
4位: [煩悩ネタ！『大学の入学を辞めただけの奴』](https://www.youtube.com/watch?v=wg-hRa2CFEs)0.167%
5位: [煩悩ネタ！『一服する奴』](https://www.youtube.com/watch?v=1TzGN4VSGR0)0.152%

# 技術
Go + YoutubeAPI で作成しているシンプルなプログラムです。

//図解

## ツール、使用方法

## 実装方法


---
//あとがき

こういう機能が欲しい等のコメント待ってる

こういう記事もある。
https://natalie.mu/owarai/pp/jarujaru02