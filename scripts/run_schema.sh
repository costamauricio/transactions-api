#!/bin/sh
sqlite3 /app/db/transactions.db < /app/scripts/database_tables.sql
