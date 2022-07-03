package utils

func Remove[T comparable](l []T, item T) []T {
    for i, other := range l {
        if other == item {
            return append(l[:i], l[i+1:]...)
        }
    }
    return l
}

func GetBigger[T comparable](a T, b T) T {
    if (a > b) {
        return a
    }
    return b
}

func GetSmaller[T comparable](a T, b T) T {
    if (a > b) {
        return b
    }
    return a
}

func StrArrContains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}