package types

type Song struct {
	label              []string
	songFolderName     string
	audioFileName      string
	backgroundFileName string
	mania              bool
	taiko              bool
	duration           int
	souceNumber        int
}

type Source struct {
	name          string
	numberOfSongs int
}
