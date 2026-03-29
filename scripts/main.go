package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	// spaceRe: 連続する空白文字を1つにまとめるために使う。
	// タブ(\t)、改行(\n)、改ページ(\f)、復帰(\r)、半角スペース、全角スペース(\u3000) の
	// 1文字以上にマッチ。Slug() でタイトル内の空白をハイフンに置換する際に使用。
	spaceRe = regexp.MustCompile(`[` + "\t\n\f\r \u3000" + `]+`)

	// dashRe: 連続するハイフン(-)にマッチ。
	// Slug() で「---」などを「-」1つに正規化し、URLスラッグをきれいにするために使用。
	dashRe = regexp.MustCompile(`-+`)

	// h2Re: Markdown の見出し2（## タイトル）の行にマッチ。
	// ^##\s+ で行頭の「##」とその後の空白、(.+)$ で見出しテキストをキャプチャ。
	// Split() で本文を ## 見出しごとに分割する際に使用。
	h2Re = regexp.MustCompile(`^##\s+(.+)$`)

	// yearRe: 4桁の数字のみ（西暦年）にマッチ。例: 2024, 2025。
	// findFiles() で YYYY/MM/DD ディレクトリ構造の「年」フォルダ名の検証に使用。
	yearRe = regexp.MustCompile(`^\d{4}$`)

	// monthRe: 2桁の数字のみ（月）にマッチ。例: 01, 12。
	// findFiles() で「月」フォルダ名の検証に使用。
	monthRe = regexp.MustCompile(`^\d{2}$`)

	// dayRe: 1桁または2桁の数字（日）にマッチ。例: 1, 09, 31。
	// findFiles() で日付ディレクトリ内のファイル名（日付部分）の検証に使用。
	dayRe = regexp.MustCompile(`^\d{1,2}$`)

	// nonDigit: 数字以外の文字（0-9 以外）の連続にマッチ。
	// 文字列から数字だけを取り出したいときに、数字以外を除去するために使用。
	nonDigit = regexp.MustCompile(`[^0-9]+`)
)

// imgExt は画像ファイルの拡張子を判定するためのマップです。
// Key に画像拡張子（例: ".png", ".jpg" など）を持ち、対応する Value を true にしています。
// これにより、あるファイル名やパスが画像かどうかを簡単に判定できるようになっています。
// 例: if imgExt[filepath.Ext(filename)] { ... } で画像ファイルかを判別できる。
var imgExt = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".webp": true,
	".svg":  true,
}

type entry struct{ Title, Body string }

type srcFile struct{ Path, Year, Month, Day string }

func Slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' ||
			unicode.IsSpace(r) || (r >= 0x3000 && r <= 0x9fff) || (r >= 0xff00 && r <= 0xffef) {
			b.WriteRune(r)
		}
	}
	out := strings.Trim(dashRe.ReplaceAllString(spaceRe.ReplaceAllString(b.String(), "-"), "-"), "-")
	if len(out) > 80 {
		out = out[:80]
	}
	return out
}

func Split(body string) []entry {
	var out []entry
	var cur *entry
	flush := func() {
		if cur != nil {
			cur.Body = strings.TrimSpace(cur.Body)
			out = append(out, *cur)
		}
	}
	for _, line := range strings.Split(body, "\n") {
		m := h2Re.FindStringSubmatch(line)
		if len(m) == 2 && !strings.HasPrefix(line, "###") {
			flush()
			cur = &entry{Title: strings.TrimSpace(m[1])}
			continue
		}
		if cur != nil {
			cur.Body += line + "\n"
		}
	}
	flush()
	return out
}

