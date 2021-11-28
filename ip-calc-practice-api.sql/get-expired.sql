SELECT
  `id`,
  `total`,
  `is_end`,
  STR_TO_DATE(`expire`, '%Y-%m-%dT%H:%i:%s') AS `expire`
FROM T_ID
WHERE `expire` < NOW();
