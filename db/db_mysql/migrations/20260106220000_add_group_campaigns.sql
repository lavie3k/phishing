-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `group_campaigns` (group_id bigint,campaign_id bigint );
CREATE INDEX `idx_group_campaigns_group_id` ON `group_campaigns` (group_id);
CREATE INDEX `idx_group_campaigns_campaign_id` ON `group_campaigns` (campaign_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `group_campaigns`;
