# ✅ User API テストケース

| No | API             | メソッド | ケース             | 入力                          | 期待レスポンス         |
|----|------------------|----------|--------------------|-------------------------------|------------------------|
| 1  | `/`              | POST     | 正常登録           | name, email, password         | 201 Created + User情報 |
| 2  | `/`              | POST     | email未入力        | name のみ                     | 400 + バリデーションエラー |
| 3  | `/:user`         | GET      | 存在ユーザー取得    | `nakachan-ing`                       | 200 + ユーザー情報     |
| 4  | `/:user`         | GET      | 存在しないユーザー  | `unknown`                     | 404 Not Found          |
| 5  | `/:user`         | PATCH    | 名前・メール更新     | 新しい name, email            | 200 OK + 更新情報      |
| 6  | `/:user`         | DELETE   | ユーザー論理削除   | `nakachan-ing`                       | 204 No Content         |
| 7  | `/login`         | POST     | ログイン成功        | name, password                | 200 OK + トークン       |
| 8  | `/login`         | POST     | ログイン失敗        | 間違った name/password        | 401 Unauthorized       |
| 9  | `/admin/users`     | GET      | 全ユーザー取得（管理）| なし                          | 200 OK + 配列          |
| 10 | `/admin/users/:id` | DELETE   | 物理削除（管理）     | userID                        | 204 No Content         |

---

# ✅ Project API テストケース

| No | API                           | メソッド | ケース           | 入力                                    | 期待レスポンス             |
|----|--------------------------------|----------|------------------|-----------------------------------------|----------------------------|
| 11 | `/:user/projects`             | POST     | 正常作成         | title, deadline など                    | 201 Created + Project情報 |
| 12 | `/:user/projects`             | POST     | 必須項目不足     | title 空など                            | 400 Bad Request           |
| 13 | `/:user/projects`             | GET      | 一覧取得         | なし                                    | 200 OK + 配列             |
| 14 | `/:user/projects/:pid`        | GET      | 詳細取得         | projectID                               | 200 OK + Project情報      |
| 15 | `/:user/projects/:pid`        | PATCH    | 更新             | 新しい値（statusなど）                 | 200 OK + 更新情報         |
| 16 | `/:user/projects/:pid`        | DELETE   | 論理削除         | projectID                               | 204 No Content            |
| 17 | `/api/projects`               | GET      | 全件取得（管理） | なし                                    | 200 OK + 配列             |
| 18 | `/api/projects/:id`           | DELETE   | 物理削除（管理） | projectID                               | 204 No Content            |
| 19 | `/:user/projects/:pid`        | GET      | 他人のプロジェクト参照 | ログインユーザー ≠ :user         | 403 Forbidden             |
| 20 | `/:user/projects`             | POST     | 他人のプロジェクト作成 | ログインユーザー ≠ :user         | 403 Forbidden             |

---

# ✅ Task API テストケース

| No | API                                           | メソッド | ケース           | 入力                                | 期待レスポンス           |
|----|-----------------------------------------------|----------|------------------|-------------------------------------|--------------------------|
| 21 | `/:user/projects/:pid/tasks`                 | POST     | 正常作成         | title, deadlineなど                | 201 Created + Task情報   |
| 22 | `/:user/projects/:pid/tasks`                 | POST     | バリデーション   | title空 など                        | 400 Bad Request          |
| 23 | `/:user/projects/:pid/tasks`                 | GET      | 一覧取得         | projectID                           | 200 OK + 配列            |
| 24 | `/:user/projects/:pid/tasks/:tid`            | GET      | 詳細取得         | taskID                              | 200 OK + Task情報        |
| 25 | `/:user/projects/:pid/tasks/:tid`            | PATCH    | 更新             | ステータス、期限など               | 200 OK + 更新情報        |
| 26 | `/:user/projects/:pid/tasks/:tid`            | DELETE   | 論理削除         | taskID                              | 204 No Content           |
| 27 | `/api/tasks`                                 | GET      | 全件取得（管理） | なし                                | 200 OK + 配列            |
| 28 | `/api/tasks/:id`                             | DELETE   | 物理削除（管理） | taskID                              | 204 No Content           |
| 29 | `/:user/projects/:pid/tasks/:tid`            | GET      | 他人のタスク参照 | ログインユーザー ≠ :user           | 403 Forbidden            |
| 30 | `/:user/projects/:pid/tasks`                 | POST     | 他人プロジェクトに追加 | ログインユーザー ≠ :user     | 403 Forbidden            |

---

✅ 今後必要に応じて：
- 認証導入後 → JWT有無、ロールによる制限ケース追加
- バリデーションの詳細（email形式、日付形式、数値の範囲など）
- 並び順、フィルター、ページネーションのテストケース
