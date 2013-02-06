# -*- coding: utf-8 -*-

# Copyright 2013 Florian von Bock (f at theplant dot jp)
#
# Base Class for the services

__author__ = "Florian von Bock"
__email__ = "f at theplant dot jp"
__version__ = "0.0.1"


import inspect


class ServiceError( Exception ):
    # TODO: refactor exceptions...
    def __init__( self, code = None, reason = None, description = None,
                  details = {} ):
        self.code = code
        self.reason = reason
        self.description = description
        self.details = details

    def __str__( self ):
        return str( self.description )


class Service( object ):
    """
    Service is a container for a set of functions or methods that are presented
    as a _named_ service.
    """
    def __init__( self, name ):
        self.name = name
        self.service_methods = {}

    def register_service_method( self, f = None, replace = False ):
        sig = inspect.getargspec( f )
        print( f.__name__, sig )
        args_len = len( sig.args )
        # if sig.defaults:
        #     defaults_len = len( sig.defaults )
        if f.__name__ in self.service_methods and not replace:
            raise ServiceError( reason = "Name already registered",
                                description = "There is alread a function or method registered with the name {}. If you want to replace that function you have to call register_service_method( self, f = {}, replace = True )".format( f.__name__, f.__name__ ) )
        self.service_methods[ f.__name__ ] = f

    def __call__( self, **kwargs ):
        # check arg spec/payload
        pass


class ServiceRegistry( object ):
    registered_services = {}

    @staticmethod
    def register_service( service = None ):
        pass
