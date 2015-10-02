FROM golang:latest

WORKDIR /opt/gol
RUN  useradd --home-dir /opt/gol gol
ADD  . /opt/gol

ENV  PATH /opt/gol/bin:$PATH
RUN  cd /opt/gol;\
     make;\
     chown -R gol:users .

USER gol

CMD [ "/opt/gol/gol" ]
