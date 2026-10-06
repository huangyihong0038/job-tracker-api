# Go 后端学习清单

适用目标：使用 Go、Gin 和 MySQL 完成求职投递管理 API，并为后端开发实习面试建立基础。

使用方法：按顺序学习。每个知识点至少做到“能解释、能独立写出最小示例、能在项目中使用”。不要等全部学完才开始项目。

## 0 开发工具与项目命令

- [ ] 理解源码、编译、运行和可执行文件的区别。
- [ ] 会使用 `go version`、`go env`。
- [ ] 会使用 `go run`、`go build`、`go install` 的基本形式。
- [ ] 会使用 `go fmt ./...`、`go vet ./...`、`go test ./...`。
- [ ] 会使用 VS Code Go 扩展、`gopls`、断点调试。
- [ ] 能读懂编译错误的文件名、行号和错误信息。

搜索关键词：

```text
Go run 和 go build 区别
Go fmt vet test 命令
VS Code Go gopls 调试
```

## 1 程序结构

- [ ] `package main` 与 `func main()`。
- [ ] `import` 单个包和分组导入。
- [ ] 注释、标识符、Go格式化规则。
- [ ] 未使用变量和未使用import为什么会编译失败。
- [ ] 名称首字母大写表示导出，小写表示包内可见。
- [ ] 同一目录中的 Go 文件通常属于同一个 package。

搜索关键词：

```text
Go package main func main
Go exported names 大写 小写
Go package 目录关系
```

掌握标准：能创建、编译并运行一个包含多个函数的命令行程序。

## 2 变量 常量与基本类型

- [ ] `var` 声明、初始化和零值。
- [ ] 短变量声明 `:=` 的作用域限制。
- [ ] `const` 与 `iota`。
- [ ] `bool`、`string`。
- [ ] `int`、`int64`、`uint64` 等整数类型。
- [ ] `float64`。
- [ ] `byte` 与 `rune`。
- [ ] 显式类型转换。
- [ ] `fmt.Print`、`Println`、`Printf`、`Sprintf`。

搜索关键词：

```text
Go 变量 零值 :=
Go const iota
Go byte rune string UTF-8
Go fmt Printf 占位符
```

掌握标准：能解释为什么数据库ID常用`int64`，能够正确格式化字符串和数字。

## 3 控制流程

- [ ] `if`及其初始化语句。
- [ ] `for`的三种常用写法。
- [ ] `range`。
- [ ] `switch`与不自动贯穿。
- [ ] `break`、`continue`。
- [ ] 理解Go没有`while`关键字。

搜索关键词：

```text
Go if for range switch
```

掌握标准：能遍历投递记录并按状态筛选或统计。

## 4 数组 切片与map

### 数组

- [ ] 数组长度属于类型的一部分。
- [ ] 数组赋值会复制整个数组。

### 切片

- [ ] `[]T`、切片字面量、`make`。
- [ ] `len`、`cap`。
- [ ] `append`、`copy`。
- [ ] 切片表达式和底层数组共享。
- [ ] `nil`切片与空切片。
- [ ] append扩容后底层数组可能改变。
- [ ] 删除元素的常用写法。

### map

- [ ] 创建、写入、读取和删除。
- [ ] `value, ok := m[key]`。
- [ ] map的零值不可直接写入。
- [ ] map不是并发安全的。
- [ ] 遍历顺序不固定。

搜索关键词：

```text
Go slice len cap append 底层数组
Go nil slice empty slice
Go map comma ok
Go map concurrent safe
```

掌握标准：能使用`map[int64]Application`实现内存版增删改查，并使用slice返回列表。

## 5 函数

- [ ] 参数与返回值。
- [ ] 多返回值。
- [ ] 命名返回值只需了解，不滥用。
- [ ] 可变参数。
- [ ] 函数作为值和回调。
- [ ] 匿名函数和闭包的基本概念。
- [ ] Go不支持函数重载和默认参数。

搜索关键词：

```text
Go function multiple return values
Go variadic function
Go closure
```

掌握标准：能够编写`validateStatus(status string) error`和分页参数校验函数。

## 6 指针 值语义与生命周期

- [ ] `&`取地址、`*`解引用。
- [ ] Go没有指针运算。
- [ ] 参数默认按值传递。
- [ ] 结构体值传递与指针传递的差异。
- [ ] slice、map本身也是值，但内部引用共享数据。
- [ ] `new`与`make`的区别。
- [ ] 垃圾回收的基本概念。
- [ ] 逃逸分析只需了解，不在第一阶段深入。

搜索关键词：

```text
Go pointer value semantics
Go new make 区别
Go slice map 值传递
```

