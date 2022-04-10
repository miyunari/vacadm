CREATE TABLE team (
    id UUID NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    created_at DATE NOT NULL,
    deleted_at DATE,
    updated_at DATE,
    PRIMARY KEY(id)
);

CREATE TABLE user (
    id UUID NOT NULL DEFAULT UUID(),
    parent_id UUID,
    team_id UUID,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at DATE NOT NULL DEFAULT time(),
    deleted_at DATE,
    updated_at DATE,
    PRIMARY KEY(id),
    FOREIGN KEY(team_id) REFERENCES team(id)
);

CREATE TABLE vaccation (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    approved_id UUID,
    `from` DATE NOT NULL,
    `to` DATE NOT NULL,
    created_at DATE NOT NULL,
    deleted_at DATE,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user(id),
    FOREIGN KEY(approved_id) REFERENCES user(id)
);

CREATE TABLE vaccation_request (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    `from` DATE NOT NULL,
    `to` DATE NOT NULL,
    created_at DATE NOT NULL,
    deleted_at DATE,
    updated_at DATE,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user(id)
);

CREATE TABLE vaccation_resource (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    yearlyDays INT NOT NULL,
    `from` DATE NOT NULL,
    `to` DATE NOT NULL,
    created_at DATE NOT NULL,
    deleted_at DATE,
    updated_at DATE,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user(id)
);