# 权限管理服务平台

该服务采用比较流行的微服务思想, 利用[**shiro**](http://shiro.apache.org/)<域，动作，实体>思想，实现权限管理平台服务， 它支持**单用户多角色** , 比RBAC的资源管理更细粒度化  

权限管理服务平台的实现，主要由三个步骤构成： 
+ 第一步：判断URL是否在白名单中，如果是，直接返回状态码：200  
+ 第二步：解析URL，获取功能ID和实体ID，服务初始化阶段，会构建多叉树  
+ 第三步：获取<域，动作，实体>, 并在redis中采用SET集合存储<UserId-FuncId, SET集合={实体1, 实体2, ... , 实体N}>

## 权限管理库表设计

[权限管理库表](table.md)

## 环境依赖

+ [beego框架](https://beego.me/)
+ [redis](https://redis.io/)

## OpenResty配置

权限管理安插在Nginx Access访问阶段，对http请求的合法性进行校验

access_by_lua_file "/data/openresty/lua_files/test_ycfm_lua_files/access_by_grbac.lua"

```lua
-- GRBAC权限管理模块
ngx.req.read_body()
local bodyData = ngx.req.get_body_data()
ngx.log(ngx.ERR, "body data:", bodyData)
local cjson = require "cjson"
local info={
        ["body"] =  bodyData,
        ["method"] = ngx.req.get_method(),
        ["uri"] = ngx.var.uri,
}
local encode = cjson.encode(info)
local res = ngx.location.capture('/v1/grbac/functions/tree_parsing', {method=ngx.HTTP_POST, body=encode})
if res.status == 403 then
        ngx.exit(ngx.HTTP_FORBIDDEN)
end
```
## DEMO
![demo](grbac_demo.jpg)
## 说明

+ `希望与大家一起成长，有任何该服务运行或者代码问题，可以及时找我沟通，喜欢开源，热爱开源, 欢迎多交流`   
+ `联系方式：cdh_cjx@163.com`
