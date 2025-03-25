# APIテストケース一覧

## ✅ USER API

### 🔹 正常系

| No  | テスト内容                      | Endpoint        | 説明                     |
|-----|--------------------------------|-----------------|--------------------------|
| U01 | ユーザー作成（正常）           | POST /signup    | 必須項目を正しく指定    |
| U02 | ログイン（成功）               | POST /login     | 正しい認証情報を指定    |
| U03 | ユーザー取得                   | GET /:user      | ログインユーザー本人     |
| U04 | ユーザー更新                   | PATCH /:user    | 名前・メール変更         |
| U05 | ユーザー削除（ソフト）         | DELETE /:user   | 対象ユーザーを論理削除   |

### 🔹 異常系

| No  | テスト内容                        | Endpoint        | 説明                             |
|-----|----------------------------------|-----------------|----------------------------------|
| U06 | 作成時バリデーションエラー       | POST /signup    | name, email, password が空       |
| U07 | 無効なemail形式                  | POST /signup    | emailが不正                      |
| U08 | 重複name/email登録               | POST /signup    | 同一ユーザー名/メールの再登録   |
| U09 | 存在しないユーザー更新/削除     | PATCH /:user, DELETE /:user | ghostuser など存在しない        |

### 🔹 バリデーション

| フィールド     | 制約内容                      |
|----------------|-------------------------------|
| name           | 必須、最大30文字              |
| email          | 必須、正しい形式、最大255文字 |
| password       | 必須、8〜64文字               |

---

## ✅ PROJECT API

### 🔹 正常系

| No  | テスト内容                      | Endpoint                    |
|-----|-------------------------------|-----------------------------|
| P01 | プロジェクト作成              | POST /:user/projects        |
| P02 | プロジェクト取得              | GET /:user/projects/:pid    |
| P03 | プロジェクト更新              | PATCH /:user/projects/:pid  |
| P04 | プロジェクト削除（ソフト）    | DELETE /:user/projects/:pid |

### 🔹 異常系

| No  | テスト内容                        | Endpoint                    | 説明                             |
|-----|----------------------------------|-----------------------------|----------------------------------|
| P05 | 他人のプロジェクト参照禁止       | GET /otheruser/projects/:pid|
| P06 | 他人のプロジェクト更新/削除禁止 | PATCH/DELETE /otheruser/... |
| P07 | 不正な日付形式                  | POST /:user/projects        | deadline = "31-12-2025"         |

### 🔹 バリデーション

| フィールド      | 制約内容                               |
|-----------------|----------------------------------------|
| title           | 必須、最大100文字                      |
| platform        | 任意、最大100文字                      |
| client          | 任意、最大100文字                      |
| estimated_fee   | 任意、数値                             |
| status          | 必須、"Open", "In progress", "Completed", "Canceled" のみ |
| deadline        | 任意、YYYY-MM-DD 形式                  |

---

## ✅ TASK API

### 🔹 正常系

| No  | テスト内容                      | Endpoint                                |
|-----|-------------------------------|-----------------------------------------|
| T01 | タスク作成                     | POST /:user/projects/:pid/tasks         |
| T02 | タスク取得                     | GET /:user/projects/:pid/tasks/:tid     |
| T03 | タスク更新                     | PATCH /:user/projects/:pid/tasks/:tid   |
| T04 | タスク削除（ソフト）           | DELETE /:user/projects/:pid/tasks/:tid  |

### 🔹 異常系

| No  | テスト内容                        | Endpoint                              | 説明                           |
|-----|----------------------------------|---------------------------------------|--------------------------------|
| T05 | 他人のタスク参照禁止             | GET /otheruser/projects/1/tasks/1     |
| T06 | 他人のタスク更新/削除禁止       | PATCH/DELETE /otheruser/...           |
| T07 | 存在しないタスクの更新/削除     | PATCH/DELETE /.../tasks/999           |
| T08 | 不正なstatus/priorityの指定      | POST /:user/projects/:pid/tasks       |

### 🔹 バリデーション

| フィールド     | 制約内容                               |
|----------------|----------------------------------------|
| title          | 必須、最大100文字                      |
| description    | 任意、最大500文字                      |
| status         | 必須、"Open", "In progress", "Completed", "Canceled" のみ |
| priority       | 必須、"Low", "Medium", "High" のみ     |
| due_date       | 任意、YYYY-MM-DD 形式                  |

---

## その他

- 認可ミドルウェアも全ルートに適用済み。
- ログインユーザーとパスの`/:user`が一致しない場合は403 Forbidden。
- `/login`, `/signup`, `/admin/...` は `ReservedPathGuard` により認可対象外。
