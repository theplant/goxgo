#!/usr/bin/python
# -*- coding: utf-8 -*-

import regex as re

from string import lower

import tinysegmenter

from nltk.stem import PorterStemmer, SnowballStemmer
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize


class Stopwords( object ):
    """
    Wrapper class for nltk stopwords and internal additional lists.
    """

    def __init__( self, locale ):
        if locale != 'default':
            # cut off any localization parts of the locale
            self.locale = lower( locale[0:2] )
        else:
            self.locale = locale

        if self.locale == 'da':
            self.s = stopwords.words( 'danish' )
        elif self.locale == 'nl':
            self.s = stopwords.words( 'dutch' )
        elif self.locale == 'en':
            self.s = stopwords.words( 'english' )
        elif self.locale == 'fi':
            self.s = stopwords.words( 'finnish' )
        elif self.locale == 'fr':
            self.s = stopwords.words( 'french' )
        elif self.locale == 'de':
            self.s = stopwords.words( 'german' )
        elif self.locale == 'hu':
            self.s = stopwords.words( 'hungarian' )
        elif self.locale == 'it':
            self.s = stopwords.words( 'italian' )
        elif self.locale == 'no':
            self.s = stopwords.words( 'norwegian' )
        elif self.locale == 'pt':
            self.s = stopwords.words( 'portuguese' )
        elif self.locale == 'ru':
            self.s = stopwords.words( 'russian' )
        elif self.locale == 'es':
            self.s = stopwords.words( 'spanish' )
        elif self.locale == 'sv':
            self.s = stopwords.words( 'swedish' )
        elif self.locale == 'tr':
            self.s = stopwords.words( 'turkish' )
        else:
            # default for languages we do not have a stopword list for yet.
            self.s = []

    def __iter__( self ):
        for sw in self.s:
            yield sw


class Tokenizer( object ):
    """
    Wrapper class for different tokenizers.
    """

    def __init__( self, locale ):
        # cut off any localization parts of the locale
        self.locale = lower( locale[0:2] )
        if self.locale == 'ja':
            self.t = tinysegmenter.TinySegmenter()
            self.tokenize_function = self.t.tokenize
        else:
            self.tokenize_function = word_tokenize

    def tokenize( self, body ):
        """
        tokenizes a given body. will attempt utf-8 decoding when a
        ``UnicodeDecodeError`` gets raised.

        :param string body: the body to be tokenized
        :rtype: list
        :returns: a list containing the tokens
        """

        try:
            return self.tokenize_function( body )
        except UnicodeDecodeError:
            return self.tokenize_function( body.decode( 'utf-8' ) )



class NonStemmer( object ):
    """
    A dummy stemmer that will just return the word as is
    """

    def stem( self, word ):
        return word


class Stemmer( object ):
    """
    Wrapper class for different Stemmers.
    """

    def __init__( self, locale ):
        # cut off any localization parts of the locale
        if locale != 'default':
            self.locale = lower( locale[0:2] )
        else:
            self.locale = locale

        if self.locale == 'da':
            self.s = SnowballStemmer( 'danish' )
        elif self.locale == 'nl':
            self.s = SnowballStemmer( 'dutch' )
        elif self.locale == 'en':
            self.s = PorterStemmer()
        elif self.locale == 'fi':
            self.s = SnowballStemmer( 'finnish' )
        elif self.locale == 'fr':
            self.s = SnowballStemmer( 'french' )
        elif self.locale == 'de':
            self.s = SnowballStemmer( 'german' )
        elif self.locale == 'hu':
            self.s = SnowballStemmer( 'hungarian' )
        elif self.locale == 'it':
            self.s = SnowballStemmer( 'italian' )
        elif self.locale == 'no':
            self.s = SnowballStemmer( 'norwegian' )
        elif self.locale == 'pt':
            self.s = SnowballStemmer( 'portuguese' )
        elif self.locale == 'ro':
            self.s = SnowballStemmer( 'romanian' )
        elif self.locale == 'ru':
            self.s = SnowballStemmer( 'russian' )
        elif self.locale == 'es':
            self.s = SnowballStemmer( 'spanish' )
        elif self.locale == 'sv':
            self.s = SnowballStemmer( 'swedish' )
        else:
            # default for languages we do not have a stemmer for yet.
            self.s = NonStemmer()

    def stem( self, word ):
        """
        stems a given word. will attempt utf-8 decoding when a
        ``UnicodeDecodeError`` gets raised.
        """
        try:
            return self.s.stem( word )
        except UnicodeDecodeError:
            return self.s.stem( word.decode( 'utf-8' ) )
