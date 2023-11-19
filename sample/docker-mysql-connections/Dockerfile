# ベースイメージの指定
FROM golang:latest


# 作業ディレクトリの設定
WORKDIR /api/
COPY go.mod go.sum /api/

# ソースコードをコピー
COPY . /api/

# ビルド
RUN go build -o main .

# アプリケーションの起動
CMD ["./main"]