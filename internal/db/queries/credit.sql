-- name: ListCreditByYearMonth :many
SELECT id
     , name
     , transacted_at
     , amount_in_microsgd 
     , created_at
FROM credit
WHERE transacted_at BETWEEN sqlc.arg(start_date) and sqlc.arg(end_date)
;

-- name: CreateCredit :one
INSERT INTO credit (
  name
  , transacted_at
  , amount_in_microsgd 
) VALUES ( $1, $2, $3 )
RETURNING id
;
