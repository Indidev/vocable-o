package main

import (
	"github.com/indidev/vocable-o/console"
	"github.com/indidev/vocable-o/lang"
	"github.com/indidev/vocable-o/util/mathutil"
	"github.com/indidev/vocable-o/util/stringutil"
	"io/ioutil"
	//"crypto/rand"
	"os"
	"strconv"
	"strings"
	//"time"
)

const replFile = "replacements.txt"

var replacements map[string]string

func main() {
	//add lanuage dir if not present
	os.Mkdir("languages", os.ModeDir|os.ModePerm)

	//rand.Seed(time.Now().Unix())
	console.Init()

	setInfoBottom("", true)

	initReplacements()
	lang.Init()

	mainMenu()

	console.Quit()
}

func mainMenu() {

	options := []string{"Learn Language", "Edit languages", "Show/Edit replacements", "Quit"}

	for repeate := true; repeate; {
		option := console.Menu(options, "What would you like to do?")

		switch option {
		// Learn Language
		case 0:
			learnLangMenu()

		// Edit Language
		case 1:
			edtLangMenu()

		// Show replacements
		case 2:
			showReplacements()

		// Quit
		case -1, 3:
			repeate = false
		}
	}
}

func showReplacements() {

	repeate := true
	curIndex := 0
	for repeate {
		console.Clear()

		setInfoBottom("Esc - Back, Enter - Edit, a - Add, d - Delete", false)

		items := make([]string, 0)
		for old, new := range replacements {
			items = append(items, stringutil.Join(stringutil.Join(old, " -> "), new))
		}
		index, char := console.ExtendedMenu(items, "", []rune{'d', 'a'}, curIndex)

		switch index {
		case -1:
			repeate = false

		default:
			key, _ := stringutil.SplitFirst(items[index], "->")
			switch true {
			// Delete replacement
			case char == 'd':
				delete(replacements, key)
				curIndex = mathutil.MaxInt(index-1, 0)

			// Add replacement
			case char == 'a':
				edtReplacement("")
			// Edit replacement
			default:
				edtReplacement(key)
			}
			saveReplacements()
		}
	}
	setInfoBottom("", true)
}

func edtReplacement(key string) {
	setInfoBottom("", true)
	console.Clear()
	old, valid := console.DisplayCenteredWithInput([]string{"", "", "", "Enter your string to replace:"}, nil, key)
	if valid {
		console.Clear()
		new, valid := console.DisplayCenteredWithInput([]string{"String to replace:", old, "", "Enter the new character/string "}, nil, replacements[key])
		if valid {
			replacements[old] = new
			delete(replacements, key)
		}
	}
}

func saveReplacements() {
	tmpStr := ""
	for key, value := range replacements {
		line := key
		line = stringutil.Join(line, " := ")
		line = stringutil.Join(line, value)
		line = stringutil.Join(line, "\n")

		tmpStr = stringutil.Join(tmpStr, line)
	}

	ioutil.WriteFile(replFile, []byte(tmpStr), os.ModePerm)
}

func edtLangMenu() {

	console.Clear()

	for repeate := true; repeate; {

		options := lang.AvailableLanguages
		options = append(options, "Add a language")
		options = append(options, "Back to main menu")

		option := console.Menu(options, "Choose your language to edit:")

		switch true {

		// Add language
		case option == len(options)-2:
			addLanguage()
		// Back to main menu
		case (option == len(options)-1) || (option == -1):
			repeate = false

		// Language selected
		default:
			err := lang.LoadLanguage(option)
			if err == nil {
				edtLang()
			} else {
				console.Clear()
				console.DisplayCentered([]string{err.Error()})
				console.WaitForAnyInput()
			}
		}
	}
}

func edtLang() {
	options := []string{"Add vocables", "Edit vocables", "Delete language", "Back to language selection"}

	info := stringutil.Join("[", stringutil.Join(lang.CurLang(), "] What would you like to do?"))

	for repeate := true; repeate; {
		console.Clear()
		option := console.Menu(options, info)

		switch option {

		// Add vocables
		case 0:
			addVocables()

		// Edit vocables
		case 1:
			edtVocables()

		// Delete Language
		case 2:

			console.Clear()

			question := []string{stringutil.Join("Realy delete: ", stringutil.Join(lang.CurLang(), "?")),
				"Type \"yes\" to do so."}
			input, valid := console.DisplayCenteredWithInput(question, nil, "")
			if ("yes" == input) && valid {
				err := lang.DeleteCurLanguage()
				if err != nil {
					console.Clear()
					console.DisplayCentered([]string{err.Error()})
					console.WaitForEnter()
				}
				repeate = false
			}

		case 3, -1:
			lang.SaveCurLanguage()
			repeate = false
		}

	}
}

func setInfoBottom(option string, standardOptions bool) {
	bottomText := ""
	if standardOptions {
		bottomText = stringutil.Join(bottomText, "Esc - Back, Enter - Confirm")
	}
	if bottomText != "" && option != "" {
		bottomText = stringutil.Join(bottomText, ", ")
	}
	bottomText = stringutil.Join(bottomText, option)
	console.SetInfoBottom(bottomText)
}

