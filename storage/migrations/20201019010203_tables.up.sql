CREATE TABLE account
(
       id uuid NOT NULL,
       login text NOT NULL,
       pass text NOT NULL,
       name text,
       created_at timestamptz default now(),
       CONSTRAINT account_pk PRIMARY KEY (id)
);

CREATE TABLE session
(
       id serial NOT NULL,
       uid uuid NOT NULL,
       session text NOT NULL,
       created_at timestamptz default now(),
       CONSTRAINT session_pk PRIMARY KEY (id)
);

CREATE TABLE story
(
       id uuid NOT NULL,
       body text NOT NULL,
       score integer NOT NULL DEFAULT 0,
       views integer NOT NULL DEFAULT 0,
       kind text,
       created_at timestamptz default now(),
       CONSTRAINT story_pk PRIMARY KEY (id)
);

CREATE TABLE comment
(
       id uuid NOT NULL,
       story_id uuid NOT NULL,
       body text NOT NULL,
       created_at timestamptz default now(),
       CONSTRAINT comment_pk PRIMARY KEY (id)
);
