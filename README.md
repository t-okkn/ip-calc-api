# ip-calc-api

IPアドレス計算アプリのAPI（バックエンド）

## 事前準備

DB：複数対応したいところですが、現状**MySQL（SSLなし、証明書なし）のみ対応**しています。
DB migrate：[golang-migrate](https://github.com/golang-migrate/migrate)をインストールしておいてください。

## インストール方法

お好きなDB名でDATABASEを作成しておいてください。
※エンジンはInnoDB、文字コードはutf8mb4を利用します。

```sql
CREATE DATABASE IF NOT EXISTS `{{ DB名 }}`;
```

次に `ip-calc-api.sql` 内にある `connect.toml` にDBの接続情報を記載してください。
接続するUSERは先程作成したDBに対してすべての権限を付与しておいてください。
※現状、**TLSの設定は触らないでください。**
　接続できません。

その後、ターミナルに戻ってインストール作業を実施してください。
インストール先等のディレクトリの変更を実施したい場合は、お手数ですが `Makefile` を直接編集してください。
configure ファイルは気がむいたら作るかもしれません。

```terminal
$ make init
$ make build
$ make db
$ sudo make install
$ sudo make start
```

待ち受けポートを変更したい場合、 `main.go` の `LISTEN_PORT` を変更してください。

