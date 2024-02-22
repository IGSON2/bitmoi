CREATE TABLE `bidding_history` (
  `user_id` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `location` varchar(255) NOT NULL,
  `tx_hash` varchar(255) NOT NULL PRIMARY KEY,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `wmoi_transaction` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `from` varchar(50) NOT NULL DEFAULT "admin",
  `to` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `title` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY (`to`) REFERENCES `users`(`user_id`)
);