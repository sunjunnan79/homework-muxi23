# Kratos入门

## 简介

Kratos 是一套由Bilibili开源的轻量级 Go 微服务框架，包含大量微服务相关框架及工具。(虽然不知道什么时候倒闭,但还在做开源,这就是开源精神啊!)

> 名字来源于:《战神》游戏以希腊神话为背景，讲述奎托斯（Kratos）由凡人成为战神并展开弑神屠杀的冒险经历。

## 特性

go-kratos 是一个基于 Protobuf 和 gRPC 的微服务框架，提供了快速开发和部署微服务应用的工具和库。它支持**快速迭代**、**高性能**和**高可用性**。(好腻害,虽然不知道怎么这么厉害)

## 适用场景

适用于构建**高性能**的微服务架构，特别是对于需要**快速迭代和高可用性**的项目。

## 社区表现

| 框架      | Fork | Star  | contributers | year |
| --------- | ---- | ----- | ------------ | ---- |
| go-Kratos | 3.9K | 22.5K | 272          | 2019 |
| go-micro  | 2.3K | 21.4K | 189          | 2015 |
| go-zero   | 3.8K | 27.7K | 259          | 2020 |
| go-kit    | 2.4K | 26.1K | 205          | 2015 |

可以看出go-zero和go-kratos两位后起之秀实力强劲并且**社区活跃度比较高**,很有潜力。

先说为什么我觉得选go-kratos更好，首先是社区活跃度来说go-kratos略高于go-micro而且是国内的框架，足以说明这个框架必然很有独到之处。(主要是**中文看得懂**)

go-kratos与go-zero的区别我觉得主要还是在于其**轻量级的定位**，由于go-zero比较重量级，会有大量的约束使得自由度相对来说比较低

## 设计哲学

**简单**：不过度设计，代码平实简单；
**通用**：通用业务开发所需要的基础库的功能；
**高效**：提高业务迭代的效率；
**稳定**：基础库可测试性高，覆盖率高，有线上实践安全可靠；
**健壮**：通过良好的基础库设计，减少错用；
**高性能**：性能高，但不特定为了性能做 hack 优化，引入 unsafe ；
**扩展性**：良好的接口设计，来扩展实现，或者通过新增基础库目录来扩展功能；
**容错性**：为失败设计，大量引入对 SRE 的理解，鲁棒性高；
**工具链**：包含大量工具链，比如 cache 代码生成，lint 工具等等；
这是kratos官方挂出的框架设计出发点，其中有几点是在现有工具中尤为宝贵，并且十分契合go开发风格的。如 **简单，高效，扩展性，容错性**。

跟gin之类的轻量级框架很像,很符合Golang简洁高效的设计哲学,就是简洁舒服。

## 项目生态

Kratos是一个Go语言实现的微服务框架，说得更准确一点，它更类似于一个使用Go构建微服务的工具箱，开发者可以按照自己的习惯选用或定制其中的组件，来打造自己的微服务。也正是由于这样的原因，Kratos并不绑定于特定的基础设施，不限定于某种注册中心，或数据库ORM等，所以您可以十分轻松地将任意库集成进项目里，与Kratos共同运作。

围绕这样的核心设计理念，设计了如下的项目生态：

