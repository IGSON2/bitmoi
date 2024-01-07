CREATE TABLE `verify_emails` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `user_id` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `secret_code` varchar(255) NOT NULL,
  `is_used` boolean NOT NULL DEFAULT false,
  `created_at` timestamp NOT NULL,
  `expired_at` timestamp NOT NULL
);

ALTER TABLE `verify_emails` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);