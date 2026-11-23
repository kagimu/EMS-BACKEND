-- Run this script against db_whole to fix the legacy transaction_number values in qb_sync_queue
-- After the previous bugfix, transaction_number stores the true UUID of the transaction
-- but old QUEUED items still have the item_receipt_no string.

UPDATE qb_sync_queue sq
JOIN invoice_item_transactions iit ON iit.item_receipt_number = sq.transaction_number
SET sq.transaction_number = iit.number
WHERE sq.action_type = 'SALES_RECEIPT_ADD' AND LENGTH(sq.transaction_number) < 36;
