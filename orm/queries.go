package orm

import "fmt"

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
	return "SELECT SUM(amount_converted) AS total_spending, COUNT(*) FROM milestracker_db.expense_table WHERE user_id = ? AND card_id = ? AND transaction_timestamp >= ? AND transaction_timestamp <= ? AND is_cancel = 0"
}

func Sql6() string {
	return "SELECT * FROM milestracker_db.expense_table WHERE user_id = ? AND is_cancel = 0 AND transaction_timestamp >= ? AND transaction_timestamp <= ? ORDER BY transaction_timestamp DESC"
}

func Sql7() string {
	return "SELECT * FROM milestracker_db.expense_table WHERE user_id = ? AND trx_id = ?"
}

func Sql8() string {
	return "SELECT * FROM milestracker_db.user_card_table WHERE user_id = ? AND card_id = ? AND card_nickname = ? AND card_expiry = ?"
}

func Sql9(userId int64) string {
	return fmt.Sprintf("SELECT user_card.*, card_info.card_name,card_info.short_card_name, card_info.card_type, card_info.card_image, card_info.card_issuer FROM milestracker_db.user_card_table AS user_card LEFT JOIN milestracker_db.card_table AS card_info ON user_card.card_id = card_info.card_id WHERE user_card.user_id = %v", userId)
}
