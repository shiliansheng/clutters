-- 关键词表
drop table if exists ccgp_keyword;
create table ccgp_keyword(
	keyword_id varchar(40) not null COMMENT '关键词编号 md5(keyword)',
	keyword varchar(128) not null COMMENT '关键词内容 关键词之间用+连接 关键词之间为AND关系',
	ts timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	primary key(keyword_id)
) ENGINE=InnoDB AUTO_INCREMENT=1000 
  DEFAULT CHARSET=utf8 COMMENT '关键词表';

-- 结果表
drop table if exists ccgp_result;
create table ccgp_result(
	result_id varchar(40) not null COMMENT '结果编号 md5(url)',
	title varchar(512) not null COMMENT '标题',
	url varchar(512) not null COMMENT '链接地址',
	keyword varchar(256) not null COMMENT '结果依据关键词',
	pubdate datetime COMMENT '公告发布时间',
	task_id varchar(40) not null '任务编号',
	ts timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(result_id)
) ENGINE=InnoDB AUTO_INCREMENT=1000 
  DEFAULT CHARSET=utf8 COMMENT '查找结果表';