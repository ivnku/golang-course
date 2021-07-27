-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               8.0.15 - MySQL Community Server - GPL
-- Server OS:                    Win64
-- HeidiSQL Version:             11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for redditclone
CREATE DATABASE IF NOT EXISTS `redditclone` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `redditclone`;

-- Dumping structure for table redditclone.comments
CREATE TABLE IF NOT EXISTS `comments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `post_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `body` varchar(500) COLLATE utf8mb4_general_ci NOT NULL,
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK2_comments_post-id_posts_id` (`post_id`),
  KEY `FK1_author_user_id` (`user_id`) USING BTREE,
  CONSTRAINT `FK1_author_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK2_comments_post-id_posts_id` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table redditclone.comments: ~2 rows (approximately)
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` (`id`, `post_id`, `user_id`, `body`, `created`) VALUES
	(1, 1, 2, 'My first comment to the first post', '2021-07-11 14:03:46'),
	(2, 7, 16, 'some comment', '2021-07-20 23:47:34');
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;

-- Dumping structure for table redditclone.posts
CREATE TABLE IF NOT EXISTS `posts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
  `category` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
  `text` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `url` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `upvote_percentage` int(10) unsigned NOT NULL DEFAULT '0',
  `score` int(11) NOT NULL DEFAULT '0',
  `views` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK1_author_id_user_id` (`user_id`) USING BTREE,
  CONSTRAINT `FK1_posts_user-id_users_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table redditclone.posts: ~6 rows (approximately)
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` (`id`, `title`, `category`, `user_id`, `type`, `text`, `url`, `upvote_percentage`, `score`, `views`, `created_at`) VALUES
	(1, 'My first post', 'programming', 1, 'text', 'My first text to my first post', NULL, 0, 0, 1, '2021-07-10 15:32:39'),
	(2, 'My link post', 'music', 1, 'link', NULL, 'https://www.youtube.com/watch?v=h_D3VFfhvs4', 0, 0, 1, '2021-07-10 15:34:41'),
	(4, 'music post', 'music', 16, 'text', 'music post text', '', 0, 0, 1, '2021-07-20 00:03:10'),
	(5, 'yet another post', 'music', 16, 'text', 'poste text2', '', 0, 0, 0, '2021-07-20 00:31:53'),
	(7, 'post to delete', 'programming', 16, 'text', 'some post to delete', '', 0, 0, 14, '2021-07-20 23:45:46'),
	(8, 'news post', 'news', 16, 'text', 'some news text', '', 0, 0, 0, '2021-07-21 00:25:08');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;

-- Dumping structure for table redditclone.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table redditclone.users: ~2 rows (approximately)
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`id`, `name`, `password`) VALUES
	(1, 'Vladimir', ''),
	(2, 'Ivan', ''),
	(16, 'someuser', '$2a$14$rUUsNaoU1824OtcRaObtrOwhTNUauUCHew23pmsYyfguCQ0sPws.a');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;

-- Dumping structure for table redditclone.votes
CREATE TABLE IF NOT EXISTS `votes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `post_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `vote` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK1_user_user_id` (`user_id`) USING BTREE,
  KEY `FK2_votes_post-id_posts_id` (`post_id`),
  CONSTRAINT `FK1_user_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK2_votes_post-id_posts_id` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table redditclone.votes: ~0 rows (approximately)
/*!40000 ALTER TABLE `votes` DISABLE KEYS */;
/*!40000 ALTER TABLE `votes` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
