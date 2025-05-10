package environment

type Node struct {
	Program    Program
	Statement  Statement
	Expression Expression
}

type Program struct {
	statements []Statement
}

type Statement struct {
	Let        LetStatement
	Expression ExpressionStatement
}

type LetStatement struct {
	name  Identifier
	value Expression
}

type Identifier struct {
	value string
}

type Expression struct {
	Identifier     Identifier
	IntegerLiteral uint
}

type ExpressionStatement struct {
	expression Expression
}
