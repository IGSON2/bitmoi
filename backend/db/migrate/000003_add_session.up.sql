CREATE TABLE `sessions` (
  `session_id` varchar(36) PRIMARY KEY DEFAULT (UUID()),
  `user_id` varchar(255) NOT NULL,
  `refresh_token` varchar(300) NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(255) NOT NULL,
  `is_blocked` boolean NOT NULL DEFAULT false,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);