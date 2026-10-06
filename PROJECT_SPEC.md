# 求职投递管理 API 项目说明

## 1 项目目标

使用 Go 开发一个可运行的 REST API，管理求职投递记录。通过这个项目走通需求分析、数据建模、接口设计、分层实现、数据库操作、错误处理、测试和交付流程。

本项目是流程训练项目，预计 3 至 4 天完成。完成验收后停止扩展，进入 AI 文档服务主项目。

## 2 技术栈

- 开发语言：Go
- Web 框架：Gin
- 数据库：MySQL 8
- 数据访问：Go 标准库 `database/sql` 与 `go-sql-driver/mysql`
- 日志：Go 标准库 `log/slog`
- 测试：Go 标准库 `testing`、`httptest`
- 配置：环境变量
- 接口调试：Apifox、Postman 或 curl
- 版本控制：Git
- 数据库运行方式：本机 MySQL 或 Docker Compose，二选一

第一版不使用 ORM，以便真正练习 SQL、参数绑定和错误处理。

## 3 功能范围

### 3.1 必须实现

1. 健康检查。
2. 新增投递记录。
3. 根据 ID 查询投递记录。
4. 分页查询投递列表。
5. 按公司名称和投递状态筛选。
6. 修改投递记录。
7. 删除投递记录。
8. 参数校验和统一错误响应。
9. 请求日志。
10. 数据库存储，服务重启后数据仍然存在。
11. 至少包含 service 单元测试和 repository 集成测试中的一种，以及 handler 接口测试。
12. README 中写清构建、配置、启动和测试方法。

### 3.2 本项目不实现

- 登录、注册、JWT 和多用户权限。
- 前端页面。
- Redis、消息队列和微服务拆分。
- AI 功能。
- 文件上传。
- 云端部署。

这些功能不属于本项目的验收范围。

## 4 数据模型

表名：`applications`

| 字段 | 类型 | 约束 | 含义 |
|---|---|---|---|
| id | BIGINT UNSIGNED | 主键，自增 | 投递记录 ID |
| company | VARCHAR(100) | NOT NULL | 公司名称 |
| position | VARCHAR(100) | NOT NULL | 岗位名称 |
| status | VARCHAR(20) | NOT NULL | 当前状态 |
| applied_at | DATE | 可为空 | 投递日期 |
| notes | VARCHAR(1000) | NOT NULL DEFAULT '' | 备注 |
| created_at | DATETIME | NOT NULL | 创建时间 |
| updated_at | DATETIME | NOT NULL | 更新时间 |

允许的状态：

- `planned`：计划投递
- `applied`：已投递
- `written_test`：笔试
- `interview`：面试
- `offer`：已录用
- `rejected`：已拒绝
- `withdrawn`：已撤回

索引：

- 普通索引：`status`
- 普通索引：`company`
- 普通索引：`applied_at`

## 5 API 设计

统一前缀：`/api/v1`

| 方法 | 路径 | 功能 | 成功状态码 |
|---|---|---|---|
| GET | `/health` | 健康检查 | 200 |
| POST | `/api/v1/applications` | 新增记录 | 201 |
| GET | `/api/v1/applications/:id` | 查询单条记录 | 200 |
| GET | `/api/v1/applications` | 分页、筛选列表 | 200 |
| PATCH | `/api/v1/applications/:id` | 部分修改记录 | 200 |
| DELETE | `/api/v1/applications/:id` | 删除记录 | 204 |

### 5.1 新增记录

请求示例：

```json
{
  "company": "示例科技",
  "position": "后端开发实习生",
  "status": "planned",
  "applied_at": null,
  "notes": "官网日常实习"
}
```

校验规则：

- `company`：去除首尾空格后长度为 1 至 100。
- `position`：去除首尾空格后长度为 1 至 100。
- `status`：必须属于允许状态。
- `applied_at`：为空或符合 `YYYY-MM-DD`。
- `notes`：最多 1000 个字符。

### 5.2 分页查询

查询参数：

- `page`：默认 1，最小为 1。
- `page_size`：默认 20，范围 1 至 100。
- `status`：可选，必须属于允许状态。
- `company`：可选，按公司名称模糊查询。

默认排序：`updated_at DESC, id DESC`。

返回分页信息：

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 0
  }
}
```

## 6 错误响应

统一结构：

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "status is invalid"
  }
}
```

至少区分：

- 400 `INVALID_ARGUMENT`：参数格式或业务校验失败。
- 404 `NOT_FOUND`：记录不存在。
- 405 `METHOD_NOT_ALLOWED`：请求方法错误。
- 500 `INTERNAL_ERROR`：未预期的服务端错误。

对外响应不得返回 SQL、堆栈或数据库连接信息；详细错误只写入服务日志。

## 7 代码结构

```text
job-tracker-api/
├── cmd/server/main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   └── middleware/
├── migrations/001_init.sql
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── PROJECT_SPEC.md
└── README.md
```

职责划分：

- handler：解析 HTTP 请求、调用 service、生成 HTTP 响应。
- service：参数与业务规则、错误转换，不直接写 SQL。
- repository：执行 SQL，不处理 HTTP。
- model：请求、响应和数据库实体。
- middleware：请求日志、异常恢复等横切逻辑。
- config：读取并校验环境变量。

## 8 配置项

`.env.example` 至少包含：

```text
APP_PORT=8080
MYSQL_DSN=user:password@tcp(127.0.0.1:3306)/job_tracker?parseTime=true&charset=utf8mb4
```

真实密码不得提交到 Git。

## 9 工程要求

- 所有数据库操作接收 `context.Context`。
- SQL 使用参数绑定，不拼接用户输入。
- 数据库连接设置合理的最大连接数和连接生命周期。
- HTTP 服务支持优雅关闭。
- 日志至少包含请求方法、路径、状态码、耗时和错误摘要。
- handler 不直接访问数据库。
- 错误必须处理，不使用空的 `if err != nil`。
- 提交前运行 `gofmt`、`go vet ./...` 和 `go test ./...`。

## 10 测试清单

- 新增合法记录成功。
- 缺少公司或岗位时返回 400。
- 非法状态返回 400。
- 查询存在的 ID 成功。
- 查询不存在的 ID 返回 404。
- 分页默认值正确。
- 非法页码返回 400。
- 状态筛选正确。
- 修改不存在的记录返回 404。
- 删除成功后再次查询返回 404。
- 数据库错误不会向客户端泄露内部信息。

## 11 Git 提交建议

1. `chore: initialize Go project`
2. `feat: add application model and create endpoint`
3. `feat: add MySQL repository and migrations`
4. `feat: add query update and delete endpoints`
5. `test: add service and handler tests`
6. `docs: add setup and API usage`

## 12 完成标准

满足以下条件即停止开发：

1. 所有必须接口可运行。
2. MySQL 持久化正常。
3. 正常路径和主要错误路径有测试。
4. 新电脑按 README 可以完成配置和启动。
5. 能在 3 分钟内讲清请求从 handler 到数据库的调用链。
6. 能解释为什么分层、为什么使用参数绑定、索引如何选择、错误如何转换。

