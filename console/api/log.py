import logging

from .dependencies import settings

if settings.LOG_FORMAT == 'JSON':
    req_fmt = ('{"timestamp":"%(asctime)s",'
               '"logger_name":"%(name)s",'
               '"log_level":"%(levelname)s",'
               '"message":"%(message)s",'
               '"type":"%(type)s",'
               '"identifier":"%(identifier)s",'
               '"path":"%(path)s",'
               '"method":"%(method)s",'
               '"headers":"%(headers)s",'
               '"query_params":"%(query_params)s",'
               '"path_params":"%(path_params)s"}')
    resp_fmt = ('{"timestamp":"%(asctime)s",'
                '"logger_name":"%(name)s",'
                '"log_level":"%(levelname)s",'
                '"message":"%(message)s",'
                '"type":"%(type)s",'
                '"identifier":"%(identifier)s",'
                '"code":"%(code)s",'
                '"headers":"%(headers)s"'
                )
    base_fmt = ('{"timestamp":"%(asctime)s",'
                '"logger_name":"%(name)s",'
                '"log_level":"%(levelname)s",'
                '"message":"%(message)s"}'
                )
elif settings.LOG_FORMAT == "PLAIN":
    req_fmt = ('%(asctime)s|%(name)s|%(levelname)s|%(type)s|%(identifier)s|%(path)s|'
               '%(method)s|%(headers)s|%(query_params)s|%(path_params)s|%(message)s')
    resp_fmt = '%(asctime)s|%(name)s|%(levelname)s|%(type)s|%(identifier)s|%(code)s|%(headers)s|%(message)s'
    base_fmt = '%(asctime)s|%(name)s|%(levelname)s|%(message)s'
else:
    raise ValueError('Wrong logging format, available JSON and PLAIN')

logger = logging.getLogger('console.api')
logger.setLevel(settings.LOG_LEVEL)


base_handler = logging.StreamHandler()
base_handler.setFormatter(logging.Formatter(base_fmt))
base_logger = logger.getChild('router')
base_logger.addHandler(base_handler)

req_handler = logging.StreamHandler()
req_handler.setFormatter(logging.Formatter(req_fmt))

resp_handler = logging.StreamHandler()
resp_handler.setFormatter(logging.Formatter(resp_fmt))

req_logger = logger.getChild('requests')
req_logger.addHandler(req_handler)

resp_logger = logger.getChild('responses')
resp_logger.addHandler(resp_handler)
