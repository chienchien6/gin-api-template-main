package service

import (
	"RCSP/global"
	"RCSP/model"
)

type UserService struct{}

func (s *UserService) GetUserByID(id string) (model.User, error) {
	var user model.User
	if err := global.GvaMysqlClient.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

/*
&user 是一個指向 user 變數的指針，這個變數通常是用來存儲查詢結果的 User 結構體實例。

*/

func (s *UserService) Create(user *model.User) (*model.User, error) {
	if err := global.GvaMysqlClient.Create(user).Error; err != nil {
		global.GvaLogger.Sugar().Errorf("Error creating user: %v", err)
		return nil, err
	}
	global.GvaLogger.Sugar().Infof("User created successfully: %#v", user)
	return user, nil
}

/*使用指針的原因是讓Create, Save(update)方法能夠直接修改 user 變數的內容，以便在查詢成功後將數據填充到這個變數中。*/

func (s *UserService) Update(user *model.User) (*model.User, error) {
	if err := global.GvaMysqlClient.Save(user).Error; err != nil {
		global.GvaLogger.Sugar().Errorf("Error updating user: %v", err)
		return nil, err
	}
	global.GvaLogger.Sugar().Infof("User updated successfully: %#v", user)
	return user, nil
}

func (s *UserService) Delete(id string) error {
	if err := global.GvaMysqlClient.Delete(&model.User{}, id).Error; err != nil {
		global.GvaLogger.Sugar().Errorf("Error deleting user: %v", err)
		return err
	}
	global.GvaLogger.Sugar().Infof("User deleted successfully: %s", id)
	return nil
}

/*&model.User{} 創建了一個 model.User 結構的指針。這個結構通常代表用戶的數據模型，並且在這裡用作類型的占位符。
在 ORM 中，你需要提供一個結構體的類型，以告訴 ORM 你要刪除哪種類型的數據。這裡使用空的結構體是因為我們只需要告訴 ORM 這是 User 模型，而不需要提供具體的用戶數據。*/
