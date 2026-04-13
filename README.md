# jikkyo

アニメ放送の「実際の開始時刻」、「A」「B」「C」パートの開始時刻を自動検出する CLI ツール

https://github.com/user-attachments/assets/9571eba9-dfa7-4dc3-841e-009037bd7593


## 概要

jikkyo は、アニメの放送タイトルとエピソード番号を指定すると、以下の情報を自動取得します：

- **ｷﾀ** - 実際の放送開始時刻（放送コメントから検出）
- **A** - Aパート開始時刻
- **B** - Bパート開始時刻
- **C** - Cパート開始時刻


### ビルド

```bash
# リポジトリをクローン

# 依存関係をインストール
go mod download

# ビルド
go build -o jikkyo.exe .
```

## 使い方

### 基本的な使用方法

```bash
./jikkyo.exe jikkyo -t "タイトル" -e エピソード番号
```

### フラグ一覧

| フラグ | 短縮形 | 説明 | 必須 | 例 |
|--------|--------|------|------|-----|
| `--title` | `-t` | アニメタイトル | ✅ | `-t "プリパラ"` |
| `--episode` | `-e` | エピソード番号 | ✅ | `-e 1` |
| `--verbose` | `-v` | 詳細ログ出力 | - | `-v` |
| `--log-file` | `-l` | ログ保存先 | - | `-l output.log` |

### 実行例

#### 1. シンプル実行

```bash
./jikkyo.exe jikkyo -t "プリパラ" -e 1
```

**結果:**
```
(title) プリパラ
(episode) 1
(subtitle) アイドル始めちゃいました！
(start) 2014-07-05 10:00:00
(real_start_time) 2014-07-05 10:00:10
(A) 2014-07-05 10:04:11
(B) 2014-07-05 10:19:40
```


### 複数マッチ時の選択

検索結果が複数ある場合、CLI で選択を促されます：

```
Multiple titles found:
1. プリパラ (TID: 3434)
2. アイドルタイムプリパラ (TID: 4541)

Select (1-2): 1   ← 番号を入力して Enter キーを押す
```

## 出力フォーマット

```
(title)          タイトル名
(episode)        エピソード番号
(subtitle)       サブタイトル
(start)          放送開始時刻（基準時刻）
(real_start_time) 実際の放送開始時刻（ｷﾀコメントから検出）
(A)              Aパート開始時刻（見つかった場合のみ表示）
(B)              Bパート開始時刻（見つかった場合のみ表示）
(C)              Cパート開始時刻（見つかった場合のみ表示）
```


### 複数エピソード処理例

```bash
# ループで複数エピソードを処理
for ep in {1..13}; do
  echo "Episode $ep..."
  echo "1" | ./jikkyo.exe jikkyo -t "プリパラ" -e $ep -l results.log
  echo "---"
done
```


## 外部 API クレジット

- [Syoboi Calendar](http://cal.syoboi.jp/) - テレビ放送予定データベース
- [Jikkyo API](https://jikkyo.tsukumijima.net/) - 放送コメント API


### その他

しょぼいカレンダーでタイトルIDを取得する際、地デジではないチャンネルしか取得できない場合があり、(例：AT-X, Youtube)正しく取得できない場合がございます。
