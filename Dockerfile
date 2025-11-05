FROM golang:1.24.2-alpine AS build-stage

WORKDIR /app

RUN go install github.com/swaggo/swag/v2/cmd/swag@latest

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all go files
COPY . .

# Set build tag
ARG BUILD_TAG

ENV GOCACHE=/root/.cache/go-build

RUN if [[ "${BUILD_TAG}" == "debug" ]]; \
    then swag i -g main.go --v3.1 ; fi
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=linux go build \
    --ldflags="-s -w" \
    --tags ${BUILD_TAG} \
    -buildvcs=false \
    -o /app/go-gin-template \
    /app/main.go


FROM scratch AS production-stage

WORKDIR /app
COPY --from=build-stage /app/go-gin-template /app/go-gin-template

CMD [ "./go-gin-template" ]