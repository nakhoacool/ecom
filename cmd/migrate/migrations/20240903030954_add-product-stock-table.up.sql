CREATE TABLE IF NOT EXISTS product_stock (
    `product_id` INT UNSIGNED PRIMARY KEY NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);