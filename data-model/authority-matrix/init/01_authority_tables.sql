-- 権限管理システムのテーブル作成（小文字版）

-- 基本テーブル（外部キー依存なし）

-- ロールテーブル
CREATE TABLE IF NOT EXISTS role (
    role_id VARCHAR(3) PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL
);

-- 画面テーブル
CREATE TABLE IF NOT EXISTS screen (
    screen_id VARCHAR(3) PRIMARY KEY,
    screen_name VARCHAR(100) NOT NULL
);

-- APIテーブル
CREATE TABLE IF NOT EXISTS api (
    api_id VARCHAR(3) PRIMARY KEY,
    api_name VARCHAR(100) NOT NULL
);

-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS "user" (
    user_id VARCHAR(10) PRIMARY KEY
);

-- 組織テーブル
CREATE TABLE IF NOT EXISTS organization (
    org_id VARCHAR(3) PRIMARY KEY,
    org_name VARCHAR(100) NOT NULL
);

-- 依存テーブル（外部キーあり）

-- 予算テーブル
CREATE TABLE IF NOT EXISTS budget (
    id TEXT PRIMARY KEY,
    department_id TEXT NOT NULL,
    org_id VARCHAR(3) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES organization(org_id)
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_budget_department_id ON budget(department_id);
CREATE INDEX IF NOT EXISTS idx_budget_org_id ON budget(org_id);

-- メニュー設定テーブル
CREATE TABLE IF NOT EXISTS menu_setting (
    menu_id VARCHAR(3) PRIMARY KEY,
    menu_name VARCHAR(100) NOT NULL,
    screen_id VARCHAR(3) NOT NULL,
    FOREIGN KEY (screen_id) REFERENCES screen(screen_id)
);

-- API画面マッピングテーブル
CREATE TABLE IF NOT EXISTS api_screen_mapping (
    screen_id VARCHAR(3) NOT NULL,
    api_id VARCHAR(3) NOT NULL,
    PRIMARY KEY (screen_id, api_id),
    FOREIGN KEY (screen_id) REFERENCES screen(screen_id),
    FOREIGN KEY (api_id) REFERENCES api(api_id)
);

-- ユーザーロールテーブル
CREATE TABLE IF NOT EXISTS user_role (
    user_id VARCHAR(10) NOT NULL,
    role_id VARCHAR(3) NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES "user"(user_id),
    FOREIGN KEY (role_id) REFERENCES role(role_id)
);

-- ロール画面権限テーブル
CREATE TABLE IF NOT EXISTS role_screen_permission (
    role_id VARCHAR(3) NOT NULL,
    screen_id VARCHAR(3) NOT NULL,
    PRIMARY KEY (role_id, screen_id),
    FOREIGN KEY (role_id) REFERENCES role(role_id),
    FOREIGN KEY (screen_id) REFERENCES screen(screen_id)
);

-- ユーザー組織所属テーブル
CREATE TABLE IF NOT EXISTS user_organization (
    user_id VARCHAR(10) NOT NULL,
    org_id VARCHAR(3) NOT NULL,
    PRIMARY KEY (user_id, org_id),
    FOREIGN KEY (user_id) REFERENCES "user"(user_id),
    FOREIGN KEY (org_id) REFERENCES organization(org_id)
);

-- 販売管理テーブル
CREATE TABLE IF NOT EXISTS sales (
    id VARCHAR(10) PRIMARY KEY,
    org_id VARCHAR(3) NOT NULL,
    product VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES organization(org_id)
);

-- 販売管理インデックス
CREATE INDEX IF NOT EXISTS idx_sales_org_id ON sales(org_id);
CREATE INDEX IF NOT EXISTS idx_sales_date ON sales(date);

-- 在庫管理テーブル
CREATE TABLE IF NOT EXISTS inventory (
    code VARCHAR(10) PRIMARY KEY,
    org_id VARCHAR(3) NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    updated DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES organization(org_id)
);

-- 在庫管理インデックス
CREATE INDEX IF NOT EXISTS idx_inventory_org_id ON inventory(org_id);
CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory(status);

-- サンプルデータの挿入

-- ロールデータ
INSERT INTO role (role_id, role_name) VALUES
('R01', '管理職'),
('R02', '営業担当'),
('R03', '在庫管理者'),
('R04', 'システム管理者')
ON CONFLICT (role_id) DO NOTHING;

-- 組織データ
INSERT INTO organization (org_id, org_name) VALUES
('001', '本社'),
('010', '北日本支店'),
('020', '東京支店'),
('999', 'その他部門')
ON CONFLICT (org_id) DO NOTHING;

-- ユーザーデータ
INSERT INTO "user" (user_id) VALUES
('user001'),
('user002'),
('user003'),
('admin')
ON CONFLICT (user_id) DO NOTHING;

-- ユーザー組織所属データ
INSERT INTO user_organization (user_id, org_id) VALUES
('user001', '001'),    -- 管理職 - 本社
('user001', '010'),    -- 管理職 - 北日本支店
('user002', '020'),    -- 営業担当 - 東京支店
('user003', '001'),    -- 在庫管理者 - 本社
('admin', '001')       -- システム管理者 - 本社
ON CONFLICT (user_id, org_id) DO NOTHING;

-- 画面データ
INSERT INTO screen (screen_id, screen_name) VALUES
('001', '予算管理画面'),
('002', '販売管理画面'),
('003', '在庫管理画面'),
('004', '権限管理画面'),
('005', 'ユーザー管理画面')
ON CONFLICT (screen_id) DO NOTHING;

-- APIデータ
INSERT INTO api (api_id, api_name) VALUES
('100', '予算管理API'),
('200', '販売管理API'),
('300', '在庫管理API'),
('400', '権限管理API'),
('500', 'ユーザー管理API')
ON CONFLICT (api_id) DO NOTHING;

-- メニュー設定データ
INSERT INTO menu_setting (menu_id, menu_name, screen_id) VALUES
('M01', '予算管理', '001'),
('M02', '販売管理', '002'),
('M03', '在庫管理', '003'),
('M04', '権限管理', '004'),
('M05', 'ユーザー管理', '005')
ON CONFLICT (menu_id) DO NOTHING;

-- API画面マッピングデータ
INSERT INTO api_screen_mapping (screen_id, api_id) VALUES
('001', '100'),
('002', '200'),
('003', '300'),
('004', '400'),
('005', '500')
ON CONFLICT (screen_id, api_id) DO NOTHING;

-- ユーザーロールデータ
INSERT INTO user_role (user_id, role_id) VALUES
('user001', 'R01'),    -- 管理職
('user002', 'R02'),    -- 営業担当
('user003', 'R03'),    -- 在庫管理者
('admin', 'R04')       -- システム管理者
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ロール画面権限データ
INSERT INTO role_screen_permission (role_id, screen_id) VALUES
('R01', '001'),    -- 管理職 → 予算管理画面
('R02', '002'),    -- 営業担当 → 販売管理画面
('R03', '003'),    -- 在庫管理者 → 在庫管理画面
('R04', '001'),    -- システム管理者 → すべての画面
('R04', '002'),
('R04', '003'),
('R04', '004'),
('R04', '005')
ON CONFLICT (role_id, screen_id) DO NOTHING;

-- 予算サンプルデータ
INSERT INTO budget (id, department_id, org_id, amount, description) VALUES 
('B001', '001', '001', 1000000.00, '本社予算'),
('B002', '010', '010', 500000.00, '北日本支店予算'),
('B003', '020', '020', 400000.00, '東京支店予算'),
('B004', '999', '999', 200000.00, 'その他部門予算')
ON CONFLICT (id) DO NOTHING;

-- 販売サンプルデータ
INSERT INTO sales (id, org_id, product, quantity, amount, date) VALUES
('S001', '001', 'ノートPC', 2, 300000, '2024-01-15'),
('S002', '020', 'モニター', 5, 250000, '2024-01-16'),
('S003', '020', 'キーボード', 10, 80000, '2024-01-17'),
('S004', '010', 'マウス', 15, 60000, '2024-01-18'),
('S005', '001', 'USBケーブル', 20, 20000, '2024-01-19')
ON CONFLICT (id) DO NOTHING;

-- 在庫サンプルデータ
INSERT INTO inventory (code, org_id, name, category, stock, status, updated) VALUES
('P001', '001', 'ノートPC', '電子機器', 15, 'in-stock', '2024-01-20'),
('P002', '020', 'モニター', '電子機器', 3, 'low-stock', '2024-01-19'),
('P003', '020', 'キーボード', '周辺機器', 0, 'out-of-stock', '2024-01-18'),
('P004', '010', 'マウス', '周辺機器', 25, 'in-stock', '2024-01-20'),
('P005', '001', 'USBケーブル', 'アクセサリ', 2, 'low-stock', '2024-01-17')
ON CONFLICT (code) DO NOTHING;