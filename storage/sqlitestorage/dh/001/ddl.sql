CREATE TABLE shortlinks (
        "from",
        "to",
        "deleted",
        "description",
        PRIMARY KEY ("from")
);

CREATE TABLE history (
        "from",
        "to",
        "when",
        "who",
        "description"
);
