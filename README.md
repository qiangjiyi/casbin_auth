# casbin_auth
A base auth go web project combined with casbin open-source authorization library.

# 项目说明
项目代码中支持两种用户状态鉴权机制,分别是beego内置的session机制和目前比较流行的token机制.配置文件中如果设置`sessionon`为开启,则系统默认会使用cookie-session机制,
否则使用token模式, token模式一方面可以避免用户禁用token的情景, 另一方面也可以避免集群环境下session失效的情况. 项目中不论是session还是token模式均采取了reids
作为存储引擎, 因此在session模式下,即使重启了服务器只要cookie中携带的sessionID还未失效,则session信息依然可以获取到,无需重新登录.  
# 运行项目步骤
- 修改conf/app.conf中数据库连接信息
- 配置conf/app.conf中是否开启session
- 安装并启动redis,如果redis为远程服务,则需要根据实际情况修改连接redis配置
- 在项目根目录下新建logs文件夹,里面存放一些日志文件,github上忽略上传需要自己来建

# casbin_auth Api接口文档使用说明

## 1.1 版本历史
日期       |版本号       |作者   |备注
------------|-----------|-----------|-----------
2019.3.31       |1.0        |Qiangjiyi      |基于casbin的auth模块1.0版本

## 1.2 文档介绍
&emsp;&emsp;本文档的接口遵循RESTful设计风格

## 1.3 全局响应状态码说明

状态码       |说明       
------------|-----------
100000| 参数不合法
100001| 缺少参数
100002| 数据库操作失败
100003| 解析请求体失败
100004| 序列化数据失败
100005| 反序列化(解析)数据失败
100006| 获取服务的上下文失败
100007| 没有权限
100008| 服务内部错误
100009| shell脚本运行失败

## 1.4 Auth模块
### 1.4.1 登录
&emsp;&emsp;基于`session`或者`token`的机制来实现用户登录
> 注意事项：如果用户在app.conf中开启了sessionOn = true，则使用session机制，否则使用token。

#### 1.4.1.1 请求URL

http://localhost:8080/v1/api/auth/user/login

#### 1.4.1.2 请求方式

POST

#### 1.4.1.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.1.4 请求Body参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
username |是   |用户名|string
password |是  | 密码    |string

**请求示例：**
```json
{
    "username": "admin",
    "password": "123456"
}
```

#### 1.4.1.5 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
.user_id | number    |用户ID
.token|  string|登录令牌

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "user_id": 1,
        "token": "gx7KDTNEA26lFcLGHaP0o645aLhe49gP"
    }
}
```

#### 1.4.1.6 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功    
100103       |用户名不正确    
100104       |密码不正确 

### 1.4.2 注册

&emsp;&emsp;根据请求参数注册一个新用户，用户名不能重复

#### 1.4.2.1 请求URL

http://localhost:8080/v1/api/auth/user/register 

#### 1.4.2.2 请求方式

POST

#### 1.4.2.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.2.4 请求Body参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
username |是   |用户名|string
password |是  | 密码    |string
real_name |否  | 真实姓名    |string
mobile |是  | 手机号    |string
email |否  | 邮箱    |string

**请求示例：**
```json
{
    "username": "qiangjiyi",
    "password": "123456",
    "real_name": "真实姓名",
    "mobile": "18611876729",
    "email": "qiangjiyi@163.com"
}
```

#### 1.4.2.5 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
.user_id | number    |用户ID
.username|  string|用户名
.create_time|  string|注册时间

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "user_id": 1,
        "username": "qiangjiyi",
        "create_time": "2019-01-31 16:43:23"
    }
}
```

#### 1.4.2.6 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功    
100102       |用户名已经被占用    

### 1.4.3 退出
&emsp;&emsp;移除当前登录用户会话信息

#### 1.4.3.1 请求URL

http://localhost:8080/v1/api/auth/user/logout 

#### 1.4.3.2 请求方式

GET

