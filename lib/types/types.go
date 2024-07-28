package types

type Song struct {
	Label              []string
	SongFolderName     string
	AudioFileName      string
	BackgroundFileName string
	Mania              bool
	Taiko              bool
	Duration           int
	SouceNumber        int
}

type Source struct {
	Name          string
	NumberOfSongs int
}
