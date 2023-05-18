CREATE TABLE `candles_1d` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `open` double NOT NULL,
  `close` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `time` bigint NOT NULL,
  `volume` double NOT NULL,
  `color` varchar(255) NOT NULL
);

CREATE TABLE `candles_4h` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `open` double NOT NULL,
  `close` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `time` bigint NOT NULL,
  `volume` double NOT NULL,
  `color` varchar(255) NOT NULL
);

CREATE TABLE `candles_1h` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `open` double NOT NULL,
  `close` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `time` bigint NOT NULL,
  `volume` double NOT NULL,
  `color` varchar(255) NOT NULL
);

CREATE TABLE `candles_15m` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `open` double NOT NULL,
  `close` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `time` bigint NOT NULL,
  `volume` double NOT NULL,
  `color` varchar(255) NOT NULL
);

CREATE UNIQUE INDEX `candles_1d_index_0` ON `candles_1d` (`name`, `time`);
CREATE UNIQUE INDEX `candles_4h_index_0` ON `candles_4h` (`name`, `time`);
CREATE UNIQUE INDEX `candles_1h_index_0` ON `candles_1h` (`name`, `time`);
CREATE UNIQUE INDEX `candles_15m_index_0` ON `candles_15m` (`name`, `time`);
