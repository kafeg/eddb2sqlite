-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up
DROP TABLE IF EXISTS `allegiance`;
CREATE TABLE `allegiance` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `allegiance_id` INTEGER UNIQUE NOT NULL,
  `allegiance` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `bodies`;
CREATE TABLE `bodies` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `bodyId` int(11) NOT NULL,
  `name` mediumtext NOT NULL,
  `type` mediumtext NOT NULL,
  `subType` mediumtext NOT NULL,
  `offset` int(11) NOT NULL,
  `distanceToArrival` int(11) NOT NULL,
  `isMainStar` tinytext NOT NULL,
  `isScoopable` tinytext NOT NULL,
  `age` int(11) NOT NULL,
  `spectralClass` tinytext NOT NULL,
  `luminosity` tinytext NOT NULL,
  `absoluteMagnitude` float NOT NULL,
  `solarMasses` float NOT NULL,
  `solarRadius` float NOT NULL,
  `surfaceTemperature` float NOT NULL,
  `orbitalPeriod` float NOT NULL,
  `semiMajorAxis` float NOT NULL,
  `orbitalEccentricity` float NOT NULL,
  `orbitalInclination` float NOT NULL,
  `argOfPeriapsis` float NOT NULL,
  `rotationalPeriod` float NOT NULL,
  `rotationalPeriodTidallyLocked` tinytext NOT NULL,
  `axialTilt` float NOT NULL,
  `updateTime` int(11) NOT NULL,
  `systemId` int(11) NOT NULL
);

DROP TABLE IF EXISTS `commodities`;
CREATE TABLE `commodities` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `name` mediumtext NOT NULL,
  `category_id` int(11) NOT NULL,
  `average_price` int(11) NOT NULL,
  `is_rare` tinyint(1) NOT NULL,
  `max_buy_price` int(11) NOT NULL,
  `max_sell_price` int(11) NOT NULL,
  `min_buy_price` int(11) NOT NULL,
  `min_sell_price` int(11) NOT NULL,
  `buy_price_lower_average` int(11) NOT NULL,
  `sell_price_upper_average` int(11) NOT NULL,
  `is_non_marketable` int(11) NOT NULL,
  `ed_id` int(11) NOT NULL
);

DROP TABLE IF EXISTS `controlling_minor_faction`;
CREATE TABLE `controlling_minor_faction` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `controlling_minor_faction_id` int(11) UNIQUE NOT NULL,
  `controlling_minor_faction` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `factions`;
CREATE TABLE `factions` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `name` mediumtext NOT NULL,
  `updated_at` int(11) NOT NULL,
  `government_id` int(11) NOT NULL,
  `allegiance_id` int(11) NOT NULL,
  `home_system_id` int(11) NOT NULL,
  `is_player_faction` tinyint(1) NOT NULL
);

DROP TABLE IF EXISTS `government`;
CREATE TABLE `government` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `government_id` int(11) UNIQUE NOT NULL,
  `government` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `graph`;
CREATE TABLE `graph` (
  `systems_id` int(11) UNIQUE NOT NULL,
  `neighbors` longtext,
  `weights` longtext
);

DROP TABLE IF EXISTS `listings`;
CREATE TABLE `listings` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `station_id` int(11) NOT NULL,
  `commodity_id` int(11) NOT NULL,
  `supply` int(11) NOT NULL,
  `supply_bracket` int(11) NOT NULL,
  `buy_price` int(11) NOT NULL,
  `sell_price` int(11) NOT NULL,
  `demand` int(11) NOT NULL,
  `demand_bracket` int(11) NOT NULL,
  `collected_at` int(11) NOT NULL
);

DROP TABLE IF EXISTS `modules`;
CREATE TABLE `modules` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `group_id` int(11) NOT NULL,
  `class` int(11) NOT NULL,
  `rating` tinytext NOT NULL,
  `price` int(11) NOT NULL,
  `weapon_mode` mediumtext NOT NULL,
  `missile_type` mediumtext NOT NULL,
  `name` mediumtext NOT NULL,
  `belongs_to` mediumtext NOT NULL,
  `ed_id` int(11) NOT NULL,
  `ed_symbol` mediumtext NOT NULL,
  `ship` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `power_state`;
