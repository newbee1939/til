- goroutineを理解する
    - バファつきチャネル
    - goroutine
    - channel
    - context
- https://httpd.apache.org/docs/current/ja/howto/htaccess.html
    - いつ .htaccess ファイルを使う(使わない)か。 ¶
- Goの特徴
    - シングルバイナリ
    - クロスコンパイル
- Cloud SQL IAM認証
    - Googleの推奨
- Go: context
- Goのゴルーチン、コンテキスト、チャネルの関係性
- Linux: プロセススケジューラ
    - 本の内容
- bashのパラメータ展開の記法について。${var:-default}みたいな
- Goルーチン
- Cloudflare AI Gateway
- GO:言語 selectブロック
- SBOM
    - フレームワークとかベストプラクティスある？
- QUERYメソッド
    - HTTPの最新の仕様
- 新卒一括採用はもう限界か　AI時代に“内定が出ない”人材とは
    - https://www.itmedia.co.jp/business/articles/2607/02/news007.html

- groupsコマンド
- 証明書の仕組み。ハンドシェイク。SNI

- go mod download

## 1. SREのコアコンセプト（マインドセットと指標）

* **SLI / SLO / SLA**
    * **SLI（サービスレベル指標）**: 信頼性を測るための具体的なメトリクス（例：リクエストの成功率、レイテンシなど）。
    * **SLO（サービスレベル目標）**: チーム内で目指す信頼性のターゲット（例：過去30日間でリクエストの99.9%が200ms以内に返る）。
    * **SLA（サービスレベル契約）**: ビジネス上の顧客との約束。下回ると返金などのペナルティが発生する。通常、SLOはSLAよりも厳しく設定します。
* **エラーバジェット（Error Budget）**
    * 「100% - SLO」で計算される、「システムが許容できる停止時間（または失敗リクエスト）」の残高。
    * これが残っているうちは新機能のリリースを攻め、バジェットを使い果たすと新機能リリースを止め、信頼性向上（バグ修正やインフラ改善）にリソースを集中させるという意思決定の枠組みです。
* **トイル（Toil）の削減**
    * 手作業、繰り返される、自動化可能、戦術的（長期的な価値がない）な作業のこと。SREは業務時間の50%以上をトイル削減（自動化や仕組み化）にあてることが求められます。

## 2. Linux / OSの基礎とトラブルシューティング
トラブル発生時にインフラの低レイヤーで何が起きているかを調査する能力が問われます。

* **Linuxカーネルとリソース監視**
    * CPUの「ロードアベレージ」が意味するもの（実行待ち・ディスクI/O待ちのプロセス数）。
    * メモリの `free`, `available` の違いや、`OOM Killer`（メモリ枯渇時にカーネルがプロセスを強制終了する仕組み）の挙動。
* **プロセスの状態とデバッグ**
    * `top`, `htop`, `vmstat`, `iostat`, `strace`, `lsof` などのコマンドをどう使ってボトルネックを特定するか。
    * 「ゾンビプロセス」や「孤児プロセス」の違い。
* **ファイルシステム**
    * iノード（Inode）の枯渇問題（ディスク容量はあるのにファイルが作れない現象）。

## 3. ネットワーク（プロトコルとデバッグ）
分散システムにおける通信のエラーや遅延の原因を特定するために、必須の知識です。

* **TCP/IP 階層モデルと3ウェイ・ハンドシェイク**
    * 接続確立（SYN, SYN-ACK, ACK）と切断（FIN, ACK）の流れ。
    * `TIME_WAIT` 状態のソケットが大量発生したときの原因と対策。
* **DNS（Domain Name System）**
    * ドメイン名からIPアドレスを引く仕組み（再帰クエリと反復クエリ）。
    * Aレコード、CNAMEレコード、TTL（キャッシュ有効期限）の意味。
* **HTTP/HTTPS**
    * HTTP/1.1 と HTTP/2, HTTP/3（QUIC）の違い（マルチプレクシングなどによる高速化）。
    * TLSハンドシェイクの仕組み。
* **トラブルシューティングコマンド**
    * `ping`, `traceroute`（どこでパケットが落ちているか）、`nslookup`/`dig`（DNS確認）、`tcpdump`/`wireshark`（パケットキャプチャ）、`curl -v`。

## 4. クラウド・コンテナ・IaC（現代のインフラ技術）
具体的なプラットフォーム（AWS, GCPなど）や、コンテナオーケストレーションの理解度です。

* **コンテナ技術（Docker / Kubernetes）**
    * Dockerの仕組み（Linuxの `Namespace` による隔離と `cgroups` によるリソース制限）。
    * Kubernetes（K8s）の基本リソース（Pod, Deployment, Service, Ingress）の役割。
    * K8sの自己修復機能（Liveness Probe / Readiness Probe の違い）。
* **IaC（Infrastructure as Code）**
    * Terraformなどを用いたインフラのコード化。
    * 「イミュータブル・インフラストラクチャ（不変のインフラ）」のメリット（環境の冪等性の担保）。
