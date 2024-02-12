CREATE TABLE `bidding_history` (
  `user_id` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `location` varchar(255) NOT NULL,
  `tx_hash` varchar(255) NOT NULL PRIMARY KEY,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

/* TODO: update transfers table */
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from" bigint NOT NULL,
  "to" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "title" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transfers" ADD FOREIGN KEY ("to") REFERENCES "users" ("id");