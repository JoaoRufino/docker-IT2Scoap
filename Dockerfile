FROM maven:3-jdk-8
WORKDIR /app/californium
CMD mvn clean install -X