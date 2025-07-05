-- 部門テーブル
CREATE TABLE IF NOT EXISTS departments (
    department_id VARCHAR(50) PRIMARY KEY,
    department_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 期間別組織属性テーブル
CREATE TABLE IF NOT EXISTS organization_attributes (
    department_id VARCHAR(50) NOT NULL,
    effective_date DATE NOT NULL,
    parent_department_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (department_id, effective_date),
    FOREIGN KEY (department_id) REFERENCES departments(department_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_department_id) REFERENCES departments(department_id)
);

-- インデックスの作成
CREATE INDEX idx_org_attr_parent ON organization_attributes(parent_department_id);
CREATE INDEX idx_org_attr_effective_date ON organization_attributes(effective_date);

-- 更新日時を自動更新するトリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_departments_updated_at BEFORE UPDATE
    ON departments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_organization_attributes_updated_at BEFORE UPDATE
    ON organization_attributes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- サンプルデータの投入
-- 部門マスタ
INSERT INTO departments (department_id, department_name) VALUES
('HQ', '本社'),
('STRATEGY', '経営企画部'),
('GENERAL', '総務部'),
('HR', '人事部'),
('FINANCE', '財務部'),
('SALES_HQ', '営業本部'),
('SALES_1', '営業1部'),
('SALES_2', '営業2部'),
('SALES_SUPPORT', '営業支援部'),
('TECH_HQ', '技術本部'),
('DEV', '開発部'),
('RESEARCH', '研究部'),
('QA', '品質管理部'),
('MFG_HQ', '製造本部'),
('MFG_1', '製造1部'),
('MFG_2', '製造2部'),
('IT', 'IT推進部'),
('DIGITAL', 'デジタル戦略部');

-- 2023年1月1日時点の組織構造
INSERT INTO organization_attributes (department_id, effective_date, parent_department_id) VALUES
('HQ', '2023-01-01', NULL),
('STRATEGY', '2023-01-01', 'HQ'),
('GENERAL', '2023-01-01', 'HQ'),
('HR', '2023-01-01', 'HQ'),
('FINANCE', '2023-01-01', 'HQ'),
('SALES_HQ', '2023-01-01', 'HQ'),
('SALES_1', '2023-01-01', 'SALES_HQ'),
('SALES_2', '2023-01-01', 'SALES_HQ'),
('SALES_SUPPORT', '2023-01-01', 'SALES_HQ'),
('TECH_HQ', '2023-01-01', 'HQ'),
('DEV', '2023-01-01', 'TECH_HQ'),
('RESEARCH', '2023-01-01', 'TECH_HQ'),
('QA', '2023-01-01', 'TECH_HQ'),
('MFG_HQ', '2023-01-01', 'HQ'),
('MFG_1', '2023-01-01', 'MFG_HQ'),
('MFG_2', '2023-01-01', 'MFG_HQ');

-- 2023年7月1日の組織改編（IT推進部を総務部配下に新設）
INSERT INTO organization_attributes (department_id, effective_date, parent_department_id) VALUES
('IT', '2023-07-01', 'GENERAL');

-- 2024年1月1日の組織改編（IT推進部を技術本部に移管）
INSERT INTO organization_attributes (department_id, effective_date, parent_department_id) VALUES
('IT', '2024-01-01', 'TECH_HQ');

-- 2024年4月1日の組織改編（デジタル戦略部を経営企画部配下に新設、営業支援部を独立）
INSERT INTO organization_attributes (department_id, effective_date, parent_department_id) VALUES
('DIGITAL', '2024-04-01', 'STRATEGY'),
('SALES_SUPPORT', '2024-04-01', 'HQ');

-- 2024年10月1日の組織改編（製造本部を技術本部に統合）
INSERT INTO organization_attributes (department_id, effective_date, parent_department_id) VALUES
('MFG_HQ', '2024-10-01', 'TECH_HQ'),
('MFG_1', '2024-10-01', 'MFG_HQ'),
('MFG_2', '2024-10-01', 'MFG_HQ');