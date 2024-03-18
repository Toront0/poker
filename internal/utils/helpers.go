package utils

import (
    "crypto/rand"
)


func Filter[T any](vs []T, f func(T) bool) []T {
    filtered := make([]T, 0)
    for _, v := range vs {
        if f(v) {
            filtered = append(filtered, v)
        }
    }
    return filtered
}

func SliceIndex(limit int, predicate func(i int) bool) int {
    for i := 0; i < limit; i++ {
        if predicate(i) {
            return i
        }
    }
    return -1
}

func Map[T any, U any](s []T, fn func(i T) U) []U {
    res := []U{}

    for _, v := range s {
        res = append(res, fn(v))
    }
    return res
}

func SliceMultipleIndex(limit int, predicate func(i int) bool) []int {
    res := []int{}

    for i := 0; i < limit; i++ {
        if predicate(i) {
            res = append(res, i)
        }
    }
    return res
}

func Contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func RemoveIndex[T any](s []T, index int) []T {
    copy(s[index:], s[index+1:])
    return s[:len(s)-1]
}

const otpChars = "123456789"

func GenerateOTP(length int) (string, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return "", err
    }

    otpCharsLength := len(otpChars)
    for i := 0; i < length; i++ {
        buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
    }

    return string(buffer), nil
}

func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
