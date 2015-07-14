// Package stringutil implements functions based on character indexes, not byte indexes.
package stringutil

import (
	"github.com/indidev/vocable-o/util/mathutil"
	"strings"
)

// merges to strings to one
func Join(a, b string) string {
	return strings.Join([]string{a, b}, "")
}

/*
	Does replacement for a map, should be used for special characters.
	Avoid using the method for stuff like a -> aa.
*/
func ReplaceMap(s string, replacements map[string]string) string {
	if replacements != nil {
		for old, new := range replacements {
			s = Replace(s, old, new, -1)
		}
	}

	return s
}

/*
	This function finds and replaces all substrings of old and replaces it with new
	/param n number of max replacements (-1 for all)
*/
func Replace(s, old, new string, n int) string {
	for i := 0; i < n || n < 0; i++ {
		index := Find(s, old)

		if index < 0 {
			break
		}

		s = RemoveAtN(s, index, Size(old))

		s = InsertAt(s, new, index)
	}

	return s
}

/*
Returns the number of characters in a string
*/
func Size(x string) int {
	i := 0
	for _ = range x {
		i++
	}
	return i
}

/*
	Splits a string after n characters and returns two substrings.
	e.q. SplitIndex("hallo", 2) returns ("ha", "llo")
*/
func SplitIndex(s string, index int) (string, string) {
	index = mathutil.MinInt(index, Size(s))
	index = mathutil.MaxInt(index, 0)

	s1 := make([]string, 0)
	s2 := make([]string, 0)

	i := 0

	for _, elem := range s {
		if i < index {
			s1 = append(s1, string(elem))
		} else {
			s2 = append(s2, string(elem))
		}

		i++
	}

	return strings.Join(s1, ""), strings.Join(s2, "")
}

/*
	Splites a string at the first occurrence of a seperator
*/
func SplitFirst(s, seperator string) (string, string) {
	l := strings.SplitN(s, seperator, 2)
	subs1 := s
	subs2 := ""

	if len(l) == 2 {
		subs1 = strings.TrimSpace(l[0])
		subs2 = strings.TrimSpace(l[1])
	}
	return subs1, subs2
}

/*
	Inserts a string into another string at a specific index.
	e.q. InsertAt("inrt", "se", 2) returns "insert"
*/
func InsertAt(s, filler string, index int) string {
	s1, s2 := SplitIndex(s, index)

	return Join(Join(s1, filler), s2)
}

/*
	Returns the character at a specific index
*/
func At(x string, index int) rune {

	if index >= Size(x) {
		return rune(0)
	}

	i := 0
	value := rune(0)
	for _, r := range x {
		if i == index {
			value = r
			break
		}
		i++
	}

	return value
}

/*
	removes n characters from the tail of the string
*/
func RemoveTail(x string, n int) string {
	size := Size(x)
	if n > size {
		return ""
	}

	for i := 0; i < n; i++ {
		x = RemoveAt(x, size-1)
	}

	return x
}

/*
	removes n characters from the given position in the string
	e.q. RemoveAtN("hallo", 2, 2) returns "hao"
*/
func RemoveAtN(x string, index, n int) string {
	for i := 0; i < n; i++ {
		x = RemoveAt(x, index)
	}

	return x
}

/*
	creates a substring from start to end (not including the end-index)
	e.q. Substring("hallo", 1, 4) returns "all"
*/
func Substring(s string, start, end int) string {

	start = mathutil.MaxInt(start, 0)
	end = mathutil.MinInt(end, Size(s))

	parts := make([]string, end-start+1)

	index := 0
	i := 0

	for _, elem := range s {
		if index >= start && index < end {
			parts[i] = string(elem)
			i++
		}
		index++
	}

	return strings.Join(parts, "")
}

/*
	removes a single character out of a string
*/
func RemoveAt(x string, index int) string {

	size := Size(x)
	if index >= size {
		return x
	}

	parts := make([]string, size-1)

	i := 0
	j := 0
	for _, elem := range x {
		if i != index {
			parts[j] = string(elem)
			j++
		}
		i++
	}

	return strings.Join(parts, "")
}

/*
	returns the first occurrence of a substring in s
*/
func Find(s, substring string) int {
	index := -1

	i := 0
	for _, char := range s {
		if char == At(substring, 0) {
			success := true
			for x := 1; x < Size(substring) && success; x++ {
				success = At(substring, x) == At(s, i+x)
			}
			if success {
				index = i
				break
			}
		}
		i++
	}

	return index
}