- [kratos](https://github.com/go-kratos/kratos) Kratos框架核心，主要包含了基础的CLI工具，内置的HTTP/gRPC接口生成和服务生命周期管理，提供链路追踪、配置文件、日志、服务发现、监控等组件能力和相关接口定义。
- [contrib](https://github.com/go-kratos/kratos/tree/main/contrib) 基于上述核心定义的基础接口，对配置文件、日志、服务发现、监控等服务进行具体实现所形成的一系列插件，可以直接使用它们，也可以参考它们的代码，做您需要的服务的适配，从而集成进kratos项目中来。
- [aegis](https://github.com/go-kratos/aegis) 我们将服务可用性相关的算法：如限流、熔断等算法放在了这个独立的项目里，几乎没有外部依赖，它更不依赖Kratos，您可以在直接在任意项目中使用。您也可以轻松将它集成到Kratos中使用，提高服务的可用性。
- [layout](https://github.com/go-kratos/kratos-layout) 我们设计的一个默认的项目模板，它包含一个参考了DDD和简洁架构设计的项目结构、Makefile脚本和Dockerfile文件。但这个项目模板不是必需的，您可以任意修改它，或使用自己设计的项目结构，Kratos依然可以正常工作。框架本身不对项目结构做任何假设和限制，您可以按照自己的想法来使用，具有很强的可定制性。
- [gateway](https://github.com/go-kratos/gateway) 这个是我们刚刚起步，用Go开发的API Gateway，后续您可以使用它来作为您Kratos微服务的网关，用于微服务API的治理，项目正在施工中，欢迎关注。

# 项目结构

 [kratos-layout](https://github.com/go-kratos/kratos-layout) 作为使用 `kratos new` 新建项目时所使用结构，其中包括了开发过程中所需的配套工具链( Makefile 等)，便于开发者更高效地维护整个项目，本项目亦可作为使用 Kratos 构建微服务的工程化最佳实践的参考。

![kratos ddd](https://go-kratos.dev/images/ddd.png)

使用如下命令即可基于 kratos-layout 创建项目：

```shell
kratos new <project-name>
```

生成的目录结构如下：

```text
  .
├── Dockerfile  
├── LICENSE
├── Makefile  
├── README.md
├── api // 下面维护了微服务使用的proto文件以及根据它们所生成的go文件
│   └── helloworld
│       └── v1
│           ├── error_reason.pb.go
│           ├── error_reason.proto
│           ├── error_reason.swagger.json
│           ├── greeter.pb.go
│           ├── greeter.proto
│           ├── greeter.swagger.json
│           ├── greeter_grpc.pb.go
│           └── greeter_http.pb.go
├── cmd  // 整个项目启动的入口文件
│   └── server
│       ├── main.go
│       ├── wire.go  // 我们使用wire来维护依赖注入
│       └── wire_gen.go
├── configs  // 这里通常维护一些本地调试用的样例配置文件
│   └── config.yaml
├── generate.go
├── go.mod
├── go.sum
├── internal  // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面，使用internal避免错误引用
│   ├── biz   // 业务逻辑的组装层，类似 DDD 的 domain 层，data 类似 DDD 的 repo，而 repo 接口在这里定义，使用依赖倒置的原则。
│   │   ├── README.md
│   │   ├── biz.go
│   │   └── greeter.go
│   ├── conf  // 内部使用的config的结构定义，使用proto格式生成
│   │   ├── conf.pb.go
│   │   └── conf.proto
│   ├── data  // 业务数据访问，包含 cache、db 等封装，实现了 biz 的 repo 接口。我们可能会把 data 与 dao 混淆在一起，data 偏重业务的含义，它所要做的是将领域对象重新拿出来，我们去掉了 DDD 的 infra层。
│   │   ├── README.md
│   │   ├── data.go
│   │   └── greeter.go
│   ├── server  // http和grpc实例的创建和配置
│   │   ├── grpc.go
│   │   ├── http.go
│   │   └── server.go
│   └── service  // 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)，同时协同各类 biz 交互，但是不应处理复杂逻辑
│       ├── README.md
│       ├── greeter.go
│       └── service.go
└── third_party  // api 依赖的第三方proto
    ├── README.md
    ├── google
    │   └── api
    │       ├── annotations.proto
    │       ├── http.proto
    │       └── httpbody.proto
    └── validate
        ├── README.md
        └── validate.proto
```

总体上还是很清晰的，虽然我还是有点不明白把infra去掉之后直接与数据库进行交互的curd的dao层放在哪里,还是放在data层下面吗?，还有greeter.go是干嘛的,欢迎吗?