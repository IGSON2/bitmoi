CREATE TABLE `users` (
  `user_id` varchar(255) PRIMARY KEY,
  `oauth_uid` varchar(255),
  `full_name` varchar(255) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_changed_at` timestamp NOT NULL DEFAULT (now()),
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `photo_url` varchar(255)
);

CREATE TABLE `score` (
  `score_id` varchar(255) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `stage` tinyint NOT NULL,
  `pairname` varchar(255) NOT NULL,
  `entrytime` varchar(255) NOT NULL,
  `position` varchar(255) NOT NULL,
  `leverage` tinyint NOT NULL,
  `outtime` tinyint NOT NULL,
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
  `final_balance` double NOT NULL,
  `comment` varchar(255) NOT NULL
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`full_name`);

CREATE INDEX `score_index_1` ON `score` (`user_id`);

CREATE UNIQUE INDEX `score_index_2` ON `score` (`user_id`, `score_id`, `stage`);

CREATE INDEX `ranking_board_index_3` ON `ranking_board` (`score_id`);

CREATE UNIQUE INDEX `ranking_board_index_4` ON `ranking_board` (`user_id`, `score_id`);

ALTER TABLE `score` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

ALTER TABLE `ranking_board` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);
