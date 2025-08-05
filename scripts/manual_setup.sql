-- ===========================================
-- MANUAL DATABASE SETUP SCRIPT
-- ===========================================
-- Chạy script này trực tiếp trên PostgreSQL để tạo bảng users
-- thay vì sử dụng GORM auto-migration

-- Kết nối tới database
-- psql -U postgres -d codebase_db -f scripts/manual_setup.sql

-- ===========================================
-- USER TABLE CREATION
-- ===========================================

/* Xóa bảng cũ nếu muốn tạo lại (cẩn thận với dữ liệu) */
-- DROP TABLE IF EXISTS users CASCADE;

/* Tạo bảng users theo đúng cấu trúc GORM */
CREATE TABLE IF NOT EXISTS users (
    -- Primary key (GORM: ID uint `gorm:"primaryKey"`)
    id SERIAL PRIMARY KEY,
    
    -- Username field (GORM: Username string `gorm:"unique;not null;size:50"`)
    username VARCHAR(50) UNIQUE NOT NULL,
    
    -- Email field (GORM: Email string `gorm:"unique;not null;size:100"`)
    email VARCHAR(100) UNIQUE NOT NULL,
    
    -- Password field (GORM: Password string `gorm:"not null;size:255"`)
    password VARCHAR(255) NOT NULL,
    
    -- First name (GORM: FirstName string `gorm:"size:50"`)
    first_name VARCHAR(50),
    
    -- Last name (GORM: LastName string `gorm:"size:50"`)
    last_name VARCHAR(50),
    
    -- Active status (GORM: IsActive bool `gorm:"default:true"`)
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps (GORM tự động thêm)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Soft delete (GORM: DeletedAt gorm.DeletedAt `gorm:"index"`)
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- ===========================================
-- INDEXES (GORM tự động tạo một số index)
-- ===========================================

/* Index cho soft delete - GORM tự động tạo */
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

/* Các index bổ sung cho performance */
CREATE INDEX IF NOT EXISTS idx_users_username_active ON users(username) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_email_active ON users(email) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_active_status ON users(is_active) 
    WHERE deleted_at IS NULL;

-- ===========================================
-- AUTO UPDATE TRIGGER
-- ===========================================

/* Function để tự động update updated_at */
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

/* Trigger cho updated_at */
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ===========================================
-- CONSTRAINTS & VALIDATIONS
-- ===========================================

/* Email format validation */
ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS chk_users_email_format 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

/* Username format validation (chỉ cho phép alphanumeric và underscore) */
ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS chk_users_username_format 
    CHECK (username ~* '^[a-zA-Z0-9_]{3,50}$');

/* Password không được rỗng */
ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS chk_users_password_not_empty 
    CHECK (LENGTH(password) > 0);

-- ===========================================
-- SAMPLE DATA
-- ===========================================

/* Dữ liệu mẫu để test (password đã được hash bằng bcrypt) */
-- Password gốc là "password123" đã được hash
INSERT INTO users (username, email, password, first_name, last_name, is_active) 
VALUES 
    ('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin', 'User', true),
    ('john_doe', 'john@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'John', 'Doe', true),
    ('jane_smith', 'jane@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Jane', 'Smith', true)
ON CONFLICT (username) DO NOTHING;

-- ===========================================
-- USEFUL QUERIES
-- ===========================================

/* Kiểm tra cấu trúc bảng */
-- \d users

/* Xem tất cả users active */
-- SELECT id, username, email, first_name, last_name, is_active, created_at 
-- FROM users WHERE deleted_at IS NULL;

/* Đếm số lượng users */
-- SELECT COUNT(*) as total_users FROM users WHERE deleted_at IS NULL;

/* Test insert user mới */
-- INSERT INTO users (username, email, password, first_name, last_name) 
-- VALUES ('test_user', 'test@example.com', 'hashed_password', 'Test', 'User');

/* Test soft delete */
-- UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = 1;

/* Test restore từ soft delete */
-- UPDATE users SET deleted_at = NULL WHERE id = 1;

-- ===========================================
-- NOTES
-- ===========================================

/*
NOTES VỀ GORM VS MANUAL SETUP:

1. GORM tự động tạo:
   - Primary key với AUTO_INCREMENT
   - Index cho deleted_at (soft delete)
   - Constraints cho UNIQUE fields
   - created_at và updated_at timestamps

2. Differences:
   - GORM sử dụng snake_case cho column names
   - GORM tự động handle soft delete với WHERE deleted_at IS NULL
   - GORM tự động update updated_at khi Save/Update

3. Để tương thích với GORM:
   - Không disable auto-migration nếu muốn GORM tự tạo
   - Hoặc sử dụng script này và set AutoMigrate = false
   - Column names phải match với GORM tags

4. Recommedations:
   - Sử dụng GORM auto-migration cho development
   - Sử dụng manual script cho production
   - Luôn backup trước khi chạy migration
*/