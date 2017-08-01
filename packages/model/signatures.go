package model

type Signatures struct {
	tableName  string
	Name       string `gorm:"primary_key;not null;size:255"`
	Value      string `gorm:not null;type:jsonb(PostgreSQL)`
	Conditions string `gorm:"not null"`
	RbID       int64  `gorm:"not null"`
}

func (s *Signatures) SetTableName(prefix string) {
	s.tableName = prefix + "_signatures"
}

func (s *Signatures) TableName() string {
	return s.tableName
}

func (s *Signatures) Get(name string) error {
	return DBConn.Where("name = ?", name).First(s).Error
}

func (s *Signatures) ExistsByName(name string) (bool, error) {
	query := DBConn.Where("name = ?", name).First(s)
	return !query.RecordNotFound(), query.Error
}

func (s *Signatures) GetAllOredered(prefix string) ([]Signatures, error) {
	var result []Signatures
	err := DBConn.Table(prefix + "_signatures").Order("name").Find(result).Error
	return result, err
}

func (s *Signatures) ToMap() map[string]string {
	var result map[string]string
	result["name"] = s.Name
	result["value"] = s.Value
	result["conditions"] = s.Conditions
	result["rb_id"] = string(s.RbID)
	return result
}

func CreateSignaturesStateTable(stateID string) error {
	return DBConn.Exec(`CREATE TABLE "` + stateID + `_signatures" (
				"name" varchar(100)  NOT NULL DEFAULT '',
				"value" jsonb,
				"conditions" text  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + stateID + `_signatures" ADD CONSTRAINT "` + stateID + `_signatures_pkey" PRIMARY KEY (name);
		    `).Error
}
