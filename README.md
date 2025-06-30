### 目录结构

```
kano/
├── Makefile                         # 构建脚本（编译、运行等）
├── server/                          # 核心服务代码
│   ├── api/dto                      # API 相关代码（数据传输对象 DTO）
         # ├── common.go                           #通用 DTO（如分页、公共字段）
         # ├── upload.go                           # 文件上传相关的 DTO（如请求参数、响应结果）
│   ├── cmd /                        # 可执行文件入口（main 函数）
         # ├── main.go                             # main 函数是程序启动的入口，负责初始化服务、加载配置、注册路由等
│   ├── configs/                     # 环境配置（开发、生产、测试环境）
         # ├── config-dev.yaml                     #  开发环境配置文件解析
         # ├── config-prod.yaml                    #  生产环境配置文件解析
         # ├── config-test.yaml                    #  测试环境配置文件解析
│   ├── docs/                        # 文档（Swagger 相关）
          # ├── docs.go                            #  用于生成文档的工具代码（通常由工具自动生成）。
          # ├── swagger.json                       #  文档内容（接口路径、参数、响应示例等）
          # ├── swagger.yaml       
│   ├── internal/                    # 内部模块（业务逻辑核心）
         # ├── config/                             # 配置解析 configs 中的 YAML 配置，提供统一的配置访问接口
                 # ├── config.go                   #  加载配置文件，映射到结构体。
                 # ├── database.go                 # 数据库配置 初始化数据库连接
         # ├──  handler/v1           # 接口处理器
                 # ├── upload.go                   # 处理 HTTP 请求，校验参数，调用 Service 层逻辑，返回响应。
         # ├──  logger/               # 日志模块
                 # ├── gin-logger.go               # 基础日志接口和实现
                 # ├── logger.go                   # Gin 框架的日志中间件，记录请求链路信息
         # ├── middleware/            # 中间件
                 # ├── login_auth.go               # 实现登录认证中间件，校验请求中的 Token。
         # ├── provider/tencent       # 第三方服务集成 //封装第三方服务（如腾讯云、阿里云）的 SDK，提供统一调用接口
                  # ├── tencent.go                 # 腾讯云服务集成
         # ├── repository/            # 数据仓储 负责数据库操作，实现数据的增删改查，隔离业务逻辑与具体数据库实现
                  # ├── model                     # 定义数据库模型（GORM 结构体映射）。
                  # ├── upload_resp.go            #实现具体的数据操作方法
         # ├── service/               # 业务逻辑层 实现具体业务逻辑（如文件上传流程、权限校验等）。
                  # ├── upload.go
├── pkg/                              # 可复用工具包
         # ├── response.go/                            # 响应处理 统一响应格式
         # ├── utils.go/                               # 工具包通用客户端工具
│   ├── logs/                         # 日志文件
│   ├── pkg/                          # 可复用工具包
│   ├── router/                       # 路由定义 注册 API 路由，绑定 Handler 函数，通常使用 Gin 框架的路由分组（如 /v1/upload）。

│   ├── go.mod/go.sum                 # Go 模块依赖管理
│   └── README.md                     # 项目说明文档

```
