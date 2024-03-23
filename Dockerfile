# Собираем в гошке
FROM golang:1.21.5 as build

ENV CODE_DIR /Applications/calendar_service/

WORKDIR ${CODE_DIR}

COPY ./ ${CODE_DIR}

RUN apt-get update
RUN apt-get -y install postgresql-client
# RUN chmod +x ./
RUN chmod +x ${CODE_DIR}wait-for-postgres.sh

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o calendar_service ./cmd/calendar/main.go





# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
#ARG LDFLAGS
#RUN CGO_ENABLED=0 go build \
#        -ldflags "$LDFLAGS" \
#       -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
#FROM alpine:3.9


#LABEL SERVICE="calendar-service"
CMD [ "./calendar-service" ]