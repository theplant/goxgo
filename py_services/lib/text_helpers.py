# -*- coding: utf-8 -*-

import regex as re


def strip_parts( body, numbers = True, punctuation = True, links = True,
                 tags = True, vars = True ):
    """
    Wrapper function to pass a body of text through the strip_x functions of
    this module.

    :param string body: the body of text that gets passed to the strip_x functions
    :param boolean numbers:
    :param boolean punctuation:
    :param boolean links:
    :param boolean vars:

    :rtype: string
    :returns: string -- body stripped
    """
    # order does matter
    if links:
        body = strip_links( body )
    if vars:
        body = strip_braced_vars( body )
    if tags:
        body = strip_tags( body )
    if numbers:
        body = strip_numbers( body )
    if punctuation:
        body = strip_punctuation( body )
    return body


def strip_numbers( body ):
    """
    Strip numbers (surrounded by white-space) from a given body

    :param string body: the body of text the numbers should be stripped from

    :rtype: string
    :returns: string -- body stripped of all numbers
    """
    # number_mask_re = re.compile( u"\s\p{N}[\p{N}.]*\s" )
    number_mask_re = re.compile( u"[0-9]+" )
    return re.sub( number_mask_re, ' ', body )


def strip_punctuation( body, strict_protocol = False ):
    """
    Strip punctuation from a given body

    :param string body: the body of text the punctuation marks should be stripped from

    :rtype: string
    :returns: string -- body stripped of all punctuation
    """
    # most of the "general punctuation" range (2000 to 206F) + other Unicode characters.
    # punct_re = re.compile( u"[‐-⁞«»\"':：「」.。,、?？¿!！¡/／]+" )
    punct_re = re.compile( u"[‐-⁞«»\"':：「」.。,;<>、?？¿!！¡/／\[\]\(\)]+" )
    return re.sub( punct_re, ' ', body )


def strip_links( body, strict_protocol = True ):
    """
    Strip links (starting with http(s) from a given body

    :param string body: the body of text comments should be stripped from

    :rtype: string
    :returns: string -- body stripped of all links
    """
    if strict_protocol:
        url_re = re.compile( """http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&#+]|[!*'\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+""" )
    else:
        url_re = re.compile( """([http://]?|[https://]?)(?:[a-zA-Z]|[0-9]|[$-_@.&#+]|[!*'\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+""" )
    stripped = re.sub( url_re, '', body );
    return stripped


def strip_tags( body ):
    """
    Strip tags from a given body

    :param string body:

    :rtype: string
    :returns: string -- body stripped of all <zyz>
    """
    var_mask_re = re.compile( u"[%$]?\<+([^<>]+)\>+" )
    return re.sub( var_mask_re, '', body )


def strip_braced_vars( body ):
    """
    Strip strings enclosed by braces from a given body

    :param string body:

    :rtype: string
    :returns: string -- body stripped of all {}
    """
    var_mask_re = re.compile( u"[%$]?\{+([^{}]+)\}+" )
    return re.sub( var_mask_re, '', body )
