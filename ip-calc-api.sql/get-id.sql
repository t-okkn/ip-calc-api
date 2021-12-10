SELECT
  `id`,
  `total`,
  CAST(`is_end` AS UNSIGNED) AS `is_end`,
  `expire`
FROM T_ID
WHERE `id` = :id;