* **CI/CD**
    * GitHub ActionsやArgoCDなどを使った、安全かつ自動化されたデプロイ戦略（カナリアデプロイ、ブルーグリーンデプロイ）。

## 5. 分散システムとオブザーバビリティ（可観測性）
大規模なシステムを安定して運用するための設計と監視の知識です。

* **オブザーバビリティの3つの柱**
    * **Metrics（メトリクス）**: CPU使用率やリクエスト数などの数値データ（Prometheus, Datadogなど）。
    * **Logs（ログ）**: アプリケーションやシステムが出力するテキストリリー（ELKスタック, CloudWatch Logsなど）。
    * **Traces（トレース）**: マイクロサービス間でリクエストがどう遷移したかの経路と時間（Jaeger, OpenTelemetryなど）。
* **分散システムの信頼性設計（レジリエンス）**
    * **サーキットブレーカー**: 依存する外部サービスが落ちた際、引きずられて自社システムが全滅するのを防ぐために通信を遮断する仕組み。
    * **リトライとエクスポネンシャルバックオフ**: 失敗したリクエストを再試行する際、間隔を徐々に広げ、さらに「Jitter（ゆらぎ）」を加えてアクセス集中（リトライの嵐）を防ぐ技術。
    * **レートリミット（流量制御）**: 特定のクライアントからの過剰なアクセスを遮断し、システムを保護する。

- bashのプロセス置換について
- >&2
- ⭐️memoのTODO.mdを全部やる
- Sentry
- eBPF
- 低レイヤー
- Product Readiness Check
    - https://tech.timee.co.jp/entry/2025/04/28/100000
- パイロットチーム
- クリティカルユーザージャーニー（CUJ: Critical User Journey）
- Redash
- ペネとレーションテスト
- VPC Service Controls
- Macを軽くする
    - Storage容量
- CORSの復習
- MDMの仕組み
    - Intune
- Claude Code: .claude/keybindings.json
- Lucky Thirteen攻撃
- IP/サブネット 読み方など整理。CIDR
- signoz
- .npmrc
- Production Readiness Check
- 証明書チェーン
    - ルートCA証明書
    - 中間証明書
- Go Wiki: Go Code Review Comments
    - https://go.dev/wiki/CodeReviewComments
- ルート証明書
- 証明書の検証の詳細な流れ
- Cloudflare TLS Inspection
- UNABLE_TO_GET_ISSUER_CERT_LOCALLY
- Cloudflare TLS Decryption
- Cloudflare Zero Trustの通信の流れ
- User-side certificates
    - https://developers.cloudflare.com/cloudflare-one/team-and-resources/devices/user-side-certificates/
- Install certificate manually
    - https://developers.cloudflare.com/cloudflare-one/team-and-resources/devices/user-side-certificates/manual-deployment/
- http request failed unable to get local issuer certificate
- PKCE
- githubのトークンの種類
- https://reisuta.com/strong-engineer/
- https://zenn.dev/kinniku_coder/articles/2025-08-06-engineer_routine
- https://qiita.com/akiralab/items/416bbd96122f9251fcf7
- https://adventar.org/calendars/12295

## ⭐️個人開発・セキュリティ

### 個人開発でシステムを作る場合の注意点その1（HTTP、クッキー、パスワードの管理、XSS、HttpOnle属性、API権限）

https://www.higlabo.ai/blog/higty-tech/indie-dev-security-mistakes-1

### 個人開発でシステムを作る場合の注意点その2（CSP、HSTS、CORS、CSRF）

https://www.higlabo.ai/blog/higty-tech/indie-dev-security-mistakes-2

## 30分でわかるデータ指向アプリケーションデザイン

https://speakerdeck.com/xerial/30fen-dewakarudetazhi-xiang-apurikesiyondezain-data-engineering-study-number-18

## 【書評】データ指向アプリケーションデザインを読了して見える世界

https://okb-shelf.hatenablog.com/entry/data_application_design

## eBPFとは？

https://ebpf.io/ja/what-is-ebpf/

## Microsoft Entra ID SSO

https://learn.microsoft.com/ja-jp/entra/identity/enterprise-apps/what-is-single-sign-on

シングルサインオンの仕組み

## Google Cloud における認証・認可の仕組みがこれを見ればおおよそわかる

https://zenn.dev/cloud_ace_jp/articles/f46fb1868249d7

## 知っておきたい！ 文字コードの基礎知識

https://gihyo.jp/book/pickup/2019/0006

## 新人さんに知ってほしい「文字コードのお話」

https://qiita.com/yuji38kwmt/items/b3a7820b4d3b544da4ff

## ⭐️アプリケーションのデフォルト認証情報の仕組み

https://docs.cloud.google.com/docs/authentication/application-default-credentials?hl=ja

## 【図解】【3分解説】UnicodeとUTF-8の違い！【今さら聞けない】

https://qiita.com/omiita/items/50814037af2fd8b2b21e

## unicodeとは？文字コードとは？UTF-8とは？

https://qiita.com/hiroyuki_mrp/items/f0b497394f3a5d8a8395
