# -*- coding: utf-8 -*-

import re
import numpy as np
cimport numpy as np

# We now need to fix a datatype for our arrays. I've used the variable
# DTYPE for this, which is assigned to the usual NumPy runtime
# type info object.
DTYPE = np.int
# "ctypedef" assigns a corresponding compile-time type to DTYPE_t. For
# every type in the numpy module there's a corresponding compile-time
# type with a _t-suffix.
ctypedef np.int_t DTYPE_t

from numpy.linalg import norm
from string import lower
from time import time
from pprint import pprint
from nltk_wrappers import Stopwords, Tokenizer, Stemmer
from lang_detect import detect_lang


class DocVector( object ):
    def __init__( self, text, lang = None, drop_stopwords = False,
                  stem_words = False, debug = False, timeit = False ):
        self.debug = debug
        self.timeit = timeit
        if self.timeit:
            st = time()
        self.words = {}
        # self.word_positions = {}
        self.lang = lang
        if not self.lang:
            self.lang = detect_lang( text )[1]
            print "***", self.lang
        self.tokenizer = Tokenizer( self.lang )
        self.stopwords = Stopwords( self.lang if drop_stopwords else 'default' )
        self.stemmer = Stemmer( self.lang if stem_words else 'default' )
        self._text = None

        # add the words from the doc to the words dict
        for w in self.tokenizer.tokenize( text ):
            w = w.lower()
            if w not in self.stopwords:
                w = self.stemmer.stem( w )
                self.words.setdefault( w, 0 )
                self.words[w] += 1
        self._text = " ".join( self.words )
        if self.timeit:
            print "__init__ took {}".format( time() - st )

    def compare( self, d ):
        if self.timeit:
            st = time()
        if not isinstance( d, DocVector ):
            raise Exception( "can only compare two DocVector objects." )
        all_words = list( set( self.words.keys() ).union( d.words.keys() ) )
        all_words.sort()

        cdef int i, w_length
        w_length = len( all_words )

        cdef np.ndarray v1 = np.zeros( w_length, dtype = DTYPE )
        cdef np.ndarray v2 = np.zeros( w_length, dtype = DTYPE )
        # cdef char* key
        for i in range( w_length ):
            # k = all_words[ i ].encode('utf-8') # makes stuff very slow - need to learn more about unicode/byte strings in cython. left for ref
            # key = k
            key = all_words[ i ]
            if key in self.words:
                v1[i] = self.words[ key ]
            if key in d.words:
                v2[i] = d.words[ key ]

        print all_words
        del all_words

        if self.debug:
            print v1
            print v2

        print v1
        print v2

        cdef float n1 = norm( v1 )
        cdef float n2 = norm( v2 )

        # due to dropped stopwords the norm of any doc vector could be 0
        # we should raise an exception then or return None when an option is set
        # TODO: raise Exception...
        cdef float doc_diff
        if n1 != 0 and n2 != 0:
            doc_diff = np.float64( np.dot( v1, v2 ) / ( n1 * n2 ) )
        else:
            doc_diff = np.float64( -1.0 )
        if self.timeit:
            print "compare took {}".format( time() - st )
        return doc_diff
