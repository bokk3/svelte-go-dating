-- Seed data for dating app
-- Note: password is 'password123' hashed with bcrypt

-- Insert admin user
INSERT INTO users (id, email, password_hash, email_verified_at, status, gdpr_consent, gdpr_consent_at) VALUES
('00000000-0000-0000-0000-000000000001', 'admin@dating-app.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW());

-- Insert diverse user profiles
INSERT INTO users (id, email, password_hash, email_verified_at, status, gdpr_consent, gdpr_consent_at) VALUES
('11111111-1111-1111-1111-111111111111', 'alice.johnson@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('22222222-2222-2222-2222-222222222222', 'mike.chen@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('33333333-3333-3333-3333-333333333333', 'sara.williams@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('44444444-4444-4444-4444-444444444444', 'david.rodriguez@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('55555555-5555-5555-5555-555555555555', 'emma.thompson@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('66666666-6666-6666-6666-666666666666', 'james.anderson@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('77777777-7777-7777-7777-777777777777', 'maria.garcia@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('88888888-8888-8888-8888-888888888888', 'alex.kim@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('99999999-9999-9999-9999-999999999999', 'jessica.brown@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW()),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ryan.taylor@email.com', '$2a$10$7qB7QXxc2vuTJDo3f5EetOkOepT1afNpwAuyzoyndZOXDD4EwlfHS', NOW(), 'active', TRUE, NOW());

-- Insert corresponding profiles
INSERT INTO profiles (user_id, display_name, bio, age, gender, interested_in, location_city, location_country, latitude, longitude, is_verified, is_premium) VALUES
('00000000-0000-0000-0000-000000000001', 'Admin', 'System Administrator', 30, 'other', '{"any"}', 'San Francisco', 'USA', 37.7749, -122.4194, TRUE, TRUE),

('11111111-1111-1111-1111-111111111111', 'Alice', 'Adventure seeker and coffee enthusiast. Love hiking, photography, and trying new restaurants. Looking for someone who shares my passion for exploring the world!', 28, 'female', '{"male"}', 'New York', 'USA', 40.7128, -74.0060, TRUE, FALSE),

('22222222-2222-2222-2222-222222222222', 'Mike', 'Software engineer by day, musician by night. I play guitar and love live music. Always up for a good conversation about technology, music, or life in general.', 32, 'male', '{"female"}', 'San Francisco', 'USA', 37.7749, -122.4194, TRUE, TRUE),

('33333333-3333-3333-3333-333333333333', 'Sara', 'Yoga instructor and wellness coach. Passionate about mindfulness, healthy living, and helping others find their inner peace. Looking for someone who values personal growth.', 26, 'female', '{"male", "female"}', 'Los Angeles', 'USA', 34.0522, -118.2437, TRUE, FALSE),

('44444444-4444-4444-4444-444444444444', 'David', 'Chef and food blogger. I believe the way to someone''s heart is through their stomach! Love cooking, traveling for food, and discovering hidden culinary gems.', 35, 'male', '{"female"}', 'Chicago', 'USA', 41.8781, -87.6298, FALSE, FALSE),

('55555555-5555-5555-5555-555555555555', 'Emma', 'Artist and creative soul. I paint, sketch, and love all forms of creative expression. Looking for someone who appreciates art and isn''t afraid to get a little paint on their hands!', 24, 'female', '{"male", "non-binary"}', 'Portland', 'USA', 45.5152, -122.6784, TRUE, FALSE),

('66666666-6666-6666-6666-666666666666', 'James', 'Personal trainer and fitness enthusiast. Love staying active, whether it''s at the gym, rock climbing, or playing sports. Looking for a workout partner and life partner!', 29, 'male', '{"female"}', 'Miami', 'USA', 25.7617, -80.1918, TRUE, TRUE),

('77777777-7777-7777-7777-777777777777', 'Maria', 'Teacher and lifelong learner. Passionate about education, reading, and making a difference in the world. Love quiet evenings with a good book and deep conversations.', 31, 'female', '{"male"}', 'Austin', 'USA', 30.2672, -97.7431, FALSE, FALSE),

('88888888-8888-8888-8888-888888888888', 'Alex', 'Non-binary data scientist who loves puzzles, board games, and analyzing interesting datasets. Looking for someone intelligent and curious about the world around them.', 27, 'non-binary', '{"any"}', 'Seattle', 'USA', 47.6062, -122.3321, TRUE, FALSE),

('99999999-9999-9999-9999-999999999999', 'Jessica', 'Marketing professional and travel enthusiast. I''ve been to 25 countries and counting! Love experiencing new cultures, trying exotic foods, and sharing adventures.', 30, 'female', '{"male"}', 'Denver', 'USA', 39.7392, -104.9903, TRUE, TRUE),

('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Ryan', 'Photographer and outdoor enthusiast. Spend my weekends hiking, camping, and capturing the perfect shot. Looking for someone who loves nature as much as I do.', 33, 'male', '{"female", "non-binary"}', 'Boulder', 'USA', 40.0150, -105.2705, FALSE, FALSE);

-- Insert sample photos
INSERT INTO photos (user_id, url, is_primary, display_order) VALUES
('11111111-1111-1111-1111-111111111111', 'https://images.unsplash.com/photo-1494790108755-2616b612bcae?w=400', TRUE, 1),
('11111111-1111-1111-1111-111111111111', 'https://images.unsplash.com/photo-1509909756405-be0199881695?w=400', FALSE, 2),

('22222222-2222-2222-2222-222222222222', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400', TRUE, 1),
('22222222-2222-2222-2222-222222222222', 'https://images.unsplash.com/photo-1560250097-0b93528c311a?w=400', FALSE, 2),

('33333333-3333-3333-3333-333333333333', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=400', TRUE, 1),
('33333333-3333-3333-3333-333333333333', 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=400', FALSE, 2),

('44444444-4444-4444-4444-444444444444', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400', TRUE, 1),

('55555555-5555-5555-5555-555555555555', 'https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?w=400', TRUE, 1),
('55555555-5555-5555-5555-555555555555', 'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=400', FALSE, 2),

('66666666-6666-6666-6666-666666666666', 'https://images.unsplash.com/photo-1552058544-f2b08422138a?w=400', TRUE, 1),

('77777777-7777-7777-7777-777777777777', 'https://images.unsplash.com/photo-1580489944761-15a19d654956?w=400', TRUE, 1),

('88888888-8888-8888-8888-888888888888', 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=400', TRUE, 1),

('99999999-9999-9999-9999-999999999999', 'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=400', TRUE, 1),
('99999999-9999-9999-9999-999999999999', 'https://images.unsplash.com/photo-1521119989659-a83eee488004?w=400', FALSE, 2),

('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400', TRUE, 1);

-- Create some sample matches between users
INSERT INTO matches (user1_id, user2_id, matched_at, is_active) VALUES
('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', NOW() - INTERVAL '2 days', TRUE),
('33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', NOW() - INTERVAL '1 day', TRUE),
('55555555-5555-5555-5555-555555555555', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '3 hours', TRUE);

-- Create some sample swipes
INSERT INTO swipes (swiper_id, swiped_id, liked, created_at) VALUES
-- Alice's swipes
('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', TRUE, NOW() - INTERVAL '2 days'),
('11111111-1111-1111-1111-111111111111', '44444444-4444-4444-4444-444444444444', TRUE, NOW() - INTERVAL '1 day'),
('11111111-1111-1111-1111-111111111111', '66666666-6666-6666-6666-666666666666', FALSE, NOW() - INTERVAL '1 day'),

-- Mike's swipes
('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', TRUE, NOW() - INTERVAL '2 days'),
('22222222-2222-2222-2222-222222222222', '33333333-3333-3333-3333-333333333333', TRUE, NOW() - INTERVAL '1 day'),
('22222222-2222-2222-2222-222222222222', '55555555-5555-5555-5555-555555555555', FALSE, NOW() - INTERVAL '12 hours'),

-- Sara's swipes
('33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', TRUE, NOW() - INTERVAL '1 day'),
('33333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222', FALSE, NOW() - INTERVAL '1 day'),

-- David's swipes
('44444444-4444-4444-4444-444444444444', '33333333-3333-3333-3333-333333333333', TRUE, NOW() - INTERVAL '1 day'),
('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', FALSE, NOW() - INTERVAL '1 day'),

-- Emma's swipes
('55555555-5555-5555-5555-555555555555', '66666666-6666-6666-6666-666666666666', TRUE, NOW() - INTERVAL '3 hours'),
('55555555-5555-5555-5555-555555555555', '88888888-8888-8888-8888-888888888888', TRUE, NOW() - INTERVAL '2 hours'),

-- James's swipes
('66666666-6666-6666-6666-666666666666', '55555555-5555-5555-5555-555555555555', TRUE, NOW() - INTERVAL '3 hours'),
('66666666-6666-6666-6666-666666666666', '99999999-9999-9999-9999-999999999999', TRUE, NOW() - INTERVAL '1 hour');

-- Sample messages between matched users
INSERT INTO messages (match_id, sender_id, message, created_at) VALUES
-- Alice and Mike conversation
((SELECT id FROM matches WHERE user1_id = '11111111-1111-1111-1111-111111111111' AND user2_id = '22222222-2222-2222-2222-222222222222'), '11111111-1111-1111-1111-111111111111', 'Hey! I saw you''re into music. What kind of shows do you usually go to?', NOW() - INTERVAL '2 days'),
((SELECT id FROM matches WHERE user1_id = '11111111-1111-1111-1111-111111111111' AND user2_id = '22222222-2222-2222-2222-222222222222'), '22222222-2222-2222-2222-222222222222', 'Hi Alice! I love indie rock and jazz. Just went to a great concert at the Fillmore last week. You?', NOW() - INTERVAL '2 days' + INTERVAL '30 minutes'),
((SELECT id FROM matches WHERE user1_id = '11111111-1111-1111-1111-111111111111' AND user2_id = '22222222-2222-2222-2222-222222222222'), '11111111-1111-1111-1111-111111111111', 'Nice! I''m more into alternative and folk. The Fillmore is amazing - such great acoustics!', NOW() - INTERVAL '1 day'),

-- Sara and David conversation
((SELECT id FROM matches WHERE user1_id = '33333333-3333-3333-3333-333333333333' AND user2_id = '44444444-4444-4444-4444-444444444444'), '33333333-3333-3333-3333-333333333333', 'I noticed you''re a chef! Any restaurant recommendations in the city?', NOW() - INTERVAL '1 day'),
((SELECT id FROM matches WHERE user1_id = '33333333-3333-3333-3333-333333333333' AND user2_id = '44444444-4444-4444-4444-444444444444'), '44444444-4444-4444-4444-444444444444', 'Hey Sara! For wellness-focused dining, I''d recommend Green Table - they have amazing plant-based options that even us chefs love!', NOW() - INTERVAL '23 hours'),

-- Emma and James conversation
((SELECT id FROM matches WHERE user1_id = '55555555-5555-5555-5555-555555555555' AND user2_id = '66666666-6666-6666-6666-666666666666'), '55555555-5555-5555-5555-555555555555', 'Love your outdoor photos! Do you have any favorite hiking spots?', NOW() - INTERVAL '3 hours'),
((SELECT id FROM matches WHERE user1_id = '55555555-5555-5555-5555-555555555555' AND user2_id = '66666666-6666-6666-6666-666666666666'), '66666666-6666-6666-6666-666666666666', 'Thanks Emma! Multnomah Falls is incredible, and I love the art scene in your area. Maybe we could explore both sometime?', NOW() - INTERVAL '2 hours');

-- Sample premium subscriptions
INSERT INTO subscriptions (user_id, plan_type, status, starts_at, ends_at) VALUES
('00000000-0000-0000-0000-000000000001', 'admin', 'active', NOW() - INTERVAL '1 year', NOW() + INTERVAL '1 year'),
('22222222-2222-2222-2222-222222222222', 'premium_monthly', 'active', NOW() - INTERVAL '15 days', NOW() + INTERVAL '15 days'),
('66666666-6666-6666-6666-666666666666', 'premium_yearly', 'active', NOW() - INTERVAL '60 days', NOW() + INTERVAL '305 days'),
('99999999-9999-9999-9999-999999999999', 'premium_monthly', 'active', NOW() - INTERVAL '5 days', NOW() + INTERVAL '25 days');

-- Update last_active timestamps to make profiles seem more realistic
UPDATE users SET last_active = NOW() - INTERVAL '5 minutes' WHERE id = '11111111-1111-1111-1111-111111111111';
UPDATE users SET last_active = NOW() - INTERVAL '2 hours' WHERE id = '22222222-2222-2222-2222-222222222222';
UPDATE users SET last_active = NOW() - INTERVAL '30 minutes' WHERE id = '33333333-3333-3333-3333-333333333333';
UPDATE users SET last_active = NOW() - INTERVAL '1 hour' WHERE id = '44444444-4444-4444-4444-444444444444';
UPDATE users SET last_active = NOW() - INTERVAL '15 minutes' WHERE id = '55555555-5555-5555-5555-555555555555';
UPDATE users SET last_active = NOW() - INTERVAL '3 hours' WHERE id = '66666666-6666-6666-6666-666666666666';
UPDATE users SET last_active = NOW() - INTERVAL '1 day' WHERE id = '77777777-7777-7777-7777-777777777777';
UPDATE users SET last_active = NOW() - INTERVAL '45 minutes' WHERE id = '88888888-8888-8888-8888-888888888888';
UPDATE users SET last_active = NOW() - INTERVAL '10 minutes' WHERE id = '99999999-9999-9999-9999-999999999999';
UPDATE users SET last_active = NOW() - INTERVAL '6 hours' WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- Add some verification for testing
UPDATE profiles SET is_verified = TRUE WHERE user_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '55555555-5555-5555-5555-555555555555',
    '66666666-6666-6666-6666-666666666666',
    '88888888-8888-8888-8888-888888888888',
    '99999999-9999-9999-9999-999999999999'
);