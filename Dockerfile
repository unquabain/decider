FROM golang:1.22-alpine AS build
ENV CGOENABLED=0

RUN mkdir /build
WORKDIR /build
COPY . /build
RUN go mod tidy
RUN go generate
RUN go build --tags server

FROM scratch

COPY --from=build /build/decider /decider
EXPOSE 8899
CMD ["/decider", "serve"]
