package models

type Product struct {
	ID       int32
	Name     string
	Price    float32
	Articles []ProductArticle
}

type ProductArticle struct {
	ID       int32
	Quantity int32
}

type ProductWithStock struct {
	Product
	Stock int32
}
