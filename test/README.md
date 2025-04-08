# Fiber-RBAC 测试指南

本目录包含 Fiber-RBAC 项目的测试文件，包括单元测试、集成测试和性能基准测试。

## 目录结构

```
test/
├── database/              # 数据库相关测试
│   ├── connection_test.go # 数据库连接测试
│   ├── benchmark_test.go  # 数据库性能基准测试
│   └── testconfig.yaml    # 测试配置文件
├── api/                   # API 测试
├── middleware/            # 中间件测试
├── repository/            # 数据访问层测试
├── service/               # 业务逻辑层测试
├── utils/                 # 工具函数测试
└── README.md              # 本文件
```

## 运行测试

### 运行所有测试

在项目根目录运行：

```bash
go test ./test/...
```

### 运行数据库测试

```bash
go test ./test/database/...
```

### 运行 API 测试

```bash
go test ./test/api/...
```

### 运行性能基准测试

```bash
# 运行所有基准测试
go test -bench=. ./test/...

# 运行数据库基准测试
go test -bench=. ./test/database/benchmark_test.go

# 运行指定的基准测试
go test -bench=BenchmarkInitDB ./test/database/benchmark_test.go
```

### 生成测试覆盖率报告

```bash
# 生成覆盖率数据
go test -coverprofile=coverage.out ./test/...

# 查看覆盖率报告
go tool cover -html=coverage.out
```

## 数据库测试说明

数据库测试使用 SQLite 内存数据库，不需要外部数据库支持。主要测试以下功能：

1. 数据库连接初始化
2. 数据库表结构迁移
3. 数据库查询性能
4. 数据库事务处理
5. 数据库连接池配置

## 编写新测试

### 单元测试

单元测试应遵循 Go 的标准测试规范，测试函数名应以 `Test` 开头：

```go
func TestSomething(t *testing.T) {
    // 测试逻辑
}
```

### 基准测试

基准测试应以 `Benchmark` 开头：

```go
func BenchmarkSomething(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 测试逻辑
    }
}
```

### 表驱动测试

对于需要多组数据测试的功能，应使用表驱动测试：

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "expected1"},
        {"case2", "input2", "expected2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionToTest(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 测试依赖

项目使用以下测试工具：

1. `testing` - Go 标准测试包
2. `testify/assert` - 断言库
3. `testify/suite` - 测试套件工具
4. `testify/mock` - 模拟工具 