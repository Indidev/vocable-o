package lang

import (
	"github.com/indidev/vocable-o/util/stringutil"
	"io/ioutil"
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type word struct {
	Name           string
	Translation    string
	successCounter int
	pocketIndex    int
}

var AvailableLanguages []string

var successPerPocket = [5]int{2, 3, 4, 5, 1000}

const wordfile = "wordlist.txt"
const dirname = "languages/"

var curLang string
var words [5][]word
var modified bool = false

type oobError int

func (x oobError) Error() string {
	return strings.Join([]string{"Index", strconv.Itoa(int(x)), "is out of Arraybounds."}, " ")
}

func CurLang() string {
	return curLang
}

func Init() {
	content, err := ioutil.ReadDir(dirname)

	AvailableLanguages = make([]string, 0)

	if err == nil {
		for _, file := range content {
			if file.IsDir() {
				AvailableLanguages = append(AvailableLanguages, file.Name())
			}
		}
	}
}

func LoadLanguageByName(langName string) error {
	return LoadLanguage(stringutil.FindInSlice(&AvailableLanguages, langName))
}

func LoadLanguage(index int) error {

	SaveCurLanguage()

	if index < 0 && index <= len(AvailableLanguages) {
		return oobError(index)
	}

	for i := range words {
		words[i] = make([]word, 0)
	}

	//replacements := make(map[string]string)

	curLang = AvailableLanguages[index]
	langFolder := stringutil.Join(curLang, "/")

	data, err := ioutil.ReadFile(stringutil.Join(dirname, stringutil.Join(langFolder, wordfile)))

	if err == nil {
		lines := strings.Split(string(data), "\n")

		for _, line := range lines {
			elements := strings.Split(line, ":=:")
			if len(elements) == 4 {
				name := strings.TrimSpace(elements[0])
				translation := strings.TrimSpace(elements[1])
				successCounter, err := strconv.Atoi(strings.TrimSpace(elements[2]))
				index, err2 := strconv.Atoi(strings.TrimSpace(elements[3]))

				if err == nil && err2 == nil {
					words[index] = append(words[index], word{name, translation, successCounter, index})
				}
			}
		}
	} else {
		return err
	}
	return nil
}

func SaveCurLanguage() {

	if modified {
		tmpStr := ""

		for _, wordlist := range words {
			for _, elem := range wordlist {

				line := elem.Name
				line = stringutil.Join(line, " :=: ")
				line = stringutil.Join(line, elem.Translation)
				line = stringutil.Join(line, " :=: ")
				line = stringutil.Join(line, strconv.Itoa(elem.successCounter))
				line = stringutil.Join(line, " :=: ")
				line = stringutil.Join(line, strconv.Itoa(elem.pocketIndex))
				line = stringutil.Join(line, "\n")

				tmpStr = stringutil.Join(tmpStr, line)
			}
		}

		langFolder := stringutil.Join(curLang, "/")
		ioutil.WriteFile(stringutil.Join(dirname, stringutil.Join(langFolder, wordfile)), []byte(tmpStr), os.ModePerm)
		modified = false
	}
}

func DeleteCurLanguage() error {
	err := os.RemoveAll(stringutil.Join(dirname, curLang))

	Init()

	return err
}

func AddVocable(name, translation string) {
	words[0] = append(words[0], word{name, translation, 0, 0})
	modified = true
}

func DeleteVocable(name, translation string) {

loop:
	for x, wordlist := range words {
		for i, elem := range wordlist {
			if elem.Name == name && elem.Translation == translation {
				words[x] = append(words[x][:i], words[x][i+1:]...)
				modified = true
				break loop
			}
		}
	}
}

func DeleteVocableSplit(compound string) {
	l := strings.SplitN(compound, " - ", 2)

	if len(l) == 2 {
		DeleteVocable(strings.TrimSpace(l[0]), strings.TrimSpace(l[1]))
	}
}

func GetAll() []string {

	all := make([]string, 0)

	for _, wordlist := range words {
		for _, elem := range wordlist {

			line := elem.Name
			line = stringutil.Join(line, " - ")
			line = stringutil.Join(line, elem.Translation)

			all = append(all, line)
		}
	}

	return stringutil.Mergesort(all)
}

func Language() (string, string) {
	return stringutil.SplitFirst(curLang, " - ")
}

func AddLanguage(lang1, lang2 string) {
	lang1 = stringutil.UpperCaseOnlyFirst(lang1)
	lang2 = stringutil.UpperCaseOnlyFirst(lang2)

	langName := stringutil.Join(stringutil.Join(lang1, " - "), lang2)

	os.Mkdir(stringutil.Join(dirname, langName), os.ModeDir|os.ModePerm)

	langFolder := stringutil.Join(langName, "/")
	ioutil.WriteFile(stringutil.Join(dirname, stringutil.Join(langFolder, wordfile)), []byte(""), os.ModePerm)

	Init()

}

//returns a random word with its index.
//words with a low succes counter are more likely to be returned.
func RandomWord(pocketIndex int) (word, int) {
	x := make([]int, 0) //index list
	for index, elem := range words[pocketIndex] {
		for i := 0; i < successPerPocket[pocketIndex] - elem.successCounter; i++ {
			x = append(x, index) //add index to the list
		}
	}
	len := big.NewInt((int64)(len(x)))
	tmpRand, _ := rand.Int(rand.Reader, len)
	index := x[tmpRand.Int64()]
	//index := rand.Int() % PocketSize(pocketIndex)
	return words[pocketIndex][index], index

}

func Right(elem word, index int) {

	pocketIndex := elem.pocketIndex

	if words[pocketIndex][index].Name == elem.Name && words[pocketIndex][index].Translation == elem.Translation {

		curWord := &words[pocketIndex][index]
		curWord.successCounter++

		if curWord.successCounter >= successPerPocket[pocketIndex] {
			moveToNextPocket(pocketIndex, index)
		}
	}
	modified = true
}

func False(elem word, index int) {
	pocketIndex := elem.pocketIndex

	if words[pocketIndex][index].Name == elem.Name && words[pocketIndex][index].Translation == elem.Translation {

		curWord := &words[pocketIndex][index]

		if curWord.successCounter == 0 && curWord.pocketIndex > 0 {

			// add word to previous pocket
			curWord.pocketIndex--
			words[pocketIndex-1] = append(words[pocketIndex-1], *curWord)

			// delete word out of current pocket
			words[pocketIndex] = append(words[pocketIndex][:index], words[pocketIndex][index+1:]...)

		} else {
			curWord.successCounter = 0
		}
	}
	modified = true
}

func moveToNextPocket(pocketIndex, wordIndex int) {
	if pocketIndex < len(words) {
		curWord := words[pocketIndex][wordIndex]
		curWord.successCounter = 0
		curWord.pocketIndex++

		words[pocketIndex] = append(words[pocketIndex][:wordIndex], words[pocketIndex][wordIndex+1:]...)
		words[pocketIndex+1] = append(words[pocketIndex+1], curWord)
	}
}

func PocketSize(pocketIndex int) int {
	size := -1

	if pocketIndex >= 0 && pocketIndex < len(words) {
		size = len(words[pocketIndex])
	}
	return size
}

func NumPockets() int {
	return len(words)
}
