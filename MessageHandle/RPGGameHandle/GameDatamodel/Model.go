package GameDatamodel

type Pet struct {
	ID                  int
	Name                string
	Type                string
	Skill               string
	HealthGrowthFactor  float32
	AtkGrowthFactor     float32
	DefenseGrowthFactor float32
	EnergyGrowthFactor  float32
	BaseHealth          int
	BaseAtk             int
	BaseDef             int
	BaseEnergy          int
}

type UserInfo struct {
	ID       int
	Name     string
	PetID    int
	PetLevel int
	Item     string
	Exp      int
	Skill    string
}
