FROM redhat/ubi9

WORKDIR /app

COPY . .

RUN dnf install -y python3-pip
RUN dnf install -y python3-devel

RUN pip3 install -r requirements.txt
ENV FLASK_APP=server.py

ENTRYPOINT ["flask", "run", "--host=0.0.0.0", "--port=5002"]
EXPOSE 5002/tcp