CREATE TABLE `power_state` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `power_state_id` int(11) UNIQUE NOT NULL,
  `power_state` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `primary_economy`;
CREATE TABLE `primary_economy` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `primary_economy_id` int(11) UNIQUE NOT NULL,
  `primary_economy` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `reserve_type`;
CREATE TABLE `reserve_type` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `reserve_type_id` int(11) UNIQUE NOT NULL,
  `reserve_type` mediumtext NOT NULL
);

DROP TABLE IF EXISTS `security`;
CREATE TABLE `security` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `security_id` int(11) UNIQUE NOT NULL,
  `security` mediumtext NOT NULL
);

--DROP TABLE IF EXISTS `state`;
--CREATE TABLE `state` (
--  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
--  `state_id` int(11) UNIQUE NOT NULL,
--  `state` mediumtext NOT NULL
--);

DROP TABLE IF EXISTS `stations`;
CREATE TABLE `stations` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `name` mediumtext NOT NULL,
  `system_id` int(11) NOT NULL,
  `updated_at` int(11) NOT NULL,
  `max_landing_pad_size` tinytext NOT NULL,
  `distance_to_star` int(11) NOT NULL,
  `government_id` int(11) NOT NULL,
  `allegiance_id` int(11) NOT NULL,
  `type_id` int(11) NOT NULL,
  `has_blackmarket` tinytext NOT NULL,
  `has_market` tinytext NOT NULL,
  `has_refuel` tinytext NOT NULL,
  `has_repair` tinytext NOT NULL,
  `has_rearm` tinytext NOT NULL,
  `has_outfitting` tinytext NOT NULL,
  `has_shipyard` tinytext NOT NULL,
  `has_docking` tinytext NOT NULL,
  `has_commodities` tinytext NOT NULL,
  `shipyard_updated_at` int(11) NOT NULL,
  `outfitting_updated_at` int(11) NOT NULL,
  `market_updated_at` int(11) NOT NULL,
  `is_planetary` tinytext NOT NULL,
  `settlement_size_id` int(11) NOT NULL,
  `settlement_size` int(11) NOT NULL,
  `settlement_security_id` int(11) NOT NULL,
  `body_id` int(11) NOT NULL,
  `controlling_minor_faction_id` int(11) NOT NULL
) ;

DROP TABLE IF EXISTS `systems`;
CREATE TABLE `systems` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `eddb_id` int(11) UNIQUE NOT NULL,
  `edsm_id` int(11) NOT NULL,
  `name` mediumtext NOT NULL,
  `x` float NOT NULL,
  `y` float NOT NULL,
  `z` float NOT NULL,
  `population` bigint(20) NOT NULL,
  `is_populated` tinyint(1) NOT NULL,
  `government_id` int(11) NOT NULL,
  `allegiance_id` int(11) NOT NULL,
  `security_id` int(11) NOT NULL,
  `primary_economy_id` int(11) NOT NULL,
  `power` mediumtext NOT NULL,
  `power_state_id` int(11) NOT NULL,
  `needs_permit` tinyint(1) NOT NULL,
  `updated_at` int(11) NOT NULL,
  `simbad_ref` mediumtext NOT NULL,
  `controlling_minor_faction_id` int(11) NOT NULL,
  `reserve_type_id` int(11) NOT NULL,
  `ed_system_address` bigint(20) NOT NULL
);

DROP TABLE IF EXISTS `type`;
CREATE TABLE `type` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `type_id` int(11) UNIQUE NOT NULL,
  `type` mediumtext NOT NULL
);

-- +migrate Down
DROP TABLE `type`;
DROP TABLE `systems`;
DROP TABLE `stations`;
DROP TABLE `state`;
DROP TABLE `security`;
DROP TABLE `reserve_type`;
DROP TABLE `primary_economy`;
DROP TABLE `power_state`;
DROP TABLE `modules`;
DROP TABLE `listings`;
DROP TABLE `graph`;
DROP TABLE `government`;
DROP TABLE `factions`;
DROP TABLE `controlling_minor_faction`;
DROP TABLE `commodities`;
DROP TABLE `bodies`;
DROP TABLE `allegiance`;
