CREATE DATABASE /*!32312 IF NOT EXISTS*/ `xnightwatch` /*!40100 DEFAULT CHARACTER SET latin1 COLLATE latin1_swedish_ci */;

USE `xnightwatch`;

--
-- Table structure for table `api_chain`
--

DROP TABLE IF EXISTS `api_chain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_chain` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `namespace` varchar(253) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '区块链名',
  `display_name` varchar(253) NOT NULL DEFAULT '' COMMENT '区块链展示名',
  `miner_type` varchar(16) NOT NULL DEFAULT '' COMMENT '区块链矿机机型',
  `image` varchar(253) NOT NULL DEFAULT '' COMMENT '区块链镜像 ID',
  `min_mine_interval_seconds` int(8) NOT NULL DEFAULT 0 COMMENT '矿机挖矿间隔',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_namespace_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='区块链表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_miner`
--

DROP TABLE IF EXISTS `api_miner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_miner` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `namespace` varchar(253) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机名',
  `display_name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机展示名',
  `phase` varchar(45) NOT NULL DEFAULT '' COMMENT '矿机状态',
  `miner_type` varchar(16) NOT NULL DEFAULT '' COMMENT '矿机机型',
  `chain_name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机所属的区块链名',
  `cpu` int(8) NOT NULL DEFAULT 0 COMMENT '矿机 CPU 规格',
  `memory` int(8) NOT NULL DEFAULT 0 COMMENT '矿机内存规格',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_namespace_name` (`namespace`,`name`),
  KEY `idx_chain_name` (`chain_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='矿机表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_minerset`
--

DROP TABLE IF EXISTS `api_minerset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_minerset` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `namespace` varchar(253) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机池名',
  `replicas` int(8) NOT NULL DEFAULT 0 COMMENT '矿机副本数',
  `display_name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机池展示名',
  `delete_policy` varchar(32) NOT NULL DEFAULT '' COMMENT '矿机池缩容策略',
  `min_ready_seconds` int(8) NOT NULL DEFAULT 0 COMMENT '矿机 Ready 最小等待时间',
  `fully_labeled_replicas` int(8) NOT NULL DEFAULT 0 COMMENT '所有标签匹配的副本数',
  `ready_replicas` int(8) NOT NULL DEFAULT 0 COMMENT 'Ready 副本数',
  `available_replicas` int(8) NOT NULL DEFAULT 0 COMMENT '可用副本数',
  `failure_reason` longtext DEFAULT NULL COMMENT '失败原因',
  `failure_message` longtext DEFAULT NULL COMMENT '失败信息',
  `conditions` longtext DEFAULT NULL COMMENT '矿机池状态',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_namespace_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='矿机池表';
/*!40101 SET character_set_client = @saved_cs_client */;