#### 1.4.3.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.3.4 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK"
}
```

#### 1.4.3.5 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功    
100106       |注销登录失败   

### 1.4.4 更新用户

&emsp;&emsp;根据请求参数更新用户资料

#### 1.4.4.1 请求URL

http://localhost:8080/v1/api/auth/user/update 

#### 1.4.4.2 请求方式

POST

#### 1.4.4.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.4.4 请求Body参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
id |是   |用户ID|number
username |是   |用户名|string
password |是  | 密码    |string
real_name |否  | 真实姓名    |string
mobile |是  | 手机号    |string
email |否  | 邮箱    |string

**请求示例：**
```json
{
    "id": 1,
    "username": "qiangjiyi1",
    "password": "123456789",
    "real_name": "真实姓名xxx",
    "mobile": "18611876729",
    "email": "1473430734@qq.com"
}
```

#### 1.4.4.5 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
.user_id | number    |用户ID
.username|  string|用户名
.update_time|  string|更新时间

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "user_id": 1,
        "username": "qiangjiyi1",
        "update_time": "2019-01-31 17:00:37"
    }
}
```

#### 1.4.4.6 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功    
100100       |用户ID不合法   
100101       |用户ID对应的用户不存在

### 1.4.5 删除用户

&emsp;&emsp;根据请求参数删除对应的用户

#### 1.4.5.1 请求URL

http://localhost:8080/v1/api/auth/user/delete/1 

#### 1.4.5.2 请求方式

GET

#### 1.4.5.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.5.4 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK"
}
```

#### 1.4.5.5 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功    
100100       |用户ID不合法   
100101       |用户ID对应的用户不存在

### 1.4.6 创建角色

&emsp;&emsp;根据请求参数新建一个角色

#### 1.4.6.1 请求URL

http://localhost:8080/v1/api/auth/role/create

#### 1.4.6.2 请求方式

POST

#### 1.4.6.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.6.4 请求Body参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
role_name |是   |角色名称|string
description |是   |角色描述|string

**请求示例：**
```json
{
    "role_name": "admin",
    "description": "管理员"
}
```

#### 1.4.6.5 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
.id | number    |角色ID
.role_name|  string|角色名称
.description|  string|角色描述
.create_time|  string|创建时间

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "id": 15,
        "role_name": "admin",
        "description": "管理员",
        "create_time": "2019-01-31 17:12:29"
    }
}
```

#### 1.4.6.6 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功   

### 1.4.7 更新角色

&emsp;&emsp;根据请求参数更新指定的角色信息

#### 1.4.7.1 请求URL

http://localhost:8080/v1/api/auth/role/update 

#### 1.4.7.2 请求方式

POST

#### 1.4.7.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.7.4 请求Body参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
id |是   |角色ID|number
role_name |是   |角色名称|string
description |是   |角色描述|string

**请求示例：**
```json
{
    "id": 1,
    "role_name": "system_admin",
    "description": "系统管理员"
}
```

#### 1.4.7.5 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
.id | number    |角色ID
.role_name|  string|角色名称
.description|  string|角色描述
.update_time|  string|更新时间

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "id": 1,
        "role_name": "system_admin",
        "description": "系统管理员",
        "update_time": "2019-01-31 17:18:01"
    }
}
```

#### 1.4.7.6 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功   
100200 | 角色ID不合法
100201 | 角色ID对应的角色不存在

### 1.4.8 删除角色

&emsp;&emsp;根据请求参数删除指定的角色信息

#### 1.4.8.1 请求URL

http://localhost:8080/v1/api/auth/role/delete/1 

#### 1.4.8.2 请求方式

GET

#### 1.4.8.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.8.4 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK"
}
```

#### 1.4.8.5 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功   
100200 | 角色ID不合法
100201 | 角色ID对应的角色不存在

### 1.4.9 添加权限

&emsp;&emsp;根据请求参数添加权限信息

#### 1.4.9.1 请求URL

http://localhost:8080/v1/api/auth/permission/add 

#### 1.4.9.2 请求方式

POST

#### 1.4.9.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.9.4 请求Body参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
type |是   |权限类型 ['p' 或 'g']|string
sub |是   |权限主体值|string
sub_type |是   |权限主体  ['u' 或 'r']|string
obj|是 |权限资源值|string
action|否 |请求方式|string

