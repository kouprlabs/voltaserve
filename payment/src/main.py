import logging.config
import os

import yaml
from fastapi import FastAPI

current_dir = os.path.dirname(__file__)
config_path = os.path.join(current_dir, '..', 'logger.yaml')

with open(config_path, 'r') as file:
    config = yaml.safe_load(file.read())
    logging.config.dictConfig(config)

app = FastAPI()