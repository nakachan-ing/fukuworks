# FukuWorks

> 副業の不安を「見える化」し、安心して働けるようにサポートする副業マネジメントツール

---

## 🎯 プロジェクト概要

**FukuWorks** は、副業をしている人が抱える不安やモヤモヤを「見える化」し、より安心して活動できるように支援するツールです。

### ✨ 主な特徴
- 稼働時間の可視化（週20時間制限の自己管理）
- 案件の管理と記録（報酬・納期・進捗）
- タスクごとの自信度・不安度の記録
- 月別の収益・作業時間の集計
- 将来的には確定申告支援も視野に

---

## 🏗️ 技術構成（MVP v0.1）

| 区分 | 使用技術 |
|------|----------|
| バックエンド | Go, Gin, sqlx, SQLite |
| フロントエンド | Appsmith（ノーコードUI） |
| 認証方式 | JWT（アクセストークン） |
| 集計処理 | GoによるSQLベースの集計 |
| レポート出力 | Python（matplotlib, pandas 等） |
| インフラ | AWS Lightsail, Docker |

---

## 📦 実装済みのMVP機能（v0.1）

- [ ] ログイン機能（JWT認証）
- [ ] ダッシュボード表示（稼働時間、収益、タスク）
- [ ] 案件管理（一覧表示／新規追加／編集／削除／CSV出力）
- [ ] タスク管理（自信度・不安の記録、作業ログ追加）
- [ ] 稼働時間の可視化（20h進捗バー）

---

## 🚀 ローカル環境での起動方法

```bash
# リポジトリをクローン
git clone https://github.com/yourname/fukuworks.git
cd fukuworks

# .envファイルを用意（認証設定など）
cp .env.example .env

# サーバ起動（make または go run 等）
make run
```

---

## 📂 ディレクトリ構成（予定）

```
fukuworks/
├── cmd/              # エントリポイント
├── internal/         # ドメインロジック
│   ├── handler/      # ハンドラ層（Gin）
│   ├── repository/   # DB層（sqlx）
│   └── model/        # データ構造体
├── scripts/          # Pythonレポートスクリプト
└── frontend-appsmith/ # Appsmith設定ファイル（オプション）
```

---

## 🔒 注意事項

- このプロジェクトは**個人開発**のものであり、実データや個人情報は含まれていません。
- 認証キーや秘密情報（JWT_SECRETなど）は `.env` に記載し、 `.gitignore` 済みです。
- セキュリティを考慮したうえで公開していますが、PRやIssue歓迎です！

---

## 🛠️ 今後の開発予定

- [ ] グラフ表示（収益や稼働時間の月次推移）
- [ ] ロール認可機能（admin / user）
- [ ] 確定申告支援（必要書類のチェックリスト等）
- [ ] Appsmith から React へのUI置き換え（shadcn/ui使用予定）

---

## 👤 作者

**nakachan-ing**  
AWSインフラ・Goバックエンド・Python運用自動化を得意とする個人開発者。  
本ツールは自らの副業体験を元に開発・改善中です。

---

## 📄 ライセンス

[MIT License](https://opensource.org/licenses/mit-license.php)
