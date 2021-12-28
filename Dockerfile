############################
# STEP 1 build executable binary
############################
FROM golang:1.17-buster@sha256:9cf49bf1e8cf0c8d995899546bea5d876e59b5a9fdeb1575cb6197f97973c20e AS builder

# Create appuser.
ENV USER=appuser
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/ github.com/raysandeep/Estimator-App
COPY . .

# # Fetch dependencies.
# RUN go mod download
# RUN go mod verify

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/main

############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/main /go/bin/main

# Use an unprivileged user.
USER appuser:appuser

# Run the binary.
ENTRYPOINT ["/go/bin/main"]