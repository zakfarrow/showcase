-- Projects table schema
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    tech_stack TEXT[] NOT NULL DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'completed',
    features TEXT[] NOT NULL DEFAULT '{}',
    challenges TEXT NOT NULL DEFAULT '',
    learnings TEXT NOT NULL DEFAULT '',
    future_plans TEXT NOT NULL DEFAULT '',
    github_url VARCHAR(512),
    live_url VARCHAR(512),
    image_url VARCHAR(512),
    featured BOOLEAN NOT NULL DEFAULT false,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add columns if they don't exist (for existing databases)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'projects' AND column_name = 'status') THEN
        ALTER TABLE projects ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'completed';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'projects' AND column_name = 'features') THEN
        ALTER TABLE projects ADD COLUMN features TEXT[] NOT NULL DEFAULT '{}';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'projects' AND column_name = 'challenges') THEN
        ALTER TABLE projects ADD COLUMN challenges TEXT NOT NULL DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'projects' AND column_name = 'learnings') THEN
        ALTER TABLE projects ADD COLUMN learnings TEXT NOT NULL DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'projects' AND column_name = 'future_plans') THEN
        ALTER TABLE projects ADD COLUMN future_plans TEXT NOT NULL DEFAULT '';
    END IF;
END $$;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);
CREATE INDEX IF NOT EXISTS idx_projects_featured ON projects(featured) WHERE featured = true;
CREATE INDEX IF NOT EXISTS idx_projects_sort_order ON projects(sort_order, created_at DESC);
