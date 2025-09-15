package usecase

// TODO: create a .go file to house each struct
// NOTE: debit  (usually left  hand side): money taken from account
// NOTE: credit (usually right hand side): money put   into account

// takes a csv file, persist each row in debit table
type IngestFinancialStatement struct{}

// create row(s) in credit table
type CategorizeTransaction struct{}
