/* plugin database 
-- Chen Qian(chinahbcq@qq.com)
-- Date: 2016.05.2 19:16 
------------------------------------------------------
*/

/*
-- create a database
-- CREATE DATABASE `smartapi` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
-- grant privileges
-- GRANT ALL PRIVILEGES ON smartapi.* TO 'im_rdqa'@'%' WITH GRANT OPTION;
*/
DROP TABLE IF EXISTS `user_info`;

CREATE TABLE `user_info` (
    `uid` bigint NOT NULL DEFAULT 0,
    `name` varchar(30) NOT NULL,
    `gender` varchar(30) NOT NULL,
    PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into user_info values(123, "smartapi", "male");
insert into user_info values(124, "test", "female");
