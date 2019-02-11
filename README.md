# wlx

## なにこれ

YAMAHA の無線 LAN アクセスポイント [WLX202](https://network.yamaha.com/products/wireless_lan/wlx202/index) の一部のシステム情報を REST API で返却するシンプルなサーバーです.

以下, 本サーバーで取得可能な情報です.

* CPU 稼働率 
* メモリ使用率
* 2.4GHz クライアント接続数
* 5GHz クライアント帯接続数

## 使い方

### WLX202 の設定

* 管理 Web コンソールのログインユーザー名とパスワードを設定する
* 機器の IP アドレスを設定する

### wlx を起動する

WLX202 の管理 Web コンソールにアクセス可能な端末でサーバーを起動します.

```sh
./wlx &
```

デフォルトでサーバーは localhost の 20200 番で Listen するので, 以下のようにアクセスすることで, 対象の機器の情報を取得することが出来ます.

```sh
curl -s -XPOST http://localhost:20200/wlx/${対象機器の IP アドレス} \
  -d "user=${ログインユーザー}" \
  -d "pass=${ログインパスワード}"
```

以下のようなレスポンスが返却されます.

```json
{
  "results": [
    {
      "ip_address": "対象機器の IP アドレス",
      "machine_name": "WLX202_XXXXXXXXXXX",
      "cpu_utilization": "10",
      "memory_usage": "39",
      "connected": "10",
      "connected_5g": "15"
    }
  ]
}
```

以下のような内容になっています.

* ip_address は WLX202 本体の機器 IP アドレス
* machine_name は機器の名称
* cpu_utilization は CPU の稼働率
* memory_usage はメモリの使用率
* connected は 2.4GHz 帯域に接続しているクライアント数
* connected_5g は 5GHz 帯域に接続しているクライアント数

## なんで作ったのか

* WLX202 は設定もシンプルで YAMAHA のルーターを使っていれば, 本当にサイコーな無線 LAN アクセスポイントなのに, SNMP でクライアントの接続数が取れないという残念なポイントを解決したかった
* Golang で API サーバーを実装してみたかった
