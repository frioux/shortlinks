CREATE TABLE IF NOT EXISTS shortlinks (
        "from",
        "to",
        "deleted",
        "description",
        PRIMARY KEY ("from")
);

CREATE TABLE IF NOT EXISTS history (
        "from",
        "to",
        "when",
        "who",
        "description"
);

