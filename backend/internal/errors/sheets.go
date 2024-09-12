package errors

func DeleteLoyaltyDataError(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка удаления данных",
		"en": "Loyalty data delete failed",
	})

}
func RefreshLoyaltyDataError(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка обновления данных",
		"en": "Refresh data failed",
	})

}
func DeactivateUsersError(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка деактивации пользователей",
		"en": "Deactivate users error",
	})

}
func RefreshUsersError(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка обновления данных пользователй",
		"en": "Refresh users error",
	})

}
func RefreshExtendDataError(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка обновления дополнительных данных",
		"en": "Extend data refresh failed",
	})

}
func DeleteExtendLoyaltyDataError(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка удаления дополнительных данных",
		"en": "Delete extend loyalty data error",
	})

}
