CREATE TABLE `bidding_history` (
  `tx_hash` varchar(255) NOT NULL PRIMARY KEY,
  `user_id` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `location` varchar(255) NOT NULL,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`)
);

CREATE TABLE `recommend_history` (
  `recommender` varchar(255) NOT NULL,
  `new_member` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`new_member`),
  FOREIGN KEY (`recommender`) REFERENCES `users`(`user_id`),
  FOREIGN KEY (`new_member`) REFERENCES `users`(`user_id`)
);

CREATE TABLE `wmoi_minting_history` (
  `to_user` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `title` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `method` VARCHAR(2) NOT NULL DEFAULT '자동',
  `giver` VARCHAR(20) NOT NULL DEFAULT '시스템',
  CONSTRAINT wmoi_chk_method CHECK(method IN ('자동', '수동')),
  PRIMARY KEY (`to_user`, `created_at`),
  FOREIGN KEY (`to_user`) REFERENCES `users`(`user_id`)
);

CREATE TABLE `accumulation_history` (
  `to_user` varchar(255) NOT NULL,
  `amount` DOUBLE NOT NULL DEFAULT 0,
  `title` VARCHAR(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `method` VARCHAR(2) NOT NULL DEFAULT '자동',
  `giver` VARCHAR(20) NOT NULL DEFAULT '시스템',
  CONSTRAINT accu_chk_method CHECK(method IN ('자동', '수동')),
  PRIMARY KEY (`to_user`, `created_at`),
  FOREIGN KEY (`to_user`) REFERENCES `users`(`user_id`)
);
