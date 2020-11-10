# 2020/10/14最新versionを取得
FROM golang:1.15.2-alpine
# アップデートとgit, bashのインストール
RUN apk update && apk add git && apk add tree
# appディレクトリの作成
RUN mkdir /go/src/app
# ワーキングディレクトリの設定
WORKDIR /go/src/app
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/src/app

RUN go get github.com/oxequa/realize \
    && go get github.com/go-sql-driver/mysql
# この設定がないとgo testが出来ない。
ENV CGO_ENABLED=0