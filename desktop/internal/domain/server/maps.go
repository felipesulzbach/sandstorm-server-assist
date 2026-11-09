package server

// Map é um mapa do jogo base: o nome que vai na URL de início e o nome usado nos cenários.
type Map struct {
	Label string `json:"label"`
	Map   string `json:"map"`
	// Scenario é o prefixo dos cenários (Scenario_<Scenario>_Checkpoint_Security).
	Scenario string `json:"scenario"`
}

// Maps são os mapas do jogo base. O mapa inicial precisa ser um deles: mapa de mod só entra pelo
// -ModDownloadTravelTo (ADR-003). Os cenários vêm do MapCycle.txt do servidor; os nomes de mapa, do guia de
// administração do jogo.
var Maps = []Map{
	{"Bab", "Bab", "Bab"},
	{"Citadel", "Citadel", "Citadel"},
	{"Crossing", "Canyon", "Crossing"},
	{"Farmhouse", "Farmhouse", "Farmhouse"},
	{"Forest", "Forest", "Forest"},
	{"Gap", "Gap", "Gap"},
	{"Hideout", "Town", "Hideout"},
	{"Hillside", "Sinjar", "Hillside"},
	{"Last Light", "LastLight", "LastLight"},
	{"Ministry", "Ministry", "Ministry"},
	{"Outskirts", "Compound", "Outskirts"},
	{"Power Plant", "PowerPlant", "PowerPlant"},
	{"Precinct", "Precinct", "Precinct"},
	{"Prison", "Prison", "Prison"},
	{"Refinery", "Oilfield", "Refinery"},
	{"Summit", "Mountain", "Summit"},
	{"Tell", "Tell", "Tell"},
	{"Tideway", "Buhriz", "Tideway"},
}

// Mode é um modo de cenário (sufixo do nome do cenário). O texto da tela vem das traduções (maps.mode.<id>).
type Mode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Modes são os modos de cenário mais usados.
var Modes = []Mode{
	{"Checkpoint_Security", "Checkpoint (Security)"},
	{"Checkpoint_Insurgents", "Checkpoint (Insurgents)"},
}

// IsKnownMap diz se o nome é de um mapa do jogo base.
func IsKnownMap(name string) bool {
	for _, m := range Maps {
		if m.Map == name {
			return true
		}
	}
	return false
}
