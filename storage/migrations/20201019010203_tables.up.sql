CREATE TABLE account
(
       id uuid NOT NULL,
       login text NOT NULL,
       pass text NOT NULL,
       created_at timestamptz default now(),
       CONSTRAINT account_pk PRIMARY KEY (id)
);

CREATE UNIQUE INDEX account_login_index ON account(login);

CREATE TABLE story
(
       id uuid NOT NULL,
       title text NOT NULL,
       kind text NOT NULL,
       body text NOT NULL,
       category text NOT NULL,
       score integer NOT NULL DEFAULT 0,
       views integer NOT NULL DEFAULT 0,
       created_by uuid NOT NULL,
       created_at timestamptz default now(),
       CONSTRAINT story_pk PRIMARY KEY (id)
);

CREATE INDEX story_author_index ON story(created_by);
CREATE INDEX story_created_index ON story(created_at);
CREATE INDEX story_category_index ON story(category);

CREATE TABLE comment
(
       id uuid NOT NULL,
       story_id uuid NOT NULL,
       body text NOT NULL,
       created_by uuid NOT NULL,
       created_at timestamptz default now(),
       CONSTRAINT comment_pk PRIMARY KEY (id)
);

CREATE INDEX comment_story_index ON comment(story_id);
CREATE INDEX comment_author_index ON comment(created_by);
CREATE INDEX comment_created_index ON comment(created_at);
