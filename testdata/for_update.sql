-- Should trigger: FOR UPDATE without SKIP LOCKED or NOWAIT
SELECT id FROM queue_items WHERE status = 'pending' ORDER BY id FOR UPDATE;

-- Should NOT trigger: FOR UPDATE SKIP LOCKED
SELECT id FROM queue_items WHERE status = 'pending' ORDER BY id FOR UPDATE SKIP LOCKED;

-- Should NOT trigger: FOR UPDATE NOWAIT
SELECT id FROM queue_items WHERE status = 'pending' ORDER BY id FOR UPDATE NOWAIT;

-- Should NOT trigger: FOR SHARE (not an exclusive lock)
SELECT id FROM users WHERE id = 1 FOR SHARE;
