FROM registry.suse.com/bci/python:3.11

WORKDIR /app

COPY . .

RUN zypper install -y gcc-c++

RUN pip3 install micropipenv

RUN micropipenv install

RUN python3 -m spacy download xx_ent_wiki_sm

ENV FLASK_APP=server.py

ENTRYPOINT ["flask", "run", "--host=0.0.0.0", "--port=5002"]

EXPOSE 5002/tcp
