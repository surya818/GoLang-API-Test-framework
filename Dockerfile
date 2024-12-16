# Copyright © 2024 Kong Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Start with an official Golang image for building
FROM golang:1.23 AS builder
WORKDIR /build

# Cache dependencies for build
COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
ARG APP_VERSION
ARG APP_COMMIT
ARG APP_OS_ARCH
ARG APP_GO_VERSION
ARG APP_DATE_FORMAT
ARG APP_BUILD_DATE
ARG APP_PACKAGE
ENV APP_DOCKER_BUILD=true
RUN make build
RUN ls /build

# Use a minimal Docker image for the actual runtime (Alpine Linux)
FROM golang:1.23
WORKDIR /app

# Copy the executable and config.yml to the final container
COPY --from=builder /build/bin/candidate-take-home-exercise-sdet /app/candidate-take-home-exercise-sdet
RUN chmod a+x /app/candidate-take-home-exercise-sdet

# Start the application
CMD ["/app/candidate-take-home-exercise-sdet"]
