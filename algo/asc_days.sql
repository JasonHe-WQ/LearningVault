SELECT MAX(consecutive_day)
FROM (
    SELECT COUNT(*) AS consecutive_day
    FROM (
        SELECT SUM(rise_mark) OVER (ORDER BY trade_date) AS days_no_gain
        FROM (
            SELECT trade_date,
                   CASE
                       WHEN closing_price > LAG(closing_price) OVER (ORDER BY trade_date) THEN 0
                       ELSE 1
                   END AS rise_mark
            FROM stock_price
        ) AS sub1
    ) AS sub2
    GROUP BY days_no_gain
) AS sub3;
