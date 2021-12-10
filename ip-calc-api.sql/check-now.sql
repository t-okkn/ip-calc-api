SELECT
  COUNT(`tq`.`id`) AS `now`,
  `tid`.`total`
FROM T_QUESTION AS `tq`
LEFT JOIN T_ID AS `tid`
  ON `tq`.`id` = `tid`.`id`
WHERE `tq`.`id` = :id;
