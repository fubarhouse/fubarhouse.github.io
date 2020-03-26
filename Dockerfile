FROM jojomi/hugo as builder
MAINTAINER Karl Hepworth

#FROM alpine:latest
#COPY --from=builder /app /app
EXPOSE 1313

# Automatically build site
ONBUILD ADD site/ /app
ONBUILD RUN hugo -d /app/

# By default, serve site
ENV HUGO_BASE_URL http://localhost:1313
CMD hugo serve -d /app/ -b ${HUGO_BASE_URL} --bind=0.0.0.0

