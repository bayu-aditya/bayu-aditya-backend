package util

func ArrayRemove[T comparable](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
