CREATE DATABASE  IF NOT EXISTS `goph-chat-db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `goph-chat-db`;
-- MySQL dump 10.13  Distrib 8.0.36, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: goph-chat-db
-- ------------------------------------------------------
-- Server version	8.0.37

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL,
  `room_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_messages_room_id` (`room_id`),
  KEY `idx_messages_user_id` (`user_id`),
  CONSTRAINT `fk_messages_room` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_messages_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (8,'Mark the beginning of the messagequeue',6,23,'2025-05-06 12:59:30','2025-05-06 12:59:30'),(9,'',6,39,'2025-05-06 13:00:05','2025-05-06 13:00:05'),(10,'I just create a monster',6,23,'2025-05-06 13:00:19','2025-05-06 13:00:19'),(11,'',6,39,'2025-05-06 13:00:31','2025-05-06 13:00:31'),(12,'',2,23,'2025-05-06 13:00:43','2025-05-06 13:00:43'),(13,'I just create a monster',2,23,'2025-05-06 13:44:24','2025-05-06 13:44:24'),(14,'Is that so ?',2,24,'2025-05-06 13:44:38','2025-05-06 13:44:38'),(15,'What\'s up',6,23,'2025-05-06 13:44:55','2025-05-06 13:44:55'),(16,'IM FINE',6,39,'2025-05-06 13:45:05','2025-05-06 13:45:05'),(17,'Can i ask who are you ?',2,24,'2025-05-16 07:39:59','2025-05-16 07:39:59'),(18,'I am user id 23',2,23,'2025-05-16 07:40:24','2025-05-16 07:40:24'),(19,'hello im user 23 from room 6',6,23,'2025-05-16 07:40:39','2025-05-16 07:40:39'),(20,'Great, you can connect to multiple websocket connection',6,39,'2025-05-16 07:41:10','2025-05-16 07:41:10'),(21,'hello im user 23 from room 6 1',6,23,'2025-05-16 07:41:47','2025-05-16 07:41:47'),(22,'hello im user 23 from room 6 2',6,23,'2025-05-16 07:41:50','2025-05-16 07:41:50'),(23,'hello im user 23 from room 6 3',6,23,'2025-05-16 07:41:53','2025-05-16 07:41:53'),(24,'hello im user 23 from room 6 4',6,23,'2025-05-16 07:41:55','2025-05-16 07:41:55'),(25,'hello im user 23 from room 6 4',6,23,'2025-05-16 07:41:56','2025-05-16 07:41:56'),(26,'hello im user 23 from room 6 4',6,23,'2025-05-16 07:41:57','2025-05-16 07:41:57'),(27,'hello im user 23 from room 6 4',6,23,'2025-05-16 07:41:57','2025-05-16 07:41:57'),(28,'hello im user 23 from room 6 4',6,23,'2025-05-16 07:41:58','2025-05-16 07:41:58'),(29,'hello im user 23 from room 6 4',6,23,'2025-05-16 07:41:58','2025-05-16 07:41:58'),(31,'hello im user 23 from room 6 4',6,23,'2025-05-19 14:09:59','2025-05-19 14:09:59'),(32,'Testout new HSET',6,39,'2025-05-19 14:13:51','2025-05-19 14:13:51'),(33,'I have change the message content',6,23,'2025-05-19 14:17:19','2025-05-20 02:04:53'),(34,'I now change the message ID 34',6,23,'2025-05-20 02:02:09','2025-05-20 02:06:43'),(35,'id 35 is what i chose',6,23,'2025-05-20 02:09:58','2025-05-20 02:10:40'),(36,'Removing Message redundant, how does everyone feels?',6,23,'2025-05-20 08:40:57','2025-05-20 08:40:57'),(37,'Its getting buggy, but im gonna fix it',6,39,'2025-05-20 08:45:59','2025-05-20 08:45:59'),(38,'tEST OUT WHERE IS THE BUG',6,39,'2025-05-20 08:49:20','2025-05-20 08:49:20'),(39,'tEST OUT WHERE IS THE BUG',6,39,'2025-05-20 08:49:38','2025-05-20 08:49:38'),(40,'tEST OUT WHERE IS THE BUG',6,23,'2025-05-20 08:50:18','2025-05-20 08:50:18'),(41,'Windsurf is trying its best',6,39,'2025-05-20 08:54:31','2025-05-20 08:54:31'),(42,'Windsurf is being eaten alive by GPT',6,39,'2025-05-20 09:03:43','2025-05-20 09:03:43'),(43,'Why id stays the same ?',6,39,'2025-05-20 09:05:02','2025-05-20 09:05:02'),(44,'Is this fix ?',6,39,'2025-05-20 09:28:33','2025-05-20 09:28:33'),(45,'Is this fix ?',6,39,'2025-05-20 09:29:15','2025-05-20 09:29:15'),(46,'Try out pubsub with update ?',2,23,'2025-05-20 09:30:59','2025-05-20 09:33:25'),(47,'Update even tho no one has said anything',2,23,'2025-05-20 09:36:56','2025-05-20 09:38:56'),(48,'This should be fixed',2,23,'2025-05-20 09:43:47','2025-05-20 09:44:18'),(49,'You just not the person to send it originally tho',2,23,'2025-05-20 09:45:43','2025-05-20 09:48:51'),(50,'Yeah maybe ? But why my message is sending to mine alone ?',2,24,'2025-05-20 09:46:01','2025-05-20 09:46:01');
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rooms`
--