掌握标准：能解释为什么修改`Application`的方法通常使用指针接收者。

## 7 结构体 方法与组合

- [ ] 定义和初始化结构体。
- [ ] 匿名结构体的简单使用。
- [ ] 结构体字段标签。
- [ ] 值接收者与指针接收者。
- [ ] 方法集的基本概念。
- [ ] 嵌入和组合。
- [ ] Go没有传统类继承。

搜索关键词：

```text
Go struct initialization
Go value receiver pointer receiver
Go struct embedding composition
Go struct tag
```

掌握标准：能够定义`Application`、创建请求和响应结构体，并实现状态修改方法。

## 8 接口

- [ ] 接口定义的是方法集合。
- [ ] 隐式实现，不写`implements`。
- [ ] 接口值包含动态类型和动态值。
- [ ] `nil`接口与“内部值为nil的接口”差异。
- [ ] 类型断言的基本语法。
- [ ] 空接口`any`只在确有必要时使用。
- [ ] 接口应由使用方按需要定义。
- [ ] 使用小接口，不提前抽象所有代码。

搜索关键词：

```text
Go interface implicit implementation
Go nil interface trap
Go type assertion
Go small interfaces consumer side
```

掌握标准：能定义`ApplicationRepository`，并分别用内存存储和MySQL存储实现它。

## 9 错误处理

- [ ] `error`接口。
- [ ] `errors.New`。
- [ ] `fmt.Errorf`与`%w`。
- [ ] `errors.Is`和`errors.As`。
- [ ] 哨兵错误，如`ErrNotFound`。
- [ ] 何时记录错误，避免每层重复打印。
- [ ] 对外错误与内部错误分离。
- [ ] `panic`与`recover`的用途和边界。

搜索关键词：

```text
Go error wrapping %w errors.Is
Go sentinel errors
Go panic recover best practices
```

掌握标准：repository返回`ErrNotFound`，service保留业务语义，handler转换为HTTP 404。

## 10 defer与资源释放

- [ ] defer在函数返回前执行。
- [ ] 多个defer按后进先出执行。
- [ ] defer参数在注册时求值。
- [ ] 使用defer关闭文件、响应体和数据库行集。
- [ ] 在长循环中滥用defer可能延迟释放资源。

搜索关键词：

```text
Go defer execution order arguments
Go defer resource cleanup
```

掌握标准：数据库查询后正确`defer rows.Close()`并检查`rows.Err()`。

## 11 Go Module与包管理

- [ ] Module、Package和Go文件的关系。
- [ ] `go mod init`。
- [ ] `go.mod`与`go.sum`。
- [ ] Module Path与内部import路径。
- [ ] `go get`、`go mod tidy`、`go mod download`。
- [ ] 直接依赖和间接依赖。
- [ ] `internal`目录限制。
- [ ] 不提交模块缓存和编译产物。

搜索关键词：

```text
Go modules go.mod go.sum
Go module path package import
Go mod tidy
Go internal directory
```

掌握标准：能够初始化本项目、添加Gin和MySQL驱动，并解释`go.mod`发生了什么变化。

## 12 常用标准库

- [ ] `fmt`：格式化。
- [ ] `strings`：字符串处理。
- [ ] `strconv`：字符串与数字转换。
- [ ] `time`：时间解析和格式化。
- [ ] `encoding/json`：JSON。
- [ ] `errors`：错误判断。
- [ ] `context`：取消、超时和值传递。
- [ ] `log/slog`：结构化日志。
- [ ] `os`：环境变量和进程信号。
- [ ] `net/http`：HTTP服务。
- [ ] `database/sql`：数据库访问。

搜索关键词：

```text
Go strings strconv time examples
Go encoding json struct tags
Go context timeout cancellation
Go slog structured logging
```

## 13 JSON

- [ ] `json.Marshal`与`json.Unmarshal`。
- [ ] 导出字段才能被编码。
- [ ] `json`结构体标签。
- [ ] `omitempty`的行为。
- [ ] 缺失字段、零值和`null`的区别。
- [ ] 使用指针字段表达“PATCH请求中未提供”和“提供零值”的差异。
- [ ] 不直接把内部错误编码给客户端。

搜索关键词：

```text
Go json Marshal Unmarshal struct tags
Go JSON omitempty pointer PATCH
```

掌握标准：能解析新增请求，并为PATCH请求正确区分未提交字段。

## 14 HTTP基础

- [ ] 请求方法、URL、请求头、请求体。
- [ ] 路径参数与查询参数。
- [ ] HTTP状态码。
- [ ] `Content-Type: application/json`。
- [ ] 幂等性的基本概念。
- [ ] REST API的资源建模。
- [ ] `net/http`中的Handler、Request、ResponseWriter。
- [ ] 请求体大小、超时等基本安全意识。

