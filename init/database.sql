CREATE TABLE test (
                      id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                      order jsonb NOT NULL
);