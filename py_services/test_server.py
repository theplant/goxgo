# -*- coding: utf-8 -*-

import gevent
from gevent import monkey; monkey.patch_all()

from service_frontend import ZmqFrontend
from service import Service
from lib.static_nltk_wrappers import tokenize, stem, vsm_compare

# name a service and add some functions that it provides
NLTKService = Service( name = 'NLTK' )
NLTKService.register_service_method( f = tokenize )
NLTKService.register_service_method( f = stem )
NLTKService.register_service_method( f = vsm_compare )

frontends = []
# register that service with a (netowrking) frontend
zmq_frontend = ZmqFrontend()
zmq_frontend.register_service( NLTKService )
frontends.append( gevent.spawn( zmq_frontend.start ) )

# this is how it should work for HTTP/REST
# http_frontend = HttpFrontend()
# http_frontend.register_service( NLTKService,
#                                 default_verb = HttpFrontend.HTTP_POST,
#                                 method_verb_map = { 'tokenize': HttpFrontend.HTTP_GET },
#                                 method_verb_map = { HttpFrontend.HTTP_GET:[ 'tokenize', '...' ]  },
#     )
# http_frontend.start()
# frontends.append( gevent.spawn( http_frontend.start ) )

gevent.joinall( frontends )
