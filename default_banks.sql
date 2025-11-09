PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
INSERT INTO banks VALUES(1,'VBank',NULL,'https://vbank.open.bankingapi.ru');
INSERT INTO banks VALUES(2,'ABank',NULL,'https://abank.open.bankingapi.ru');
INSERT INTO banks VALUES(3,'SBank',NULL,'https://sbank.open.bankingapi.ru');
COMMIT;
