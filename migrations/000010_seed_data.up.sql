BEGIN;

INSERT INTO users (
    id,
    name,
    email,
    role,
    hashed_password
) VALUES (
    'f02de4cd-411e-4680-b99e-5d8485a78165',
    'admin',
    'admin@gmail.com',
    'admin',
    '$argon2id$v=19$m=65536,t=1,p=24$jIMk6lin0oTl32mjiOJeNQ$tkCAIMqcT37RmFOyhy8td2ITm8grv5BNOfQVtG7C5zk'
);

INSERT INTO address (
    id,
    user_id,
    first_name,
    last_name,
    email,
    phone,
    address1,
    city,
    state,
    zip
) VALUES (
    'a9e5cdb9-5e37-484b-909f-66f85847045e',
    'f02de4cd-411e-4680-b99e-5d8485a78165',
    'admin',
    'admin',
    'admin@gmail.com',
    '9256157821',
    'church street',
    'bangalore',
    'karnataka',
    '560078'
);

COMMIT;

BEGIN;

INSERT INTO product (
    id,
    title,
    description,
    category,
    price
) VALUES (
    'e72eabd7-b953-40f7-9d64-a221a6f00377',
    'Nike Jordan Air 1',
    'Nike Shoes',
    'Shoes',
    99.99
);

INSERT INTO inventory (
    id,
    product_id,
    small,
    medium,
    large,
    extra_large
) VALUES (
    'e630e71a-e7cd-4de9-95c6-c8e90542614d',
    'e72eabd7-b953-40f7-9d64-a221a6f00377',
    20,
    20,
    20,
    20
);

COMMIT;