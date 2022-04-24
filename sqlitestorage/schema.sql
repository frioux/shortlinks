CREATE TABLE IF NOT EXISTS shortlinks (
        "from",
        "to",
        PRIMARY KEY ("from")
);

CREATE TABLE IF NOT EXISTS history (
        "from",
        "to",
        "when"
);
