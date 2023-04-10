ALTER TABLE user_contacts ADD COLUMN mobile_phone VARCHAR(255);
ALTER TABLE user_contacts ADD COLUMN show_phone BOOLEAN;
ALTER TABLE user_contacts DROP COLUMN chat_id;