-- Schema cơ bản cho backend-go (SQL Server)
-- Giữ nguyên cấu trúc bảng tương ứng với project Spring Boot gốc.

CREATE TABLE users (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    email NVARCHAR(255) NOT NULL UNIQUE,
    password_hash NVARCHAR(500) NOT NULL,
    active BIT NOT NULL DEFAULT 1,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE refresh_tokens (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER NOT NULL REFERENCES users(id),
    token_hash NVARCHAR(500) NOT NULL,
    expires_at DATETIME2 NOT NULL,
    revoked BIT NOT NULL DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE user_profiles (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER NOT NULL UNIQUE REFERENCES users(id),
    display_name NVARCHAR(255),
    avatar_url NVARCHAR(500),
    guitar_type NVARCHAR(50),
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE courses (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    title_en NVARCHAR(500),
    title_vi NVARCHAR(500),
    slug NVARCHAR(255),
    description_en NVARCHAR(MAX),
    description_vi NVARCHAR(MAX),
    guitar_type NVARCHAR(50),
    level NVARCHAR(50),
    sort_order INT DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE modules (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    course_id UNIQUEIDENTIFIER NOT NULL REFERENCES courses(id),
    title_en NVARCHAR(500),
    title_vi NVARCHAR(500),
    description_en NVARCHAR(MAX),
    description_vi NVARCHAR(MAX),
    sort_order INT DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE lessons (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    module_id UNIQUEIDENTIFIER NOT NULL REFERENCES modules(id),
    title_en NVARCHAR(500),
    title_vi NVARCHAR(500),
    content_en NVARCHAR(MAX),
    content_vi NVARCHAR(MAX),
    video_url NVARCHAR(500),
    sort_order INT DEFAULT 0,
    duration_min INT DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE learning_methods (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name_en NVARCHAR(500),
    name_vi NVARCHAR(500),
    description_en NVARCHAR(MAX),
    description_vi NVARCHAR(MAX),
    difficulty NVARCHAR(50),
    sort_order INT DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE lesson_blocks (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    lesson_id UNIQUEIDENTIFIER NOT NULL REFERENCES lessons(id),
    block_type NVARCHAR(50),
    sort_order INT DEFAULT 0,
    media_asset_id NVARCHAR(500),
    metadata_json NVARCHAR(MAX),
    content_en NVARCHAR(MAX),
    content_vi NVARCHAR(MAX),
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE chords (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name NVARCHAR(255),
    guitar_type NVARCHAR(50),
    category NVARCHAR(50),
    description NVARCHAR(MAX)
);

CREATE TABLE chord_positions (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    chord_id UNIQUEIDENTIFIER NOT NULL REFERENCES chords(id),
    position_number INT,
    base_fret INT,
    fingers_json NVARCHAR(MAX)
);

CREATE TABLE strumming_patterns (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name NVARCHAR(500),
    formula NVARCHAR(MAX),
    description NVARCHAR(MAX),
    difficulty NVARCHAR(50),
    guitar_type NVARCHAR(50),
    bpm_min INT,
    bpm_max INT,
    is_formula BIT DEFAULT 0,
    status NVARCHAR(50) DEFAULT 'PUBLISHED',
    sort_order INT DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE strumming_instructions (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    pattern_id UNIQUEIDENTIFIER NOT NULL REFERENCES strumming_patterns(id),
    step_number INT,
    instruction_en NVARCHAR(MAX),
    instruction_vi NVARCHAR(MAX)
);

CREATE TABLE lesson_progress (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER NOT NULL REFERENCES users(id),
    lesson_id UNIQUEIDENTIFIER NOT NULL REFERENCES lessons(id),
    status NVARCHAR(50) DEFAULT 'IN_PROGRESS',
    completed_at DATETIME2,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE songs (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    title_en NVARCHAR(500),
    title_vi NVARCHAR(500),
    artist NVARCHAR(500),
    source NVARCHAR(100),
    difficulty NVARCHAR(50),
    language NVARCHAR(10),
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE quizzes (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    lesson_id UNIQUEIDENTIFIER NOT NULL REFERENCES lessons(id),
    title_en NVARCHAR(500),
    title_vi NVARCHAR(500)
);

CREATE TABLE quiz_attempts (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    quiz_id UNIQUEIDENTIFIER NOT NULL REFERENCES quizzes(id),
    user_id UNIQUEIDENTIFIER NOT NULL REFERENCES users(id),
    score INT DEFAULT 0,
    status NVARCHAR(50) DEFAULT 'IN_PROGRESS',
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);

CREATE TABLE user_preferences (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER NOT NULL UNIQUE REFERENCES users(id),
    guitar_type NVARCHAR(50),
    completed BIT DEFAULT 0,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);
