-- Should trigger: UPDATE without WHERE updates all rows
UPDATE users SET active = false;

-- Should NOT trigger: UPDATE with WHERE
UPDATE users SET active = false WHERE id = 1;

-- Should NOT trigger: UPDATE with FROM and WHERE
UPDATE orders SET status = 'shipped' FROM shipments WHERE orders.id = shipments.order_id;
