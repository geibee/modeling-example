-- 品目基本属性テーブル
CREATE TABLE IF NOT EXISTS 品目基本属性 (
    品目ID VARCHAR(10) PRIMARY KEY,
    品目名 VARCHAR(100) NOT NULL,
    品種区分 VARCHAR(10) NOT NULL CHECK (品種区分 IN ('A', 'B')),
    品目コード VARCHAR(20) NOT NULL UNIQUE
);

-- A品種品目属性テーブル
CREATE TABLE IF NOT EXISTS A品種品目属性 (
    品目ID VARCHAR(10) PRIMARY KEY,
    容量 DECIMAL(10, 2),
    材質 VARCHAR(50),
    FOREIGN KEY (品目ID) REFERENCES 品目基本属性(品目ID) ON DELETE CASCADE
);

-- B品種品目属性テーブル
CREATE TABLE IF NOT EXISTS B品種品目属性 (
    品目ID VARCHAR(10) PRIMARY KEY,
    内径 DECIMAL(10, 2),
    外径 DECIMAL(10, 2),
    FOREIGN KEY (品目ID) REFERENCES 品目基本属性(品目ID) ON DELETE CASCADE
);

-- インデックスの作成
CREATE INDEX idx_品目コード ON 品目基本属性(品目コード);
CREATE INDEX idx_品種区分 ON 品目基本属性(品種区分);

-- テストデータの挿入
-- A品種のデータ
INSERT INTO 品目基本属性 (品目ID, 品目名, 品種区分, 品目コード) VALUES
('A001', 'プラスチックボトル500ml', 'A', 'PBOT-500'),
('A002', 'ガラスボトル1000ml', 'A', 'GBOT-1000'),
('A003', 'アルミ缶350ml', 'A', 'ACAN-350');

INSERT INTO A品種品目属性 (品目ID, 容量, 材質) VALUES
('A001', 500.00, 'PET'),
('A002', 1000.00, 'ガラス'),
('A003', 350.00, 'アルミニウム');

-- B品種のデータ
INSERT INTO 品目基本属性 (品目ID, 品目名, 品種区分, 品目コード) VALUES
('B001', 'ステンレスパイプ20mm', 'B', 'SPIPE-20'),
('B002', '塩ビパイプ50mm', 'B', 'VPIPE-50'),
('B003', '銅管15mm', 'B', 'CPIPE-15');

INSERT INTO B品種品目属性 (品目ID, 内径, 外径) VALUES
('B001', 15.00, 20.00),
('B002', 44.00, 50.00),
('B003', 13.00, 15.00);