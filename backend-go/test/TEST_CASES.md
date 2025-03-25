# ✅ FukuWorks API テストケース一覧

## 🧑‍💼 User エンドポイント

### ✅ 正常系
| No. | テスト名                        | 説明 |
|-----|----------------------------------|------|
| 1   | `TestPostUser_Success`          | ユーザー新規登録（成功） |
| 2   | `TestLogin_Success`             | ログイン成功（仮トークン発行） |
| 3   | `TestGetUser`                   | ユーザー情報取得 |
| 4   | `TestUpdateUser`                | ユーザー情報更新 |
| 5   | `TestSoftDeleteUser`            | ユーザー論理削除 |

### ❌ 異常系
| No. | テスト名                            | 説明 |
|-----|--------------------------------------|------|
| 1   | `TestPostUser_ValidationError`       | バリデーションエラー（name や email 未入力） |
| 2   | `TestPostUser_InvalidEmailFormat`    | email 形式不正 |
| 3   | `TestPostUser_Duplicate`             | name/email の重複エラー |
| 4   | `TestLogin_Failure`                  | ログイン失敗（ユーザー未登録） |
| 5   | `TestUpdateUser_NotFound`            | 存在しないユーザーの更新 |
| 6   | `TestSoftDeleteUser_NotFound`        | 存在しないユーザーの削除 |

---

## 📁 Project エンドポイント

### ✅ 正常系
| No. | テスト名                            | 説明 |
|-----|--------------------------------------|------|
| 1   | `TestPostProject_Success`           | プロジェクト作成 |
| 2   | `TestUpdateProject_Success`         | プロジェクト更新 |
| 3   | `TestSoftDeleteProject_Success`     | プロジェクト論理削除 |
| 4   | `TestGetProject_Success`            | プロジェクト取得 |

### ❌ 異常系
| No. | テスト名                                       | 説明 |
|-----|------------------------------------------------|------|
| 1   | `TestPostProject_ValidationError`              | 必須項目未入力 |
| 2   | `TestPostProject_InvalidDeadlineFormat`        | deadline の日付形式不正 |
| 3   | `TestGetProject_ForbiddenForOtherUser`         | 他人のプロジェクト取得（403）|
| 4   | `TestUpdateProject_ForbiddenForOtherUser`      | 他人のプロジェクト更新（403）|
| 5   | `TestSoftDeleteProject_ForbiddenForOtherUser`  | 他人のプロジェクト削除（403）|
| 6   | `TestUpdateProject_NotFound`                   | 存在しないプロジェクトの更新（404）|
| 7   | `TestSoftDeleteProject_NotFound`               | 存在しないプロジェクトの削除（404）|

---

## ✅ Task エンドポイント

### ✅ 正常系
| No. | テスト名                         | 説明 |
|-----|----------------------------------|------|
| 1   | `TestPostTask_Success`          | タスク作成成功 |
| 2   | `TestUpdateTask_Success`        | タスク更新成功 |
| 3   | `TestSoftDeleteTask_Success`    | タスク削除成功 |
| 4   | `TestGetTask_Success`           | タスク取得成功 |

### ❌ 異常系
| No. | テスト名                            | 説明 |
|-----|--------------------------------------|------|
| 1   | `TestPostTask_ValidationError`       | 必須フィールド未入力で作成失敗 |
| 2   | `TestGetTask_Forbidden`              | 他人のタスク取得（403） |
| 3   | `TestUpdateTask_Forbidden`           | 他人のタスク更新（403） |
| 4   | `TestSoftDeleteTask_Forbidden`       | 他人のタスク削除（403） |
| 5   | `TestGetTask_NotFound`               | 存在しないタスク取得（404） |
| 6   | `TestUpdateTask_NotFound`            | 存在しないタスク更新（404） |
| 7   | `TestSoftDeleteTask_NotFound`        | 存在しないタスク削除（404） |
