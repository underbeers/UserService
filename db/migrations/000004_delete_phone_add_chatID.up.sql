ALTER TABLE user_contacts DROP COLUMN mobile_phone;
ALTER TABLE user_contacts DROP COLUMN show_phone;
ALTER TABLE user_contacts ADD COLUMN chat_id TEXT DEFAULT '';