# Voltaserve Language

## Getting Started

Install [Pipenv](https://pipenv.pypa.io/en/latest/installation/#installing-pipenv).

Install `clang` and `python3-devel`. For RHEL based operating systems, run:

```shell
sudo dnf install -y clang python3-pip python3-devel
```

Install dependencies:

```shell
pipenv install
```

Run:

```shell
pipenv shell
FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug
```
