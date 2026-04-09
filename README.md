# TIL - Today I Learned

日々の学びを記録するリポジトリ。`go run ./scripts` で `content/` を生成し、Hugo で表示する。

## 書き方

**ファイル配置**: `YYYY/MM/DD.md`（例: `2024/01/12.md`）。1ファイルに複数エントリを書く場合は `##` 見出しで区切る。

**見出し**: `## タイトル` で区切る。探すときは検索を使う。

## 日報（TIL）の内容に書くべきこと

- 普段の仕事の中で、何を感じ、どんな観点で判断したのかを自分の言葉で説明する
- //

### 画像の使い方

画像ファイルは Markdown と同じ `YYYY/MM/` ディレクトリに置く。`go run ./scripts` 実行時に `static/images/YYYY/MM/` へ自動コピーされる。

対応フォーマット: `.png` `.jpg` `.jpeg` `.gif` `.webp` `.svg`

```
2024/10/25.md        ← 記事
2024/10/image.png    ← 画像（同じ YYYY/MM/ に配置）
```

Markdown 内では以下のように参照する:

```markdown
![説明](/images/2024/10/image.png)
```

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
