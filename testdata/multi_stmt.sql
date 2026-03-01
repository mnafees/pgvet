-- Should trigger: two statements in one block
DELETE FROM task_slots WHERE task_id IN (SELECT task_id FROM expired);
DELETE FROM tasks WHERE task_id IN (SELECT task_id FROM expired);
