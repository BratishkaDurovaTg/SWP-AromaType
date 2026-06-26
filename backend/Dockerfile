FROM golang:1.23-alpine AS build

WORKDIR /src/backend

COPY backend/go.mod ./
COPY backend/ ./

RUN go build -o /out/aromatype-api ./cmd/api \
	&& go build -o /out/aromatype-catalog-bot ./cmd/catalogbot

FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/aromatype-api /app/aromatype-api
COPY --from=build /out/aromatype-catalog-bot /app/aromatype-catalog-bot
COPY docs/ /app/docs/

EXPOSE 8080

CMD ["/app/aromatype-api"]
