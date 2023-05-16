CREATE TABLE `candles_4h` (
  `name` varchar(255) NOT NULL,
  `open` double NOT NULL,
  `close` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `time` bigint NOT NULL,
  `volume` double NOT NULL,
  `color` varchar(255) NOT NULL
);

CREATE UNIQUE INDEX `candles_4h_index_0` ON `candles_4h` (`name`, `time`);
