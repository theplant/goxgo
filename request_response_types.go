// Request and Response structures for - in this first case - python services
// providing NLP/scientific math functions

package goxgo

import (
	"fmt"
	"encoding/json"
)


/*
CallTarget contains the target information for the payload by naming services to
be called and the desired version of that service

Services - a <[]string> of "service name/procedure name" pairs that the payload
should be "piped through". ie:

 	[ "NLTK/stem" ] or
 	[ "NLTK/tokenize", "NLTK/stem" ]

If you list more than one the result from the first gets passed as arguments
into the next. Make sure they are passable - and raise exceptions on the
server-side if they don't.

NB: this pipelining is not yet implemented

Version - a <string> representing the API version
*/
type CallTarget struct {
	Services []string
	Version	 string
}


/* TokenizeRequest - request structure
*/
type TokenizeRequest struct {
	Target *CallTarget
	Body   string
	Locale string
}

/* TokenizeResponse - response structure
*/
type TokenizeResponse struct {
	Locale string	`json:"locale"`
	Tokens []string `json:"tokens"`
}


/* StemRequest - request structure
*/
type StemRequest struct {
	Target *CallTarget
	Words  []string
	Locale string
}

/* StemResponse - response structure
*/
type StemResponse struct {
	Locale string	`json:"locale"`
	Words  []string `json:"tokens"`
}


/* VsmDiffRequest - Request structure for the cosine difference function.


*/
type VsmDiffRequest struct {
	Target  		*CallTarget
	Docs            []string
	Drop_Stopwords	bool
	Stem_Words		bool
}

/* VsmDiffRequest - Response structure
*/
type VsmDiffResponse struct {
	Diff	 float64 `json:"cosine_diff"`
	DocLangA string  `json:doc1_detected_lang`
	DocLangB string  `json:doc2_detected_lang`
}


// Serialization function for the request structures into JSON
func Serialize(i interface{}) (payload []byte, err error) {
	payload, err = json.Marshal(i)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

// Unserialization function for the structures from JSON to a
// response structure
func Unserialize(payload []byte, i interface{}) (err error){
	err = json.Unmarshal(payload, &i)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}
