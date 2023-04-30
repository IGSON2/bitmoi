CREATE TABLE `Candles` (
  `name` varchar(255) NOT NULL,
  `open` double NOT NULL,
  `close` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `time` integer NOT NULL,
  `volume` double NOT NULL,
  `color` varchar(255) NOT NULL,
  `interval` varchar(255) NOT NULL
);

CREATE UNIQUE INDEX `Candles_index_0` ON `Candles` (`Name`, `Time`, `Interval`);
