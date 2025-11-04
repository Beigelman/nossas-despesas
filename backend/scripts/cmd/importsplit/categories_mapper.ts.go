package importsplit

func SplitCategoryToCategory(category string) int { //nolint:funlen
	switch category {
	case "Geral":
		return 64
	case "Mercado":
		return 16
	case "Jantar fora":
		return 17
	case "Táxi":
		return 34
	case "Seguro":
		return 45
	case "Hotel":
		return 57
	case "Filmes":
		return 22
	case "Bebidas alcoólicas":
		return 19
	case "Entretenimento - Outros":
		return 29
	case "Combustível":
		return 30
	case "Presents":
		return 43
	case "Transporte - Outros":
		return 39
	case "Estacionamento":
		return 31
	case "Casa - Outros":
		return 14
	case "Vida - Outros":
		return 43
	case "Produtos de limpeza":
		return 15
	case "Aluguel":
		return 1
	case "TV/Telefone/Internet":
		return 6
	case "Manutenção":
		return 8
	case "Eletricidade":
		return 4
	case "Aquecimento/gás":
		return 5
	case "Vestuário":
		return 44
	case "Despesas médicas":
		return 51
	case "Animais de estimação":
		return 9
	case "Móveis":
		return 11
	case "Carro":
		return 39
	case "Eletrônicos":
		return 10
	case "Comidas e bebidas - Outros":
		return 21
	case "Esports":
		return 40
	case "Serviços":
		return 12
	case "Serviços públicos - Outros":
		return 15
	case "Ônibus/trem":
		return 33
	case "Avião":
		return 56
	case "Educação":
		return 42
	case "Limpeza":
		return 13
	case "Música":
		return 26
	case "Jogos":
		return 24
	case "Impostos":
		return 63
	default:
		return 64
	}
}
