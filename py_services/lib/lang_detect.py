# -*- coding: utf-8 -*-

import cld


def is_supported( lang_code ):
    """
    :param string lang_code:

    :rtype: bool
    :returns: whether or not the given language code is supported by cld
    """
    for lang in cld.LANGUAGES:
        if lang[ 1 ].lower() == lang_code.lower():
            return True
    return False

def detect_lang( text, lang_code_hint = None ):
    """
    Try to detect the language of a given text.

    :param ustring text: Body of text of which we want to detect the language
    :param string lang_code_hint: Expected language code. May give bias towards the expected language

    :rtype: tuple
    :returns:
    detected language name in English,
    language code,
    reliability flag,
    list of possible languages - lists details for other possible matches of the form
    [('ENGLISH', 'en', text percentage, normalized score), ...]
    """
    try:
        detected_lang, detected_lang_code, is_reliable, textBytesFound, details = \
            cld.detect( text, hintLanguageCode = lang_code_hint,
                        pickSummaryLanguage = False, isPlainText = False,
                        removeWeakMatches=False )
    except UnicodeEncodeError:
        detected_lang, detected_lang_code, is_reliable, textBytesFound, details = \
            cld.detect( text.encode( 'utf-8' ), hintLanguageCode = lang_code_hint,
                        pickSummaryLanguage = False, isPlainText = False,
                        removeWeakMatches=False )
    except cld.error as e:
        raise Exception( str( err ) )

    possible_matches = []
    for d in details:
        possible_matches.append( ( d[0].title(), d[1], d[2] ) )
    return ( detected_lang.title(), detected_lang_code, is_reliable, possible_matches )
