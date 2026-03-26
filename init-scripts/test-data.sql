-- Test data for Sync-OUT testing
-- Insert into Oracle (SLAVE), then sync-out to PostgreSQL (MASTER)

-- Customer (OUT table) - This should trigger the log and sync-out
INSERT INTO customer (store_id, first_name, last_name, email, address_id, activebool, create_date, last_update, active)
VALUES (1, 'John', 'Doe', 'john@test.com', 1, true, CURRENT_DATE, CURRENT_TIMESTAMP, 1);

-- Rental (OUT table)
INSERT INTO rental (rental_date, inventory_id, customer_id, return_date, staff_id, last_update)
VALUES (CURRENT_TIMESTAMP, 1, 1, NULL, 1, CURRENT_TIMESTAMP);

-- Payment (OUT table)  
INSERT INTO payment (customer_id, staff_id, rental_id, amount, payment_date, last_update)
VALUES (1, 1, 1, 9.99, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