> 说明：如果type为g，则说明添加的是权限组，那么sub_type如果为u，则sub为用户名，否则sub_type为r，则sub为角色ID，而obj此时必须为角色ID；如果type为p，则说明添加的是普通的权限，sub_type和sub的值跟type为g时要求类似，只是obj此时就是我们要访问的restful api地址，action为我们要访问该api地址的请求方式。

**请求示例：**
```json
{
    "type": "g",
    "sub": "qiangjiyi",
    "sub_type": "u",
    "obj": "1",
    "action": "GET"
}
```

#### 1.4.9.5 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
.policy | string    |casbin格式的权限策略

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "policy": "g, 2, role_1"
    }
}
```

#### 1.4.9.6 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功   
100300| 权限执行模块加载失败
100301| 权限类型不合法 ['p' 或 'g']
100302| 权限主体类型不合法 ['u' 或 'r']
100303| 权限主体值不合法(不存在对应的主体)
100304| 权限资源值不合法(不存在对应的角色)
100305| 权限已存在

### 1.4.10 查询用户的角色和权限

&emsp;&emsp;根据请求参数查询指定用户的所有角色和权限信息

#### 1.4.10.1 请求URL

http://localhost:8080/v1/api/auth/permission/query/1 

#### 1.4.10.2 请求方式

GET

#### 1.4.10.3 请求Headers参数

参数名|必选|说明|默认值
:---    |:-: |:-----   |-----
Content-Type |是   |接受的Content-Type类型|application/json;charset=UTF-8
token |是  | 登录令牌    |空

> 说明：*如果是token模式的话，请求Headers中必须包含token参数键值对，否则需要开启cookie。*

#### 1.4.10.4 响应Response参数

参数名|类型|描述
:---    | :-----   |-----
err_code |number|响应码
err_message |string    |响应描述
payload |object|响应体
&emsp;.user_id |number    |用户ID
&emsp;.username |string|用户名
&emsp;.real_name |string    |真实姓名
&emsp;.mobile |string|手机号
&emsp;.email |string    |电子邮箱
&emsp;.role_list |array list|角色列表
&emsp;&emsp;..id |number    |角色ID
&emsp;&emsp;..role_name |string|角色名称
&emsp;&emsp;..description |string    |角色描述
&emsp;&emsp;..update_time |string    |角色最近更新时间
&emsp;.permission_list| []string| 权限列表

**响应示例：**
```json
{
    "err_code": 0,
    "err_message": "OK",
    "payload": {
        "user_id": 2,
        "username": "qiangjiyi",
        "real_name": "强吉义",
        "mobile": "18611876729",
        "email": "qiangjiyi@163.com",
        "role_list": [
            {
                "id": 5,
                "role_name": "组织管理员",
                "description": "作为一个组织内部的管理员",
                "update_time": "2019-01-21 20:18:55"
            },
            {
                "id": 2,
                "role_name": "角色name",
                "description": "角色description",
                "update_time": "2019-01-23 16:24:48"
            },
            {
                "id": 3,
                "role_name": "新角色name",
                "description": "新角色description",
                "update_time": "2019-01-29 11:11:48"
            }
        ],
        "permission_list": [
            "2,/v1/api/auth/permission/add,POST",
            "2,/v1/api/auth/role/delete/*,GET",
            "2,/v1/api/auth/role/delete/:roleId,GET",
            "role_3,/v2/api/auth/permission/query,GET"
        ]
    }
}
```

#### 1.4.10.5 响应状态码说明

状态码       |说明       
------------|-----------
0       |请求成功   
100100       |用户ID不合法   
100101       |用户ID对应的用户不存在

# 写在最后
通过以上的接口文档,相信有过go web开发经验的朋友们应该可以将该项目run起来,该部分代码是我之前做的一个比较通用的用户角色权限模块,最近才决定抽取并开源出来,不足之处也希望大家能够指正.
如果在使用过程中遇到问题,可以单独与我交流.  
![image](https://github.com/qiangjiyi/casbin_auth/blob/master/images/qjy1473430734.jpg)