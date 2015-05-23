package main

import (
	"github.com/indidev/vocable-o/util/stringutil"
	"github.com/indidev/vocable-o/util/mathutil"
	"github.com/indidev/vocable-o/console"
	"github.com/indidev/vocable-o/lang"
	"io/ioutil"
	"strings"
	"strconv"
	"time"
	"math/rand"
)

var replacements map[string]string

func main() {

	rand.Seed(time.Now().Unix())
	console.Init()

	console.SetInfoBottom("Esc - Back, Enter - Confirm")

	initReplacements()
	lang.Init()

	mainMenu()

	console.Quit()

}

func mainMenu() {

	options := []string{"Learn Language", "Edit languags", "Show replacements", "Quit"}

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

	console.Clear()

	elements := make([]string, len(replacements))

	i := 0
	for old, new := range replacements {
		elements[i] = stringutil.Join(stringutil.Join(old, " -> "), new)
		i++
	}

	console.DisplayCentered(elements)

	console.WaitForAnyInput()
}

func edtLangMenu() {

	console.Clear()

	for repeate := true; repeate; {

		options := lang.AvailableLanguages
		options = append(options, "Add a language")
		options = append(options, "Back to main menu")

		option := console.Menu(options, "Choose your language to edit:")

		switch true {

		case option == len(options)-2:
			addLanguage()

		case (option == len(options)-1) || (option == -1):
			repeate = false

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

func edtVocables() {

	curIndex := 0

	for repeate := true; repeate; {
		words := lang.GetAll()
		console.SetInfoBottom("Esc - Back, Enter - Confirm, d - Delete")
		index, char := console.ExtendedMenu(words, "", []rune{'d'}, curIndex)

		switch (index) {
		case -1:
			repeate = false

		default:
			if char == 'd' {
				curIndex = mathutil.MaxInt(index - 1, 0)

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
						[]string{stringutil.Join(lang1, ":"),word, "                    ",
						stringutil.Join(lang2, ":")}, replacements, translation)

					if valid {
						lang.DeleteVocableSplit(words[index])
						lang.AddVocable(word, translation)
					}
				}
			}
		}
	}
	console.SetInfoBottom("Esc - Back, Enter - Confirm")

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
				[]string{stringutil.Join(lang1, ":"),word, "                    ",
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

	data, err := ioutil.ReadFile("replacements.txt")

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

	lang1, lang2 := lang.Language()

	r := "  ✔"
	f := "  ✘"

	for repeate := true; repeate; {

		if lang.PocketSize(pocketIndex) > 0 {
			console.Clear()
			randWord, index := lang.RandomWord(pocketIndex)

			text := []string{stringutil.Join(lang1, ":"), randWord.Name,
				"                            ", stringutil.Join(lang2, ":"), ""}

			input, valid := console.DisplayCenteredWithInput(text, replacements, "")

			if valid {
				mark := ""
				if strings.ToLower(input) == strings.ToLower(randWord.Translation) {
					lang.Right(randWord, index)
					mark = r
				} else {
					lang.False(randWord, index)
					mark = f
				}


				text := []string{stringutil.Join(lang1, ":"), randWord.Name,
					"                            ", stringutil.Join(lang2, ":"), randWord.Translation, "",
					stringutil.Join(input, mark)}
				console.Clear()
				console.DisplayCentered(text)
				console.WaitForAnyInput()

			} else {
				repeate = false
			}
		} else {

			text := []string{"No vocables left"}

			console.Clear()
			console.DisplayCentered(text)
			console.WaitForAnyInput()

			repeate = false;
		}
	}
}
