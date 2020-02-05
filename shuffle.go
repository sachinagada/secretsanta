package secretsanta

// Shuffle interface is used to shuffle the participants to assign the secret santa
type Shuffle interface {
	Shuffle(emails []string) []string
	Type() string
}
