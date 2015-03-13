FROM golang:1.4

WORKDIR /opt/gol
RUN  useradd --home-dir /opt/gol gol
ADD  . /opt/gol

ENV  PATH /opt/gol/bin:$PATH
RUN  cd /opt/gol;\
     make;\
     chown -R gol:users .

USER gol

# todo: change executable name
CMD [ "/opt/gol/main" ]

