package language

import "fmt"

func ExampleLanguage_Human() {
	fmt.Println(French.Human())
	// Output: Français
}

func ExampleLanguage_String() {
	fmt.Println(French.String())
	// Output: fr
}

func ExampleLanguage_Path() {
	fmt.Println(French.Path("/eu/ec/eci/"))
	fmt.Println(French.Path("/eu/ec/eci/scheme."))
	fmt.Println(AllEnglish.Path("/release/"))
	// Output: /eu/ec/eci/fr.html
	// /eu/ec/eci/scheme.fr.html
	// /release/
}
