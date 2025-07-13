-- 権限レベルを追加するためのテーブル変更

-- ROLE_SCREEN_PERMISSIONテーブルに権限レベルカラムを追加
ALTER TABLE ROLE_SCREEN_PERMISSION 
ADD COLUMN IF NOT EXISTS permission_level VARCHAR(20) DEFAULT 'VIEWER';

-- 既存のデータを更新（管理職と管理者にはEDITOR権限を付与）
UPDATE ROLE_SCREEN_PERMISSION 
SET permission_level = 'EDITOR' 
WHERE ROLE_ID IN ('R01', 'R04');

-- 権限レベルの制約を追加
ALTER TABLE ROLE_SCREEN_PERMISSION 
ADD CONSTRAINT check_screen_permission_level 
CHECK (permission_level IN ('VIEWER', 'EDITOR'));