package util

func ArrayRemove[T comparable](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func ArrayAppendHead[T comparable](slice []T, newElement T) []T {
	return append([]T{newElement}, slice...)
}
