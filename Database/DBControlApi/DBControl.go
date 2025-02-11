package DBControlApi

import "gorm.io/gorm"

// 用户服务结构体
type DBcontrol struct {
	DB *gorm.DB
}

func (bcontrol *DBcontrol) Where(tablename string, model interface{}, query interface{}, args ...interface{}) (interface{}, error) {
	// 使用 GORM 的 Where 方法执行查询
	result := bcontrol.DB.Table(tablename).Where(query, args...).Find(model)
	if result.Error != nil {
		return nil, result.Error
	}
	// 返回查询到的结果，保持原始的 model 类型
	return model, nil
}

// Get 查询数据库中所有记录
func (bcontrol *DBcontrol) Get(Model interface{}, tablename string) (interface{}, error) {
	result := bcontrol.DB.Table(tablename).Find(Model)
	if result.Error != nil {
		return nil, result.Error
	}
	return Model, nil
}

// Create 插入记录
func (bcontrol *DBcontrol) Create(model interface{}, tablename string) error {
	if result := bcontrol.DB.Table(tablename).Create(model); result.Error != nil {
		return result.Error
	}
	return nil
}

// Update 更新记录
func (bcontrol *DBcontrol) Update(updates interface{}, tablename string) error {
	if result := bcontrol.DB.Table(tablename).Updates(updates); result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete 删除记录
func (bcontrol *DBcontrol) Delete(model interface{}, tablename string) error {
	if result := bcontrol.DB.Table(tablename).Delete(model); result.Error != nil {
		return result.Error
	}
	return nil
}
