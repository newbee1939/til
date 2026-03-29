# TIL - Today I Learned

日々の学びを記録するリポジトリ。`go run ./scripts` で `content/` を生成し、Hugo で表示する。

## 書き方

**ファイル配置**: `YYYY/MM/DD.md`（例: `2024/01/12.md`）。1ファイルに複数エントリを書く場合は `##` 見出しで区切る。

**見出し**: `## タイトル` で区切る。画像は同じ `YYYY/MM/` に置くと `static/images/YYYY/MM/` へコピーされる。探すときは検索を使う。

## セットアップ

- **Go**: リポジトリに `.tool-versions` あり。`mise install` で入る（未導入時は [mise.run](https://mise.run/) 参照）。

## コマンド

```sh
# content/ を再生成して Hugo でプレビュー
go run ./scripts && hugo server
# → http://localhost:1313/til/
```

- 生成のみ: `go run ./scripts`
- テスト: `go test ./scripts`

## デプロイ

`main` に push すると GitHub Actions で GitHub Pages に自動デプロイされる。