DROP TABLE IF EXISTS `rooms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rooms` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(45) DEFAULT NULL,
  `description` text,
  `user_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_user` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rooms`
--

LOCK TABLES `rooms` WRITE;
/*!40000 ALTER TABLE `rooms` DISABLE KEYS */;
INSERT INTO `rooms` VALUES (1,'Room Chat No.1','In the beninging',39,'2025-04-19 04:21:09','2025-04-19 04:21:09',NULL),(2,'No.2 Room (Websocket)','In the beninging',39,'2025-04-19 08:49:50','2025-04-19 08:49:50',NULL),(3,'No.3 Room (ws-fixing)','In the beninging',39,'2025-04-19 16:02:44','2025-04-19 16:02:44',NULL),(4,'Room 4','In the beninging',21,'2025-04-25 05:08:44','2025-04-25 05:08:44',NULL),(5,'Room 5: Fixing seed','In the beninging',21,'2025-04-25 05:27:19','2025-04-25 05:27:19',NULL),(6,'Room 6: Testing id with seed','In the beninging',21,'2025-04-25 05:28:18','2025-04-25 05:28:18',NULL);
/*!40000 ALTER TABLE `rooms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(100) DEFAULT NULL,
  `role` enum('user','admin','shipper','mod') DEFAULT 'user',
  `status` int DEFAULT '1',
  `password` varchar(100) DEFAULT NULL,
  `salt` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) DEFAULT NULL,
  `first_name` varchar(100) DEFAULT NULL,
  `phone` varchar(40) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (15,'fir2s2t111@gmail.com','user',1,'','Qzelaixtq1rLYl49clkB1pXJzZsZJHNhzsvC-zQ-QV70R8-4wsyJ2qr2AodKUHe_11w','Phuc','Hoang','','2025-04-14 16:08:41','2025-04-14 16:08:41'),(16,'fixpass@gmail.com','user',1,'$2a$12$Pk0fR03AOoI3b9dVnwblruM0XeBGSzPu86F5sD1bePpD1kHkxtXtu','0F8gMRhw9rBC3Lslk3dKE5H1xNShzlYZ1p-9DAoqzn-VsnBbf0DN6A','Phuc','Hoang','','2025-04-14 16:14:30','2025-04-14 16:14:30'),(17,'first1@gmail.com','user',1,'$2a$12$Qym0lM0S.DPgZuKYjfi5x.HU8tPTgcTpa9n4sDWNGS0VyLVfqmN86','jp8y_MMbKCJVv3Cvi7mu3CBEGP3JxZiNOtjzSqQ0qpNA3sCB78AKRQ','Phuc','Hoang','','2025-04-15 05:43:24','2025-04-15 05:43:24'),(21,'first123@gmail.com','user',1,'$2a$12$wC9HJd97750kR113kzx9q.5r4M.A/atXqw8tXsDS3T9Y3oKyn69Z.','LEq3Wp6TlGF38EUzNoIfOsDZBrtCKCpapi3YYwwZFTVEkgO6zt58SQ','Phuc','Hoang','','2025-04-15 08:20:30','2025-04-15 08:20:30'),(22,'first1233@gmail.com','user',1,'$2a$12$l6tndyRMFZffTCy5T2awV.8NnKlAnPJ5mYOe8erDPwzz/jXNWx53m','GEgP55WWeZLR85p1yPV5OzyKRbCRMrxaxHNBq8ceII4sdE0HlZ_BRA','Phuc','Hoang','','2025-04-15 08:29:18','2025-04-15 08:29:18'),(23,'first111233@gmail.com','user',1,'$2a$12$uuUT6hPyL/EU8WlxG2IdFOvwPKj.Qv81TJVEOO7NLAqSe6fQI5Lle','qZCI9S61z5o0Y7Zq2abURZy_M-MuuI5vLoZwRN-IMn9U7HwTjAPGbg','Phuc','Hoang','','2025-04-15 08:33:07','2025-04-15 08:33:07'),(24,'fixnumber2@gmail.com','user',1,'$2a$12$wBPdJrrqsq9l.WAxKjfxI.tmBv5njPkXuT3RULJyPg1DVHUZGzhFy','ZI9gGrBIoSxdM4s0VE6znkexE1TWbo3egDTkqrodAwJVd7SLWAL7ag','Phuc','Hoang','','2025-04-15 08:51:23','2025-04-15 08:51:23'),(39,'newone@gmail.com','user',1,'$2a$12$Ej04L4vOVpLc5.aIfxeHTu296fpLQ31NQun1RuzuXIqDPXNfO/Ed2','lXODVwYhMlidP-sqJIqmg2D7gvLXKq4h6lvV0yqdsLmsMXIlEYW_3A','Phuc','Hoang','','2025-04-19 04:18:10','2025-04-19 04:18:10');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-05-24 14:55:59
