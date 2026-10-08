# Job Tracker API

一个使用 Go、Gin 和 MySQL 实现的求职投递管理 REST API。

项目支持投递记录的创建、查询、分页筛选、部分更新和删除，并包含统一错误响应、请求日志、参数校验、MySQL 持久化、优雅关闭以及分层自动测试。

## 技术栈

- Go 1.27
- Gin
- MySQL 8
- `database/sql`
- `go-sql-driver/mysql`
- `testing`
- `httptest`
- Git

项目没有使用 ORM，Repository 层直接编写 SQL，以练习参数绑定、数据库错误处理和数据映射。

## 功能

- 健康检查
- 创建投递记录
- 根据 ID 查询记录
- 分页查询记录
- 按公司名称模糊筛选
- 按投递状态筛选
- 部分更新记录
- 删除记录
- 统一 JSON 错误响应
- 请求日志
- MySQL 持久化
- HTTP 服务优雅关闭

允许的投递状态：

- `planned`
- `applied`
- `written_test`
- `interview`
- `offer`
- `rejected`
- `withdrawn`

## 项目结构

```text
job-tracker-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── migrations/
│   └── 001_init.sql
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── PROJECT_SPEC.md
└── README.md
```

各层职责：

- `handler`：解析 HTTP 请求并生成 HTTP 响应。
- `service`：执行业务校验、数据标准化和错误转换。
- `repository`：访问 MySQL 或内存存储。
- `model`：定义请求、响应和实体结构。
- `middleware`：记录请求方法、路径、状态码、耗时和客户端地址。
- `config`：读取并校验数据库配置。

请求调用链：

```text
HTTP Request
    ↓
Gin Router
    ↓
Middleware
    ↓
Handler
    ↓
Service
    ↓
Repository
    ↓
MySQL
```

## 环境要求

- Go 1.27 或兼容版本
- MySQL 8
- Git

检查环境：

```powershell
go version
git --version
mysql --version
```

## 初始化数据库

在项目根目录执行：

```powershell
mysql -u root -p -e "source migrations/001_init.sql"
```

输入 MySQL 密码后，将创建：

```text
job_tracker
```

数据库以及：

```text
applications
```

数据表。

也可以在 MySQL Workbench 中打开：

```text
migrations/001_init.sql
```

并执行全部 SQL。

## 配置环境变量

项目使用两个环境变量：

```text
APP_PORT
MYSQL_DSN
```

PowerShell 示例：

```powershell
$env:APP_PORT = '8080'
$env:MYSQL_DSN = 'root:你的MySQL密码@tcp(127.0.0.1:3306)/job_tracker?parseTime=true&charset=utf8mb4&loc=Local'
```

`APP_PORT` 未设置时默认使用 `8080`。

环境变量只对当前 PowerShell 终端及其子进程有效。

## 启动服务

进入项目目录：

```powershell
cd job-tracker-api
```

启动：

```powershell
go run ./cmd/server
```

服务默认监听：

```text
http://127.0.0.1:8080
```

健康检查：

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri 'http://127.0.0.1:8080/health'
```

预期响应：

```json
{
  "status": "ok"
}
```

按下 `Ctrl+C` 时，服务器停止接收新请求，并在最多 10 秒内等待正在处理的请求完成。

## API

| 方法 | 路径 | 功能 | 成功状态码 |
|---|---|---|---:|
| GET | `/health` | 健康检查 | 200 |
| POST | `/api/v1/applications` | 创建记录 | 201 |
| GET | `/api/v1/applications/:id` | 查询单条记录 | 200 |
| GET | `/api/v1/applications` | 分页与筛选 | 200 |
| PATCH | `/api/v1/applications/:id` | 部分更新 | 200 |
| DELETE | `/api/v1/applications/:id` | 删除记录 | 204 |

### 创建记录

```powershell
$body = @{
    company    = 'ByteDance'
    position   = 'Backend Intern'
    status     = 'planned'
    applied_at = $null
    notes      = '官网日常实习'
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri 'http://127.0.0.1:8080/api/v1/applications' `
  -ContentType 'application/json' `
  -Body $body
```

### 查询单条记录

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri 'http://127.0.0.1:8080/api/v1/applications/1'
```

### 分页与筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri 'http://127.0.0.1:8080/api/v1/applications?page=1&page_size=20&status=interview&company=Byte'
```

查询参数：

- `page`：最小为 1。
- `page_size`：范围为 1 至 100。
- `status`：可选，精确匹配。
- `company`：可选，模糊匹配。

默认排序：

```text
updated_at DESC, id DESC
```

### 部分更新

只需要发送要修改的字段：

```powershell
$body = @{
    status = 'interview'
    notes  = '进入一面'
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Patch `
  -Uri 'http://127.0.0.1:8080/api/v1/applications/1' `
  -ContentType 'application/json' `
  -Body $body
```

清空投递日期：

```json
{
  "applied_at": null
}
```

未出现在 PATCH 请求中的字段保持不变。

### 删除记录

```powershell
Invoke-RestMethod `
  -Method Delete `
  -Uri 'http://127.0.0.1:8080/api/v1/applications/1'
```

删除成功返回：

```text
204 No Content
```

## 错误响应

所有业务错误使用统一 JSON 结构：

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "invalid status"
  }
}
```

主要错误：

| HTTP 状态码 | 错误码 | 含义 |
|---:|---|---|
| 400 | `INVALID_ARGUMENT` | 请求格式或参数校验失败 |
| 404 | `NOT_FOUND` | 路由或记录不存在 |
| 405 | `METHOD_NOT_ALLOWED` | 路径存在，但 HTTP 方法不支持 |
| 500 | `INTERNAL_ERROR` | 未预期的服务器错误 |

数据库错误等内部信息不会直接返回给客户端。

## 测试

运行全部测试：

```powershell
go test ./... -count=1
```

运行静态检查：

```powershell
go vet ./...
```

MySQL 集成测试依赖 `MYSQL_DSN`。如果没有设置，该测试会跳过。

运行 MySQL Repository 测试：

```powershell
go test ./internal/repository -v -count=1
```

运行 PATCH 集成测试：

```powershell
go test ./internal/repository `
  -run '^TestMySQLApplicationRepositoryUpdate$' `
  -v `
  -count=1
```

运行 Handler 测试：

```powershell
go test ./internal/handler -v -count=1
```

运行 Service 测试：

```powershell
go test ./internal/service -v -count=1
```

## 构建

```powershell
go build -o job-tracker-api.exe ./cmd/server
```

启动编译后的程序：

```powershell
.\job-tracker-api.exe
```

## 设计说明

Service 层通过接口依赖 Repository，因此可以在不修改业务逻辑的情况下切换内存实现和 MySQL 实现。

所有 SQL 值均通过占位符传入：

```sql
WHERE id = ?
```

用户输入不会直接拼接到 SQL 中。

PATCH 使用自定义可选字段类型区分三种状态：

```text
字段未出现       → 不修改
字段值为 null    → 清空可为空字段
字段值为字符串   → 校验并更新
```

数据库为 `status`、`company` 和 `applied_at` 建立了索引。当前列表接口支持按 `status` 精确筛选和按 `company` 包含匹配。由于公司筛选使用 `LIKE '%keyword%'`，普通 B-Tree 索引通常无法优化该条件；数据规模扩大后可以改用前缀匹配或全文索引。

HTTP 服务支持请求日志、panic 恢复和优雅关闭。