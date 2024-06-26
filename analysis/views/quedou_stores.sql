CREATE TABLE quedou_stores AS
SELECT
	NAME,
	id,
	rooms,
	address
FROM
	(
		SELECT
			store_name AS NAME,
			store_id AS id,
			store_address AS address,
			COUNT(*) OVER (PARTITION BY store_id) AS rooms,
			ROW_NUMBER() OVER (PARTITION BY store_id) AS rn
		FROM
			quedou_data
		WHERE
			YEAR = '2024'
			AND MONTH = '6'
			AND DAY = '18'
			AND HOUR = '13'
	) AS foo
WHERE
	rn = 1;