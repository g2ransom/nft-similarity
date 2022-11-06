package query

import "fmt"

type EIP721QueryResponse struct {
	Data ERC721Contracts `json:"data"`
}

type ERC721Contracts struct {
	Contracts []ERC721Contract `json:"erc721Contracts"`
}

type ERC721Contract struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Tokens []Token `json:"tokens"`
}
type Token struct {
	Identifier string `json:"identifier"`
	Owner struct {
		ID string `json:"id"`
	} `json:"owner"`
	Uri string `json:"uri"`
}

type ERC721Token struct {
	ContractID string
	Name string
	Symbol string
	Identifier string
	Owner string
	Uri string
}

func NewERC721Token(contractId string, name string, symbol string, identifier string, owner string, uri string) *ERC721Token {
	return &ERC721Token{
		ContractID: contractId,
		Name: name,
		Symbol: symbol,
		Identifier: identifier,
		Owner: owner,
		Uri: uri,
	}
}

func EIP721Query(pageLen int, timeStart int, timeEnd int, tokenLen int) map[string]string {
	queryStatement := fmt.Sprintf(`{
	  erc721Contracts(first: %d, 
	    where: 
	    {
	      transfers_: 
	      	{timestamp_gte: %d, timestamp_lt: %d}
	    }) {
	    id
	    name
	    symbol
	    tokens(first: %d) {
	      identifier
	      owner {
	        id
	      }
	      uri
	    }
	  }
	}`, pageLen, timeStart, timeEnd, tokenLen)
	return map[string]string{"query": queryStatement}
}