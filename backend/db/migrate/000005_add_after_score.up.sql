CREATE TABLE `prac_after_score` (
  `score_id` varchar(50) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `max_roe` double NOT NULL,
  `min_roe` double NOT NULL,
  `after_outtime` bigint NOT NULL,
  PRIMARY KEY (`score_id`, `user_id`),
  FOREIGN KEY (`score_id`, `user_id`) REFERENCES `prac_score`(`score_id`, `user_id`)
);