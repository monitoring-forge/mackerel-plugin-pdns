# mackerel-plugin-pdns

PowerDNS Authoritative Nameserver のメトリクスを Mackerel に投稿するプラグインです。
[hanazuki/mackerel-plugin-pdns](https://github.com/hanazuki/mackerel-plugin-pdns)（Ruby 製）の Go 移植版です。

## インストール

以下のいずれかの方法でインストールしてください。

- GitHub Releases ページからバイナリをダウンロードする
- `mkr plugin install monitoring-forge/mackerel-plugin-pdns`

## 使い方

```sh
mackerel-plugin-pdns [オプション]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-v`, `--version` | - | バージョン情報を表示して終了します |
| `--prefix` | `pdns` | メトリクス名の接頭辞を指定します |
| `--control-command` | `/usr/bin/pdns_control` | `pdns_control` コマンドへのパスを指定します |

### 実行例

```sh
# デフォルト設定で実行
mackerel-plugin-pdns

# 接頭辞を変更する
mackerel-plugin-pdns --prefix powerdns

# pdns_control のパスを変更する
mackerel-plugin-pdns --control-command /usr/local/bin/pdns_control
```

## mackerel-agent への設定例

`/etc/mackerel-agent/mackerel-agent.conf` に以下のように記述します。

```toml
[plugin.metrics.pdns]
command = ["/usr/local/bin/mackerel-plugin-pdns", "--control-command", "/usr/bin/pdns_control"]
```

接頭辞を変更する場合は、`plugin.metrics.<任意の名前>` のセクション名と `--prefix` オプションを合わせてください。

```toml
[plugin.metrics.powerdns]
command = ["/usr/local/bin/mackerel-plugin-pdns", "--prefix", "powerdns"]
```

## 取得できるメトリクス

`pdns_control show '*'` の出力からメトリクスを取得します。

### グラフ定義

| グラフ名 | 単位 | 内容 |
|----------|------|------|
| `dnsupdate` | integer | Dynamic DNS Update の応答数・変更数・クエリ数・拒否数 |
| `notifications` | integer | DNS Notification の受信数 |
| `packetcache` | integer | Packet Cache のヒット数・ミス数 |
| `query-cache` | integer | Query Cache のヒット数・ミス数 |
| `cache-size` | integer | Packet cache / Key cache / Signature cache / Metadata cache のサイズ |
| `fails` | integer | SERVFAIL / Corrupt / Timeout / バックエンド過負荷によるドロップ数 |
| `backend` | integer | バックエンドへのクエリ数 |
| `tcp-connection` | integer | TCP 接続数とファイルディスクリプタ使用量 |
| `signatures` | integer | 作成した DNSSEC Signature 数 |
| `latency` | integer | レイテンシ（マイクロ秒） |
| `qsize` | integer | キューサイズ |
| `answers` | integer | TCP / UDP / TCP4 / UDP4 / TCP6 / UDP6 別の応答数 |
| `queries` | integer | TCP / UDP / TCP4 / UDP4 / TCP6 / UDP6 別のクエリ数 |
| `answer-bytes` | integer | TCP / UDP / TCP4 / UDP4 / TCP6 / UDP6 別の応答バイト数 |
| `memory` | bytes | 実メモリ使用量 |
| `cpu` | integer | ユーザー時間・システム時間（ミリ秒） |

### メトリクス名の例

接頭辞が `pdns` の場合、メトリクス名は以下のようになります。

- `custom.pdns.dnsupdate.dnsupdate-answers`
- `custom.pdns.packetcache.packetcache-hit`
- `custom.pdns.tcp-connection.open-tcp-connections`
- `custom.pdns.memory.real-memory-usage`

## 注意事項

- 本プラグインは PowerDNS Authoritative Nameserver を対象としています。Recursor には対応していません。
- `pdns_control` を実行する権限が必要です。
- 監視対象の PowerDNS サーバー上で実行するか、`pdns_control` が実行可能な環境で実行してください。

## License

MIT