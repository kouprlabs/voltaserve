FROM opensuse/leap:15.5

WORKDIR /app

COPY . .

RUN zypper install -y python311 python311-devel python311-pip

RUN pip3 install pipenv==2023.6.12

ENV PATH="/root/.local/bin:${PATH}"

RUN pipenv install

ENV FLASK_APP=server.py

ENTRYPOINT ["pipenv", "run", "flask", "run", "--host=0.0.0.0", "--port=5002"]

EXPOSE 5002/tcp