package tools

func Insert[T any](slice []T, target T, index int) []T {
	if len(slice) == 0 {
		slice = append(slice, target)
		return slice
	}
  if len(slice) == index {
    slice = append(slice, target)
    return slice
  }
	slice = append(slice[:index+1], slice[index:]...)
	slice[index] = target
	return slice
}
