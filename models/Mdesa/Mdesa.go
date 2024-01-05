package Mdesa

type Desa struct {
	Id               int    `gorm:"PrimaryKey; column:id" json:"id"`
	Nama             string `gorm:"column:nama; type:varchar(30)" json:"nama"`
	Nama_Kepala_Desa string `gorm:"column:nama_kepala_desa; type:varchar(30)" json:"nama_kepala_desa"`
	Id_mukim         int    `gorm:"column:id_mukim; type:bigint(20)" json:"id_mukim"`
}
type DesaResult struct {
	Id               int    `json:"id"`
	Nama             string `json:"nama"`
	Nama_Kepala_Desa string `json:"nama_kepala_desa"`
	Id_mukim         int    `json:"id_mukim"`
	Nama_Mukim       string `json:"nama_mukim"`
	Nama_Imum_Mukim  string `json:"nama_imum_mukim"`
}
