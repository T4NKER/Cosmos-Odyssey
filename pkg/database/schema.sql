CREATE TABLE pricelists (
    id TEXT PRIMARY KEY,
    valid_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE routes (
    id TEXT PRIMARY KEY,
    pricelist_id TEXT REFERENCES pricelists(id) ON DELETE CASCADE,
    from_planet VARCHAR(255),
    to_planet VARCHAR(255),
    distance DECIMAL
);
CREATE TABLE providers (
    id TEXT PRIMARY KEY,
    route_id TEXT REFERENCES routes(id) ON DELETE CASCADE,
    company_name VARCHAR(255),
    price DECIMAL,
    flight_start TIMESTAMP,
    flight_end TIMESTAMP
);
CREATE TABLE reservations (
    id TEXT PRIMARY KEY,
    pricelist_id TEXT REFERENCES pricelists(id) ON DELETE CASCADE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    routes  VARCHAR(255),
    total_quoted_price DECIMAL,
    total_quoted_travel_time INTEGER,
    travel_companies varchar(255),
    flight_Ids TEXT
);