func edtVocables() {

	curIndex := 0

	for repeate := true; repeate; {
		words := lang.GetAll()
		setInfoBottom("d - Delete", true)
		index, char := console.ExtendedMenu(words, "", []rune{'d'}, curIndex)

		switch index {
		case -1:
			repeate = false

		default:
			if char == 'd' {
				curIndex = mathutil.MaxInt(index-1, 0)

				lang.DeleteVocableSplit(words[index])
			} else {
				lang1, lang2 := lang.Language()

				word, translation := stringutil.SplitFirst(words[index], " - ")
				console.Clear()
				word, valid := console.DisplayCenteredWithInput(
					[]string{"", "", "                    ",
						stringutil.Join(lang1, ":")}, replacements, word)

				if valid {

					console.Clear()
					translation, valid := console.DisplayCenteredWithInput(
						[]string{stringutil.Join(lang1, ":"), word, "                    ",
							stringutil.Join(lang2, ":")}, replacements, translation)

					if valid {
						lang.DeleteVocableSplit(words[index])
						lang.AddVocable(word, translation)
					}
				}
			}
		}
	}
	setInfoBottom("", true)
}

func addVocables() {

	lang1, lang2 := lang.Language()
	for repeate := true; repeate; {
		console.Clear()
		word, valid := console.DisplayCenteredWithInput(
			[]string{"", "", "                    ", stringutil.Join(lang1, ":")}, replacements, "")
		repeate = valid

		if valid {

			console.Clear()
			translation, valid := console.DisplayCenteredWithInput(
				[]string{stringutil.Join(lang1, ":"), word, "                    ",
					stringutil.Join(lang2, ":")}, replacements, "")
			repeate = valid

			if valid {
				lang.AddVocable(word, translation)
			}
		}
	}
}

func initReplacements() {

	replacements = make(map[string]string)

	data, err := ioutil.ReadFile(replFile)

	if err == nil {
		lines := strings.Split(string(data), "\n")

		for _, line := range lines {
			elements := strings.Split(line, ":=")
			if len(elements) == 2 {
				replacements[strings.TrimSpace(elements[0])] = strings.TrimSpace(elements[1])
			}
		}
	} else {
		console.Write(2, 2, err.Error())
	}
}

func addLanguage() {
	console.Clear()

	info := []string{"Enter the language you want to learn:"}

	lang2, valid := console.DisplayCenteredWithInput(info, replacements, "")

	if valid {
		info[0] = "Enter the language you already know:"

		console.Clear()

		lang1, valid := console.DisplayCenteredWithInput(info, replacements, "")

		if valid {
			lang.AddLanguage(lang1, lang2)
		}
	}
}

func learnLangMenu() {
	console.Clear()

	for repeate := true; repeate; {

		options := lang.AvailableLanguages
		options = append(options, "Back to main menu")

		option := console.Menu(options, "Choose your language to edit:")

		switch true {

		case (option == len(options)-1) || (option == -1):
			repeate = false

		default:
			err := lang.LoadLanguage(option)
			if err == nil {
				pocketSelect()
			} else {
				console.Clear()
				console.DisplayCentered([]string{err.Error()})
				console.WaitForAnyInput()
			}
		}
		lang.SaveCurLanguage()
	}
}

func pocketSelect() {

	for repeate := true; repeate; {
		options := make([]string, 0)

		for i := 0; i < lang.NumPockets(); i++ {
			item := strings.Join([]string{"Pocket ", strconv.Itoa(i + 1), " (",
				strconv.Itoa(lang.PocketSize(i)), " vocables)"}, "")
			options = append(options, item)
		}

		options = append(options, "Back to main menu")

		option := console.Menu(options, "Choose your Pocket:")

		switch true {
		case (option == len(options)-1) || (option == -1):
			repeate = false

		default:
			learn(option)
		}
	}
}

func learn(pocketIndex int) {

	for wordGuess(pocketIndex){}
}

func wordGuess(pocketIndex int) bool {

	lang1, lang2 := lang.Language()

	r := " /green  ✔"
	f := " /red  ✘"

	if lang.PocketSize(pocketIndex) > 0 {
		console.Clear()
		//get random word
		randWord, index := lang.RandomWord(pocketIndex)

		text := []string{stringutil.Join(lang1, ":"), randWord.Name,
			"                            ", stringutil.Join(lang2, ":"), ""}

		//display text and get user input
		input, valid := console.DisplayCenteredWithInput(text, replacements, "")

		//check if input is valid
		if valid {
			mark := ""
			//check if input is equal
			if stringutil.CheckEqual(input, randWord.Translation, true) {
				lang.Right(randWord, index)
				mark = r
			} else {
				if stringutil.Levenshtein(input, randWord.Translation, true, true, 1) <= 1 {
					mark = " /orange <- almost right"
					text = []string{stringutil.Join(lang1, ":"), randWord.Name,
						"                            ", stringutil.Join(lang2, ":"), stringutil.Join(input, mark)}

					console.Clear()
					input, valid = console.DisplayCenteredWithInput(text, replacements, input)
					if valid && stringutil.CheckEqual(input, randWord.Translation, true) {
						mark = r
					} else {
						lang.False(randWord, index)
						mark = f
					}
				} else {
					lang.False(randWord, index)
					mark = f
				}
			}

			text = []string{stringutil.Join(lang1, ":"), randWord.Name,
				"                            ", stringutil.Join(lang2, ":"), randWord.Translation, "",
				stringutil.Join(input, mark)}
			console.Clear()
			console.DisplayCentered(text)
			console.WaitForAnyInput()

		} else {
			return false
		}
	} else {

		text := []string{"No vocables left"}

		console.Clear()
		console.DisplayCentered(text)
		console.WaitForAnyInput()

		return false
	}
	return true
}
