# -*- coding: utf-8 -*-

# Copyright 2013 Florian von Bock (f at theplant dot jp)
#
# ZMQ frontend for a service


from gevent import spawn
import zmq.green as zmq


# import os
import sys
from time import time
import traceback
import json

import service


class ServiceFrontendError( Exception ):
    # TODO: refactor exceptions...
    def __init__( self, code = None, reason = None, description = None,
                  details = {} ):
        self.code = code
        self.reason = reason
        self.description = description
        self.details = details

    def __str__( self ):
        return str( self.description )



class ServiceFrontend( object ):
    """
    Base class for service frontends
    """
    def __init__( self, proto = None, host = None, port = None ):
        self.dsn = '{}://{}:{}'.format( proto, host, port )
        self.server = None
        self.services = {}
        # TODO: signal handling to stop - maybe reload?

    def register_service( self, service ):
        self.services[ service.name ] = service



class ZmqFrontend( ServiceFrontend ):
    """
    RPC style frontend for arbitrary services using ZMQ as transport
    """
    def __init__( self, proto = 'tcp', host = '*', port = 4242 ):
        super( ZmqFrontend, self ).__init__( proto, host, port )
        # server
        self.context = zmq.Context()
        self.socket = self.context.socket( zmq.REP )
        # self.socket
        self.socket.bind( self.dsn )

    def start( self ):
        print "Starting server. Listening on {}...".format( self.dsn )
        self.server = spawn( self._serve )
        self.server.join()

    def stop( self ):
        self.server.kill()

    def _serve( self ):
        """the main server loop"""
        print "Start serving."
        while True:
            message = self.socket.recv()
            start_time = time()
            try:
                msg_data = json.loads( message )
                # print msg_data
                # print "payload keys: {}".format( msg_data.keys() )
                # TODO: pipelining NOT yet implemented...
                service_name, f_name = msg_data['Target']['Services'][0].split('/')
                response = json.dumps(
                    self.services[ service_name ].service_methods[ f_name ](
                        **{ k.lower(): msg_data[k] for k in msg_data.keys() if k not in [ 'Target' ] } ) )
            except Exception, e:
                traceback.print_exc( file = sys.stdout )
                response = "shit hit the fan: {}".format( e )
            self.socket.send( response )
            print "req took: {}".format( time() - start_time )

    def register_service( self, s ):
        if not isinstance( s, service.Service ):
            raise ServiceFrontendError( reason = "Could not retister {}".format( s.__name__ ),
                                        description = "Can only register service.Service types." )
        self.services[ s.name ] = s


class HttpFrontend( ServiceFrontend ):
    """
    HTTP REST style frontend for arbitrary services
    """
    pass