/*
	retruns all occurrences of a substring in s
*/
func FindAll(s, substring string) []int {
	indexes := make([]int, 0)

	i := 0
	for _, char := range s {
		if char == At(substring, 0) {
			success := true
			for x := 1; x < Size(substring) && success; x++ {
				success = At(substring, x) == At(s, i+x)
			}
			if success {
				indexes = append(indexes, i)
			}
		}
		i++
	}

	return indexes
}

/*
	changes the whole string to lower case and the first character of the string to upper case
*/
func UpperCaseOnlyFirst(s string) string {
	x1, x2 := SplitIndex(s, 1)

	x1 = strings.ToUpper(x1)
	x2 = strings.ToLower(x2)

	return Join(x1, x2)
}

/*
	sorts a slice of strings parallel accoarding to alphabetical order,
	does not distinghuis between lower and upper case
*/
func Mergesort(slice []string) []string {
	fin := make(chan []string)
	go mergehelp(slice, fin)
	return <-fin
}

func mergehelp(slice []string, fin chan []string) {
	length := len(slice)
	if length > 1 {
		slice1 := slice[:length/2]
		slice2 := slice[length/2:]

		fin1 := make(chan []string)
		fin2 := make(chan []string)
		go mergehelp(slice1, fin1)
		go mergehelp(slice2, fin2)

		slice1 = <-fin1
		slice2 = <-fin2

		slice = merge(slice1, slice2)
	}
	fin <- slice
}

/*
	merges to sorted string slices
*/
func merge(slice1, slice2 []string) []string {
	sorted := make([]string, 0)

	for len(slice1) > 0 && len(slice2) > 0 {
		if Compare(slice1[0], slice2[0]) {
			sorted = append(sorted, slice1[0])
			slice1 = slice1[1:]
		} else {
			sorted = append(sorted, slice2[0])
			slice2 = slice2[1:]
		}
	}

	for _, x := range slice1 {
		sorted = append(sorted, x)
	}
	for _, x := range slice2 {
		sorted = append(sorted, x)
	}

	return sorted
}

/*
	compares to strings and returns tue if x1 < x2
*/
func Compare(x1, x2 string) bool {
	x1 = strings.ToLower(x1)
	x2 = strings.ToLower(x2)

	x1LessX2 := true

	for i, x := range x1 {
		if x != At(x2, i) {
			x1LessX2 = x < At(x2, i)
			break
		}
	}
	return x1LessX2
}

/*
	finds a string in a string slice and returns its index
*/
func FindInSlice(wordlist *[]string, needle string) int {
	for index, elem := range *wordlist {
		if elem == needle {
			return index
		}
	}
	return -1
}

/*
	check wheter two strings are equal or not, include option to ignore puctuation
*/
func CheckEqual(x1, x2 string, ignorePunctuation bool) bool {
	if ignorePunctuation {
		x1 = removePunctiation(x1)
		x2 = removePunctiation(x2)
	}
	return strings.ToLower(x1) == strings.ToLower(x2)
}

/*
	returns the levenshtein distance between two strings,ignoring cases and option to ignore puctuation
*/
func Levenshtein(str1, str2 string, ignorePunctuation bool, maxLD int) int {
	if ignorePunctuation {
		str1 = removePunctiation(str1)
		str2 = removePunctiation(str2)
	}
	str1 = strings.ToLower(str1)
	str2 = strings.ToLower(str2)
	return levenshteinBT(strings.Split(str1, ""), strings.Split(str2, ""), 0, maxLD)
}

/*
	levenshtein back tracking helping function
*/
func levenshteinBT(str1, str2 []string, curLD, maxLD int) int {
	if (curLD > maxLD) {
		return curLD
	}
	if len(str1) == 0 || len(str2) == 0 {
		return curLD + mathutil.AbsInt(len(str1) - len(str2))
	}
	result := levenshteinBT(str1[1:], str2, curLD + 1, maxLD) //delete
	result = mathutil.MinInt(levenshteinBT(str1, str2[1:], curLD + 1, maxLD), result) //insert
	if (str1[0] == str2[0]) {
	result = mathutil.MinInt(levenshteinBT(str1[1:], str2[1:], curLD, maxLD), result) //equality
	} else {
	result = mathutil.MinInt(levenshteinBT(str1[1:], str2[1:], curLD + 1, maxLD), result) //substitution
	}
	return result
}

func removePunctiation(x string) string {
	replacements := make(map[string]string)
	replacements[","] = ""
	replacements["."] = ""
	replacements["?"] = ""
	replacements["!"] = ""
	replacements[":"] = ""
	x = ReplaceMap(x, replacements)
	x = Replace(x, "  ", " ", -1);
	x = strings.TrimSpace(x)

	return x
}
