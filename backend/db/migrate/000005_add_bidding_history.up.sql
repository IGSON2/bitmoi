CREATE TABLE `bidding_history` (
  `user_id` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `location` varchar(255) NOT NULL,
  `tx_hash` varchar(255) NOT NULL PRIMARY KEY,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE `bidding_history` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);