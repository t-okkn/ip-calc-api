SELECT
  COUNT(`tq`.`id`) AS `now`,
  `mid`.`total`
FROM T_QUESTION AS `tq`
LEFT JOIN M_ID AS `mid`
  ON `tq`.`id` = `mid`.`id`
WHERE `tq`.`id` = :id;
