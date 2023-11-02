# Data Stream Project

## Overview
This project is a data pipeline designed to ingest data from a CSV file, stream it to Apache Kafka, store it in MySQL, and replicate it to ClickHouse using Debezium and ClickHouse's ReplicatingMergeTree for efficient querying. After setting up and running the pipeline, you can execute SQL queries on the stored data to extract meaningful insights.


## Project Structure
Here's an overview of the project's directory structure:

- `main.go`: The main entry point of the application.
- `api/`: Contains API-related code.
  - `handler.go`: Handles file upload and processing.
- `services/`: Manages database interactions.
   - `clickhouse.go`: Responsible for ClickHouse Database operations.
   - `mysql.go`: Responsible for MySQL Database operations.
   - `kafka.go`: Responsible for Kafka operations.

- `routes/`: Defines API routes.
  - `routes.go`: Configures HTTP routes for the API.
- `types/`: Contains custom data types.
  - `types.go`: Defines custom data structures and types.
- `templates/`: Stores HTML templates for the web interface.
  - `HomePage.html`: Home page template.
  - `QueryPage.html`: Query page template.
  - `ResultPage.html`: Output page template.
- `static/`: Contains static assets like CSS files.
  - `css/`: Stylesheets.
    - `HomePage.css`: Custom styles for the web interface.
    - `ResultPage.css`: Custom styles for the web interface.
- `config/`: Configuration settings for the project.
  - `db.go`: Responsible for Database configurations.
  - `kafka.go`: Responsible for Kafka configurations.
- `logs/`: Handles project logging.
  - `log.go`: Configures logging for the application.
- `dataprocess/`: Handles project logging.
  - `dataprocess.go`: Prepare the data before insertion.
   - `datageneration.go`: generate random data.


## Prerequisites
Before you begin, ensure you have met the following requirements:
- Go (Golang) installed
- Apache Kafka, MySQL, and ClickHouse set up and running
- Debezium configured for CDC (Change Data Capture) from MySQL to ClickHouse


## Installation
1. Clone the repository:
   ```bash
   git clone git@github.com:shamnas10/Go_Induvidual_Project.git
   cd your-repo
