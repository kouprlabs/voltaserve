FROM gradle:jdk21-alpine as build-env
WORKDIR /build
COPY . .
RUN gradle assemble

FROM gradle:jdk21-alpine
WORKDIR /app
COPY --from=build-env /build/build/libs .

ENTRYPOINT ["java", "-jar", "./voltaserve-watermark-2.0.0.jar"]
EXPOSE 8086/tcp