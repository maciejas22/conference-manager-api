# Conference Manager API

This repository hosts the backend API of the Conference Management System, developed to facilitate seamless management and organization of conferences.

## Overview

The API is structured as a microservices architecture using modern web technologies. It efficiently handles business logic, user authentication, payments, and data management.

## Implementation and Architecture

The API follows a microservices architecture for modular and scalable backend management. The architecture is designed to support expanding functionalities with ease, leveraging the efficiency of Go for service implementation.

-  **Microservices Architecture**: Composed of distinct services handling authentication, conference management, and information retrieval.
-  **GraphQL Gateway**: Serves as a single entry point for client interactions, consolidating data requests across services.
-  **Security and Session Management**: Employs OAuth tokens and session management for secure interactions.
-  **Database Design**: Utilizes PostgreSQL with a focus on normalized data storage for efficient querying and data integrity.
-  **External Service Integration**: Integrates with Stripe for payments and Amazon S3 for file management.

## Features

-    User authentication and session management
-    Conference creation, update, and management
-    Ticket handling and validation
-    Secure payment processing using external services

![be-diagram](https://github.com/user-attachments/assets/0bb2a785-a553-4e76-a744-c78faa2d5cf4)


## Technologies

-    **Go**
-    **GraphQL**
-    **gRPC**
-    **PostgreSQL**
-    **Amazon S3**
-    **Stripe**

## Setup Instructions

1. **Clone Repository**:
   ```bash
   git clone https://github.com/maciejas22/conference-manager-api.git
   cd conference-manager-api
   
2. **Create docker network**:
   ```bash
   docker network create cm-network

3. ** Start Databases via Docker
   ```bash
   docker-compose up -d cm-conferences-db cm-auth-db cm-info-db localstack

4. **Apply Migrations**:
   ```bash
   make -C cm-auth migrate-up
   make -C cm-conferences migrate-up
   make -C cm-info migrate-up

5. **Start Backend Services**:
   ```bash
    docker-compose up -d cm-conferences-api cm-auth-api cm-info-api
