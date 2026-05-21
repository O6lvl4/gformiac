# gformiac

[English](README.md)

Google Forms Infrastructure as Code — YAML定義からGoogle Formsを宣言的に管理するCLIツール。

## Install

```bash
go install github.com/O6lvl4/gformiac@latest
```

## Quick Start

```yaml
# form.yaml
title: "顧客満足度アンケート"
description: "サービス改善のため、ご意見をお聞かせください。"

items:
  - title: "お名前"
    type: short_answer
    required: true

  - title: "満足度を教えてください"
    type: choice
    required: true
    choice:
      type: radio
      options:
        - "とても満足"
        - "満足"
        - "普通"
        - "不満"
        - "とても不満"

  - title: "自由記述"
    type: paragraph
```

```bash
gformiac plan                  # 差分プレビュー
gformiac apply                 # フォーム作成/更新
gformiac import <formID>       # 既存フォームをYAMLに変換
```

## Setup

1. [Google Cloud Console](https://console.cloud.google.com/) でプロジェクトを作成
2. Google Forms API を有効化
3. OAuth 2.0 クライアントIDを作成し `credentials.json` として保存
4. 初回 `apply` 時にブラウザ認証 → `token.json` が自動生成される

サービスアカウントの場合は `GOOGLE_APPLICATION_CREDENTIALS` 環境変数を設定するだけでOK。

## Commands

| コマンド | 説明 |
|---|---|
| `plan` | ローカル定義とリモートの差分をプレビュー |
| `apply` | フォームを作成または更新 |
| `import <formID>` | 既存フォームをYAML定義にインポート |

### Flags

```
-f, --file string          フォーム定義ファイル (default "form.yaml")
    --credentials string   OAuth2認証情報ファイル (default "credentials.json")
    --token string         OAuthトークンファイル (default "token.json")
    --state string         状態ファイル (default "gformiac.state.json")
    --auto-approve         確認プロンプトをスキップ (apply時)
```

## Item Types

| type | 説明 | 追加フィールド |
|---|---|---|
| `short_answer` | 短文回答 | — |
| `paragraph` | 長文回答 | — |
| `choice` | 選択式 | `choice.type`: `radio` / `checkbox` / `dropdown` |
| `scale` | スケール | `scale.low`, `scale.high`, `scale.low_label`, `scale.high_label` |
| `date` | 日付 | — |
| `time` | 時刻 | — |
| `page_break` | セクション区切り | — |

## YAML Schema

```yaml
title: "フォームタイトル"           # 必須
description: "フォームの説明"       # 任意

items:
  - title: "質問文"                # 必須
    type: short_answer             # 必須
    description: "補足説明"         # 任意
    required: true                 # 任意 (default: false)

  - title: "選択質問"
    type: choice
    choice:
      type: radio                  # radio / checkbox / dropdown
      options:
        - "選択肢1"
        - "選択肢2"

  - title: "評価"
    type: scale
    scale:
      low: 1
      high: 10
      low_label: "低い"
      high_label: "高い"
```

## State Management

`apply` 実行時に `gformiac.state.json` が生成され、フォームIDとアイテムIDのマッピングを保持する。このファイルにより次回以降の `plan` / `apply` でリモートとの差分検出が可能になる。

```
credentials.json        ← .gitignoreに追加推奨
token.json              ← .gitignoreに追加推奨
gformiac.state.json     ← .gitignoreに追加推奨
```

## License

MIT
