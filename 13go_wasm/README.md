# 始めに
これはGoでWebAssemblyを試したものです。

# ディレクトリの説明
ディレクトリ構成
```
.
├── README.md
├── server
│   ├── public
│   │   ├── go.wasm           // web assembly(Goビルド)
│   │   ├── index.html
│   │   ├── tinygo.wasm       // web assembly(TinyGoビルド)
│   │   ├── tinywasm_exec.js  // TinyGo用のjsファイル
│   │   └── wasm_exec.js      // Go用のjsファイル
│   └── server.go
├── wasm
│   └── main.go
└── wasm-tiny
    └── main.go
```
## server
ウェブサーバと静的ファイルのホスティング
ウェブサーバ→server.go
静的ファイル→public/

## wasm
Goでバイナリファイルの生成

## wasm-tiny
TonyGoでバイナリファイルの生成

1. wasm(wasm-tiny)/main.goの記述
2. ビルド
    ```
    通常:GOOS=js GOARCH=wasm go build -o go.wasm
    Tiny:GOOS=js GOARCH=wasm tinygo build -o tinygo.wasm
    ```
    出力されたファイルを`server/public/`に移動
3. public/index.htmlの作成
4. wasm_exec.jsを取得
    ```
    # GOのwasmを実行する場合は、以下のファイルをダウンロードする
    $ wget https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js

    # TinyGOのwasmを実行する場合は、以下のファイルをダウンロードする
    $ wget https://raw.githubusercontent.com/tinygo-org/tinygo/master/targets/wasm_exec.js
    ```
5. index.htmlスクリプトタグ内のfetch関数の引数を変更(***の部分)
    ```js
    WebAssembly.instantiateStreaming(fetch("***.wasm"), go.importObject).then((result) => {
    ```
6. サーバを起動

# バイナリファイルのサイズ
```
$ ls -l
total 4928
-rwxr-xr-x  1 mbp  staff  2265398 11 10 10:48 go.wasm*
-rw-r--r--  1 mbp  staff      338 11 10 11:08 index.html
-rwxr-xr-x  1 mbp  staff   210803 11 10 10:48 tinygo.wasm*
-rw-r--r--  1 mbp  staff    15466 11 10 11:06 tinywasm_exec.js
-rw-r--r--  1 mbp  staff    17255 11 10 10:24 wasm_exec.js
```
大規模なアプリケーションでもない限りGoでそのままビルドしたwasmはファイルはサイズが大きいと思う。
今のところTinyGoを使う前提でGoでWebAssemblyを使った方が良さそう。