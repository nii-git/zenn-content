---
title: "YoutubeAPIを使用してジャルジャルの人気ネタを解析する奴"
emoji: "😎"
type: "idea" # tech: 技術記事 / idea: アイデア
topics: ["ジャルジャル","GCP","Youtube","YoutubeAPI","Go"]
published: false
---

皆様、[ジャルジャル](https://profile.yoshimoto.co.jp/talent/detail?id=179)というコンビ芸人はご存知でしょうか？

![ジャルジャル](/images/others/jarujaru.jpg)

独特なキャラクターと無限のネタで笑わせてくれるお二人のネタが大好きで、よく開発の作業用BGMとして利用しています。

そんな中でも私のお気に入りなのが、2019年の大晦日に開催した「大晦日に108回もジャルってんじゃねえよ」というライブです。
除夜の鐘が**108回**鳴ることにかけて、披露するネタの数が**108回**、開演時間も**10時間8分**と正気の沙汰でないライブとなっています（料金も10800円というこだわりっぷりです）。

2023年5月現在、この108本のコントがYoutube上で無料で見ることができます。
なんていい時代なのでしょうか。ジャルジャルさんには頭が上がりませんね。

(この108件の動画を一気にアップしたせいでチャンネル登録欄がジャルジャルで埋め尽くされ、登録者が1.8万人減ったニュースを見た時は爆笑しました)

https://www.youtube.com/playlist?app=desktop&list=PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz


---

さて、この煩悩プレーリストを周回している時にふと思いました。
どのコントが人気あるのでしょうか？

というわけで、解析するためにツールを作成しました。
再生リスト内の動画の再生回数や高評価、コメント数を取得するツールです。

https://github.com/nii-git/youtube-app

使用技術に関してはランキングの下に記載しています。

# ランキング
[ジャルジャルタワー JARUJARU TOWER](https://www.youtube.com/@jarujarutower365)チャンネルにアップされている[１０８本！煩悩ネタ！](https://www.youtube.com/playlist?app=desktop&list=PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz)再生リストの中から、下記の項目で5位までを記載します。
- 再生回数
- 高評価数
- コメント数

データは2023年5月3日のものです。

## 再生回数Top5

|  URL  |  タイトル  |  再生回数  |
| ---- | ---- | ---- |
| https://www.youtube.com/watch?v=dKB8Z_KJG3I |  <a href="">煩悩ネタ！『嘘つき通す奴』</a>  |  1398004  |
|  <a href="https://www.youtube.com/watch?v=f_r4F6QTZx4">煩悩ネタ！『第一声、間違えて、友達失う奴』</a>  |  973753  |
|  <a href="https://www.youtube.com/watch?v=hFnhrRqDrA0">煩悩ネタ！『オキンタマデカ男って奴』</a>  |  969133  |
|  <a href="https://www.youtube.com/watch?v=TJDTMyVX4fw">煩悩ネタ！『あだ名コロコロ変わる奴』</a>  |  510030  |
|  <a href="https://www.youtube.com/watch?v=wg-hRa2CFEs">煩悩ネタ！『大学の入学を辞めただけの奴』</a>  |  460008  |


## 高評価数Top5

|  タイトル  |  高評価数  |
| ---- | ---- |
|  <a href="https://www.youtube.com/watch?v=hFnhrRqDrA0">煩悩ネタ！『オキンタマデカ男って奴』</a>  |  16882  |
|  <a href="https://www.youtube.com/watch?v=dKB8Z_KJG3I">煩悩ネタ！『嘘つき通す奴』</a>  |  14939  |
|  <a href="https://www.youtube.com/watch?v=f_r4F6QTZx4">煩悩ネタ！『第一声、間違えて、友達失う奴』</a>  |  12595  |
|  <a href="https://www.youtube.com/watch?v=Au-mjjkqe7g">煩悩ネタ！『声の周波数のせいで離れたら英語に聞こえる奴』</a>  |  7613  |
|  <a href="https://www.youtube.com/watch?v=TJDTMyVX4fw">煩悩ネタ！『あだ名コロコロ変わる奴』</a>  |  6528  |

## コメント数Top5
|  タイトル  |  コメント数  |
| ---- | ---- |
|  <a href="https://www.youtube.com/watch?v=xZDBDH8eseY">煩悩ネタ！『107本を通して感じたことをコントにする奴』</a>  |  862  |
|  <a href="https://www.youtube.com/watch?v=wg-hRa2CFEs">煩悩ネタ！『大学の入学を辞めただけの奴』</a>  |  763  |
|  <a href="https://www.youtube.com/watch?v=f_r4F6QTZx4">煩悩ネタ！『第一声、間違えて、友達失う奴』</a>  |  746  |
|  <a href="https://www.youtube.com/watch?v=hFnhrRqDrA0">煩悩ネタ！『オキンタマデカ男って奴』</a>  |  720  |
|  <a href="https://www.youtube.com/watch?v=dKB8Z_KJG3I">煩悩ネタ！『嘘つき通す奴』</a>  |  605  |



# 技術
Go + YoutubeAPI で作成しているシンプルなプログラムです。

## フローチャート図

![処理フロー図](/images/others/youtube-playlist-analyze.drawio.png)

## バリデーションチェック
configファイルの読み込みや、引数チェックを行います。
ここでエラーが発生した場合、そのままプログラムが終了します。

## playlist情報を取得
まずは引数のプレーリストIDを基に、リストに含まれる動画情報を抽出していきます。
[Playlists: list API](https://developers.google.com/youtube/v3/docs/playlists/list?hl=ja)を使用します。

必須パラメータ`part`にはsnippetを指定します。
これにより、item.snippet.recourceId.videoId(動画ID)が返されます。

1回のリクエストにつき最大50件までしか動画情報を取得できないため、51件以上の動画を含む動画リストを指定する場合は複数回Playlists list APIを呼ぶ必要があります。この時、レスポンス`nextPageToken`を使用することでページネーションと最終ページ処理を合わせて行うことができます。

ここでYoutube list APIから200以外のレスポンスが返された場合、エラーとしてプログラムを終了します。


:::details Response例
```JSON
{
    "kind": "youtube#playlistItemListResponse",
    "etag": "w0REtXALenja70i0bGWZOdyycVo",
    "nextPageToken": "EAAaBlBUOkNBVQ",
    "items": [
        {
            "kind": "youtube#playlistItem",
            "etag": "ZGt5Se0wR8z49JKiKPJo60yGFcg",
            "id": "UExSZGlhYW5LQUZRbDVFUkRnSkh4MlpSS0NjSWwtSThmei41NkI0NEY2RDEwNTU3Q0M2",
            "snippet": {
                "publishedAt": "2020-03-28T03:27:55Z",
                "channelId": "UChwgNUWPM-ksOP3BbfQHS5Q",
                "title": "煩悩ネタ！『一服する奴』",
                "description": "【福徳】\n〈１０８本ライブをアップした翌日に結局１０８本を１本ずつアップする奴〉です。\n本当は１０８本連続で見て欲しいけど個別で見れるようにもなりました。\n川上くんというスタッフさんが地獄の編集を乗り越えて個別のアップが実現しました。\n川上くんお疲れさま。\nありがとう。\n\n【後藤】\n結局個別にも観てほしいんです。\nだって何時間もある動画観るなんて辛いじゃないですか。\n皆さんにはそんな辛い思いしてほしくないんです！\n分刻みの自動アップができないので川上くんが手動でアップしていきます。\n皆さま是非観て下さい！\n川上くんありがとう！\n\n単独ライブツアー「JARUJARU TOWER 2020－ジャルってんじゃねぇよ－」\n大阪公演（5月7日～10日 シアタードラマシティ）\n福岡公演（5月16日、5月17日 西鉄ホール）\n名古屋公演（5月23日 ウインクあいち）\n仙台公演（5月29日 電力ホール）\n東京公演（5月30日～6月2日　ルミネtheよしもと）\n\n【ジャルジャル公式HP】http://jarujaru.com\n【ジャルジャル公式LINE】https://live.line.me/r/channels/121\n【ジャルジャルTwitter】https://twitter.com/jarujaru12th?lang=ja\n【ジャルジャルに興味ある奴】https://www.jarujaru-kyomi.com",
                "thumbnails": {
                    "default": {
                        "url": "https://i.ytimg.com/vi/1TzGN4VSGR0/default.jpg",
                        "width": 120,
                        "height": 90
                    },
                    "medium": {
                        "url": "https://i.ytimg.com/vi/1TzGN4VSGR0/mqdefault.jpg",
                        "width": 320,
                        "height": 180
                    },
                    "high": {
                        "url": "https://i.ytimg.com/vi/1TzGN4VSGR0/hqdefault.jpg",
                        "width": 480,
                        "height": 360
                    },
                    "standard": {
                        "url": "https://i.ytimg.com/vi/1TzGN4VSGR0/sddefault.jpg",
                        "width": 640,
                        "height": 480
                    },
                    "maxres": {
                        "url": "https://i.ytimg.com/vi/1TzGN4VSGR0/maxresdefault.jpg",
                        "width": 1280,
                        "height": 720
                    }
                },
                "channelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "playlistId": "PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz",
                "position": 0,
                "resourceId": {
                    "kind": "youtube#video",
                    "videoId": "1TzGN4VSGR0"
                },
                "videoOwnerChannelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "videoOwnerChannelId": "UChwgNUWPM-ksOP3BbfQHS5Q"
            }
        },
        {
            "kind": "youtube#playlistItem",
            "etag": "gde9WUM8i8GK9V6FskQna7q0_m8",
            "id": "UExSZGlhYW5LQUZRbDVFUkRnSkh4MlpSS0NjSWwtSThmei4yODlGNEE0NkRGMEEzMEQy",
            "snippet": {
                "publishedAt": "2020-03-28T03:41:01Z",
                "channelId": "UChwgNUWPM-ksOP3BbfQHS5Q",
                "title": "煩悩ネタ！『本名、ふざけてると思われる奴』",
                "description": "【福徳】\n〈１０８本ライブをアップした翌日に結局１０８本を１本ずつアップする奴〉です。\n本当は１０８本連続で見て欲しいけど個別で見れるようにもなりました。\n川上くんというスタッフさんが地獄の編集を乗り越えて個別のアップが実現しました。\n川上くんお疲れさま。\nありがとう。\n\n【後藤】\n結局個別にも観てほしいんです。\nだって何時間もある動画観るなんて辛いじゃないですか。\n皆さんにはそんな辛い思いしてほしくないんです！\n分刻みの自動アップができないので川上くんが手動でアップしていきます。\n皆さま是非観て下さい！\n川上くんありがとう！\n\n単独ライブツアー「JARUJARU TOWER 2020－ジャルってんじゃねぇよ－」\n大阪公演（5月7日～10日 シアタードラマシティ）\n福岡公演（5月16日、5月17日 西鉄ホール）\n名古屋公演（5月23日 ウインクあいち）\n仙台公演（5月29日 電力ホール）\n東京公演（5月30日～6月2日　ルミネtheよしもと）\n\n【ジャルジャル公式HP】http://jarujaru.com\n【ジャルジャル公式LINE】https://live.line.me/r/channels/121\n【ジャルジャルTwitter】https://twitter.com/jarujaru12th?lang=ja\n【ジャルジャルに興味ある奴】https://www.jarujaru-kyomi.com",
                "thumbnails": {
                    "default": {
                        "url": "https://i.ytimg.com/vi/oD6jlChFRvs/default.jpg",
                        "width": 120,
                        "height": 90
                    },
                    "medium": {
                        "url": "https://i.ytimg.com/vi/oD6jlChFRvs/mqdefault.jpg",
                        "width": 320,
                        "height": 180
                    },
                    "high": {
                        "url": "https://i.ytimg.com/vi/oD6jlChFRvs/hqdefault.jpg",
                        "width": 480,
                        "height": 360
                    },
                    "standard": {
                        "url": "https://i.ytimg.com/vi/oD6jlChFRvs/sddefault.jpg",
                        "width": 640,
                        "height": 480
                    },
                    "maxres": {
                        "url": "https://i.ytimg.com/vi/oD6jlChFRvs/maxresdefault.jpg",
                        "width": 1280,
                        "height": 720
                    }
                },
                "channelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "playlistId": "PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz",
                "position": 1,
                "resourceId": {
                    "kind": "youtube#video",
                    "videoId": "oD6jlChFRvs"
                },
                "videoOwnerChannelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "videoOwnerChannelId": "UChwgNUWPM-ksOP3BbfQHS5Q"
            }
        },
        {
            "kind": "youtube#playlistItem",
            "etag": "pLFhiYaQleEQGGMsYjYOSjc_aw8",
            "id": "UExSZGlhYW5LQUZRbDVFUkRnSkh4MlpSS0NjSWwtSThmei4wMTcyMDhGQUE4NTIzM0Y5",
            "snippet": {
                "publishedAt": "2020-03-28T03:58:52Z",
                "channelId": "UChwgNUWPM-ksOP3BbfQHS5Q",
                "title": "煩悩ネタ！『声の周波数のせいで離れたら英語に聞こえる奴』",
                "description": "【福徳】\n〈１０８本ライブをアップした翌日に結局１０８本を１本ずつアップする奴〉です。\n本当は１０８本連続で見て欲しいけど個別で見れるようにもなりました。\n川上くんというスタッフさんが地獄の編集を乗り越えて個別のアップが実現しました。\n川上くんお疲れさま。\nありがとう。\n\n【後藤】\n結局個別にも観てほしいんです。\nだって何時間もある動画観るなんて辛いじゃないですか。\n皆さんにはそんな辛い思いしてほしくないんです！\n分刻みの自動アップができないので川上くんが手動でアップしていきます。\n皆さま是非観て下さい！\n川上くんありがとう！\n\n単独ライブツアー「JARUJARU TOWER 2020－ジャルってんじゃねぇよ－」\n大阪公演（5月7日～10日 シアタードラマシティ）\n福岡公演（5月16日、5月17日 西鉄ホール）\n名古屋公演（5月23日 ウインクあいち）\n仙台公演（5月29日 電力ホール）\n東京公演（5月30日～6月2日　ルミネtheよしもと）\n\n【ジャルジャル公式HP】http://jarujaru.com\n【ジャルジャル公式LINE】https://live.line.me/r/channels/121\n【ジャルジャルTwitter】https://twitter.com/jarujaru12th?lang=ja\n【ジャルジャルに興味ある奴】https://www.jarujaru-kyomi.com",
                "thumbnails": {
                    "default": {
                        "url": "https://i.ytimg.com/vi/Au-mjjkqe7g/default.jpg",
                        "width": 120,
                        "height": 90
                    },
                    "medium": {
                        "url": "https://i.ytimg.com/vi/Au-mjjkqe7g/mqdefault.jpg",
                        "width": 320,
                        "height": 180
                    },
                    "high": {
                        "url": "https://i.ytimg.com/vi/Au-mjjkqe7g/hqdefault.jpg",
                        "width": 480,
                        "height": 360
                    },
                    "standard": {
                        "url": "https://i.ytimg.com/vi/Au-mjjkqe7g/sddefault.jpg",
                        "width": 640,
                        "height": 480
                    },
                    "maxres": {
                        "url": "https://i.ytimg.com/vi/Au-mjjkqe7g/maxresdefault.jpg",
                        "width": 1280,
                        "height": 720
                    }
                },
                "channelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "playlistId": "PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz",
                "position": 2,
                "resourceId": {
                    "kind": "youtube#video",
                    "videoId": "Au-mjjkqe7g"
                },
                "videoOwnerChannelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "videoOwnerChannelId": "UChwgNUWPM-ksOP3BbfQHS5Q"
            }
        },
        {
            "kind": "youtube#playlistItem",
            "etag": "v0Op1479WPP-oa4lwHoRSF2g760",
            "id": "UExSZGlhYW5LQUZRbDVFUkRnSkh4MlpSS0NjSWwtSThmei41MjE1MkI0OTQ2QzJGNzNG",
            "snippet": {
                "publishedAt": "2020-03-28T04:16:27Z",
                "channelId": "UChwgNUWPM-ksOP3BbfQHS5Q",
                "title": "煩悩ネタ！『ちょけた出しもんして後悔する奴』",
                "description": "【福徳】\n〈１０８本ライブをアップした翌日に結局１０８本を１本ずつアップする奴〉です。\n本当は１０８本連続で見て欲しいけど個別で見れるようにもなりました。\n川上くんというスタッフさんが地獄の編集を乗り越えて個別のアップが実現しました。\n川上くんお疲れさま。\nありがとう。\n\n【後藤】\n結局個別にも観てほしいんです。\nだって何時間もある動画観るなんて辛いじゃないですか。\n皆さんにはそんな辛い思いしてほしくないんです！\n分刻みの自動アップができないので川上くんが手動でアップしていきます。\n皆さま是非観て下さい！\n川上くんありがとう！\n\n単独ライブツアー「JARUJARU TOWER 2020－ジャルってんじゃねぇよ－」\n大阪公演（5月7日～10日 シアタードラマシティ）\n福岡公演（5月16日、5月17日 西鉄ホール）\n名古屋公演（5月23日 ウインクあいち）\n仙台公演（5月29日 電力ホール）\n東京公演（5月30日～6月2日　ルミネtheよしもと）\n\n【ジャルジャル公式HP】http://jarujaru.com\n【ジャルジャル公式LINE】https://live.line.me/r/channels/121\n【ジャルジャルTwitter】https://twitter.com/jarujaru12th?lang=ja\n【ジャルジャルに興味ある奴】https://www.jarujaru-kyomi.com",
                "thumbnails": {
                    "default": {
                        "url": "https://i.ytimg.com/vi/-FCCCaYJQgU/default.jpg",
                        "width": 120,
                        "height": 90
                    },
                    "medium": {
                        "url": "https://i.ytimg.com/vi/-FCCCaYJQgU/mqdefault.jpg",
                        "width": 320,
                        "height": 180
                    },
                    "high": {
                        "url": "https://i.ytimg.com/vi/-FCCCaYJQgU/hqdefault.jpg",
                        "width": 480,
                        "height": 360
                    },
                    "standard": {
                        "url": "https://i.ytimg.com/vi/-FCCCaYJQgU/sddefault.jpg",
                        "width": 640,
                        "height": 480
                    },
                    "maxres": {
                        "url": "https://i.ytimg.com/vi/-FCCCaYJQgU/maxresdefault.jpg",
                        "width": 1280,
                        "height": 720
                    }
                },
                "channelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "playlistId": "PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz",
                "position": 3,
                "resourceId": {
                    "kind": "youtube#video",
                    "videoId": "-FCCCaYJQgU"
                },
                "videoOwnerChannelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "videoOwnerChannelId": "UChwgNUWPM-ksOP3BbfQHS5Q"
            }
        },
        {
            "kind": "youtube#playlistItem",
            "etag": "F-ay-XPK1aVPYaxNO8LcIU2eLzs",
            "id": "UExSZGlhYW5LQUZRbDVFUkRnSkh4MlpSS0NjSWwtSThmei4wOTA3OTZBNzVEMTUzOTMy",
            "snippet": {
                "publishedAt": "2020-03-28T04:40:10Z",
                "channelId": "UChwgNUWPM-ksOP3BbfQHS5Q",
                "title": "煩悩ネタ！『体操走りな奴』",
                "description": "【福徳】\n〈１０８本ライブをアップした翌日に結局１０８本を１本ずつアップする奴〉です。\n本当は１０８本連続で見て欲しいけど個別で見れるようにもなりました。\n川上くんというスタッフさんが地獄の編集を乗り越えて個別のアップが実現しました。\n川上くんお疲れさま。\nありがとう。\n\n【後藤】\n結局個別にも観てほしいんです。\nだって何時間もある動画観るなんて辛いじゃないですか。\n皆さんにはそんな辛い思いしてほしくないんです！\n分刻みの自動アップができないので川上くんが手動でアップしていきます。\n皆さま是非観て下さい！\n川上くんありがとう！\n\n単独ライブツアー「JARUJARU TOWER 2020－ジャルってんじゃねぇよ－」\n大阪公演（5月7日～10日 シアタードラマシティ）\n福岡公演（5月16日、5月17日 西鉄ホール）\n名古屋公演（5月23日 ウインクあいち）\n仙台公演（5月29日 電力ホール）\n東京公演（5月30日～6月2日　ルミネtheよしもと）\n\n【ジャルジャル公式HP】http://jarujaru.com\n【ジャルジャル公式LINE】https://live.line.me/r/channels/121\n【ジャルジャルTwitter】https://twitter.com/jarujaru12th?lang=ja\n【ジャルジャルに興味ある奴】https://www.jarujaru-kyomi.com",
                "thumbnails": {
                    "default": {
                        "url": "https://i.ytimg.com/vi/QtmM8l8E3HM/default.jpg",
                        "width": 120,
                        "height": 90
                    },
                    "medium": {
                        "url": "https://i.ytimg.com/vi/QtmM8l8E3HM/mqdefault.jpg",
                        "width": 320,
                        "height": 180
                    },
                    "high": {
                        "url": "https://i.ytimg.com/vi/QtmM8l8E3HM/hqdefault.jpg",
                        "width": 480,
                        "height": 360
                    },
                    "standard": {
                        "url": "https://i.ytimg.com/vi/QtmM8l8E3HM/sddefault.jpg",
                        "width": 640,
                        "height": 480
                    },
                    "maxres": {
                        "url": "https://i.ytimg.com/vi/QtmM8l8E3HM/maxresdefault.jpg",
                        "width": 1280,
                        "height": 720
                    }
                },
                "channelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "playlistId": "PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz",
                "position": 4,
                "resourceId": {
                    "kind": "youtube#video",
                    "videoId": "QtmM8l8E3HM"
                },
                "videoOwnerChannelTitle": "ジャルジャルタワー JARUJARU TOWER",
                "videoOwnerChannelId": "UChwgNUWPM-ksOP3BbfQHS5Q"
            }
        }
    ],
    "pageInfo": {
        "totalResults": 108,
        "resultsPerPage": 5
    }
}

```
:::

## playlist情報からvideoリストを取得
上記で取得した動画IDを基に、高評価数や視聴回数を取得していきます。
[Videos: list API](https://developers.google.com/youtube/v3/docs/videos/list?hl=ja)を使用します。

必須のクエリパラメータ`part`にはstatisticsを指定します。
これにより、items.statistics(高評価数、視聴回数、お気に入り回数、コメント数)が返されます。

また、クエリパラメータ`id`には複数idを指定することができます。
playlistidで取得したvideoIdを全て指定し、 VideosListAPIを呼び出します。

ここでも同様に、VideosListAPIから200以外のレスポンスが返された場合、エラーとしてプログラムを終了します。

:::details Response例
```JSON
{
    "kind": "youtube#videoListResponse",
    "etag": "xbXG2Hba3mcYnrNlb8gGLOysSpQ",
    "items": [
        {
            "kind": "youtube#video",
            "etag": "pEaHUXOSgzZokZJMj4XRulBm2HY",
            "id": "1yUoU1Gr_zE",
            "statistics": {
                "viewCount": "58979",
                "likeCount": "687",
                "favoriteCount": "0",
                "commentCount": "46"
            }
        },
        {
            "kind": "youtube#video",
            "etag": "VYbCsTUy_BfVJhQNPP-ehN0J5yM",
            "id": "xZDBDH8eseY",
            "statistics": {
                "viewCount": "113829",
                "likeCount": "2238",
                "favoriteCount": "0",
                "commentCount": "862"
            }
        },
        {
            "kind": "youtube#video",
            "etag": "pEaHUXOSgzZokZJMj4XRulBm2HY",
            "id": "1yUoU1Gr_zE",
            "statistics": {
                "viewCount": "58979",
                "likeCount": "687",
                "favoriteCount": "0",
                "commentCount": "46"
            }
        },
        {
            "kind": "youtube#video",
            "etag": "VYbCsTUy_BfVJhQNPP-ehN0J5yM",
            "id": "xZDBDH8eseY",
            "statistics": {
                "viewCount": "113829",
                "likeCount": "2238",
                "favoriteCount": "0",
                "commentCount": "862"
            }
        }
    ],
    "pageInfo": {
        "totalResults": 4,
        "resultsPerPage": 4
    }
}
```
:::

## video情報を取得
VideosListAPIのレスポンスからvidwCount, likeCount, commentCountを取得します。
favoriteCountに関しては常に0を返すようです。公式ソースを見つけることができませんでしたが、favorite機能はリスト機能に置き換わった可能性が高いです。

一応ChatGPT-4に聞いてみました。
```
User:
YoutubeAPI v3のvideo list api(https://www.googleapis.com/youtube/v3/videos)に関して、レスポンス内にあるfavoriteCountが常に0となるのはなぜですか？

ChatGPT:
YouTube API v3のvideo list API (https://www.googleapis.com/youtube/v3/videos) に関して、レスポンス内にあるfavoriteCountが常に0となる理由は、YouTubeがお気に入り機能を廃止し、代わりに「いいね」や「プレイリストに追加」といった機能に置き換えたためです。
かつては、ユーザーは動画をお気に入りに追加することができましたが、現在はその機能が利用できなくなっています。そのため、APIのレスポンスでfavoriteCountが常に0を返すようになっています。
代わりに、動画に関する他の詳細情報を取得したい場合は、APIレスポンス内の「likeCount」、「dislikeCount」、「viewCount」、「commentCount」などの項目を使用してください。これらは、引き続き有効なデータを提供しています。
```


# おわりに
思ったより簡単にYoutubeの情報を取得することができて驚きました。
YoutubePlaylistAnalyzeツールを使用し、改善点や感想をいただけると嬉しいです。

YoutubeAPIで取得できないジャルジャルのデータはこちらも見てみてください。
https://natalie.mu/owarai/pp/jarujaru02
