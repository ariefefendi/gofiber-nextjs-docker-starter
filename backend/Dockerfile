# ==========================================================
# BACKEND (DEVELOPMENT): Go + Fiber + Air (hot reload)
# ==========================================================
FROM golang:1.26-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Versi Air dikunci agar build stabil. Jika dinaikkan, pastikan
# kompatibel dengan versi Go di FROM.
ARG AIR_VERSION=v1.61.7
RUN go install github.com/air-verse/air@${AIR_VERSION}

WORKDIR /app

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
# sed: buang CRLF (file dari Windows) supaya script jalan di Linux
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh \
    && chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]
