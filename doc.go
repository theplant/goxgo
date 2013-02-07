// goXgo - Intro

/*
They basic idea for/behind goXgo is to be able to have external services
provide functionality that might not yet exist in go or would take too much
time to implement yourself - in an RPCish manner.

These services could be written in any other language and be running on
the same machine as the go code or be somewhere else.

The first idea for this came up when we talked about Natural Language
Processing (NLP) for which there are a couple of nice libraries - most
notably the NLTK (http://www.nltk.org).

Another thing we will need when dealing with NLP (/Pattern Matching/
Information Retreival/etc.) will be a solid scientific mathematics library:
I found some starts in go but nothing compared to numPy (http://www.numpy.org/).
R (http://www.r-project.org/) would come to mind too.

To start playing and testing NLP and statistical features i think it makes
a lot of sense to use stuff that's out there (and was written by linguists
or mathematicians specialized in scientific computing ;)


To network those services i chose
ZMQ(http://zguide.zeromq.org/page:all#-MQ-in-a-Hundred-Words).
The python code i built some examples with already has some stubs to also
provide an HTTP frontend. It also uses gevent (http://gevent.org/) which is a
coroutine implementation for python using libev in the background.

The idea is to basically have a bunch of (static) functions that you want to
call from go:

1. You group those functions into Services

2. Give the service a name

3. And register the service with a (or more) "networking frontend"


On the python side this would look somthing like this

 from service_frontend import ZmqFrontend
 from service import Service
 from lib.static_nltk_wrappers import tokenize, stem, vsm_compare

 NLTKService = Service( name = 'NLTK' )
 NLTKService.register_service_method( f = tokenize )
 NLTKService.register_service_method( f = stem )
 NLTKService.register_service_method( f = vsm_compare )

 zmq_frontend = ZmqFrontend()
 zmq_frontend.register_service( NLTKService )
 zmq_frontend.start()

 # this is just an example - does _not_ work yet
 http_frontend = HttpFrontend()
 http_frontend.register_service( NLTKService,
 				default_verb = HttpFrontend.HTTP_POST,
 				method_verb_map = { HttpFrontend.HTTP_GET:[ 'tokenize', '...' ]  }, )
 http_frontend.start()


Try it:

First you need to install some stuff:

You need python2.7+, pip, libev4, libev-dev

 [OPTIONAL] python2.7 dev-headers, cython, std buildtools [/OPTIONAL]

Then you can install the needed python packages:

 ~/data/dev/go/src/goxgo [goxgo] $sudo pip install -r py_services/requirements.txt

 [OPTIONAL]
 One library in the py_services uses cython so you have to compile it. This is a bit annoying for a small test but i think the function is kind of cool/interesting.
 The following assumes your system puts header files in /usr/include/ ... (debian/ubuntu)

 ~/data/dev/go/src/goxgo [goxgo] $ cd py_services/lib/
 ~/data/dev/go/src/goxgo [goxgo] $ cython -a vsm_diff.pyx
 ~/data/dev/go/src/goxgo [goxgo] $ gcc -c -lpthread -O2 -Wall -fPIC -I /usr/include/python2.7/ vsm_diff.c
 ~/data/dev/go/src/goxgo [goxgo] $ gcc -shared -lpthread -O2 -Wall vsm_diff.o -o vsm_diff.so
 ~/data/dev/go/src/goxgo [goxgo] $ cd ../..
 [/OPTIONAL]

 If this does not work or is too tedious to do but you still wanna test comment the following line

 	NLTKService.register_service_method( f = vsm_compare )

 from the file goxgo/py_services/test_server.py

 and comment everything after

 	// test vsm diff

 in tests/goxgo_test.go (currently line 69)


To run the python service:

 ~/data/dev/go/src/goxgo [goxgo] $ python py_services/test_server.py
 ('tokenize', ArgSpec(args=['body', 'locale'], varargs=None, keywords=None, defaults=('en',)))
 ('stem', ArgSpec(args=['words', 'locale'], varargs=None, keywords=None, defaults=('en',)))
 ('vsm_compare', ArgSpec(args=['docs', 'drop_stopwords', 'stem_words'], varargs=None, keywords=None, defaults=(True, True)))
 Starting server. Listening on tcp://*:4242...
 Start serving.

Now the service runs and you can hit it.

 ~/data/dev/go/src/goxgo [goxgo] $ go test -v tests/goxgo_test.go
 === RUN TestGoXGo
 --- PASS: TestGoXGo (0.60 seconds)
 goxgo_test.go:36: {Locale:en Tokens:[Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied]}
 goxgo_test.go:46: {Locale:en Words:[giv me a tok vert of thi unoptimz body of text pls. ont success don we wil try to stem the word too. test try embody embody]}
 goxgo_test.go:14: result: [Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied]
 goxgo_test.go:14: result: [Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied]
 goxgo_test.go:14: result: [Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied]
 goxgo_test.go:98: Compare cats and dogs:
 {Diff:0.16651983559131622 DocLangA: DocLangB:}
 goxgo_test.go:127: Compare cats and engines:
 {Diff:0.020799094811081886 DocLangA: DocLangB:}
 PASS
 ok      command-line-arguments  0.611s


Service/Function/Argument Naming/Mapping/Case convention

In python i registered a function tokenize:

 NLTKService = Service( name = 'NLTK' )
 NLTKService.register_service_method( f = tokenize )

In go i need to call it with this target:

 gxg.CallTarget { Services: []string{"NLTK/tokenize"}, Version: "0.1" }

The arguments for that function in python are:

 def tokenize( body, locale = 'en' ):

BUT in go my payload keys starts with uppercase letters:

 tokenizePayload := gxg.TokenizeRequest {
 	Target: &gxg.CallTarget { Services: []string{"NLTK/tokenize"}, Version: "0.1" },
 	Body: "Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied",
 	Locale: "en",
 }

To keep case cannonical in both languages the python frontend lower-cases every
parameter it finds. Unserialization from the JSON response does not need
something like that because I can rely on encoding/json/Unmarshal to find the
struct definition...

*/
package goxgo
