package goxgo

import (
	"fmt"
	gxg "github.com/theplant/goxgo"
	"sync"
	"testing"
)

func TestGoXGoSingleCalls(t *testing.T) {
	dsn := gxg.DSN{
		Protocol: "tcp",
		Host:     "localhost",
		Port:     4242,
	}

	tokenizePayload := gxg.TokenizeRequest{
		Target: &gxg.CallTarget{Services: []string{"NLTK/tokenize"}, Version: "0.1"},
		Body:   "Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied",
		Locale: "en",
	}

	var tokenizeResponse gxg.TokenizeResponse
	gxg.Call(&dsn, &tokenizePayload, &tokenizeResponse)
	t.Log("tokenizeResponse", tokenizeResponse)

	stemPayload := gxg.StemRequest{
		Target: &gxg.CallTarget{Services: []string{"NLTK/stem"}, Version: "0.1"},
		Words:  tokenizeResponse.Tokens,
		Locale: tokenizeResponse.Locale,
	}

	var stemResponse gxg.StemResponse
	gxg.Call(&dsn, &stemPayload, &stemResponse)
	t.Log("stemResponse", stemResponse)
}

func TestGoXGoMultiCalls(t *testing.T) {
	dsn := gxg.DSN{
		Protocol: "tcp",
		Host:     "localhost",
		Port:     4242,
	}

	compl := sync.WaitGroup{}
	tokenizingJobs := func() chan *gxg.TokenizeRequest {
		tChan := make(chan *gxg.TokenizeRequest, 32)
		go func(s *sync.WaitGroup) {
			for i := 0; i < 100; i++ {
				s.Add(1)
				tChan <- &gxg.TokenizeRequest{
					Target: &gxg.CallTarget{Services: []string{"NLTK/tokenize"}, Version: "0.1"},
					Body:   fmt.Sprintf("Job Number %v. Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied", i),
					Locale: "en",
				}
			}
			close(tChan)
		}(&compl)
		return tChan
	}()

	for {
		payload, open := <-tokenizingJobs
		if !open {
			break
		}
		go func(p *gxg.TokenizeRequest, s *sync.WaitGroup) {
			var tokenizeResponse gxg.TokenizeResponse
			gxg.Call(&dsn, p, &tokenizeResponse)
			t.Log("tokenizeResponse", tokenizeResponse)
			s.Done()
		}(payload, &compl)
	}
	compl.Wait()
}
