SELECT
  `id`,
  `question_number`,
  `source`,
  `cidr_bits`,
  CAST(`is_cidr` AS UNSIGNED) AS `is_cidr`,
  `correct_nw`,
  `answer_nw`,
  `correct_bc`,
  `answer_bc`,
  `elapsed`,
  `created`,
  `updated`
FROM T_QUESTION
WHERE `id` IN (:ids)
FOR UPDATE;
