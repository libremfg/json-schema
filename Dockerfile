# Stage 1
FROM alpine:latest AS build

# Install the Hugo go app.
RUN apk add --update hugo jq

# Set Working directory
WORKDIR /opt/b2mml-batchml

# Copy Hugo config into the container Workdir.
COPY . .

# Set the Host URL for the JSON Schema files & Hugo build.
ARG HOST=http://localhost/

# Copy in the schema files & replace JSON Schema Location
RUN mkdir ./Docs/static/schemas && \
    cp ./Schema/* ./Docs/static/schemas/ && \
    sed -i 's,\("$id"\: *"\),\1'"$HOST"'schemas/,g' ./Docs/static/schemas/*.json && \
    sed -i 's,\("$ref"\: *"\)\./,\1'"$HOST"'schemas/,g' ./Docs/static/schemas/*.json

# Minimize the JSON Schema files.
RUN for f in ./Docs/static/schemas/*.json; do \
path=${f%.json}.min.json && \
(jq -c . < ./$f) > $path; \
done

# Run Hugo in the Workdir to generate HTML.
RUN hugo -s ./Docs -b ${HOST}

# Stage 2
FROM nginx:1.25-alpine

# Set workdir to the NGINX default dir.
WORKDIR /usr/share/nginx/html

# Copy HTML from previous build into the Workdir.
COPY --from=build /opt/b2mml-batchml/Docs/public .
COPY --from=build /opt/b2mml-batchml/Config/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80/tcp