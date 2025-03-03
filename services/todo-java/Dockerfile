# Use an official OpenJDK image
FROM eclipse-temurin:17-jdk

# Set the working directory inside the container
WORKDIR /app

# Copy the Gradle wrapper and project files
COPY gradlew gradlew
COPY gradle gradle
COPY build.gradle settings.gradle ./
COPY src src

# Give execution permissions to the Gradle wrapper
RUN chmod +x gradlew

# Build the application
RUN ./gradlew clean build -x test

EXPOSE 3000

# Run the application
CMD ["java", "-jar", "build/libs/todo-java-0.0.1-SNAPSHOT.jar"]