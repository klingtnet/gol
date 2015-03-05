FROM golang:1.4

WORKDIR /opt/gol
RUN  useradd --home-dir /opt/gol gol
ADD  . /opt/gol

RUN  cd /opt/gol;\
     make;\
     chown -R gol:users .

USER gol

# todo: change executable name
CMD [ "/opt/gol/main" ]

