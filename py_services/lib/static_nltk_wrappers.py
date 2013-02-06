#!/usr/bin/python
# -*- coding: utf-8 -*-

from string import lower
from types import StringType, ListType

import regex as re
import tinysegmenter

from nltk.stem import PorterStemmer, LancasterStemmer, SnowballStemmer
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize

import vsm_diff

from text_helpers import strip_parts

def tokenize( body, locale = 'en' ):
    """
    tokenizes a given body. will attempt utf-8 decoding when a
    ``UnicodeDecodeError`` gets raised.

    :param string body: the body to be tokenized
    :rtype: list
    :returns: a list containing the tokens
    """

    # cut off any localization parts of the locale
    locale = lower( locale[0:2] )
    if locale == 'ja':
        t = tinysegmenter.TinySegmenter()
        tokenize_function = t.tokenize
    else:
        tokenize_function = word_tokenize

    try:
        res = tokenize_function( body )
    except UnicodeDecodeError:
        res = tokenize_function( body.decode( 'utf-8' ) )

    return { 'locale': locale,
             'tokens': res }


class NonStemmer( object ):
    """
    A dummy stemmer that will just return the word as is
    """

    def stem( self, word ):
        return word


def stem( words, locale = 'en'  ):
    """
    Wrapper class for different Stemmers.

    stems a given word. will attempt utf-8 decoding when a
    ``UnicodeDecodeError`` gets raised.
    """

    # cut off any localization parts of the locale
    if locale != 'default':
        locale = lower( locale[0:2] )
    else:
        locale = locale

    # print "STEMMING: Loacle:{} Words:{}".format(locale, words)

    if locale == 'da':
        s = SnowballStemmer( 'danish' )
    elif locale == 'nl':
        s = SnowballStemmer( 'dutch' )
    elif locale == 'en':
        # s = PorterStemmer()
        s = LancasterStemmer()
    elif locale == 'fi':
        s = SnowballStemmer( 'finnish' )
    elif locale == 'fr':
        s = SnowballStemmer( 'french' )
    elif locale == 'de':
        s = SnowballStemmer( 'german' )
    elif locale == 'hu':
        s = SnowballStemmer( 'hungarian' )
    elif locale == 'it':
        s = SnowballStemmer( 'italian' )
    elif locale == 'no':
        s = SnowballStemmer( 'norwegian' )
    elif locale == 'pt':
        s = SnowballStemmer( 'portuguese' )
    elif locale == 'ro':
        s = SnowballStemmer( 'romanian' )
    elif locale == 'ru':
        s = SnowballStemmer( 'russian' )
    elif locale == 'es':
        s = SnowballStemmer( 'spanish' )
    elif locale == 'sv':
        s = SnowballStemmer( 'swedish' )
    else:
        # default for languages we do not have a stemmer for yet.
        s = NonStemmer()

    def do_stem( word ):
        try:
            return s.stem( word )
        except UnicodeDecodeError:
            return s.stem( word.decode( 'utf-8' ) )


    # if isinstance( words, StringType ):
    #     return do_stem( words )
    # el

    if isinstance( words, ListType ):
        return { 'locale': locale,
                 'words': map( do_stem, words ) }
    else:
        print "nope:", type(words), words
        raise Exception( "The parameter words needs to be a list." )



def vsm_compare( docs, drop_stopwords = True, stem_words = True ):
    if not isinstance( docs, ListType ):
        raise Exception( "You need to submit a list of documents to compare." )
    d1 = vsm_diff.DocVector( strip_parts( docs[0] ), drop_stopwords = drop_stopwords, stem_words = stem_words )
    d2 = vsm_diff.DocVector( strip_parts( docs[1] ), drop_stopwords = drop_stopwords, stem_words = stem_words )
    response =  { 'doc1_detected_lang': d1.lang,
                  'doc2_detected_lang': d2.lang,
                  'cosine_diff': d1.compare( d2 ) }
    print response
    return response
