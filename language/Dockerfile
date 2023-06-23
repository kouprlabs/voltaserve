FROM opensuse/leap:15.5

WORKDIR /app

COPY . .

# Install dependencies required for pyenv and Python build
RUN zypper --non-interactive install git curl make gcc gcc-c++ glibc-devel libopenssl-devel zlib-devel bzip2 libbz2-devel readline-devel libffi-devel sqlite3 sqlite3-devel tk-devel xz xz-devel tar

# Install pyenv
RUN curl -L https://github.com/pyenv/pyenv-installer/raw/master/bin/pyenv-installer | bash

# Set up pyenv in .bashrc
RUN echo 'export PYENV_ROOT="/root/.pyenv"' >> /root/.bashrc \
    && echo 'export PATH="$PYENV_ROOT/bin:$PATH"' >> /root/.bashrc \
    && echo 'eval "$(pyenv init --path)"' >> /root/.bashrc \
    && echo 'eval "$(pyenv virtualenv-init -)"' >> /root/.bashrc

# Install Python 3.9.14 and set it as global
RUN /bin/bash -c "source /root/.bashrc && pyenv install 3.9.14 && pyenv global 3.9.14"

# Install Pipenv
RUN /bin/bash -c "source /root/.bashrc && pip install pipenv"

# Install setuptools and wheel
RUN /bin/bash -c "source /root/.bashrc && pip install setuptools wheel"

# Install project dependencies with Pipenv
RUN /bin/bash -c "source /root/.bashrc && pipenv install --system"

# Download Spacy model
RUN /bin/bash -c "source /root/.bashrc && python -m spacy download xx_ent_wiki_sm"

ENV FLASK_APP=server.py

ENTRYPOINT ["/bin/bash", "-c", "source /root/.bashrc && pipenv run flask run --host=0.0.0.0 --port=5002"]

EXPOSE 5002/tcp