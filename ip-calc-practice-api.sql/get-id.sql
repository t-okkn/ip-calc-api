SELECT
  `id`,
  `total`,
  CAST(`is_end` AS UNSIGNED) AS `is_end`,
  `expire`
FROM M_ID
WHERE `id` = :id;