搜索关键词：

```text
HTTP request response methods status codes
REST API idempotency
Go net http Handler Request ResponseWriter
```

掌握标准：不使用Gin也能写出`GET /health`并返回JSON。

## 15 Gin

- [ ] 创建engine与启动服务。
- [ ] 注册路由和路由组。
- [ ] handler与`gin.Context`。
- [ ] `ShouldBindJSON`。
- [ ] 路径参数`Param`与查询参数`Query`。
- [ ] 返回JSON和状态码。
- [ ] 中间件执行顺序。
- [ ] Recovery和请求日志。
- [ ] 参数校验错误的统一转换。
- [ ] handler测试。

搜索关键词：

```text
Gin getting started routes groups
Gin ShouldBindJSON validation
Gin middleware execution order
Gin httptest handler test
```

掌握标准：完成健康检查、新增、单条查询接口，handler中不直接写SQL。

## 16 MySQL与database/sql

- [ ] 导入MySQL驱动。
- [ ] DSN组成与`parseTime=true`。
- [ ] `sql.Open`不会立即验证连接。
- [ ] `PingContext`验证连接。
- [ ] 连接池不是单个连接。
- [ ] `SetMaxOpenConns`、`SetMaxIdleConns`、`SetConnMaxLifetime`。
- [ ] `ExecContext`。
- [ ] `QueryRowContext`。
- [ ] `QueryContext`。
- [ ] `Scan`。
- [ ] `sql.ErrNoRows`。
- [ ] `RowsAffected`与`LastInsertId`。
- [ ] `rows.Close()`与`rows.Err()`。
- [ ] 参数绑定，禁止拼接用户输入。
- [ ] `NULL`字段与`sql.Null*`或指针。
- [ ] 事务`BeginTx`、`Commit`、`Rollback`。

搜索关键词：

```text
Go database sql MySQL tutorial
Go sql.Open PingContext connection pool
Go QueryRowContext QueryContext Scan
Go sql.ErrNoRows
Go database transaction BeginTx rollback defer
```

掌握标准：独立实现applications表的CRUD、分页和筛选。

## 17 SQL与数据库基础

- [ ] `CREATE TABLE`。
- [ ] `INSERT`、`SELECT`、`UPDATE`、`DELETE`。
- [ ] `WHERE`、`ORDER BY`、`LIMIT`、`OFFSET`。
- [ ] `COUNT(*)`。
- [ ] 主键、唯一索引、普通索引。
- [ ] 联合索引和最左前缀只需初步理解。
- [ ] 事务ACID。
- [ ] 隔离级别和常见并发异常。
- [ ] SQL注入。
- [ ] 使用`EXPLAIN`查看查询计划的基础。

搜索关键词：

```text
MySQL CRUD limit offset count
MySQL index B+ tree
MySQL transaction isolation levels
MySQL SQL injection prepared statement
MySQL explain basic
```

## 18 分层与依赖方向

- [ ] handler、service、repository、model的职责。
- [ ] 请求DTO、响应DTO与数据库实体的区别。
- [ ] 依赖从外层指向业务接口。
- [ ] 构造函数注入依赖。
- [ ] 避免全局数据库变量。
- [ ] 不为小项目堆砌复杂设计模式。

搜索关键词：

```text
Go handler service repository pattern
Go dependency injection constructor
Go project layout cmd internal
```

掌握标准：能画出请求从handler到MySQL再返回的调用链，并说明每层为何存在。

## 19 context

- [ ] context用于取消、截止时间和请求级信息。
- [ ] `context.Background`与请求context。
- [ ] `WithTimeout`、`WithCancel`及调用cancel。
- [ ] context作为第一个参数传递。
- [ ] 不把context存进结构体长期保存。
- [ ] 数据库操作使用`*Context`方法。
- [ ] 不用context传递普通业务参数。

搜索关键词：

```text
Go context cancellation timeout best practices
Go HTTP request context database
```

掌握标准：客户端取消请求或超时后，数据库操作能够尽快停止。

## 20 并发基础

- [ ] goroutine。
- [ ] channel、缓冲和关闭。
- [ ] `select`。
- [ ] `sync.Mutex`和`RWMutex`。
- [ ] `sync.WaitGroup`。
- [ ] 数据竞争。
- [ ] `go test -race`。
- [ ] 不要通过“加锁”掩盖错误的共享状态设计。
- [ ] HTTP服务器会并发处理请求。

搜索关键词：

```text
Go goroutine channel select
Go Mutex WaitGroup
Go data race race detector
Go HTTP handler concurrent requests
```

