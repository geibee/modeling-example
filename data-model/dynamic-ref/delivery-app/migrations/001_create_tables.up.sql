-- 地域別配達方法マスタ
CREATE TABLE region_delivery_methods (
    region VARCHAR(50) NOT NULL,
    delivery_method VARCHAR(50) NOT NULL,
    standard_delivery_days INTEGER NOT NULL,
    standard_delivery_fee DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (region, delivery_method)
);

-- 都道府県マスタ
CREATE TABLE prefectures (
    prefecture VARCHAR(50) PRIMARY KEY,
    region VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 顧客マスタ
CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    customer_name VARCHAR(100) NOT NULL,
    prefecture VARCHAR(50) NOT NULL,
    delivery_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (prefecture) REFERENCES prefectures(prefecture)
);

-- インデックスの作成
CREATE INDEX idx_prefectures_region ON prefectures(region);
CREATE INDEX idx_customers_prefecture ON customers(prefecture);
CREATE INDEX idx_customers_delivery_method ON customers(delivery_method);

-- サンプルデータの挿入
INSERT INTO region_delivery_methods (region, delivery_method, standard_delivery_days, standard_delivery_fee) VALUES
('関西', '通常', 1, 300),
('中部', '通常', 3, 300),
('中部', '特急', 1, 600);

INSERT INTO prefectures (prefecture, region) VALUES
('大阪府', '関西'),
('三重県', '中部');

INSERT INTO customers (customer_name, prefecture, delivery_method) VALUES
('Aさん', '大阪府', '通常'),
('Bさん', '三重県', '通常'),
('Cさん', '三重県', '特急');