// findFiles は、root 以下の「年/月」フォルダ（例: 2024/01）をたどり、
// その中の .md ファイルを集めて srcFile のリストで返します。
// あわせて画像ファイルを static/images/年/月/ にコピーします。
func findFiles(root string) ([]srcFile, error) {
	// 見つかった .md ファイルの情報を入れるスライス（可変長のリスト）
	var list []srcFile

	// --- ステップ1: 「年」フォルダを取得 ---
	// subdirs(root, yearRe) で、root 直下で「4桁の数字」の名前のフォルダだけを取得。
	// 例: 2024, 2025 など。years には ["2024", "2025"] のような文字列のスライスが入る。
	years, err := subdirs(root, yearRe)
	if err != nil {
		return nil, fmt.Errorf("failed to find year directories in %s: %w", root, err)
	}

	// 各「年」フォルダの中を見ていく
	for _, y := range years {
		// --- ステップ2: 「月」フォルダを取得 ---
		// filepath.Join(root, y) で「root/2024」のようなパスを作る。
		// subdirs(..., monthRe) で「2桁の数字」のフォルダだけを取得。例: 01, 12。
		monthDir := filepath.Join(root, y)
		months, err := subdirs(monthDir, monthRe)
		if err != nil {
			return nil, fmt.Errorf("failed to find month directories in %s: %w", monthDir, err)
		}

		for _, m := range months {
			// --- ステップ3: 年/月 の中身（ファイル一覧）を読む ---
			// path は「root/2024/01」のような「日付フォルダ」のパス。
			path := filepath.Join(root, y, m)
			dir, err := os.ReadDir(path)
			if err != nil {
				// フォルダが読めない（存在しない・権限なし等）場合はエラーで終了
				return nil, err
			}

			// このフォルダ内の「各要素」（ファイル or サブフォルダ）を1つずつ処理
			for _, e := range dir {
				// フォルダの場合はスキップ（.md や画像はファイルだけを対象にする）
				if e.IsDir() {
					continue
				}

				name := e.Name() // 例: "15.md", "image.png"
				// 拡張子を小文字で取得。".md", ".MD" どちらも ".md" で扱うため
				ext := strings.ToLower(filepath.Ext(name))

				// --- Markdown ファイルの場合 ---
				if ext == ".md" {
					// ファイル名から「日」の数字を取り出す。
					// "15.md" → "15", "9.md" → "9"
					// TrimSuffix(name, ext) で "15.md" → "15", nonDigit で数字以外を除去。
					day := nonDigit.ReplaceAllString(strings.TrimSuffix(name, ext), "")
					// 1〜2桁の数字でないもの（例: "abc.md"）は無視
					if !dayRe.MatchString(day) {
						continue
					}
					// "9" を "09" のように2桁にそろえる（日付の表記を統一するため）
					if len(day) == 1 {
						day = "0" + day
					}
					// リストに追加: フルパス、年、月、日 を srcFile として保存
					list = append(list, srcFile{filepath.Join(path, name), y, m, day})
					continue
				}

				// --- 画像ファイルの場合 ---
				// imgExt は .png, .jpg などのマップ。この拡張子なら画像として扱う。
			if imgExt[ext] {
				// コピー先: static/images/年/月/元のファイル名
				dest := filepath.Join(root, "static", "images", y, m, name)
				// コピー先のフォルダが無ければ作成（0755 は一般的なディレクトリの権限）
				if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
					return nil, fmt.Errorf("failed to create directory for image %s: %w", dest, err)
				}
				data, err := os.ReadFile(filepath.Join(path, name))
				if err != nil {
					return nil, fmt.Errorf("failed to read image %s: %w", name, err)
				}
				if err := os.WriteFile(dest, data, 0644); err != nil {
					return nil, fmt.Errorf("failed to write image %s: %w", dest, err)
				}
			}
			}
		}
	}

	// --- ステップ4: 日付順（年→月→日）でソート ---
	// 文字列を連結して比較すると "20240109" < "20240115" のように日付順になる。
	sort.Slice(list, func(i, j int) bool {
		return list[i].Year+list[i].Month+list[i].Day < list[j].Year+list[j].Month+list[j].Day
	})

	return list, nil
}

func subdirs(dir string, re *regexp.Regexp) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && re.MatchString(e.Name()) {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func Run(root, contentDir string) (count int, fileCount int, err error) {
	if err := os.RemoveAll(contentDir); err != nil {
		return 0, 0, fmt.Errorf("failed to remove content directory %s: %w", contentDir, err)
	}
	postDir := filepath.Join(contentDir, "posts")
	if err := os.MkdirAll(postDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create posts directory: %w", err)
	}
	files, err := findFiles(root)
	if err != nil {
		return 0, 0, err
	}
	for _, f := range files {
		body, err := os.ReadFile(f.Path)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read file %s: %w", f.Path, err)
		}
		date := f.Year + "-" + f.Month + "-" + f.Day
		daily := f.Year + "/" + f.Month + "/" + f.Day
		for i, e := range Split(string(body)) {
			s := Slug(e.Title)
			if s == "" {
				s = "entry-" + strconv.Itoa(i+1)
			}
			fm := "---\ntitle: \"" + strings.ReplaceAll(e.Title, `"`, `\"`) + "\"\ndate: " + date + "\ndaily: \"" + daily + "\"\n---"
			outPath := filepath.Join(postDir, date+"-"+s+".md")
			if err := os.WriteFile(outPath, []byte(fm+"\n\n"+e.Body+"\n"), 0644); err != nil {
				return 0, 0, fmt.Errorf("failed to write post %s: %w", outPath, err)
			}
			count++
		}
	}
	if err := writeIndex(postDir, "すべてのエントリ"); err != nil {
		return 0, 0, err
	}
	searchDir := filepath.Join(contentDir, "search")
	if err := os.MkdirAll(searchDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create search directory: %w", err)
	}
	if err := writeIndex(searchDir, "検索"); err != nil {
		return 0, 0, err
	}
	return count, len(files), nil
}

func writeIndex(dir, title string) error {
	if err := os.WriteFile(filepath.Join(dir, "_index.md"), []byte("---\ntitle: "+title+"\n---\n"), 0644); err != nil {
		return fmt.Errorf("failed to write index file in %s: %w", dir, err)
	}
	return nil
}

func main() {
	root, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("failed to get absolute path: %v", err)
	}
	n, files, err := Run(root, filepath.Join(root, "content"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d entries from %d files\n", n, files)
}
