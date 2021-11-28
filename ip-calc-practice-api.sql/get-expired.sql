SELECT
  `id`,
  `total`,
  CAST(`is_end` AS UNSIGNED) AS `is_end`,
  `expire`
FROM T_ID
WHERE STR_TO_DATE(`expire`, '%Y-%m-%dT%H:%i:%s') < NOW();
