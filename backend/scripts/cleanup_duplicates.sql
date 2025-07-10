-- 清理重复数据的SQL脚本

-- 删除重复的app_types数据（保留最新的）
WITH ranked_app_types AS (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY name ORDER BY created_at DESC) as rn
    FROM app_types
)
DELETE FROM app_types WHERE id IN (
    SELECT id FROM ranked_app_types WHERE rn > 1
);

-- 删除重复的api_keys数据（保留最新的）
WITH ranked_api_keys AS (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY name, provider ORDER BY created_at DESC) as rn
    FROM api_keys
)
DELETE FROM api_keys WHERE id IN (
    SELECT id FROM ranked_api_keys WHERE rn > 1
);

-- 删除重复的chat_models数据（保留最新的）
WITH ranked_chat_models AS (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY provider, value ORDER BY created_at DESC) as rn
    FROM chat_models
)
DELETE FROM chat_models WHERE id IN (
    SELECT id FROM ranked_chat_models WHERE rn > 1
);

-- 添加唯一约束（如果不存在）
ALTER TABLE app_types ADD CONSTRAINT IF NOT EXISTS app_types_name_unique UNIQUE (name);
ALTER TABLE api_keys ADD CONSTRAINT IF NOT EXISTS api_keys_name_provider_unique UNIQUE (name, provider);
ALTER TABLE chat_models ADD CONSTRAINT IF NOT EXISTS chat_models_provider_value_unique UNIQUE (provider, value);