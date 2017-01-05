# 创建权限管理库
```
CREATE DATABASE IF NOT EXISTS ycfm_grbcs DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```
## 创建实体表
```
CREATE TABLE IF NOT EXISTS `entities` (
  `entity_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `region_id` int(11) DEFAULT NULL COMMENT '域ID',
  `name` varchar(300) DEFAULT NULL COMMENT '名称',
  `third_id` int(11) DEFAULT NULL COMMENT '第三方ID',
  `status` smallint(6) DEFAULT NULL COMMENT '状态：-20:逻辑删除；10:正常; 20:无效',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`entity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

### 创建功能表
```
CREATE TABLE IF NOT EXISTS  `functions` (
  `function_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(60) DEFAULT NULL COMMENT '名称',
  `region_id` int(11) DEFAULT NULL COMMENT '域ID',
  `uri` varchar(500) DEFAULT NULL COMMENT '功能的统一资源标识符',
  `third_region_mark_key` varchar(300) DEFAULT NULL COMMENT ':xxx: 表示在URL的正则表达参数中，名称为xxx; GET:xxx 表示在URL的kv参数中，名称为xxx； body_kv:xxx 表示在请求参数body的kv参数中，名称为xxx；session:xxx表示在session中，名 称为xxx',
  `method_type` smallint(6) DEFAULT NULL COMMENT '10:GET;20:POST;30:PUT;40:PATCH;50:DELETE',
  `status` smallint(6) DEFAULT NULL COMMENT '状态：-20:逻辑删除；10正常；20: 已废弃',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`function_id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4
```

###  创建角色与功能关系表
```
CREATE TABLE `role_function_relationships` (
  `role_function_relationship_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `role_id` int(11) DEFAULT NULL COMMENT '角色ID',
  `function_id` int(11) DEFAULT NULL COMMENT '功能ID',
  `region_id` int(11) DEFAULT NULL COMMENT '域ID',
  `status` smallint(6) DEFAULT NULL COMMENT '状态：-20：逻辑删除；10：有效；20：无效',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`role_function_relationship_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4
```

###  创建角色表
```
CREATE TABLE IF NOT EXISTS `roles` (
  `role_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `region_id` int(11) DEFAULT NULL COMMENT '域ID',
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `status` smallint(6) DEFAULT NULL COMMENT '状态：-20: 逻辑删除;10: 正常; 20:冻结',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4
```

### 创建用户与角色表
```
CREATE TABLE IF NOT EXISTS  `user_roles` (
  `user_role_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` int(11) DEFAULT NULL COMMENT '用户ID',
  `role_id` int(11) DEFAULT NULL COMMENT '角色ID',
  `region_id` int(11) DEFAULT NULL COMMENT '域ID',
  `status` smallint(6) DEFAULT NULL COMMENT '状态：-20:逻辑删除;10:正常; 20:冻结',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`user_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

### 创建白名单表
```
CREATE TABLE IF NOT EXISTS `white_list` (
  `white_list_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `url` varchar(100) DEFAULT NULL COMMENT '白名单URL',
  `desc` varchar(300) DEFAULT NULL COMMENT '备注',
  `status` smallint(6) DEFAULT NULL COMMENT '状态：-20:逻辑删除；10:正常; 20:无效',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`white_list_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```