掌握标准：内存repository在并发请求下没有map并发读写问题，并能解释锁保护的范围。

## 21 日志 配置与进程生命周期

- [ ] 环境变量配置。
- [ ] 启动时校验必需配置。
- [ ] 不提交真实密码。
- [ ] 使用`slog`记录键值字段。
- [ ] 日志级别概念。
- [ ] 请求ID的作用。
- [ ] Windows/Linux进程信号基本概念。
- [ ] HTTP服务优雅关闭。
- [ ] 启动失败应快速退出并打印明确原因。

搜索关键词：

```text
Go os.Getenv config
Go slog logger
Go HTTP graceful shutdown signal NotifyContext
```

## 22 测试

- [ ] `_test.go`文件和`TestXxx`。
- [ ] 表驱动测试。
- [ ] `t.Run`子测试。
- [ ] `httptest`。
- [ ] 用接口替换repository测试service。
- [ ] 单元测试与集成测试区别。
- [ ] 测试正常路径、边界和错误路径。
- [ ] `go test ./...`、`-v`、`-cover`、`-race`。

搜索关键词：

```text
Go table driven tests
Go httptest Gin
Go mock interface without framework
Go test cover race
```

掌握标准：新增接口至少覆盖成功、缺字段、非法状态和repository失败。

## 23 代码质量与安全基础

- [ ] 输入长度和枚举校验。
- [ ] SQL参数绑定。
- [ ] 不向客户端泄露内部错误。
- [ ] 密码与密钥放环境变量。
- [ ] HTTP与数据库超时。
- [ ] 限制分页大小。
- [ ] 正确使用HTTP状态码。
- [ ] 错误日志避免包含密码等敏感数据。
- [ ] `gofmt`、`go vet`与依赖整理。

搜索关键词：

```text
Go secure coding SQL injection error disclosure
REST API input validation pagination limit
Go HTTP server timeouts
```

## 24 Git与交付

- [ ] `git init`、`status`、`add`、`commit`、`log`。
- [ ] `.gitignore`。
- [ ] 一次提交表达一个清晰变化。
- [ ] README包含环境、配置、建库、启动、测试和接口示例。
- [ ] `.env.example`不含真实密码。
- [ ] migration脚本可以从空数据库建表。
- [ ] 构建产物不提交。

搜索关键词：

```text
Git basic workflow commit message
Go gitignore
README backend API project
database migration SQL
```

## 25 完成项目后再学

- [ ] Redis数据类型、过期、缓存穿透/击穿/雪崩。
- [ ] 后台任务和消息队列。
- [ ] 限流、重试、退避和幂等。
- [ ] Dockerfile与Docker Compose。
- [ ] pprof性能分析。
- [ ] Benchmark。
- [ ] Prometheus与可观测性。
- [ ] RAG、模型API、流式响应、Agent工具调用。

这些属于主项目阶段，不要阻塞当前CRUD项目。

## 26 当前无需深入

- 编译器源码。
- runtime调度器源码。
- 垃圾回收器源码。
- `unsafe`和CGO。
- 反射高级用法。
- 高级泛型技巧。
- 微服务框架和服务治理全家桶。
- Kubernetes。
- 复杂设计模式。
- 自己实现Web框架或ORM。

## 第一阶段开始项目的最低门槛

完成以下内容即可开始写内存版API，无需等待整份清单完成：

- [ ] 程序结构。
- [ ] 变量与控制流程。
- [ ] slice与map。
- [ ] 函数。
- [ ] 指针。
- [ ] 结构体和方法。
- [ ] error。
- [ ] Go Module。
- [ ] JSON基础。
- [ ] HTTP请求响应基本概念。

## 第二阶段接入MySQL的最低门槛

- [ ] interface基础。
- [ ] context基础。
- [ ] `database/sql`。
- [ ] SQL CRUD与索引基础。
- [ ] repository分层。
- [ ] 配置与日志。

## 自测问题

1. Module、package和Go文件分别是什么？
2. slice与数组有什么区别？append后为什么旧切片可能不受影响？
3. map读取不存在的key会发生什么？
4. 值接收者和指针接收者有什么区别？
5. Go接口为什么不需要显式声明实现？
6. `error`如何包装和判断？
7. `defer rows.Close()`解决什么问题？
8. 为什么`sql.Open`成功不代表数据库可连接？
9. handler、service、repository分别负责什么？
10. 如何避免SQL注入？
11. 为什么HTTP handler需要考虑并发？
12. context取消如何传递到数据库操作？
13. `go.mod`与`go.sum`分别有什么作用？
14. 如何测试一个依赖数据库repository的service？
15. 程序收到退出信号时如何优雅关闭HTTP服务？

