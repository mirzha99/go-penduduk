package Mkepalakularga

type KepalaKeluarga struct {
	Id         int    `gorm:"PrimaryKey" column:"id" json:"id"`
	Nama       string `gorm:"type:varchar(30);column:nama" json:"nama" form:"nama"`
	Nik        int    `gorm:"type:int; column:nik" json:"nik" form:"nik"`
	Id_Desa    int    `gorm:"type:int;column:id_desa" json:"id_desa" form:"id_desa"`
	Created_at string `gorm:"type:varchar(30);column:created_at" json:"created_at"`
	Change_at  string `gorm:"type:varchar(30);column:change_at" json:"change_at"`
	Gambar     string `gorm:"type:varchar(250);column:gambar" json:"gambar"`
}
type KepalaKeluargaResult struct {
	Id               int    `json:"id"`
	Nama             string `json:"nama"`
	Nik              int    `json:"nik"`
	Id_Desa          int    `json:"id_desa"`
	Nama_Desa        string `json:"nama_desa"`
	Nama_Kepala_Desa string `json:"nama_kepala_desa"`
	Id_Mukim         string `json:"id_mukim"`
	Nama_Mukim       string `json:"nama_mukim"`
	Nama_Imum_Mukim  string `json:"nama_imum_mukim"`
	Created_at       string `json:"created_at"`
	Change_at        string `json:"change_at"`
	Gambar           string `gorm:"type:varchar(250);column:gambar" json:"gambar"`
}
