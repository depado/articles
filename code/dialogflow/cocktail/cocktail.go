package cocktail

// FullDrink is a complete description and recipe of a single Drink
type FullDrink struct {
	IDDrink         string      `json:"idDrink"`
	StrDrink        string      `json:"strDrink"`
	StrVideo        interface{} `json:"strVideo"`
	StrCategory     string      `json:"strCategory"`
	StrIBA          string      `json:"strIBA"`
	StrAlcoholic    string      `json:"strAlcoholic"`
	StrGlass        string      `json:"strGlass"`
	StrInstructions string      `json:"strInstructions"`
	StrDrinkThumb   string      `json:"strDrinkThumb"`
	StrIngredient1  string      `json:"strIngredient1"`
	StrIngredient2  string      `json:"strIngredient2"`
	StrIngredient3  string      `json:"strIngredient3"`
	StrIngredient4  string      `json:"strIngredient4"`
	StrIngredient5  string      `json:"strIngredient5"`
	StrIngredient6  string      `json:"strIngredient6"`
	StrIngredient7  string      `json:"strIngredient7"`
	StrIngredient8  string      `json:"strIngredient8"`
	StrIngredient9  string      `json:"strIngredient9"`
	StrIngredient10 string      `json:"strIngredient10"`
	StrIngredient11 string      `json:"strIngredient11"`
	StrIngredient12 string      `json:"strIngredient12"`
	StrIngredient13 string      `json:"strIngredient13"`
	StrIngredient14 string      `json:"strIngredient14"`
	StrIngredient15 string      `json:"strIngredient15"`
	StrMeasure1     string      `json:"strMeasure1"`
	StrMeasure2     string      `json:"strMeasure2"`
	StrMeasure3     string      `json:"strMeasure3"`
	StrMeasure4     string      `json:"strMeasure4"`
	StrMeasure5     string      `json:"strMeasure5"`
	StrMeasure6     string      `json:"strMeasure6"`
	StrMeasure7     string      `json:"strMeasure7"`
	StrMeasure8     string      `json:"strMeasure8"`
	StrMeasure9     string      `json:"strMeasure9"`
	StrMeasure10    string      `json:"strMeasure10"`
	StrMeasure11    string      `json:"strMeasure11"`
	StrMeasure12    string      `json:"strMeasure12"`
	StrMeasure13    string      `json:"strMeasure13"`
	StrMeasure14    string      `json:"strMeasure14"`
	StrMeasure15    string      `json:"strMeasure15"`
	DateModified    string      `json:"dateModified"`
}

// FullDrinkList contains a slice of FullDrink
type FullDrinkList struct {
	Drinks []FullDrink `json:"drinks"`
}

// Drink is a minimal representation of a FullDrink
type Drink struct {
	Name     string `json:"strDrink"`
	Thumnail string `json:"strDrinkThumb"`
	ID       string `json:"idDrink"`
}

// DrinkList contains a slice of Drink
type DrinkList struct {
	Drinks []Drink `json:"drinks"`
}
