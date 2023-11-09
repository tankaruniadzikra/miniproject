CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    deposit_amount DECIMAL
);

CREATE TABLE equipments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    availability INT NOT NULL,
    daily_rental_cost DECIMAL NOT NULL,
    category VARCHAR(255) NOT NULL
);

CREATE TABLE rental_histories (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    equipment_id INT NOT NULL,
    rental_date DATE NOT NULL, 
    return_date DATE, 
    total_cost DECIMAL, 
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (equipment_id) REFERENCES equipments(id)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    rental_history_id INT NOT NULL,
    payment_date DATE NOT NULL,
    is_deposit BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (rental_history_id) REFERENCES rental_histories(id)
);

--
INSERT INTO equipments (name, availability, daily_rental_cost, category)
VALUES 
    ('Bulldozer', 3, 10.50, 'Heavy Machinery'),
    ('Excavator', 2, 12.75, 'Construction Equipment'),
    ('Loader', 4, 15.00, 'Heavy Machinery'),
    ('Crane', 7, 8.99, 'Construction Equipment');