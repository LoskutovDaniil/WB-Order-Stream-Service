CREATE TABLE orders (	
	id SERIAL PRIMARY KEY,
	order_uid VARCHAR(255) UNIQUE,
	delivery_id INTEGER REFERENCES delivery(id),
	payment_id INTEGER REFERENCES payment(id),
	track_number VARCHAR(255),
	entry VARCHAR(255),
	locale VARCHAR(255),
	internal_signature VARCHAR(255),
	customer_id VARCHAR(255),
	delivery_service VARCHAR(255),
	shardkey VARCHAR(255),
	sm_id INTEGER,
	date_created VARCHAR(255),
	oof_shard VARCHAR(255)
);

CREATE TABLE delivery (	
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	phone VARCHAR(255),
	zip VARCHAR(255),
	city VARCHAR(255),
	address VARCHAR(255),
	region VARCHAR(255),
	email VARCHAR(255)
);

CREATE TABLE payment (	
	id SERIAL PRIMARY KEY,
	transaction VARCHAR(255),
	request_id VARCHAR(255),
	currency VARCHAR(255),
	provider VARCHAR(255),
	amount INTEGER,
	payment_dt INTEGER,
	bank VARCHAR(255),
	delivery_cost INTEGER,
	goods_total INTEGER,
	custom_fee INTEGER
);

CREATE TABLE items (	
	id SERIAL PRIMARY KEY,
	chrt_id INTEGER,
	track_number VARCHAR(255),
	price INTEGER,
	rid VARCHAR(255),
	name VARCHAR(255),
	sale INTEGER,
	size VARCHAR(255),
	total_price INTEGER,
	nm_id INTEGER,
	brand VARCHAR(255),
	status INTEGER,
	order_uid VARCHAR(50) REFERENCES orders(order_uid)
);
