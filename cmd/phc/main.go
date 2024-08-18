package main

import (
	"fmt"
	"github.com/CyberGefest/Pc_how_conosle.git/internal/helper"
)

func main() {
	GamesData, err := helper.ExtractGamesData("C:\\Users\\admin\\GolangProjects\\Pc_how_conosle\\data")
	fmt.Println(GamesData, err)
}
