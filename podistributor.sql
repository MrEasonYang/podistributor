CREATE DATABASE /*!32312 IF NOT EXISTS*/ `podistributor` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `podistributor`;

CREATE TABLE `episodes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `podcast_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'Related podcast ID',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT 'Episode unique name',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT 'Episode nickname',
  `main_uri_list` text NOT NULL COMMENT 'Episode streaming uri list in JSON array format',
  `backup_url_list` text NOT NULL COMMENT 'Episode streaming backup uri list in JSON array format',
  `analysis_url_list` text NOT NULL COMMENT 'Analysis url list in JSON array format',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Create timestamp',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update timestamp',
  PRIMARY KEY (`id`),
  UNIQUE KEY `pod_ep` (`podcast_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `podcasts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT 'Podcast unique name',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT 'Podcast nickname',
  `enabled` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Whether the record of podcast enabled or not',
  `rss` varchar(256) NOT NULL DEFAULT '' COMMENT 'The RSS url of the podcast',
  `domain` varchar(32) NOT NULL DEFAULT '' COMMENT 'The domain of the podcast',
  `episode_url_domain` varchar(32) NOT NULL DEFAULT '' COMMENT 'The domain of every episode',
  `episode_main_url_level` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Episode resource URI array index',
  `episode_backup_url_enabled` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Whether the backup url are enabled or not',
  `episode_backup_url_level` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Backup episode resource URI array index',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Create timestamp',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update timestamp',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
