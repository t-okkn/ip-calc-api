SELECT
  `id`,
  `total`,
  STR_TO_DATE(`expire`, '%Y-%m-%dT%H:%i:%s') AS `expire`
FROM M_ID
WHERE `expire` < NOW();
