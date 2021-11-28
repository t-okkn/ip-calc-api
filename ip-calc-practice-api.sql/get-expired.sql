SELECT
  `id`,
  `total`,
  CAST(`is_end` AS UNSIGNED) AS `is_end`,
  STR_TO_DATE(`expire`, '%Y-%m-%dT%H:%i:%s') AS `expire`
FROM T_ID
WHERE `expire` < NOW();
