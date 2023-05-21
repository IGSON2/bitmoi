CREATE TABLE `users` (
  `user_id` varchar(255) PRIMARY KEY,
  `fullname` varchar(255) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_changed_at` timestamp NOT NULL DEFAULT (now()),
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `photo_url` varchar(255) NOT NULL
);

CREATE TABLE `score` (
  `score_id` varchar(255) PRIMARY KEY,
  `user_id` varchar(255) NOT NULL,
  `stage` tinyint NOT NULL,
  `pairname` varchar(255) NOT NULL,
  `entrytime` varchar(255) NOT NULL,
  `position` varchar(255) NOT NULL,
  `leverage` tinyint NOT NULL,
  `outtime` bigint NOT NULL,
  `entryprice` double NOT NULL,
  `endprice` double NOT NULL,
  `pnl` double NOT NULL,
  `roe` double NOT NULL
);

CREATE TABLE `ranking_board` (
  `user_id` varchar(255) PRIMARY KEY,
  `photo_url` varchar(255) NOT NULL,
  `score_id` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `final_balance` double,
  `comment` varchar(255)
);

CREATE INDEX `score_index_0` ON `score` (`user_id`);

CREATE INDEX `ranking_board_index_1` ON `ranking_board` (`score_id`);

CREATE INDEX `ranking_board_index_2` ON `ranking_board` (`photo_url`);

ALTER TABLE `score` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

ALTER TABLE `ranking_board` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

ALTER TABLE `ranking_board` ADD FOREIGN KEY (`photo_url`) REFERENCES `users` (`user_id`);

ALTER TABLE `score` ADD FOREIGN KEY (`score_id`) REFERENCES `ranking_board` (`score_id`);
