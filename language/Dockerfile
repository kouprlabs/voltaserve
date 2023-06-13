FROM redhat/ubi9

WORKDIR /app

COPY . .

RUN dnf install -y clang python3-pip python3-devel

RUN pip install --user --upgrade pipenv
RUN pipenv install

ENV FLASK_APP=server.py

ENTRYPOINT ["flask", "run", "--host=0.0.0.0", "--port=5002"]

EXPOSE 5002/tcp