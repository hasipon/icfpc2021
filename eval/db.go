type Solution struct {
    id           text,
    json         text,
    problem_id   text,
    valid        text,
    dislike      real,
    dislike_s    text,
    use_bonus    text,
    unlock_bonus text,
    created_at   timestamp,
    updated_at   timestamp,
    PRIMARY KEY (id)
}
