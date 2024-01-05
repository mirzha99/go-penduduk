package Mmukim

type Mukim struct {
	Id              int    `gorm:"PrimaryKey column:id" json:"id"`
	Nama            string `gorm:"column:nama; type:varchar(30)" json:"nama"`
	Nama_Imum_Mukim string `gorm:"column:nama_imum_mukim; type:varchar(30)" json:"nama_imum_mukim"`
}
