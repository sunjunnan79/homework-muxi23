# gRPC

### gRPC协议架构

gRPC是一种用于实现RPC API的技术。由于gRPC是开源框架，通信双方都基于该框架进行二次开发，从而使得通信双方聚焦在业务，无需关注由gRPC软件框架实现的底层通信。如下图，DATA部分即为业务层面内容，DATA下面所有的信息都由gRPC进行封装。

![gRPC协议架构](https://download.huawei.com/mdl/image/download?uuid=1914371252454b70bf3fbfcc59e5d175)
*gRPC协议架构*

从这里可以看出grpc是支持http的,所以大部分的语言在使用grpc的同时也**支持restful**,所以既能够用于微服务**内部通信**也能够和客户端**直接通信**。

### gRPC支持的操作

设备在网络架构里支持Dial-in和Dial-out两种对接模式。

1. Dial-in模式：设备作为gRPC服务器，采集器作为gRPC客户端。由采集器主动向设备发起gRPC连接并获取需要采集的数据信息或下发配置。Dial-in模式适用于小规模网络和采集器需要向设备下发配置的场景。

   Dial-in模式支持以下操作：

   - **Subscribe操作**：高速采集设备的接口流量统计、CPU和内存等数据信息。当前仅支持基于[Telemetry](https://info.support.huawei.com/info-finder/encyclopedia/zh/Telemetry.html)技术的Subscribe操作。
   - **Get操作**：获取设备运行状态和运行配置。当前仅支持基于gNMI（gRPC Network Management Interface）规范的Get操作。
   - **Capabilities操作**：获取设备能力数据。当前仅支持基于gNMI规范的Capabilities操作。
   - **Set操作**：向设备下发配置。当前仅支持基于gNMI规范的Set操作。

2. Dial-out模式：设备作为gRPC客户端，采集器作为gRPC服务器。设备主动和采集器建立gRPC连接，将设备上配置的订阅数据推送给采集器。Dial-out模式适用于网络设备较多的情况下，由**设备主动向采集器**提供设备数据信息。Dial-out模式只支持基于Telemetry技术的Subscribe操作。(听不懂思密达)

### gRPC交互过程

如下图，gRPC采用**客户端和服务器**模型，以网络设备为gRPC客户端，采集器为gRPC服务器为例，说明gRPC的交互过程：

1. 设备在开启gRPC功能后作为gRPC客户端，采集器作为gRPC服务器。
2. 设备会根据应用服务（如**订阅的事件**）**构建对应数据**的格式（GPB/JSON），通过ProtoBuf（Protocol Buffers）编写Proto文件。然后，设备与采集器**建立gRPC通道**，通过gRPC协议向采集器发送请求消息。
3. 采集器收到请求消息后，会**通过ProtoBuf解译Proto**文件，还原出事先定义好的数据结构，进行业务处理。
4. 采集器处理完数据后，需要使用**ProtoBuf重新编译**应答数据，通过gRPC协议向设备发送应答消息。
5. 设备收到应答消息后，结束本次的gRPC交互。

简单地说，**设备主动和采集器建立**gRPC连接，将设备上配置的订阅数据推送给采集器。在整个gRPC交互的过程中，设备和采集器都需要**使用ProtoBuf来定义Proto文件**。

![gRPC交互过程](https://download.huawei.com/mdl/image/download?uuid=422b463a2a604a4c9fae8deec8186b6b)
*gRPC交互过程*

### gRPC的应用（用于跟硬件交互？）

gRPC支持通过**Telemetry技术**实现订阅功能（Subscribe操作）。Telemetry是一项远程的从物理设备或虚拟设备上高速采集数据的技术。设备通过推模式（Push Mode）周期性地主动向采集器上送设备的接口流量统计、CPU和内存数据等信息。

如下图所示，网络设备和网络管理系统建立gRPC连接后，网络管理系统可以订阅设备上指定模块的数据信息。Telemetry有动态订阅和静态订阅两种方式，**动态订阅基于Dial-in模式建立，静态订阅基于Dial-out模式建立。**

![基于gRPC的Telemetry技术](https://download.huawei.com/mdl/image/download?uuid=e644dfe5c67a41e99df6f1ddbd80c9b6)
*基于gRPC的Telemetry技术*

Telemetry的实现流程：

1. 用户定义Telemetry**静态订阅**或Telemetry**动态订阅。**
   - Telemetry静态订阅：在huawei-grpc-dialout.proto文件中定义。
   - Telemetry动态订阅：在huawei-grpc-dialin.proto文件中定义。
2. 用户将采集到的信息通过GPB或JSON格式进行编码，在huawei-telemetry.proto文件里定义采样路径、采样时间戳等重要信息。
   - GPB编码时，huawei-telemetry.proto文件中的encoding字段为Encoding_GPB（值为0），data_gpb字段承载GPB编码格式的采样数据，data_str字段为空。
   - JSON编码时，huawei-telemetry.proto文件中的encoding字段为Encoding_JSON（值为1），data_str字段承载JSON编码格式的采样数据，data_gpb字段为空。
3. 设备传输数据到采集器，解码数据并分析结果。
   - huawei-telemetry.proto文件中data_gpb字段内容需要相应的业务proto文件进行解码，由huawei-telemetry.proto文件中的sensor_path字段标识对应哪个具体的业务proto文件，例如，当sensor_path取值为huawei-ifm:ifm/interfaces/interface时，其数据结构定义在huawei-ifm.proto文件中。
   - 当采用纯JSON编码格式（编码层和数据模型层均为JSON编码格式。）时，用户只需要对huawei-grpc-dialout.proto文件或huawei-grpc-dialin.proto文件进行解码。当采用混合JSON编码格式（编码层为GPB编码格式，数据模型层为JSON编码格式。）时，用户只需要对huawei-grpc-dialout.proto文件或huawei-grpc-dialin.proto文件和huawei-telemetry.proto文件进行解码，不需要相应的业务proto文件。

## 什么是gRPC ProtoBuf?

gRPC ProtoBuf是gRPC协议的接口描述语言，是一种与语言无关、平台无关、扩展性好的用于通信协议、数据存储的序列化结构数据格式。gRPC ProtoBuf编码格式也称为GPB（Google Protocol Buffers）编码格式。GPB提供了一种灵活、高效、自动序列化结构数据的机制。GPB与XML、JSON编码类似，也是一种编码方式，但不同的是，它是一种**二进制编码，性能好，效率高。**

目前，GPB包括v2和v3两个版本，设备当前支持的GPB版本是v3。

GPB在gRPC的框架中主要有三个作用：

- 定义数据结构

  **xxx.proto**:

  ```
  syntax = "proto3"; //proto版本定义为v3版本。
  
  message serviceArgs { //消息格式描述。
      int64 ReqId = 1; //请求ID。
      oneof MessageData {
          bytes data = 2; //表示承载GPB编码格式的采样数据。
          string data_json = 4; //表示承载JSON编码格式的采样数据。
          }
      string errors = 3; //产生错误时的描述信息。
  }
  ```

- 定义服务接口

  **xxx.proto**:

  ```
  syntax = "proto3"; //proto版本定义为v3版本。
  package huawei_dialout; //本包名称为huawei_dialout。
  service gRPCDataservice { //服务名称为gRPCDataservice。
      rpc dataPublish(stream serviceArgs) 
      returns(stream serviceArgs) {}; //方法为dataPublish，双向流，提供数据推送方法。入参是serviceArgs数据流。
  }
  ```

  

- 通过**序列化和反序列化**提升传输效率

  GPB编码格式的内容只是提供给操作者阅读的，实际上并不会以这种文本形式进行传输，而是以序列化后的**二进制数据**进行传输。而JSON编码格式则以数据文本形式呈现，传输时也以数据文本形式传输，所以GPB编码格式的传输效率相对JSON、XML编码格式有着天然的优势。

  **GPB编码**:

  ```
  {
      1:"HUAWEI"
      2:"s4"
      3:"huawei-ifm:ifm/interfaces/interface"
      4:46
  }
  ```

  **JSON**:

  ```
  {
      "node_id_str":"HUAWEI",
      "subscription_id_str":"s4",
      "sensor_path":"huawei-ifm:ifm/interfaces/interface",
      "collection_id":46,
  }
  ```

## 什么是Proto文件?

gRPC协议用**GPB编码**格式承载数据，GPB编码格式的文件名后缀为.proto，即为Proto文件。

GPB通过“.proto”文件描述编码使用的字典，即数据结构描述。采集器可以利用**Protoc等工具软件**根据“.proto”文件自动生成代码（例如java代码），然后用户基于自动生成的代码进行二次开发对获取到的数据进行解析，从而实现与设备的数据对接。

Proto文件包含公共Proto文件和业务数据Proto文件。

### 公共Proto文件

[Telemetry](https://info.support.huawei.com/info-finder/encyclopedia/zh/Telemetry.html)提供3个公共的proto文件，支持数据上送和订阅功能：

- huawei-grpc-dialout.proto文件：定义了设备作为gRPC客户端对外推送数据。
- huawei-grpc-dialin.proto文件：定义了设备作为gRPC服务端对外推送数据。
- huawei-telemetry.proto文件：定义了Telemetry采样数据上送时的数据头，包括采样路径，采样时间戳等重要信息。

### 业务数据Proto文件

设备提供多个业务数据Proto文件，用于定义具体业务数据的GPB编码，采集器侧需要根据实际要监控的业务选择对应proto文件。

### 一个简单的Proto文件示例

Telemetry静态订阅功能：设备作为gRPC客户端，采集器作为gRPC服务端，由设备主动发起到采集器的连接，进行数据采集上送。

**表1-1** proto文件内容:**huawei-grpc-dialout.proto**

```
syntax = "proto3"; //proto版本定义为v3版本。

package huawei_dialout; //本包名称为huawei_dialout。

service gRPCDataservice { //服务名称为gRPCDataservice。
	rpc dataPublish(stream serviceArgs) returns(stream serviceArgs) {}; //方法为dataPublish，双向流，提供数据推送方法。入参是serviceArgs数据流。
}

message serviceArgs { //消息格式描述。
	int64 ReqId = 1; //请求ID。
    oneof MessageData {
        bytes data = 2; //表示承载GPB编码格式的采样数据。
        string data_json = 4; //表示承载JSON编码格式的采样数据。
    }
    string errors = 3; //产生错误时的描述信息。
}
```

