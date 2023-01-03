package orm

func Sql1() string {
	return "SELECT * FROM milestracker_db.card_table WHERE card_name = ?"
}

func Sql2() string {
	return "SELECT * FROM milestracker_db.card_table WHERE card_id = ?"
}

func Sql3() string {
	return "SELECT * FROM milestracker_db.user_table WHERE user_id = ?"
}

func Sql4() string {
	return "SELECT * FROM milestracker_db.user_card_table WHERE user_id = ? AND card_id = ?"
}

func Sql5() string {
	return "SELECT SUM(amount_converted) AS total_spending, COUNT(*) FROM milestracker_db.expense_table " +
		"WHERE user_id = ? AND card_id = ? AND transaction_timestamp >= ? AND transaction_timestamp <= ? AND is_cancel = 0"
}
