-- System accounts for special purposes
-- This should be run after the initial schema

-- System account for platform fees
INSERT INTO accounts (id, owner_id, account_type, currency, balance, tier, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000000',
    'system',
    'NGN',
    0,
    'tier2',
    'active'
) ON CONFLICT DO NOTHING;

-- Escrow master account
INSERT INTO accounts (id, owner_id, account_type, currency, balance, tier, status)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000000',
    'escrow',
    'NGN',
    0,
    'tier2',
    'active'
) ON CONFLICT DO NOTHING;