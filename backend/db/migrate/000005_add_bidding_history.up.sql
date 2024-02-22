CREATE TABLE `bidding_history` (
  `user_id` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `location` varchar(255) NOT NULL,
  `tx_hash` varchar(255) NOT NULL PRIMARY KEY,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `recommend_history` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `from_user` varchar(255) NOT NULL,
  `to_user` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY (`from_user`) REFERENCES `users`(`user_id`),
  FOREIGN KEY (`to_user`) REFERENCES `users`(`user_id`)
);

CREATE TABLE `wmoi_transaction` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `from_user` varchar(50) NOT NULL DEFAULT "admin",
  `to_user` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `title` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY (`to_user`) REFERENCES `users`(`user_id`)
);

CREATE UNIQUE INDEX `recommend_history_1` ON `recommend_history` (`from_user`, `to_user`);
