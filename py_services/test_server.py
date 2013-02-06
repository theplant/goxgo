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


# http_frontend = HttpFrontend()
# http_frontend.register_service( NLTKService,
#                                 default_verb = HttpFrontend.HTTP_POST,
#                                 method_verb_map = { 'tokenize': HttpFrontend.HTTP_GET },
#                                 method_verb_map = { HttpFrontend.HTTP_GET:[ 'tokenize', '...' ]  },
#     )
# http_frontend.start()




# from service_frontend import ZmqFrontend
# from service import Service
# from somelib import func1, func2

# SomeService = Service( name = 'SomeService' )
# SomeService.register_service_method( f = func1 )
# SomeService.register_service_method( f = func2 )

# zmq_frontend = ZmqFrontend()
# zmq_frontend.register_service( SomeService )
# zmq_frontend.start()

# http_frontend = HttpFrontend()
# http_frontend.register_service( SomeService,
#                                 default_verb = HttpFrontend.HTTP_POST,
#                                 method_verb_map = { HttpFrontend.HTTP_GET:[ 'func2', '...' ]  },
#     )
# http_frontend.start()
