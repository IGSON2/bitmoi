CREATE TABLE `users` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `user_id` varchar(255) NOT NULL UNIQUE,
  `oauth_uid` varchar(255) UNIQUE,
  `nickname` varchar(50) NOT NULL UNIQUE,
  `hashed_password` varchar(255) UNIQUE,
  `email` varchar(255) NOT NULL UNIQUE,
  `metamask_address` varchar(100) UNIQUE,
  `photo_url` varchar(100),
  `prac_balance` double NOT NULL DEFAULT 0,
  `comp_balance` double NOT NULL DEFAULT 0,
  `wmoi_balance` double NOT NULL DEFAULT 0,
  `recommender_code` varchar(50) NOT NULL UNIQUE,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `last_accessed_at` timestamp,
  `password_changed_at` timestamp NOT NULL DEFAULT (now()),
  `address_changed_at` timestamp
);

CREATE TABLE `prac_score` (
  `score_id` varchar(50) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `stage` tinyint NOT NULL,
  `pairname` varchar(50) NOT NULL,
  `entrytime` varchar(50) NOT NULL,
  `position` varchar(20) NOT NULL,
  `leverage` tinyint NOT NULL,
  `outtime` varchar(50) NOT NULL,
  `entryprice` double NOT NULL,
  `quantity` double NOT NULL,
  `endprice` double NOT NULL,
  `pnl` double NOT NULL,
  `roe` double NOT NULL,
  `settled_at` timestamp,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`score_id`, `user_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE
);

CREATE TABLE `comp_score` (
  `score_id` varchar(50) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `stage` tinyint NOT NULL,
  `pairname` varchar(50) NOT NULL,
  `entrytime` varchar(50) NOT NULL,
  `position` varchar(20) NOT NULL,
  `leverage` tinyint NOT NULL,
  `outtime` varchar(50) NOT NULL,
  `entryprice` double NOT NULL,
  `quantity` double NOT NULL,
  `endprice` double NOT NULL,
  `pnl` double NOT NULL,
  `roe` double NOT NULL,
  `settled_at` timestamp,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`score_id`, `user_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE
);

CREATE TABLE `ranking_board` (
  `user_id` varchar(255) PRIMARY KEY,
  `score_id` varchar(255) NOT NULL,
  `final_balance` double NOT NULL,
  `comment` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`)
);

CREATE TABLE `used_token` (
  `score_id` varchar(255) PRIMARY KEY,
  `user_id` varchar(255) NOT NULL,
  `metamask_address` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`)
);

CREATE UNIQUE INDEX `ranking_board_index_2` ON `ranking_board` (`user_id`, `score_id`);

CREATE UNIQUE INDEX `used_token_index_1` ON `used_token` (`user_id`, `score_id`);