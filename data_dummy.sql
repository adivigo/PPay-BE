-- Dummy data for wallets table
INSERT INTO wallets (user_id, balance, is_deleted)
VALUES
(1, 1000.00, FALSE),
(2, 2500.00, FALSE);

-- Dummy data for transactions table
INSERT INTO transactions (user_id, amount, transaction_type, notes, is_deleted)
VALUES
(1, 200.00, 'transfer', 'Groceries', FALSE),
(2, 500.00, 'top_up', 'Salary', FALSE),
(1, 100.00, 'transfer', 'Utilities', FALSE),
(3, 50.00, 'top_up', 'Refund', TRUE);

-- Dummy data for payment_methods table
INSERT INTO payment_methods (name, tax, is_deleted)
VALUES
('Credit Card', 5.00, FALSE),
('Bank Transfer', 0.00, FALSE),
('E-Wallet', 1.50, FALSE),
('Cash', 0.00, TRUE);

-- Dummy data for transfer_transactions table
INSERT INTO transfer_transactions (transaction_id, target_user_id, is_deleted)
VALUES
(1, 2, FALSE),
(3, 4, FALSE);

-- Dummy data for topup_transactions table
INSERT INTO topup_transactions (transaction_id, payment_method_id, is_deleted)
VALUES
(2, 1, FALSE),
(4, 3, TRUE);
