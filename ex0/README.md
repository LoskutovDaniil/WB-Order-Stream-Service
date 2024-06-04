# Order Management System

## Overview
This project is an Order Management System that uses NATS Streaming for message handling and PostgreSQL for data storage. The system includes a web server that allows viewing orders stored in the database. Orders are subscribed to from a NATS Streaming channel and are stored in the database.

## Features
1. NATS Streaming Integration: Subscribes to a NATS Streaming channel to receive order data.
2. PostgreSQL Storage: Stores order data in a PostgreSQL database.
3. Cache Layer: Implements caching to improve performance.
4. Web Server: Serves a web interface to view order details.

## Access the web interface
http://localhost:8080